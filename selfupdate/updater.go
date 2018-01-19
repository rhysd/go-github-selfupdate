package selfupdate

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	gitconfig "github.com/tcnksm/go-gitconfig"
	"golang.org/x/oauth2"
)

// Updater is responsible for managing the context of self-update.
// It contains GitHub client and its context.
type Updater struct {
	api    *github.Client
	apiCtx context.Context
}

// Config represents the configuration of self-update.
type Config struct {
	// APIToken represents GitHub API token. If it's not empty, it will be used for authentication of GitHub API
	APIToken string
	// TODO: Add host URL for API endpoint
}

// NewUpdater crates a new detector instance. It initializes GitHub API client.
func NewUpdater(config Config) *Updater {
	token := config.APIToken
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		token, _ = gitconfig.GithubToken()
	}
	ctx := context.Background()

	var auth *http.Client
	if token != "" {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		auth = oauth2.NewClient(ctx, src)
	}

	client := github.NewClient(auth)
	return &Updater{client, ctx}
}
