package selfupdate

import "runtime"

var (
	useOS   = runtime.GOOS
	useArch = runtime.GOARCH
)

// SetArch forces the use of a different architecture than the default runtime.GOARCH
//
// For example SetArch("arm_v6") instead of the default "arm"
func SetArch(arch string) {
	if arch == "" {
		// Back to the default
		arch = runtime.GOARCH
	}
	useArch = arch
}

// GetOSArch returns the OS and Architecture currently used to detect a new version
func GetOSArch() (string, string) {
	return useOS, useArch
}
