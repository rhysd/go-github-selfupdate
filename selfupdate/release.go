package selfupdate

import (
	"time"

	"github.com/blang/semver"
	"github.com/google/go-github/v30/github"
)

// Release represents a release asset for current OS and arch.
type Release struct {
	// Version is the version of the release
	Version semver.Version
	// AssetURL is a URL to the uploaded file for the release
	AssetURL string
	// AssetSize represents the size of asset in bytes
	AssetByteSize int
	// AssetID is the ID of the asset on GitHub
	AssetID int64
	// ValidationAssetID is the ID of additional validaton asset on GitHub
	ValidationAssetID int64
	// URL is a URL to release page for browsing
	URL string
	// ReleaseNotes is a release notes of the release
	ReleaseNotes string
	// Name represents a name of the release
	Name string
	// PublishedAt is the time when the release was published
	PublishedAt *time.Time
	// RepoOwner is the owner of the repository of the release
	RepoOwner string
	// RepoName is the name of the repository of the release
	RepoName string
}

func newRelease(repo *repoInfo, release *github.RepositoryRelease,
	asset *github.ReleaseAsset, version *semver.Version) *Release {

	publishedAt := release.GetPublishedAt().Time

	return &Release{
		Version:           *version,
		AssetURL:          asset.GetBrowserDownloadURL(),
		AssetByteSize:     asset.GetSize(),
		AssetID:           asset.GetID(),
		ValidationAssetID: -1,
		URL:               release.GetHTMLURL(),
		ReleaseNotes:      release.GetBody(),
		Name:              release.GetName(),
		PublishedAt:       &publishedAt,
		RepoOwner:         repo.owner,
		RepoName:          repo.name,
	}
}
