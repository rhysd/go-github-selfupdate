package selfupdate

import (
	"github.com/blang/semver"
	"os"
	"strings"
	"testing"
)

func TestGitHubTokenEnv(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("because $GITHUB_TOKEN is not set")
	}
	_ = NewDetector()
}

func TestDetectReleaseWithVersionPrefix(t *testing.T) {
	r, ok, err := DetectLatest("rhysd/github-clone-all")
	if err != nil {
		t.Fatal("Fetch failed:", err)
	}
	if !ok {
		t.Fatal("Failed to detect latest")
	}
	if r == nil {
		t.Fatal("Release detected but nil returned for it")
	}
	if r.Version.LE(semver.MustParse("2.0.0")) {
		t.Error("Incorrect version:", r.Version)
	}
	if !strings.HasSuffix(r.AssetURL, ".zip") && !strings.HasSuffix(r.AssetURL, ".tar.gz") {
		t.Error("Incorrect URL for asset:", r.AssetURL)
	}
	if r.URL == "" {
		t.Error("Document URL should not be empty")
	}
	if r.ReleaseNotes == "" {
		t.Error("Description should not be empty for this repo")
	}
}

func TestDetectReleaseButNoAsset(t *testing.T) {
	_, ok, err := DetectLatest("rhysd/clever-f.vim")
	if err != nil {
		t.Fatal("Fetch failed:", err)
	}
	if ok {
		t.Fatal("When no asset found, result should be marked as 'not found'")
	}
}

func TestDetectNoRelease(t *testing.T) {
	_, ok, err := DetectLatest("rhysd/clever-f.vim")
	if err != nil {
		t.Fatal("Fetch failed:", err)
	}
	if ok {
		t.Fatal("When no release found, result should be marked as 'not found'")
	}
}

func TestInvalidSlug(t *testing.T) {
	d := NewDetector()

	for _, slug := range []string{
		"foo",
		"/",
		"foo/",
		"/bar",
		"foo/bar/piyo",
	} {
		_, _, err := d.DetectLatest(slug)
		if err == nil {
			t.Error(slug, "should be invalid slug")
		}
	}
}

func TestNonExistingRepo(t *testing.T) {
	d := NewDetector()
	v, ok, err := d.DetectLatest("rhysd/non-existing-repo")
	if err != nil {
		t.Fatal("Non-existing repo should not cause an error:", v)
	}
	if ok {
		t.Fatal("Release for non-existing repo should not be found")
	}
}

func TestNoReleaseFound(t *testing.T) {
	d := NewDetector()
	_, ok, err := d.DetectLatest("rhysd/misc")
	if err != nil {
		t.Fatal("Repo having no release should not cause an error:", err)
	}
	if ok {
		t.Fatal("Repo having no release should not be found")
	}
}
