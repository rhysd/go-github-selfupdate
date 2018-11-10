Like `go get`, but it downloads and installs the latest release binary from GitHub instead.

Please download a binary from [release page](https://github.com/rhysd/go-github-selfupdate/releases/tag/go-get-release)
and put it in `$PATH` or build from source with `go get`.

```
$ go get -u github.com/rhysd/go-github-selfupdate/cmd/go-get-release
```

Usage is quite similar to `go get`. But `{package}` must be hosted on GitHub. So it needs to start with `github.com/`.

```
$ go-get-release {package}
```

Please note that this command assumes that specified package is following Git tag naming rules and
released binaries naming rules described in [README](../../README.md).

For example, following command downloads and installs the released binary of [ghr](https://github.com/tcnksm/ghr)
to `$GOPATH/bin`.

```
$ go-get-release github.com/tcnksm/ghr
Command was updated to the latest version 0.5.4: /Users/you/.go/bin/ghr

$ ghr -version
ghr version v0.5.4 (a12ff1c)
```

