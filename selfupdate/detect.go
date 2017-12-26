package selfupdate

import (
	"context"
	"fmt"
	"github.com/blang/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var reVersion = regexp.MustCompile(`\d+\.\d+\.\d+`)

// ReleaseDetector is responsible for detecting the latest release using GitHub Releases API.
type ReleaseDetector struct {
	api    *github.Client
	apiCtx context.Context
}

// NewDetector crates a new detector instance. It initializes GitHub API client.
func NewDetector() *ReleaseDetector {
	token := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()

	var auth *http.Client
	if token != "" {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		auth = oauth2.NewClient(ctx, src)
	}

	client := github.NewClient(auth)
	return &ReleaseDetector{client, ctx}
}

// DetectLatest tries to get the latest version of the repository on GitHub. 'slug' means 'owner/name' formatted string.
func (d *ReleaseDetector) DetectLatest(slug string) (ver semver.Version, found bool, err error) {
	repo := strings.Split(slug, "/")
	if len(repo) != 2 || repo[0] == "" || repo[1] == "" {
		err = fmt.Errorf("Invalid slug format. It should be 'owner/name': %s", slug)
		return
	}

	rel, res, err := d.api.Repositories.GetLatestRelease(d.apiCtx, repo[0], repo[1])
	if err != nil {
		log.Println("API returned an error response:", err)
		if res.StatusCode == 404 {
			// 404 means repository not found or release not found. It's not an error here.
			found = false
			err = nil
			log.Println("API returned 404. Repository or release not found")
		}
		return
	}

	tag := rel.GetTagName()
	log.Println("Successfully fetched the latest release. tag:", tag, ", name:", rel.GetName(), ", URL:", rel.GetURL())

	// Strip version prefix
	if indices := reVersion.FindStringIndex(tag); indices != nil && indices[0] > 0 {
		tag = tag[indices[0]:]
	}

	ver, err = semver.Make(tag)
	if err == nil {
		found = true
	}

	return
}

// DetectLatest detects the latest release of the slug (owner/repo).
func DetectLatest(slug string) (semver.Version, bool, error) {
	return NewDetector().DetectLatest(slug)
}
