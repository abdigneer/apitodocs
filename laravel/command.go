package laravel

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func execute() []byte {
	executable := os.Getenv("EXECUTABLE")

	executables := []string{}
	switch Version {
	case PHP73_LARAVEL8:
		executables = append(
			strings.Split(executable, " "),
			"route:list", "--columns=name,method,uri,middleware,action", "--json", "--sort=uri",
		)
	case PHP81_LARAVEL9:
		executables = append(
			strings.Split(executable, " "),
			"route:list", "--json", "--sort=uri",
		)
	}

	output, err := exec.Command(executables[0], executables[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	return output
}
