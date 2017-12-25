Self-Update Mechanism for Go Commands using GitHub
==================================================

[go-github-selfupdate][] is a Go library to provide self-update mechanism to command line tools.

Go does not provide the way to install/update the specific version of tools. By default, Go command line tools are updated

- using `go get -u` (updating to HEAD)
- using system's package manager (depending on the platform)
- downloading executables from GitHub release page

By using this library, you will get 4th choice:

- from your command line tool directly

go-github-selfupdate retrieves the information of released binaries via [GitHub Releases API][] and check the current version. If newer version than itself is detected, it updates binary by downloading from GitHub and replacing itself.

Please note that this library assumes you adopt [semantic versioning][]. It is necessary for comparing versions systematically.

[go-github-selfupdate]: https://github.com/rhysd/go-github-selfupdate
[semantic versioning]: https://semver.org/
[GitHub Releases API]: https://developer.github.com/v3/repos/releases/

## Usage

### Code

It provides `selfupdate` package.

```go
import (
    "log"
    "github.com/rhysd/go-github-selfupdate"
)

func doUpdate(version string) {
    up, err := selfupdate.TryUpdate(version, "myname/myrepo", nil)
    if err != nil {
        log.Println("Binary update failed", err)
        return
    }
    if up.Version == version {
        log.Println("Current binary is the latest version", version)
    } else {
        log.Println("Update successfully done to version", up.Version)
    }
}
```

- `selfupdate.TryUpdate()`
- `selfupdate.TryUpdateTo()`
- `selfupdate.DetectLatest()`

Please see [the documentation page][GoDoc] for more detail.

[GoDoc]: https://godoc.org/github.com/rhysd/go-github-selfupdate

### Naming Rules of Released Binaries

go-github-selfupdate assumes that released binaries are put for each combination of platforms and archs.

For example, if your command name is `foo-bar`, one of followings is expected to be put in release page on GitHub as binary for platform `linux` and arch `amd64`.

- `foo-bar-linux-amd64` (executable)
- `foo-bar-linux-amd64.zip` (zip file)
- `foo-bar-linux-amd64.tar.gz` (gzip file)
- `foo-bar_linux_amd64` (executable)
- `foo-bar_linux_amd64.zip` (zip file)
- `foo-bar_linux_amd64.tar.gz` (gzip file)

### Naming Rules of Git Tags

go-github-selfupdate searches binaries' versions via corresponding Git tag names. When your binary's version is `1.2.3`, you should use `1.2.3` or `v1.2.3` (prefixed with `v`) as Git tag name.

### Structure of Releases

In summary, structure of releases on GitHub looks like:

- `v1.2.0`
  - `foo-bar-linux-amd64.zip`
  - `foo-bar-linux-386.zip`
  - `foo-bar-darwin-amd64.zip`
  - `foo-bar-windows-amd64.zip`
  - ... (Other binaries for v1.2.0)
- `v1.1.3`
  - `foo-bar-linux-amd64.zip`
  - `foo-bar-linux-386.zip`
  - `foo-bar-darwin-amd64.zip`
  - `foo-bar-windows-amd64.zip`
  - ... (Other binaries for v1.1.3)
- ... (older versions)


## Dependencies

This library utilizes [go-github][] to retrieve the information of releases and [go-update][] to replace current binary and [semver][] to compare versions.

> Copyright (c) 2013 The go-github AUTHORS. All rights reserved.

> Copyright 2015 Alan Shreve

> Copyright (c) 2014 Benedikt Lang <github at benediktlang.de>

[go-github]: https://github.com/google/go-github
[go-upadte]: https://github.com/inconshreveable/go-update
[semver]: https://github.com/blang/semver
