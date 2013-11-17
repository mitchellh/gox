package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

func main() {
	// Call realMain so that defers work properly, since os.Exit won't
	// call defers.
	os.Exit(realMain())
}

func realMain() int {
	if _, err := exec.LookPath("go"); err != nil {
		fmt.Fprintf(os.Stderr, "go executable must be on the PATH\n")
		return 1
	}

	version, err := GoVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading Go version: %s", err)
		return 1
	}

	var outputTpl string
	var parallel int
	flag.StringVar(&outputTpl, "output", "{{.Dir}}_{{.OS}}_{{.Arch}}", "output path")
	flag.IntVar(&parallel, "parallel", -1, "parallelization factor")
	flag.Parse()

	// Determine the packages that we want to compile. We have to be sure
	// to turn any absolute paths into relative paths so that they work
	// properly with `go list`.
	packages := flag.Args()
	if len(packages) == 0 {
		packages = []string{"."}
	}

	// Determine what amount of parallelism we want
	if parallel <= 0 {
		parallel = runtime.NumCPU()
	}
	fmt.Printf("Number of parallel builds: %d\n", parallel)

	// Get the packages that are in the given paths
	mainDirs, err := GoMainDirs(packages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading packages: %s", err)
		return 1
	}

	// Determine the platforms we're building for
	platforms := SupportedPlatforms(version)

	// Build in parallel!
	var wg sync.WaitGroup
	semaphore := make(chan int, parallel)
	for _, platform := range platforms {
		for _, path := range mainDirs {
			// Start the goroutine that will do the actual build
			wg.Add(1)
			go func(path string, platform Platform) {
				defer wg.Done()
				semaphore <- 1
				fmt.Printf("--> %s: %s\n", platform.String(), path)
				if err := GoCrossCompile(path, platform, outputTpl); err != nil {
					fmt.Fprintf(os.Stderr, "%s error: %s", platform.String(), err)
				}
				<-semaphore
			}(path, platform)
		}
	}
	wg.Wait()

	return 0
}
