package selfupdate

import (
	"fmt"
	"regexp"
	"runtime"

	"github.com/google/go-github/v30/github"
)

func findReleases(repo *repoInfo, rels []*github.RepositoryRelease, filters []*regexp.Regexp) ([]*Release, error) {
	var (
		suffixes = assetSuffixes()
		found    = make([]*Release, 0)
	)
	for _, rel := range rels {
		asset, ver, ok := findAssetFromRelease(rel, suffixes, "", filters)
		if ok {
			found = append(found, newRelease(repo, rel, asset, &ver))
		}
	}
	if len(found) == 0 {
		return nil, fmt.Errorf("could not find any release for %s and %s", runtime.GOOS, runtime.GOARCH)
	}
	return found, nil
}

func (up *Updater) ListReleases(slug string) ([]*Release, bool, error) {
	repo, err := parseRepo(slug)
	if err != nil {
		return nil, false, fmt.Errorf("parse slug: %v", err)
	}

	rels, _, err := up.api.Repositories.ListReleases(up.apiCtx, repo.owner, repo.name, nil)
	if err != nil {
		return nil, false, fmt.Errorf("list GitHub releases (owner=%q, name=%q): %v", repo.owner, repo.name, err)
	}

	releases, err := findReleases(repo, rels, up.filters)
	if err != nil {
		return nil, false, fmt.Errorf("find releases: %v", err)
	}

	return releases, true, nil
}

func ListReleases(slug string) ([]*Release, bool, error) {
	return DefaultUpdater().ListReleases(slug)
}
