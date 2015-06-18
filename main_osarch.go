package main

import (
	"fmt"
)

func mainListOSArch(version string) int {
	for _, p := range SupportedPlatforms(version) {
		fmt.Printf("%s\t(default: %v)\n", p.String(), p.Default)
	}

	return 0
}
