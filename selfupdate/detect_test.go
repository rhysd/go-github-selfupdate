package selfupdate

import (
	"fmt"
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

func TestDetectReleasesForVariousArchives(t *testing.T) {
	for _, repo := range []string{
		"rhysd-test/test-release-zip",
		"rhysd-test/test-release-tar",
		"rhysd-test/test-release-gzip",
	} {
		t.Run(repo, func(t *testing.T) {
			r, ok, err := DetectLatest(repo)
			if err != nil {
				t.Fatal("Fetch failed:", err)
			}
			if !ok {
				t.Fatal(repo, "not found")
			}
			if r == nil {
				t.Fatal("Release not detected")
			}
			if !r.Version.Equals(semver.MustParse("1.2.3")) {
				t.Error("")
			}
			url := fmt.Sprintf("https://github.com/%s/releases/tag/v1.2.3", repo)
			if r.URL != url {
				t.Error("URL is not wrong. Want", url, "but got", r.URL)
			}
			if r.ReleaseNotes == "" {
				t.Error("Release note is unexpectedly empty")
			}
			if !strings.HasPrefix(r.AssetURL, fmt.Sprintf("https://github.com/%s/releases/download/v1.2.3/", repo)) {
				t.Error("Unexpected asset URL:", r.AssetURL)
			}
		})
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
		if !strings.Contains(err.Error(), "Invalid slug format") {
			t.Error("Unexpected error for", slug, ":", err)
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
