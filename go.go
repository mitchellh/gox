package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type OutputTemplateData struct {
	Dir  string
	OS   string
	Arch string
}

// GoCrossCompile
func GoCrossCompile(packagePath string, platform Platform, outputTpl string, ldflags string) error {
	env := append(os.Environ(),
		"GOOS="+platform.OS,
		"GOARCH="+platform.Arch)

	var outputPath bytes.Buffer
	tpl, err := template.New("output").Parse(outputTpl)
	if err != nil {
		return err
	}
	tplData := OutputTemplateData{
		// If compiling a file, use that file's name without extension.
		// Otherwise use the package name
		Dir:  filepath.Base(strings.TrimSuffix(packagePath, ".go")),
		OS:   platform.OS,
		Arch: platform.Arch,
	}
	if err := tpl.Execute(&outputPath, &tplData); err != nil {
		return nil
	}

	if platform.OS == "windows" {
		outputPath.WriteString(".exe")
	}

	// Determine the full path to the output so that we can change our
	// working directory when executing go build.
	outputPathReal := outputPath.String()
	outputPathReal, err = filepath.Abs(outputPathReal)
	if err != nil {
		return err
	}

	// Go prefixes the import directory with '_' when it is outside
	// the GOPATH.For this, we just drop it since we move to that
	// directory to build.  Only do this for packages, not files.
	chdir := ""
	if packagePath[0] == '_' && filepath.Ext(packagePath) != ".go" {
		chdir = packagePath[1:]
		packagePath = ""
	}

	_, err = execGo(env, chdir, "build",
		"-ldflags", ldflags,
		"-o", outputPathReal,
		packagePath)
	return err
}

// GoMainDirs returns the file paths to the packages that are "main"
// packages, from the list of packages given. The list of packages can
// include relative paths, the special "..." Go keyword, etc.
func GoMainDirs(compileThese []string) ([]string, error) {
	// Separate the packages from the files
	files := make([]string, 0)
	packages := make([]string, 0)
	for _, p := range compileThese {
		if filepath.Ext(p) == ".go" {
			files = append(files, p)
		} else {
			packages = append(packages, p)
		}
	}

	results := make([]string, 0)
	output := ""

	if len(packages) != 0 {
		args := make([]string, 0, len(packages)+3)
		args = append(args, "list", "-f", "{{.Name}}|{{.ImportPath}}")
		args = append(args, packages...)

		var err error
		output, err = execGo(nil, "", args...)
		if err != nil {
			return nil, err
		}
	}

	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			log.Printf("Bad line reading packages: %s", line)
			continue
		}

		if parts[0] == "main" {
			results = append(results, parts[1])
		}
	}

	results = append(results, files...)

	return results, nil
}

// GoRoot returns the GOROOT value for the compiled `go` binary.
func GoRoot() (string, error) {
	output, err := execGo(nil, "", "env", "GOROOT")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}

// GoVersion reads the version of `go` that is on the PATH. This is done
// instead of `runtime.Version()` because it is possible to run gox against
// another Go version.
func GoVersion() (string, error) {
	// NOTE: We use `go run` instead of `go version` because the output
	// of `go version` might change whereas the source is guaranteed to run
	// for some time thanks to Go's compatibility guarantee.

	td, err := ioutil.TempDir("", "gox")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(td)

	// Write the source code for the program that will generate the version
	sourcePath := filepath.Join(td, "version.go")
	if err := ioutil.WriteFile(sourcePath, []byte(versionSource), 0644); err != nil {
		return "", err
	}

	// Execute and read the version, which will be the only thing on stdout.
	return execGo(nil, "", "run", sourcePath)
}

func execGo(env []string, dir string, args ...string) (string, error) {
	var stderr, stdout bytes.Buffer
	cmd := exec.Command("go", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if env != nil {
		cmd.Env = env
	}
	if dir != "" {
		cmd.Dir = dir
	}
	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("%s\nStderr: %s", err, stderr.String())
		return "", err
	}

	return stdout.String(), nil
}

const versionSource = `package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Print(runtime.Version())
}`
