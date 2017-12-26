package selfupdate

import (
	"context"
	"fmt"
	"github.com/blang/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"strings"
)

// ReleaseDetector is responsible for detecting the latest release using GitHub Releases API.
type ReleaseDetector struct {
	verPrefix string
	api       *github.Client
	apiCtx    context.Context
}

// NewDetector crates a new detector instance. It initializes GitHub API client.
func NewDetector(versionPrefix string) *ReleaseDetector {
	token := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()

	var auth *http.Client
	if token != "" {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		auth = oauth2.NewClient(ctx, src)
	}

	client := github.NewClient(auth)
	return &ReleaseDetector{versionPrefix, client, ctx}
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
		if res.StatusCode == 404 {
			// 404 means repository not found or release not found. It's not an error here.
			found = false
			err = nil
			log.Println("API returned 404. Repository or release not found")
		} else {
			log.Println("API returned an error:", err)
		}
		return
	}

	tag := rel.GetTagName()
	log.Println("Successfully fetched the latest release. tag:", tag, ", name:", rel.GetName(), ", URL:", rel.GetURL())

	ver, err = semver.Make(strings.TrimPrefix(tag, d.verPrefix))
	if err == nil {
		found = true
	}

	return
}

// DetectLatest detects the latest release of the slug (owner/repo). verPrefix is a prefix of version in tag name (i.e. 'v' for 'v1.2.3').
func DetectLatest(slug, verPrefix string) (semver.Version, bool, error) {
	return NewDetector(verPrefix).DetectLatest(slug)
}
