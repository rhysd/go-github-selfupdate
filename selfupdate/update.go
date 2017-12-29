package selfupdate

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/inconshreveable/go-update"
	"net/http"
	"os"
	"path/filepath"
)

// UpdateTo download an executable from assetURL and replace the current binary with the downloaded one. cmdPath is a file path to command executable.
func UpdateTo(assetURL, cmdPath string) error {
	res, err := http.Get(assetURL)
	if err != nil {
		return fmt.Errorf("Failed to download a release file from %s: %s", assetURL, err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to download a release file from %s", assetURL)
	}

	defer res.Body.Close()
	_, cmd := filepath.Split(cmdPath)
	asset, err := uncompress(res.Body, assetURL, cmd)
	if err != nil {
		return err
	}

	return update.Apply(asset, update.Options{
		TargetPath: cmdPath,
	})
}

// UpdateCommand updates a given command binary to the latest version.
// 'slug' represents 'owner/name' repository on GitHub and 'current' means the current version.
func UpdateCommand(cmdPath string, current semver.Version, slug string) (*Release, error) {
	rel, ok, err := DetectLatest(slug)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &Release{Version: current}, nil
	}
	if current.Equals(rel.Version) {
		return rel, nil
	}
	if err := UpdateTo(rel.AssetURL, cmdPath); err != nil {
		return nil, err
	}
	return rel, nil
}

// UpdateSelf updates the running executable itself to the latest version.
// 'slug' represents 'owner/name' repository on GitHub and 'current' means the current version.
func UpdateSelf(current semver.Version, slug string) (*Release, error) {
	return UpdateCommand(os.Args[0], current, slug)
}
