package main

import (
	"fmt"
	"log"
	"strings"

	version "github.com/hashicorp/go-version"
)

// Platform is a combination of OS/arch that can be built against.
type Platform struct {
	OS   string
	Arch string

	// Default, if true, will be included as a default build target
	// if no OS/arch is specified. We try to only set as a default popular
	// targets or targets that are generally useful. For example, Android
	// is not a default because it is quite rare that you're cross-compiling
	// something to Android AND something like Linux.
	Default bool
}

func (p *Platform) String() string {
	return fmt.Sprintf("%s/%s", p.OS, p.Arch)
}

var (
	OsList = []string{
		"darwin",
		"dragonfly",
		"freebsd",
		"linux",
		"netbsd",
		"openbsd",
		"plan9",
		"solaris",
		"windows",
	}

	ArchList = []string{
		"386",
		"amd64",
		"arm",
		"arm64",
		"ppc64",
		"ppc64le",
		"mips64",
		"mips64le",
		"mips",
		"mipsle",
	}

	Platforms_1_0 = []Platform{
		{"darwin", "386", true},
		{"darwin", "amd64", true},
		{"linux", "386", true},
		{"linux", "amd64", true},
		{"linux", "arm", true},
		{"freebsd", "386", true},
		{"freebsd", "amd64", true},
		{"openbsd", "386", true},
		{"openbsd", "amd64", true},
		{"windows", "386", true},
		{"windows", "amd64", true},
	}

	Platforms_1_1 = append(Platforms_1_0, []Platform{
		{"freebsd", "arm", true},
		{"netbsd", "386", true},
		{"netbsd", "amd64", true},
		{"netbsd", "arm", true},
		{"plan9", "386", false},
	}...)

	Platforms_1_3 = append(Platforms_1_1, []Platform{
		{"dragonfly", "386", false},
		{"dragonfly", "amd64", false},
		{"nacl", "amd64", false},
		{"nacl", "amd64p32", false},
		{"nacl", "arm", false},
		{"solaris", "amd64", false},
	}...)

	Platforms_1_4 = append(Platforms_1_3, []Platform{
		{"android", "arm", false},
		{"plan9", "amd64", false},
	}...)

	Platforms_1_5 = append(Platforms_1_4, []Platform{
		{"darwin", "arm", false},
		{"darwin", "arm64", false},
		{"linux", "arm64", false},
		{"linux", "ppc64", false},
		{"linux", "ppc64le", false},
	}...)

	Platforms_1_6 = append(Platforms_1_5, []Platform{
		{"android", "386", false},
		{"linux", "mips64", false},
		{"linux", "mips64le", false},
	}...)

	Platforms_1_7 = append(Platforms_1_6, []Platform{
		{"linux", "s390x", false},
		{"plan9", "arm", false},
		{"linux", "mips64", true},
		{"linux", "mips64le", true},
	}...)

	Platforms_1_8 = append(Platforms_1_7, []Platform{
		{"linux", "mips", true},
		{"linux", "mipsle", true},
	}...)
)

// SupportedPlatforms returns the full list of supported platforms for
// the version of Go that is
func SupportedPlatforms(v string) []Platform {
	// Use latest if we get an unexpected version string
	if !strings.HasPrefix(v, "go") {
		return Platforms_1_5
	}
	// go-version only cares about version numbers
	v = v[2:]

	current, err := version.NewVersion(v)
	if err != nil {
		log.Printf("Unable to parse current go version: %s\n%s", v, err.Error())

		// Default to 1.5
		return Platforms_1_5
	} else if strings.HasPrefix(v, "go1.6") {
		return Platforms_1_6
	} else if strings.HasPrefix(v, "go1.7") {
		return Platforms_1_7
	} else if strings.HasPrefix(v, "go1.8") {
		return Platforms_1_8
	}

	var platforms = []struct {
		constraint string
		plat       []Platform
	}{
		{"<= 1.0", Platforms_1_0},
		{">= 1.1, < 1.3", Platforms_1_1},
		{">= 1.3, < 1.4", Platforms_1_3},
		{">= 1.4, < 1.5", Platforms_1_4},
		{">= 1.5", Platforms_1_5},
	}

	for _, p := range platforms {
		constraints, err := version.NewConstraint(p.constraint)
		if err != nil {
			panic(err)
		}
		if constraints.Check(current) {
			return p.plat
		}
	}

	// Assume latest
	return Platforms_1_8
}
