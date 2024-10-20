package helpers

import (
	"os"
	"path/filepath"
)

func IsInDevEnv() bool {
	_, err := os.Stat("go.mod")

	return err == nil
}

func GetCurrentPath() string {
	if IsInDevEnv() {
		currentDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		return currentDir + "/"
	}

	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(exePath) + "/"
}
