package main

import (
	"fmt"
	"strings"
)

// Platform is a combination of OS/arch that can be built against.
type Platform struct {
	OS   string
	Arch string
}

func (p *Platform) String() string {
	return fmt.Sprintf("%s/%s", p.OS, p.Arch)
}

var (
	OsList = []string{
		"darwin",
		"linux",
		"freebsd",
		"netbsd",
		"openbsd",
		"plan9",
		"windows",
	}

	ArchList = []string{
		"386",
		"amd64",
		"arm",
	}

	Platforms_1_0 = []Platform{
		{"darwin", "386"},
		{"darwin", "amd64"},
		{"linux", "386"},
		{"linux", "amd64"},
		{"linux", "arm"},
		{"freebsd", "386"},
		{"freebsd", "amd64"},
		{"openbsd", "386"},
		{"openbsd", "amd64"},
		{"windows", "386"},
		{"windows", "amd64"},
	}

	Platforms_1_1 = append(Platforms_1_0, []Platform{
		{"freebsd", "arm"},
		{"netbsd", "386"},
		{"netbsd", "amd64"},
		{"netbsd", "arm"},
		{"plan9", "386"},
	}...)
)

// SupportedPlatforms returns the full list of supported platforms for
// the version of Go that is
func SupportedPlatforms(v string) []Platform {
	if strings.HasPrefix(v, "go1.0") {
		return Platforms_1_0
	}

	return Platforms_1_1
}
