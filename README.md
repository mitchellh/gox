# Gox - Simple Go Cross Compilation

Gox is a simple, no-frills tool for Go cross compilation that behaves a
lot like standard `go build`. Gox will parallelize builds for multiple
platforms.

## Usage

```
$ gox
```

## Versus Other Cross-Compile Tools

A big thanks to these other options for existing. They each paved the
way in many aspects to make Go cross-compilation approachable.

* [Dave Cheney's go-cross-compile](#) -
  Gox compiles for multiple platforms and can therefore easily run on
  any platform Go supports. Dave's scripts require a shell.

* [goxc](#) -
  A very richly featured tool that can even do things such as build system
  packages, upload binaries, generate download webpages, etc. Gox is a
  super slim alternative that only cross-compiles binaries.

