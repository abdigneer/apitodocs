package main

import (
	"apitodocs/helpers"
	"apitodocs/laravel"
	"apitodocs/postman"
	"flag"
	"strings"
)

// currently support laravel 8 only
const PHP73_LARAVEL8 string = "php73-laravel8"
const BASE_URL = "http://localhost:8080"
const PROTOCOL = "http"

var fromFlag *string
var baseUrlFlag *string
var useRouteParam *bool

func main() {
	helpers.LoadEnv()
	flagParser()
	flagModifier()

	postmanCollection := postman.Collection{}

	collectionFrom(&postmanCollection)

	helpers.ExportToFileAsJson(postmanCollection)
}

func flagParser() {
	baseUrlFlag = flag.String("base-url", BASE_URL, "Custom base url")
	fromFlag = flag.String("from", PHP73_LARAVEL8, "Laravel version")
	flag.Parse()
}

func collectionFrom(postmanCollection *postman.Collection) {
	switch *fromFlag {
	case PHP73_LARAVEL8:
		*postmanCollection = laravel.MakeCollection()
	}
}

func flagModifier() {
	if !strings.Contains(*baseUrlFlag, "http://") && !strings.Contains(*baseUrlFlag, "https://") {
		baseUrl := PROTOCOL + "://" + *baseUrlFlag
		*baseUrlFlag = baseUrl
	}
}
