package main

import (
	"flag"
	"fmt"
	"strings"
)

// PlatformFlag is a flag.Value (and flag.Getter) implementation that
// is used to track the os/arch flags on the command-line.
type PlatformFlag struct {
	OS     []string
	Arch   []string
	OSArch []Platform
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
	ignoreOSArch := make(map[Platform]struct{})
	includeOSArch := make(map[Platform]struct{})
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
	for _, v := range p.OSArch {
		if v.OS[0] == '!' {
			v = Platform{
				OS:   v.OS[1:],
				Arch: v.Arch,
			}

			ignoreOSArch[v] = struct{}{}
		} else {
			includeOSArch[v] = struct{}{}
		}
	}

	// We're building a list of new platforms, so build the list
	// based only on the configured OS/arch pairs.
	var prefilter []Platform = nil
	if len(includeOSArch) > 0 {
		prefilter = make([]Platform, 0, len(p.Arch)*len(p.OS)+len(includeOSArch))
		for k, _ := range includeOSArch {
			prefilter = append(prefilter, k)
		}
	}

	if len(includeOS) > 0 && len(includeArch) > 0 {
		// Build up the list of prefiltered by what is specified
		if prefilter == nil {
			prefilter = make([]Platform, 0, len(p.Arch)*len(p.OS))
		}

		for _, os := range p.OS {
			if _, ok := includeOS[os]; !ok {
				continue
			}

			for _, arch := range p.Arch {
				if _, ok := includeArch[arch]; !ok {
					continue
				}

				prefilter = append(prefilter, Platform{os, arch})
			}
		}
	}

	if prefilter != nil {
		// Remove any that aren't supported
		result := make([]Platform, 0, len(prefilter))
		for _, pending := range prefilter {
			found := false
			for _, platform := range supported {
				if pending == platform {
					found = true
					break
				}
			}

			if found {
				result = append(result, pending)
			}
		}

		prefilter = result
	}

	if prefilter == nil {
		prefilter = supported
	}

	// Go through each default platform and filter out the bad ones
	result := make([]Platform, 0, len(prefilter))
	for _, platform := range prefilter {
		if len(ignoreOSArch) > 0 {
			if _, ok := ignoreOSArch[platform]; ok {
				continue
			}
		}

		// We only want to check the components (OS and Arch) if we didn't
		// specifically ask to include it via the osarch.
		checkComponents := true
		if len(includeOSArch) > 0 {
			if _, ok := includeOSArch[platform]; ok {
				checkComponents = false
			}
		}

		if checkComponents {
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
		}

		result = append(result, platform)
	}

	return result
}

// ArchFlagValue returns a flag.Value that can be used with the flag
// package to collect the arches for the flag.
func (p *PlatformFlag) ArchFlagValue() flag.Value {
	return (*appendStringValue)(&p.Arch)
}

// OSFlagValue returns a flag.Value that can be used with the flag
// package to collect the operating systems for the flag.
func (p *PlatformFlag) OSFlagValue() flag.Value {
	return (*appendStringValue)(&p.OS)
}

// OSArchFlagValue returns a flag.Value that can be used with the flag
// package to collect complete os and arch pairs for the flag.
func (p *PlatformFlag) OSArchFlagValue() flag.Value {
	return (*appendPlatformValue)(&p.OSArch)
}

// appendPlatformValue is a flag.Value that appends a full platform (os/arch)
// to a list where the values from space-separated lines. This is used to
// satisfy the -osarch flag.
type appendPlatformValue []Platform

func (s *appendPlatformValue) String() string {
	return ""
}

func (s *appendPlatformValue) Set(value string) error {
	if *s == nil {
		*s = make([]Platform, 0, 1)
	}

	if value == "" {
		return nil
	}

	for _, v := range strings.Split(value, " ") {
		parts := strings.Split(v, "/")
		if len(parts) != 2 {
			return fmt.Errorf(
				"Invalid platform syntax: %s should be os/arch", v)
		}

		platform := Platform{
			OS:   strings.ToLower(parts[0]),
			Arch: strings.ToLower(parts[1]),
		}

		s.appendIfMissing(&platform)
	}

	return nil
}

func (s *appendPlatformValue) appendIfMissing(value *Platform) {
	for _, existing := range *s {
		if existing == *value {
			return
		}
	}

	*s = append(*s, *value)
}

// appendStringValue is a flag.Value that appends values to the list,
// where the values come from space-separated lines. This is used to
// satisfy the -os="windows linux" flag to become []string{"windows", "linux"}
type appendStringValue []string

func (s *appendStringValue) String() string {
	return strings.Join(*s, " ")
}

func (s *appendStringValue) Set(value string) error {
	if *s == nil {
		*s = make([]string, 0, 1)
	}

	for _, v := range strings.Split(value, " ") {
		if v != "" {
			s.appendIfMissing(strings.ToLower(v))
		}
	}

	return nil
}

func (s *appendStringValue) appendIfMissing(value string) {
	for _, existing := range *s {
		if existing == value {
			return
		}
	}

	*s = append(*s, value)
}
