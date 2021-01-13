## [v1.2.3] - 2021-01-13

- Fix security issues in dependencies; CVE-2020-16845, CVE-2019-11840, CVE-2020-14040 (Thanks to [@bhamail](https://github.com/bhamail)).

## [v1.2.2] - 2020-04-10

- Update `go-github` dependency to v30.1.0

## [v1.2.1] - 2019-12-19

- Fix `.tgz` file was not handled as `.tar.gz`.


## [v1.2.0] - 2019-12-19

- New Feature: Filtering releases by matching regular expressions to release names (Thanks to [@fredbi](https://github.com/fredbi)).
  Regular expression strings specified at `Filters` field in `Config` struct are used on detecting the
  latest release. Please read [documentation](https://godoc.org/github.com/rhysd/go-github-selfupdate/selfupdate#Config)
  for more details.
- Allow `{cmd}_{os}_{arch}` format for executable names.
- `.tgz` file name suffix was supported.


## [v1.1.0] - 2018-11-10

- New Feature: Signature validation for release assets (Thanks to [@tobiaskohlbau](https://github.com/tobiaskohlbau)).
  Please read [the instruction](https://github.com/rhysd/go-github-selfupdate#hash-or-signature-validation) for usage.


## [v1.0.0] - 2018-09-23

First release! :tada:


[v1.2.3]: https://github.com/rhysd/go-github-selfupdate/compare/v1.2.2...v1.2.3
[v1.2.2]: https://github.com/rhysd/go-github-selfupdate/compare/v1.2.1...v1.2.2
[v1.2.1]: https://github.com/rhysd/go-github-selfupdate/compare/v1.2.0...v1.2.1
[v1.2.0]: https://github.com/rhysd/go-github-selfupdate/compare/go-get-release...v1.2.0
[v1.1.0]: https://github.com/rhysd/go-github-selfupdate/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/rhysd/go-github-selfupdate/compare/example-1.2.4...v1.0.0
