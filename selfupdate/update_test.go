package selfupdate

import (
	"github.com/blang/semver"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func setupTestBinary() {
	if err := exec.Command("go", "build", "./testdata/github-release-test/").Run(); err != nil {
		panic(err)
	}
}

func teardownTestBinary() {
	bin := "github-release-test"
	if runtime.GOOS == "windows" {
		bin = "github-release-test.exe"
	}
	if err := os.Remove(bin); err != nil {
		panic(err)
	}
}

func TestUpdateCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skip tests in short mode.")
	}

	for _, slug := range []string{
		"rhysd-test/test-release-zip",
		"rhysd-test/test-release-tar",
		"rhysd-test/test-release-gzip",
	} {
		t.Run(slug, func(t *testing.T) {
			setupTestBinary()
			defer teardownTestBinary()
			latest := semver.MustParse("1.2.3")
			prev := semver.MustParse("1.2.2")
			rel, err := UpdateCommand("github-release-test", prev, slug)
			if err != nil {
				t.Fatal(err)
			}
			if rel.Version.NE(latest) {
				t.Error("Version is not latest", rel.Version)
			}
			bytes, err := exec.Command(filepath.FromSlash("./github-release-test")).Output()
			if err != nil {
				t.Fatal("Failed to run test binary after update:", err)
			}
			out := string(bytes)
			if out != "v1.2.3\n" {
				t.Error("Output from test binary after update is unexpected:", out)
			}
		})
	}
}

func TestNoReleaseFoundForUpdate(t *testing.T) {
	v := semver.MustParse("1.0.0")
	rel, err := UpdateCommand("foo", v, "rhysd/misc")
	if err != nil {
		t.Fatal("No release should not make an error:", err)
	}
	if rel.Version.NE(v) {
		t.Error("No release should return the current version as the latest:", rel.Version)
	}
	if rel.URL != "" {
		t.Error("Browse URL should be empty when no release found:", rel.URL)
	}
	if rel.AssetURL != "" {
		t.Error("Asset URL should be empty when no release found:", rel.AssetURL)
	}
	if rel.ReleaseNotes != "" {
		t.Error("Release notes should be empty when no release found:", rel.ReleaseNotes)
	}
}

func TestCurrentIsTheLatest(t *testing.T) {
	v := semver.MustParse("1.2.3")
	rel, err := UpdateCommand("github-release-test", v, "rhysd-test/test-release-zip")
	if err != nil {
		t.Fatal(err)
	}
	if rel.Version.NE(v) {
		t.Error("v1.2.3 should be the latest:", rel.Version)
	}
	if rel.URL == "" {
		t.Error("Browse URL should not be empty when release found:", rel.URL)
	}
	if rel.AssetURL == "" {
		t.Error("Asset URL should not be empty when release found:", rel.AssetURL)
	}
	if rel.ReleaseNotes == "" {
		t.Error("Release notes should not be empty when release found:", rel.ReleaseNotes)
	}
}

func TestBrokenBinaryUpdate(t *testing.T) {
	_, err := UpdateCommand("foo", semver.MustParse("1.2.2"), "rhysd-test/test-incorrect-release")
	if err == nil {
		t.Fatal("Error should occur for broken package")
	}
	if !strings.Contains(err.Error(), "Failed to uncompress .tar.gz file") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestInvalidSlugForUpdate(t *testing.T) {
	_, err := UpdateCommand("foo", semver.MustParse("1.0.0"), "rhysd/")
	if err == nil {
		t.Fatal("Unknown repo should cause an error")
	}
	if !strings.Contains(err.Error(), "Invalid slug format") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestInvalidAssetURL(t *testing.T) {
	err := UpdateTo("https://github.com/rhysd/non-existing-repo/releases/download/v1.2.3/foo.zip", "foo")
	if err == nil {
		t.Fatal("Error should occur for URL not found")
	}
	if !strings.Contains(err.Error(), "Failed to download a release file") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestBrokenAsset(t *testing.T) {
	asset := "https://github.com/rhysd-test/test-incorrect-release/releases/download/invalid/broken-zip.zip"
	err := UpdateTo(asset, "foo")
	if err == nil {
		t.Fatal("Error should occur for URL not found")
	}
	if !strings.Contains(err.Error(), "Failed to uncompress zip file") {
		t.Fatal("Unexpected error:", err)
	}
}
