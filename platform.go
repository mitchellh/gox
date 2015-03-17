package main

import "fmt"

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

	Platforms_1_3 = append(Platforms_1_1, []Platform{
		{"dragonfly", "386"},
		{"dragonfly", "amd64"},
		{"solaris", "amd64"},
	}...)

	Platforms_1_4 = append(Platforms_1_3, []Platform{
		{"android", "arm"},
		{"plan9", "amd64"},
	}...)
)

// SupportedPlatforms returns the full list of supported platforms for
// the version of Go that is
func SupportedPlatforms(v string) []Platform {
	supportMap := map[string][]Platform{
		"go1.0": Platforms_1_0,
		"go1.1": Platforms_1_1,
		"go1.2": Platforms_1_1,
		"go1.3": Platforms_1_3,
		"go1.4": Platforms_1_4,
	}

	supported, ok := supportMap[v[0:5]]

	if ok == true {
		return supported
	} else {
		return Platforms_1_1

	}
}
