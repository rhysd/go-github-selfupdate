package selfupdate

import (
	"github.com/blang/semver"
)

// Release represents a release asset for current OS and arch.
type Release struct {
	Version  semver.Version
	AssetURL string
}
