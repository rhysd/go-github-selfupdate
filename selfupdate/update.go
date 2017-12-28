package selfupdate

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/blang/semver"
	"github.com/inconshreveable/go-update"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func uncompress(src io.ReadCloser, url, cmd string) (io.Reader, error) {
	if strings.HasSuffix(url, ".zip") {
		// Zip format requires its file size for uncompressing.
		// So we need to read the HTTP response into a buffer at first.
		buf, err := ioutil.ReadAll(src)
		if err != nil {
			return nil, err
		}

		r := bytes.NewReader(buf)
		z, err := zip.NewReader(r, r.Size())
		if err != nil {
			return nil, err
		}

		for _, file := range z.File {
			if file.Name == cmd {
				return file.Open()
			}
		}

		return nil, fmt.Errorf("File '%s' for the command is not found in %s", cmd, url)
	} else if strings.HasSuffix(url, ".gzip") {
		return gzip.NewReader(src)
	} else if strings.HasSuffix(url, ".tar.gz") {
		gz, err := gzip.NewReader(src)
		if err != nil {
			return nil, err
		}

		t := tar.NewReader(gz)
		for {
			h, err := t.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			if h.Name == cmd {
				return t, nil
			}
		}

		return nil, fmt.Errorf("File '%s' for the command is not found in %s", cmd, url)
	}

	return src, nil
}

// UpdateTo download an executable from assetURL and replace the current binary with the downloaded one. cmdPath is a file path to command executable.
func UpdateTo(assetURL, cmdPath string) error {
	res, err := http.Get(assetURL)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Failed to download a release file from %s ", assetURL)
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
