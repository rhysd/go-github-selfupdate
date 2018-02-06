package selfupdate

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/blang/semver"
	"github.com/google/go-github/github"
)

var reVersion = regexp.MustCompile(`\d+\.\d+\.\d+`)

func findAssetFromReleasse(rel *github.RepositoryRelease, suffixes []string) (*github.ReleaseAsset, bool) {
	if rel.GetDraft() {
		log.Println("Skip draft version", rel.GetTagName())
		return nil, false
	}
	if rel.GetPrerelease() {
		log.Println("Skip pre-release version", rel.GetTagName())
		return nil, false
	}
	if !reVersion.MatchString(rel.GetTagName()) {
		log.Println("Skip version not adopting semver", rel.GetTagName())
		return nil, false
	}
	for _, asset := range rel.Assets {
		name := asset.GetName()
		for _, s := range suffixes {
			if strings.HasSuffix(name, s) {
				return &asset, true
			}
		}
	}
	log.Println("No suitable asset was found in release", rel.GetTagName())
	return nil, false
}

func findReleaseAndAsset(rels []*github.RepositoryRelease) (*github.RepositoryRelease, *github.ReleaseAsset, bool) {
	// Generate valid suffix candidates to match
	suffixes := make([]string, 0, 2*7*2)
	for _, sep := range []rune{'_', '-'} {
		for _, ext := range []string{".zip", ".tar.gz", ".gzip", ".gz", ".tar.xz", ".xz", ""} {
			suffix := fmt.Sprintf("%s%c%s%s", runtime.GOOS, sep, runtime.GOARCH, ext)
			suffixes = append(suffixes, suffix)
			if runtime.GOOS == "windows" {
				suffix = fmt.Sprintf("%s%c%s.exe%s", runtime.GOOS, sep, runtime.GOARCH, ext)
				suffixes = append(suffixes, suffix)
			}
		}
	}

	for _, rel := range rels {
		if asset, ok := findAssetFromReleasse(rel, suffixes); ok {
			return rel, asset, true
		}
	}

	log.Println("Could not find any release for", runtime.GOOS, "and", runtime.GOARCH)
	return nil, nil, false
}

// DetectLatest tries to get the latest version of the repository on GitHub. 'slug' means 'owner/name' formatted string.
// It fetches releases information from GitHub API and find out the latest release with matching the tag names and asset names.
// Drafts and pre-releases are ignored. Assets whould be suffixed by the OS name and the arch name such as 'foo_linux_amd64'
// where 'foo' is a command name. '-' can also be used as a separator. File can be compressed with zip, gzip, zxip, tar&zip or tar&zxip.
// So the asset can have a file extension for the corresponding compression format such as '.zip'.
// On Windows, '.exe' also can be contained such as 'foo_windows_amd64.exe.zip'.
func (up *Updater) DetectLatest(slug string) (release *Release, found bool, err error) {
	repo := strings.Split(slug, "/")
	if len(repo) != 2 || repo[0] == "" || repo[1] == "" {
		err = fmt.Errorf("Invalid slug format. It should be 'owner/name': %s", slug)
		return
	}

	rels, res, err := up.api.Repositories.ListReleases(up.apiCtx, repo[0], repo[1], nil)
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

	rel, asset, found := findReleaseAndAsset(rels)
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
		AssetID:       asset.GetID(),
		URL:           rel.GetHTMLURL(),
		ReleaseNotes:  rel.GetBody(),
		Name:          rel.GetName(),
		PublishedAt:   &publishedAt,
		RepoOwner:     repo[0],
		RepoName:      repo[1],
	}

	release.Version, err = semver.Make(tag)
	return
}

// DetectLatest detects the latest release of the slug (owner/repo).
// This function is a shortcut version of updater.DetectLatest() method.
func DetectLatest(slug string) (*Release, bool, error) {
	return DefaultUpdater().DetectLatest(slug)
}
