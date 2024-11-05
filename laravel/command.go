package laravel

import (
	"log"
	"os"
	"os/exec"
)

func artisan(phpExec string) []string {
	_, err := os.Stat(Location + "/artisan")
	if err != nil {
		log.Fatal("laravel artisan not found")
	}

	return []string{phpExec, Location + "/artisan"}
}

func execute() []byte {
	phpExec := os.Getenv("PHP")

	executables := []string{}
	switch Version {
	case PHP73_LARAVEL8:
		if os.Getenv("PHP73") != "" {
			phpExec = os.Getenv("PHP73")
		}

		executables = append(
			artisan(phpExec),
			"route:list", "--columns=name,method,uri,middleware,action", "--json", "--sort=uri",
		)
	case PHP81_LARAVEL9:
		if os.Getenv("PHP81") != "" {
			phpExec = os.Getenv("PHP81")
		}

		executables = append(
			artisan(phpExec),
			"route:list", "--json", "--sort=uri",
		)
	case PHP81_LARAVEL10:
		if os.Getenv("PHP81") != "" {
			phpExec = os.Getenv("PHP81")
		}

		executables = append(
			artisan(phpExec),
			"route:list", "--json", "--sort=uri",
		)
	}

	output, err := exec.Command(executables[0], executables[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	return output
}
