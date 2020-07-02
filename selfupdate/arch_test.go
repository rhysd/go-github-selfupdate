package selfupdate

import (
	"runtime"
	"testing"
)

func TestDefaultOSAndArch(t *testing.T) {
	os, arch := GetOSArch()
	if os != runtime.GOOS {
		t.Errorf("OS should be %s but found %s", runtime.GOOS, os)
	}
	if arch != runtime.GOARCH {
		t.Errorf("Arch should be %s but found %s", runtime.GOARCH, arch)
	}
}

func TestForcesArch(t *testing.T) {
	testArch := "test"
	SetArch(testArch)
	_, arch := GetOSArch()

	if arch != testArch {
		t.Errorf("Arch should be %s but found %s", testArch, arch)
	}

	SetArch("")
	_, arch = GetOSArch()

	if arch != runtime.GOARCH {
		t.Errorf("Arch should be %s but found %s", runtime.GOARCH, arch)
	}
}
