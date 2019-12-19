package selfupdate

import (
	"os"
	"strings"
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

func TestGitHubTokenIsNotSet(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		defer os.Setenv("GITHUB_TOKEN", token)
	}
	os.Setenv("GITHUB_TOKEN", "")
	_ = DefaultUpdater()
	if _, err := NewUpdater(Config{}); err != nil {
		t.Error("Failed to initialize updater with empty config")
	}
}

func TestGitHubEnterpriseClient(t *testing.T) {
	url := "https://github.company.com/api/v3/"
	up, err := NewUpdater(Config{APIToken: "hogehoge", EnterpriseBaseURL: url})
	if err != nil {
		t.Fatal(err)
	}
	if up.api.BaseURL.String() != url {
		t.Error("Base URL was set to", up.api.BaseURL, ", want", url)
	}
	if up.api.UploadURL.String() != url {
		t.Error("Upload URL was set to", up.api.UploadURL, ", want", url)
	}

	url2 := "https://upload.github.company.com/api/v3/"
	up, err = NewUpdater(Config{
		APIToken:            "hogehoge",
		EnterpriseBaseURL:   url,
		EnterpriseUploadURL: url2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if up.api.BaseURL.String() != url {
		t.Error("Base URL was set to", up.api.BaseURL, ", want", url)
	}
	if up.api.UploadURL.String() != url2 {
		t.Error("Upload URL was set to", up.api.UploadURL, ", want", url2)
	}
}

func TestGitHubEnterpriseClientInvalidURL(t *testing.T) {
	_, err := NewUpdater(Config{APIToken: "hogehoge", EnterpriseBaseURL: ":this is not a URL"})
	if err == nil {
		t.Fatal("Invalid URL should raise an error")
	}
}

func TestCompileRegexForFiltering(t *testing.T) {
	filters := []string{
		"^hello$",
		"^(\\d\\.)+\\d$",
	}
	up, err := NewUpdater(Config{
		Filters: filters,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(up.filters) != 2 {
		t.Fatalf("Wanted 2 regexes but got %d", len(up.filters))
	}
	for i, r := range up.filters {
		want := filters[i]
		got := r.String()
		if want != got {
			t.Errorf("Compiled regex is %q but specified was %q", got, want)
		}
	}
}

func TestFilterRegexIsBroken(t *testing.T) {
	_, err := NewUpdater(Config{
		Filters: []string{"(foo"},
	})
	if err == nil {
		t.Fatal("Error unexpectedly did not occur")
	}
	msg := err.Error()
	if !strings.Contains(msg, "Could not compile regular expression \"(foo\" for filtering releases") {
		t.Fatalf("Error message is unexpected: %q", msg)
	}
}
