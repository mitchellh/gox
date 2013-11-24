## 0.3.0 (November 23, 2013)

FEATURES:

  - Use `-osarch` to specify complete os/arch pairs to build for.

## 0.2.0 (November 19, 2013)

FEATURES:

  - Can specify `-ldflags` for the go build in order to get things like
    variables injected into the compile.

IMPROVEMENTS:

  - Building toolchain no longer requires the GOROOT env var. It is
    now automatically detected using `go env`
  - On Windows builds, ".exe" is appended to the output path.

BUG FIXES:

  - When building toolchains with verbose mode, wait until output is fully
    read before moving on to next compilation.
  - Work with `-os` or `-arch` is an empty string.
  - Building toolchain doesn't output "plan9" for all platforms.
  - Don't parallelize toolchain building, because it doesn't work.

## 0.1.0 (November 17, 2013)

Initial release
