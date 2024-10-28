package helpers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func IsInDevEnv() bool {
	_, err := os.Stat("go.mod")

	return err == nil
}

func getCurrentPath() string {
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

func ExportToFileAsJson(v any, filename string) {
	result, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(getCurrentPath()+filename, result, 0644)
}
