package selfupdate

import (
	"fmt"
	"strings"
)

type repoInfo struct {
	owner string
	name  string
}

func parseRepo(slug string) (*repoInfo, error) {
	parts := strings.Split(slug, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid slug %q: it should be in the format %q", slug, "owner/name")
	}
	return &repoInfo{
		owner: parts[0],
		name:  parts[1],
	}, nil
}
