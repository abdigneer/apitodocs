package helpers

import "runtime"

func GetOsArch() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	return goos + "-" + goarch
}
