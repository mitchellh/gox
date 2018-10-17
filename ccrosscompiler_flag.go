package main

import (
	"fmt"
	"strings"
)

const (
	osArchPos    = 0
	cCompilerPos = 1
)

// CCrossCompilerFlag stores the platforms and their selected C cross-compilers.
type CCrossCompilerFlag struct {
	selected map[string]string
	raw      string
}

// Set parses the input comma-separated list of platforms and compilers.
// Valid format: "linux/arm=arm-linux-gnueabi-gcc-6"
func (c *CCrossCompilerFlag) Set(s string) error {
	pairs := strings.Split(s, ",")
	c.selected = make(map[string]string, len(pairs))

	for _, pair := range pairs {
		elems := strings.Split(pair, "=")
		if len(elems) != 2 {
			return fmt.Errorf("invalid format: requires two elements separated by =")
		}
		c.selected[elems[osArchPos]] = elems[cCompilerPos]
	}

	c.raw = s
	return nil
}

func (c *CCrossCompilerFlag) String() string {
	return c.raw
}

func (c *CCrossCompilerFlag) Get() map[string]string {
	return c.selected
}
