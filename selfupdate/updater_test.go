package selfupdate

import (
	"os"
	"testing"
)

func TestGitHubTokenEnv(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("because $GITHUB_TOKEN is not set")
	}
	_ = NewUpdater(Config{})
	_ = NewUpdater(Config{APIToken: token})
}
