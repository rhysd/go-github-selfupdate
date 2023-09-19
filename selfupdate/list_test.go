package selfupdate

import "testing"

func TestListReleases(t *testing.T) {
	r, ok, err := ListReleases("rhysd/github-clone-all")
	if err != nil {
		t.Fatal("List releases:", err)
	}
	if !ok {
		t.Fatalf("Failed to list releases")
	}
	if r == nil {
		t.Fatal("Release detected but nil returned for it")
	}
	if len(r) != 7 {
		t.Fatalf("The of releases does not match: got %d, expected %d", len(r), 7)
	}
	t.Logf("Found releases: %#v", r)
}
