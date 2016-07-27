package main

import (
	"flag"
	"reflect"
	"testing"
)

func TestPlatformFlagPlatforms(t *testing.T) {
	cases := []struct {
		OS        []string
		Arch      []string
		OSArch    []Platform
		Supported []Platform
		Result    []Platform
	}{
		// Building a new list of platforms
		{
			[]string{"foo", "bar"},
			[]string{"baz"},
			[]Platform{},
			[]Platform{
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "bar", Arch: "baz", Default: true},
				{OS: "boo", Arch: "bop", Default: true},
			},
			[]Platform{
				{OS: "foo", Arch: "baz", Default: false},
				{OS: "bar", Arch: "baz", Default: false},
			},
		},

		// Skipping platforms
		{
			[]string{"!foo"},
			[]string{},
			[]Platform{},
			[]Platform{
				{OS: "foo", Arch: "bar", Default: true},
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "bar", Arch: "bar", Default: true},
			},
			[]Platform{
				{OS: "bar", Arch: "bar", Default: false},
			},
		},

		// Specifying only an OS
		{
			[]string{"foo"},
			[]string{},
			[]Platform{},
			[]Platform{
				{OS: "foo", Arch: "bar", Default: true},
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "bar", Arch: "bar", Default: true},
			},
			[]Platform{
				{OS: "foo", Arch: "bar", Default: false},
				{OS: "foo", Arch: "baz", Default: false},
			},
		},

		// Building a new list, but with some skips
		{
			[]string{"foo", "bar", "!foo"},
			[]string{"baz"},
			[]Platform{},
			[]Platform{
				{OS: "foo", Arch: "bar", Default: true},
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "bar", Arch: "baz", Default: true},
				{OS: "baz", Arch: "bar", Default: true},
			},
			[]Platform{
				{OS: "bar", Arch: "baz", Default: false},
			},
		},

		// Unsupported pairs
		{
			[]string{"foo", "bar"},
			[]string{"baz"},
			[]Platform{},
			[]Platform{
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "bar", Arch: "what", Default: true},
			},
			[]Platform{
				{OS: "foo", Arch: "baz", Default: false},
			},
		},

		// OSArch basic
		{
			[]string{},
			[]string{},
			[]Platform{
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "foo", Arch: "bar", Default: true},
			},
			[]Platform{
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "bar", Arch: "what", Default: true},
			},
			[]Platform{
				{OS: "foo", Arch: "baz", Default: false},
			},
		},

		// Negative OSArch
		{
			[]string{},
			[]string{},
			[]Platform{
				{OS: "!foo", Arch: "baz", Default: true},
			},
			[]Platform{
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "bar", Arch: "what", Default: true},
			},
			[]Platform{
				{OS: "bar", Arch: "what", Default: false},
			},
		},

		// Mix it all
		{
			[]string{"foo", "bar"},
			[]string{"bar"},
			[]Platform{
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "!bar", Arch: "bar", Default: true},
			},
			[]Platform{
				{OS: "foo", Arch: "bar", Default: true},
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "bar", Arch: "bar", Default: true},
			},
			[]Platform{
				{OS: "foo", Arch: "baz", Default: false},
				{OS: "foo", Arch: "bar", Default: false},
			},
		},

		// Ignores non-default
		{
			[]string{},
			[]string{},
			[]Platform{},
			[]Platform{
				{OS: "foo", Arch: "bar", Default: true},
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "!bar", Arch: "bar", Default: false},
			},
			[]Platform{
				{OS: "foo", Arch: "bar", Default: false},
				{OS: "foo", Arch: "baz", Default: false},
			},
		},

		// Adds non-default by OS
		{
			[]string{"bar"},
			[]string{},
			[]Platform{},
			[]Platform{
				{OS: "foo", Arch: "bar", Default: true},
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "bar", Arch: "bar", Default: false},
			},
			[]Platform{
				{OS: "bar", Arch: "bar", Default: false},
			},
		},

		// Adds non-default by both
		{
			[]string{"bar"},
			[]string{"bar"},
			[]Platform{},
			[]Platform{
				{OS: "foo", Arch: "bar", Default: true},
				{OS: "foo", Arch: "baz", Default: true},
				{OS: "bar", Arch: "bar", Default: false},
			},
			[]Platform{
				{OS: "bar", Arch: "bar", Default: false},
			},
		},
	}

	for i, tc := range cases {
		f := PlatformFlag{
			OS:     tc.OS,
			Arch:   tc.Arch,
			OSArch: tc.OSArch,
		}

		result := f.Platforms(tc.Supported)
		if !reflect.DeepEqual(result, tc.Result) {
			t.Errorf("Index: %d. input: %#v\nresult: %#v", i, f, result)
		}
	}
}

func TestPlatformFlagArchFlagValue(t *testing.T) {
	var f PlatformFlag
	val := f.ArchFlagValue()
	if err := val.Set("foo bar"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{"foo", "bar"}
	if !reflect.DeepEqual(f.Arch, expected) {
		t.Fatalf("bad: %#v", f.Arch)
	}
}

func TestPlatformFlagOSArchFlagValue(t *testing.T) {
	var f PlatformFlag
	val := f.OSArchFlagValue()
	if err := val.Set("foo/bar"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []Platform{{OS: "foo", Arch: "bar", Default: false}}
	if !reflect.DeepEqual(f.OSArch, expected) {
		t.Fatalf("bad: %#v", f.OSArch)
	}
}

func TestPlatformFlagOSFlagValue(t *testing.T) {
	var f PlatformFlag
	val := f.OSFlagValue()
	if err := val.Set("foo bar"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{"foo", "bar"}
	if !reflect.DeepEqual(f.OS, expected) {
		t.Fatalf("bad: %#v", f.OS)
	}
}

func TestAppendPlatformValue_impl(t *testing.T) {
	var _ flag.Value = new(appendPlatformValue)
}

func TestAppendPlatformValue(t *testing.T) {
	var value appendPlatformValue

	if err := value.Set(""); err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(value) > 0 {
		t.Fatalf("bad: %#v", value)
	}

	if err := value.Set("windows/arm/bad"); err == nil {
		t.Fatal("should err")
	}

	if err := value.Set("windows"); err == nil {
		t.Fatal("should err")
	}

	if err := value.Set("windows/arm windows/386"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []Platform{
		{OS: "windows", Arch: "arm", Default: false},
		{OS: "windows", Arch: "386", Default: false},
	}
	if !reflect.DeepEqual([]Platform(value), expected) {
		t.Fatalf("bad: %#v", value)
	}
}

func TestAppendStringValue_impl(t *testing.T) {
	var _ flag.Value = new(appendStringValue)
}

func TestAppendStringValue(t *testing.T) {
	var value appendStringValue

	if err := value.Set(""); err != nil {
		t.Fatalf("err: %s", err)
	}

	if len(value) > 0 {
		t.Fatalf("bad: %#v", value)
	}

	if err := value.Set("windows LINUX"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []string{"windows", "linux"}
	if !reflect.DeepEqual([]string(value), expected) {
		t.Fatalf("bad: %#v", value)
	}

	if err := value.Set("darwin"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected = []string{"windows", "linux", "darwin"}
	if !reflect.DeepEqual([]string(value), expected) {
		t.Fatalf("bad: %#v", value)
	}

	if err := value.Set("darwin"); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected = []string{"windows", "linux", "darwin"}
	if !reflect.DeepEqual([]string(value), expected) {
		t.Fatalf("bad: %#v", value)
	}
}
