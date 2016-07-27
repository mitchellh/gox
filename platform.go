package main

import (
	"fmt"
	"strings"
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
	ARM     string
}

func PlatformFromString(os, arch string) Platform {
	if strings.HasPrefix(arch, "arm") && len(arch) >= 5 {
		return Platform{
			OS:   os,
			Arch: "arm",
			ARM:  arch[4:],
		}
	}
	return Platform{
		OS:   os,
		Arch: arch,
	}
}

func (p *Platform) String() string {
	return fmt.Sprintf("%s/%s", p.OS, p.GetArch())
}

func (p *Platform) GetArch() string {
	return fmt.Sprintf("%s%s", p.Arch, p.GetARMVersion())
}

func (p *Platform) GetARMVersion() string {
	if len(p.ARM) > 0 {
		return "v" + p.ARM
	}
	return ""
}

var (
	OsList = []string{
		"darwin",
		"dragonfly",
		"linux",
		"android",
		"solaris",
		"freebsd",
		"nacl",
		"netbsd",
		"openbsd",
		"plan9",
		"windows",
	}

	ArchList = []string{
		"386",
		"amd64",
		"amd64p32",
		"arm",
		"arm64",
		"mips64",
		"mips64le",
		"ppc64",
		"ppc64le",
	}

	Platforms_1_0 = []Platform{
		{OS: "darwin", Arch: "386", Default: true},
		{OS: "darwin", Arch: "amd64", Default: true},
		{OS: "linux", Arch: "386", Default: true},
		{OS: "linux", Arch: "amd64", Default: true},
		{OS: "linux", Arch: "arm", Default: true},
		{OS: "freebsd", Arch: "386", Default: true},
		{OS: "freebsd", Arch: "amd64", Default: true},
		{OS: "openbsd", Arch: "386", Default: true},
		{OS: "openbsd", Arch: "amd64", Default: true},
		{OS: "windows", Arch: "386", Default: true},
		{OS: "windows", Arch: "amd64", Default: true},
	}

	Platforms_1_1 = append(Platforms_1_0, []Platform{
		{OS: "freebsd", Arch: "arm", Default: true},
		{OS: "linux", Arch: "arm", Default: false, ARM: "5"},
		{OS: "linux", Arch: "arm", Default: false, ARM: "6"},
		{OS: "linux", Arch: "arm", Default: false, ARM: "7"},
		{OS: "netbsd", Arch: "386", Default: true},
		{OS: "netbsd", Arch: "amd64", Default: true},
		{OS: "netbsd", Arch: "arm", Default: true},
		{OS: "plan9", Arch: "386", Default: false},
	}...)

	Platforms_1_3 = append(Platforms_1_1, []Platform{
		{OS: "dragonfly", Arch: "386", Default: false},
		{OS: "dragonfly", Arch: "amd64", Default: false},
		{OS: "nacl", Arch: "amd64", Default: false},
		{OS: "nacl", Arch: "amd64p32", Default: false},
		{OS: "nacl", Arch: "arm", Default: false},
		{OS: "solaris", Arch: "amd64", Default: false},
	}...)

	Platforms_1_4 = append(Platforms_1_3, []Platform{
		{OS: "android", Arch: "arm", Default: false},
		{OS: "plan9", Arch: "amd64", Default: false},
	}...)

	Platforms_1_5 = append(Platforms_1_4, []Platform{
		{OS: "darwin", Arch: "arm", Default: false},
		{OS: "darwin", Arch: "arm64", Default: false},
		{OS: "linux", Arch: "arm64", Default: false},
		{OS: "linux", Arch: "ppc64", Default: false},
		{OS: "linux", Arch: "ppc64le", Default: false},
	}...)

	// Nothing changed from 1.5 to 1.6
	Platforms_1_6 = Platforms_1_5
)

// SupportedPlatforms returns the full list of supported platforms for
// the version of Go that is
func SupportedPlatforms(v string) []Platform {
	if strings.HasPrefix(v, "go1.0") {
		return Platforms_1_0
	} else if strings.HasPrefix(v, "go1.1") {
		return Platforms_1_1
	} else if strings.HasPrefix(v, "go1.3") {
		return Platforms_1_3
	} else if strings.HasPrefix(v, "go1.4") {
		return Platforms_1_4
	} else if strings.HasPrefix(v, "go1.5") {
		return Platforms_1_5
	} else if strings.HasPrefix(v, "go1.6") {
		return Platforms_1_6
	}

	// Assume latest
	return Platforms_1_6
}
