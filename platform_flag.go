package main

import (
	"flag"
	"strings"
)

// PlatformFlag is a flag.Value (and flag.Getter) implementation that
// is used to track the os/arch flags on the command-line.
type PlatformFlag struct {
	OS   []string
	Arch []string
}

// Platforms returns the list of platforms that were set by this flag.
// The default set of platforms must be passed in.
func (p *PlatformFlag) Platforms(def []Platform) []Platform {
	// NOTE: Reading this method alone is a bit hard to understand. It
	// is much easier to understand this method if you pair this with the
	// table of test cases it has.

	// Determine if we're building off of a new list or if we're
	// building a brand new list.
	isNew := false
	for _, v := range p.Arch {
		if v[0] != '-' {
			isNew = true
			break
		}
	}
	for _, v := range p.OS {
		if v[0] != '-' {
			isNew = true
			break
		}
	}

	// Build a list of OS and archs NOT to build
	ignoreArch := make(map[string]struct{})
	ignoreOS := make(map[string]struct{})
	for _, v := range p.Arch {
		ignoreArch[v[1:]] = struct{}{}
	}
	for _, v := range p.OS {
		ignoreOS[v[1:]] = struct{}{}
	}

	if isNew {
		// We're building a list of new platforms, so build the list
		// based only on the configured OS/arch pairs.
		def = make([]Platform, 0, len(p.Arch)*len(p.OS))
		for _, os := range p.OS {
			if os[0] == '-' {
				continue
			}

			for _, arch := range p.Arch {
				if arch[0] == '-' {
					continue
				}

				def = append(def, Platform{os, arch})
			}
		}
	}

	// Go through each default platform and filter out the bad ones
	result := make([]Platform, 0, len(def))
	for _, platform := range def {
		if _, ok := ignoreArch[platform.Arch]; ok {
			continue
		}
		if _, ok := ignoreOS[platform.OS]; ok {
			continue
		}

		result = append(result, platform)
	}

	return result
}

// ArchFlagValue returns a flag.Value that can be used with the flag
// package to collect the arches for the flag.
func (p *PlatformFlag) ArchFlagValue() flag.Value {
	return (*appendPlatformValue)(&p.Arch)
}

// OSFlagValue returns a flag.Value that can be used with the flag
// package to collect the operating systems for the flag.
func (p *PlatformFlag) OSFlagValue() flag.Value {
	return (*appendPlatformValue)(&p.OS)
}

// appendPlatformValue is a flag.Value that appends values to the list,
// where the values come from space-separated lines. This is used to
// satisfy the -os="windows linux" flag to become []string{"windows", "linux"}
type appendPlatformValue []string

func (s *appendPlatformValue) String() string {
	return strings.Join(*s, " ")
}

func (s *appendPlatformValue) Set(value string) error {
	if *s == nil {
		*s = make([]string, 0, 1)
	}

	for _, v := range strings.Split(value, " ") {
		s.appendIfMissing(v)
	}

	return nil
}

func (s *appendPlatformValue) appendIfMissing(value string) {
	for _, existing := range *s {
		if existing == value {
			return
		}
	}

	*s = append(*s, value)
}
