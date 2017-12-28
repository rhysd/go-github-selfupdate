package selfupdate

import (
	"os/exec"
	"testing"
)

func init() {
	if err := exec.Command("go", "build", "../cmd/selfupdate-example/main.go").Run(); err != nil {
		panic(err)
	}
}

func TestRunSelfUpdateExample(t *testing.T) {
	t.Fatal("TODO")
}
