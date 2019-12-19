## [Unreleased]


## [v1.2.0] - 2019-12-19

- Filtering releases by matching regular expressions to release names was added (Thanks to @fredbi).
Regular expression strings specified at `Filters` field in `Config` struct are used on detecting the
latest release. Please read [documentation](https://godoc.org/github.com/rhysd/go-github-selfupdate/selfupdate#Config)
for more details.
- Allow `{cmd}_{os}_{arch}` format for executable names
- `.tgz` file name suffix was supported


## [v1.1.0] - 2018-11-10

- Signature validation for release assets was added. (Thanks to @tobiaskohlbau). Please read
[the instruction](https://github.com/rhysd/go-github-selfupdate#hash-or-signature-validation) for usage.


## [v1.0.0] - 2018-09-23

First release! :tada:


[Unreleased]: https://github.com/rhysd/go-github-selfupdate/compare/v1.2.0...HEAD
[v1.2.0]: https://github.com/rhysd/go-github-selfupdate/compare/go-get-release...v1.2.0
[v1.1.0]: https://github.com/rhysd/go-github-selfupdate/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/rhysd/go-github-selfupdate/compare/example-1.2.4...v1.0.0
