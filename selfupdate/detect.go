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
	"runtime"
	"strings"
)

var reVersion = regexp.MustCompile(`\d+\.\d+\.\d+`)

// ReleaseDetector is responsible for detecting the latest release using GitHub Releases API.
type ReleaseDetector struct {
	api    *github.Client
	apiCtx context.Context
}

func findSuitableReleaseAndAsset(rels []*github.RepositoryRelease) (*github.RepositoryRelease, *github.ReleaseAsset, bool) {
	// Generate candidates
	cs := make([]string, 0, 8)
	for _, sep := range []rune{'_', '-'} {
		for _, ext := range []string{".zip", ".tar.gz", ".gzip", ".gz", ".tar.xz", ".xz", ""} {
			suffix := fmt.Sprintf("%s%c%s%s", runtime.GOOS, sep, runtime.GOARCH, ext)
			cs = append(cs, suffix)
			if runtime.GOOS == "windows" {
				suffix = fmt.Sprintf("%s%c%s.exe%s", runtime.GOOS, sep, runtime.GOARCH, ext)
				cs = append(cs, suffix)
			}
		}
	}

	for _, rel := range rels {
		if rel.GetDraft() {
			log.Println("Skip draft version", rel.GetTagName())
			continue
		}
		if rel.GetPrerelease() {
			log.Println("Skip pre-release version", rel.GetTagName())
			continue
		}
		if !reVersion.MatchString(rel.GetTagName()) {
			log.Println("Skip version not adopting semver", rel.GetTagName())
			continue
		}
		for _, asset := range rel.Assets {
			name := asset.GetName()
			for _, c := range cs {
				if strings.HasSuffix(name, c) {
					return rel, &asset, true
				}
			}
		}
	}

	log.Println("Could no find any release for", runtime.GOOS, "and", runtime.GOARCH)
	return nil, nil, false
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
func (d *ReleaseDetector) DetectLatest(slug string) (release *Release, found bool, err error) {
	repo := strings.Split(slug, "/")
	if len(repo) != 2 || repo[0] == "" || repo[1] == "" {
		err = fmt.Errorf("Invalid slug format. It should be 'owner/name': %s", slug)
		return
	}

	rels, res, err := d.api.Repositories.ListReleases(d.apiCtx, repo[0], repo[1], nil)
	if err != nil {
		log.Println("API returned an error response:", err)
		if res != nil && res.StatusCode == 404 {
			// 404 means repository not found or release not found. It's not an error here.
			found = false
			err = nil
			log.Println("API returned 404. Repository or release not found")
		}
		return
	}

	rel, asset, found := findSuitableReleaseAndAsset(rels)
	if !found {
		return
	}

	tag := rel.GetTagName()
	url := asset.GetBrowserDownloadURL()
	log.Println("Successfully fetched the latest release. tag:", tag, ", name:", rel.GetName(), ", URL:", rel.GetURL(), ", Asset:", url)

	// Strip version prefix
	if indices := reVersion.FindStringIndex(tag); indices != nil && indices[0] > 0 {
		log.Println("Strip prefix of version:", tag[:indices[0]])
		tag = tag[indices[0]:]
	}

	publishedAt := rel.GetPublishedAt().Time
	release = &Release{
		AssetURL:      url,
		AssetByteSize: asset.GetSize(),
		URL:           rel.GetHTMLURL(),
		ReleaseNotes:  rel.GetBody(),
		Name:          rel.GetName(),
		PublishedAt:   &publishedAt,
	}

	release.Version, err = semver.Make(tag)
	return
}

// DetectLatest detects the latest release of the slug (owner/repo).
func DetectLatest(slug string) (*Release, bool, error) {
	return NewDetector().DetectLatest(slug)
}
