package selfupdate

import (
	"os/exec"
	"testing"
)

func TestRunSelfUpdateExample(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	if err := exec.Command("go", "build", "../cmd/selfupdate-example").Run(); err != nil {
		t.Fatal(err)
	}
	t.Fatal("TODO")
}
