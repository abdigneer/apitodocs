package main

import (
	"apitodocs/helpers"
	"apitodocs/laravel"
	"apitodocs/postman"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// currently support laravel 8 only
const PHP73_LARAVEL8 string = "php73-laravel8"
const BASE_URL = "http://localhost:8080"
const PROTOCOL = "http"

var fromFlag *string
var baseUrlFlag *string

func main() {
	loadEnv()
	flagParserAndModifier()

	postmanCollection := collectionFrom()

	exportToFile(postmanCollection)
}

func loadEnv() {
	err := godotenv.Load(helpers.GetCurrentPath() + ".env")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func flagParserAndModifier() {
	baseUrlFlag = flag.String("base-url", BASE_URL, "Custom base url")
	fromFlag = flag.String("from", PHP73_LARAVEL8, "Laravel version")
	flag.Parse()
	flagModifier()
	fmt.Println(*baseUrlFlag)
}

func collectionFrom() postman.Collection {
	postmanCollection := postman.Collection{}
	switch *fromFlag {
	case PHP73_LARAVEL8:
		postmanCollection = laravel.MakeCollection()
	}

	return postmanCollection
}

func flagModifier() {
	if !strings.Contains(*baseUrlFlag, "http://") && !strings.Contains(*baseUrlFlag, "https://") {
		baseUrl := PROTOCOL + "://" + *baseUrlFlag
		*baseUrlFlag = baseUrl
	}
}

func exportToFile(postmanCollection postman.Collection) {
	result, err := json.MarshalIndent(postmanCollection, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(helpers.GetCurrentPath()+"collection.json", result, 0644)
}
