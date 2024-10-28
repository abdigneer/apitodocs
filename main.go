package main

import (
	"apitodocs/helpers"
	"apitodocs/laravel"
	"apitodocs/postman"
	"flag"
	"log"
	"strings"
)

const BASE_URL = "http://localhost:8080"
const PROTOCOL = "http"

var fromFlag *string
var baseUrlFlag *string
var useRouteParam *bool
var sanitizeRouteParam *bool

func main() {
	helpers.LoadEnv()
	flagParser()
	flagModifier()

	postmanCollection := postman.Collection{}

	collectionFrom(&postmanCollection)

	helpers.ExportToFileAsJson(postmanCollection, "collection.json")
}

func flagParser() {
	baseUrlFlag = flag.String("base-url", BASE_URL, "Custom base url")
	fromFlag = flag.String("from", laravel.PHP73_LARAVEL8, "Laravel version \nSupported: php73-laravel8, php81-laravel9")
	useRouteParam = flag.Bool("use-route-param", false, "Use route parameter")
	sanitizeRouteParam = flag.Bool("sanitize-route-param", false, "Sanitize route parameter")

	flag.Parse()

	if *useRouteParam && *sanitizeRouteParam {
		log.Fatal("cannot use both use-route-param and sanitize-route-param")
	}
}

func collectionFrom(postmanCollection *postman.Collection) {
	switch *fromFlag {
	case laravel.PHP73_LARAVEL8:
		laravel.Version = laravel.PHP73_LARAVEL8
		*postmanCollection = laravel.MakeCollection(useRouteParam, sanitizeRouteParam)
	case laravel.PHP81_LARAVEL9:
		laravel.Version = laravel.PHP81_LARAVEL9
		*postmanCollection = laravel.MakeCollection(useRouteParam, sanitizeRouteParam)
	}
}

func flagModifier() {
	if !strings.Contains(*baseUrlFlag, "http://") && !strings.Contains(*baseUrlFlag, "https://") {
		baseUrl := PROTOCOL + "://" + *baseUrlFlag
		*baseUrlFlag = baseUrl
	}
}
