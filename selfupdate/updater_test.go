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
	_ = DefaultUpdater()
	if _, err := NewUpdater(Config{}); err != nil {
		t.Error("Failed to initialize updater with empty config")
	}
	if _, err := NewUpdater(Config{APIToken: token}); err != nil {
		t.Error("Failed to initialize updater with API token config")
	}
}
