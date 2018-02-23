package main

import (
	"reflect"
	"testing"
)

func TestSupportedPlatforms(t *testing.T) {
	var ps []Platform

	ps = SupportedPlatforms("go1.0")
	if !reflect.DeepEqual(ps, Platforms_1_0) {
		t.Fatalf("bad: %#v", ps)
	}

	ps = SupportedPlatforms("go1.1")
	if !reflect.DeepEqual(ps, Platforms_1_1) {
		t.Fatalf("bad: %#v", ps)
	}

	ps = SupportedPlatforms("go1.2")
	if !reflect.DeepEqual(ps, Platforms_1_1) {
		t.Fatalf("bad: %#v", ps)
	}

	ps = SupportedPlatforms("go1.3")
	if !reflect.DeepEqual(ps, Platforms_1_3) {
		t.Fatalf("bad: %#v", ps)
	}

	ps = SupportedPlatforms("go1.4")
	if !reflect.DeepEqual(ps, Platforms_1_4) {
		t.Fatalf("bad: %#v", ps)
	}

	ps = SupportedPlatforms("go1.5")
	if !reflect.DeepEqual(ps, Platforms_1_5) {
		t.Fatalf("bad: %#v", ps)
	}

	ps = SupportedPlatforms("go1.6")
	if !reflect.DeepEqual(ps, Platforms_1_6) {
		t.Fatalf("bad: %#v", ps)
	}

	ps = SupportedPlatforms("go1.7")
	if !reflect.DeepEqual(ps, Platforms_1_7) {
		t.Fatalf("bad: %#v", ps)
	}

	ps = SupportedPlatforms("go1.8")
	if !reflect.DeepEqual(ps, Platforms_1_8) {
		t.Fatalf("bad: %#v", ps)
	}

	ps = SupportedPlatforms("go1.9")
	if !reflect.DeepEqual(ps, Platforms_1_9) {
		t.Fatalf("bad: %#v", ps)
	}

	ps = SupportedPlatforms("go1.10")
	if !reflect.DeepEqual(ps, Platforms_1_10) {
		t.Fatalf("bad: %#v", ps)
	}

	// Unknown
	ps = SupportedPlatforms("foo")
	if !reflect.DeepEqual(ps, PlatformsLatest) {
		t.Fatalf("bad: %#v", ps)
	}
}

func TestMIPS(t *testing.T) {
	g16 := SupportedPlatforms("go1.6")
	for _, p := range g16 {
		if p.Arch == "mips64" && p.Default {
			t.Fatal("mips64 should not be default for 1.6")
		}
	}

	g17 := SupportedPlatforms("go1.7")
	for _, p := range g17 {
		if p.Arch == "mips64" && !p.Default {
			t.Fatal("mips64 should be default for 1.7")
		}
	}

}
