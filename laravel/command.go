package laravel

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func execute() []byte {
	executable := os.Getenv("EXECUTABLE")
	executables := append(
		strings.Split(executable, " "),
		"route:list", "--columns=name,method,uri,middleware,action", "--json", "--sort=uri",
	)

	output, err := exec.Command(executables[0], executables[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	return output
}
