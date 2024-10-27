package helpers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(getCurrentPath() + ".env")
	if err != nil {
		log.Fatal(err.Error())
	}
}
