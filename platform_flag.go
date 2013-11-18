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
func (p *PlatformFlag) Platforms(supported []Platform) []Platform {
	// NOTE: Reading this method alone is a bit hard to understand. It
	// is much easier to understand this method if you pair this with the
	// table of test cases it has.

	// Build a list of OS and archs NOT to build
	ignoreArch := make(map[string]struct{})
	includeArch := make(map[string]struct{})
	ignoreOS := make(map[string]struct{})
	includeOS := make(map[string]struct{})
	for _, v := range p.Arch {
		if v[0] == '!' {
			ignoreArch[v[1:]] = struct{}{}
		} else {
			includeArch[v] = struct{}{}
		}
	}
	for _, v := range p.OS {
		if v[0] == '!' {
			ignoreOS[v[1:]] = struct{}{}
		} else {
			includeOS[v] = struct{}{}
		}
	}

	// We're building a list of new platforms, so build the list
	// based only on the configured OS/arch pairs.
	var prefilter []Platform = supported
	if len(includeOS) > 0 && len(includeArch) > 0 {
		// Build up the list of prefiltered by what is specified
		pendings := make([]Platform, 0, len(p.Arch)*len(p.OS))
		for _, os := range p.OS {
			if _, ok := includeOS[os]; !ok {
				continue
			}

			for _, arch := range p.Arch {
				if _, ok := includeArch[arch]; !ok {
					continue
				}

				pendings = append(pendings, Platform{os, arch})
			}
		}

		// Remove any that aren't supported
		prefilter = make([]Platform, 0, len(pendings))
		for _, pending := range pendings {
			found := false
			for _, platform := range supported {
				if pending == platform {
					found = true
					break
				}
			}

			if found {
				prefilter = append(prefilter, pending)
			}
		}
	}

	// Go through each default platform and filter out the bad ones
	result := make([]Platform, 0, len(prefilter))
	for _, platform := range prefilter {
		if len(ignoreArch) > 0 {
			if _, ok := ignoreArch[platform.Arch]; ok {
				continue
			}
		}
		if len(ignoreOS) > 0 {
			if _, ok := ignoreOS[platform.OS]; ok {
				continue
			}
		}
		if len(includeArch) > 0 {
			if _, ok := includeArch[platform.Arch]; !ok {
				continue
			}
		}
		if len(includeOS) > 0 {
			if _, ok := includeOS[platform.OS]; !ok {
				continue
			}
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
