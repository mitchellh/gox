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
	var outputTpl string
	var parallel int
	var platformFlag PlatformFlag
	flags := flag.NewFlagSet("gox", flag.ExitOnError)
	flags.Usage = func() { printUsage() }
	flags.Var(platformFlag.ArchFlagValue(), "arch", "arch to build for or skip")
	flags.Var(platformFlag.OSFlagValue(), "os", "os to build for or skip")
	flags.StringVar(&outputTpl, "output", "{{.Dir}}_{{.OS}}_{{.Arch}}", "output path")
	flags.IntVar(&parallel, "parallel", -1, "parallelization factor")
	if err := flags.Parse(os.Args[1:]); err != nil {
		flags.Usage()
		return 1
	}

	if _, err := exec.LookPath("go"); err != nil {
		fmt.Fprintf(os.Stderr, "go executable must be on the PATH\n")
		return 1
	}

	version, err := GoVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading Go version: %s", err)
		return 1
	}

	// Determine the packages that we want to compile. Default to the
	// current directory if none are specified.
	packages := flags.Args()
	if len(packages) == 0 {
		packages = []string{"."}
	}

	// Determine what amount of parallelism we want Default to the current
	// number of CPUs is <= 0 is specified.
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
	platforms := platformFlag.Platforms(SupportedPlatforms(version))

	// Build in parallel!
	var errorLock sync.Mutex
	var wg sync.WaitGroup
	errors := make([]string, 0)
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
					errorLock.Lock()
					defer errorLock.Unlock()
					errors = append(errors,
						fmt.Sprintf("%s error: %s", platform.String(), err))
				}
				<-semaphore
			}(path, platform)
		}
	}
	wg.Wait()

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "\n%d errors occurred:\n", len(errors))
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "--> %s\n", err)
		}
		return 1
	}

	return 0
}

func printUsage() {
	fmt.Fprintf(os.Stderr, helpText)
}

const helpText = `Usage: gox [options] [packages]

  Gox cross-compiles Go applications in parallel.

  If no specific operating systes or architectures are specified, Gox
  will build for all pairs supported by your version of Go.

Options:

  -arch=""            Space-separated list of architectures to build for.
  -os=""              Space-separated list of operating systems to build for.
  -output="foo"       Output path template. See below for more info.
  -parallel=-1        Amount of parallelism, defaults to number of CPUs.

Output path template:

  The output path for the compiled binaries is specified with the
  "-output" flag. The value is a string that is a Go text template.
  The default value is "{{.Dir}}_{{.OS}}_{{.Arch}}". The variables and
  their values should be self-explanatory.

Platforms (OS/Arch):

  The operating systems and architectures to cross-compile for may be
  specified with the "-arch" and "-os" flags. These are space separated lists
  of valid GOOS/GOARCH values to build for, respectively. You may prefix an
  OS or Arch with "-" to negate and not build for that platform. If the list
  is made up of only negations, then the negations will come from the default
  list.

`
