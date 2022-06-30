package main

import (
	"strings"
	"testing"
)

func TestGoVersion(t *testing.T) {
	v, err := GoVersion()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	acceptable := []string{
		"devel",
		"go1.0",
		"go1.1",
		"go1.2",
		"go1.3",
		"go1.4",
		"go1.5",
		"go1.6",
		"go1.7",
		"go1.8",
		"go1.9",
		"go1.10",
		"go1.11",
		"go1.12",
		"go1.13",
		"go1.14",
		"go1.15",
		"go1.16",
		"go1.17",
		"go1.18",
	}
	found := false
	for _, expected := range acceptable {
		if strings.HasPrefix(v, expected) {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("bad: %#v", v)
	}
}
