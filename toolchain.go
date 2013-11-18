package main

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/iochan"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// The "main" method for when the toolchain build is requested.
func mainBuildToolchain(parallel int, platformFlag PlatformFlag, verbose bool) int {
	if _, err := exec.LookPath("go"); err != nil {
		fmt.Fprintf(os.Stderr, "You must have Go already built for your native platform\n")
		fmt.Fprintf(os.Stderr, "and the `go` binary on the PATH to build toolchains.\n")
		return 1
	}

	version, err := GoVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading Go version: %s", err)
		return 1
	}

	root := os.Getenv("GOROOT")
	if root == "" {
		fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(gorootErrorText))
		return 1
	}

	if verbose {
		fmt.Println("Verbose mode enabled. Output from building each toolchain will be")
		fmt.Println("outputted to stdout as they are built.\n")
	}

	// Determine the platforms we're building the toolchain for.
	platforms := platformFlag.Platforms(SupportedPlatforms(version))

	// Tell the user how much parallelization we will use
	fmt.Printf("Number of parallel builds: %d\n\n", parallel)
	var errorLock sync.Mutex
	var wg sync.WaitGroup
	errs := make([]error, 0)
	semaphore := make(chan int, parallel)
	for _, platform := range platforms {
		wg.Add(1)
		go func() {
			err := buildToolchain(&wg, semaphore, platform, verbose)
			if err != nil {
				errorLock.Lock()
				defer errorLock.Unlock()
				errs = append(errs, fmt.Errorf("%s: %s", platform.String(), err))
			}
		}()
	}
	wg.Wait()

	if len(errs) > 0 {
		fmt.Fprintf(os.Stderr, "\n%d errors occurred:\n", len(errs))
		for _, err := range errs {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		return 1
	}

	return 0
}

func buildToolchain(wg *sync.WaitGroup, semaphore chan int, platform Platform, verbose bool) error {
	defer wg.Done()
	semaphore <- 1
	defer func() { <-semaphore }()
	fmt.Printf("--> Toolchain: %s\n", platform.String())

	scriptName := "make.bash"
	if runtime.GOOS == "windows" {
		scriptName = "make.bat"
	}

	var stderr bytes.Buffer
	scriptDir := filepath.Join(os.Getenv("GOROOT"), "src")
	scriptPath := filepath.Join(scriptDir, scriptName)
	cmd := exec.Command(scriptPath, "--no-clean")
	cmd.Dir = scriptDir
	cmd.Env = append(os.Environ(),
		"GOARCH="+platform.Arch,
		"GOOS="+platform.OS)
	cmd.Stderr = &stderr

	if verbose {
		// In verbose mode, we output all stdout to the console.
		r, w := io.Pipe()
		cmd.Stdout = w
		go func() {
			for line := range iochan.DelimReader(r, '\n') {
				fmt.Printf("%s: %s", platform.String(), line)
			}
		}()
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error building '%s': %s", platform.String(), err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Error building '%s': %s",
			platform.String(), stderr.String())
	}

	return nil
}

const gorootErrorText string = `
You must set GOROOT to build the cross-compile toolchain. GOROOT must point
to the directory containing the checkout of the Go source code. This only
needs to be set while building the toolchain for cross-compilation with gox,
and doesn't need to be set when using gox otherwise.

Note that you probably should NOT set this value globally. Read the blog
post below for more information on why that is:

http://dave.cheney.net/2013/06/14/you-dont-need-to-set-goroot-really
`
