package main

import (
	"apitodocs/helpers"
	"apitodocs/laravel"
	"apitodocs/postman"
	"flag"
	"fmt"
	"log"
	"strings"
)

const BASE_URL = "http://localhost:8080"
const PROTOCOL = "http"

var fromFlag *string
var baseUrlFlag *string
var useRouteParam *bool
var removeRouteParam *bool
var projectLocation *string

func main() {
	helpers.LoadEnv()
	flagParser()
	flagModifier()

	postmanCollection := postman.Collection{}

	collectionFrom(&postmanCollection)

	printCollection(postmanCollection)

	helpers.ExportToFileAsJson(postmanCollection, "collection.json")
}

func flagParser() {
	baseUrlFlag = flag.String("base-url", BASE_URL, "Custom base url")
	fromFlag = flag.String("from", laravel.PHP73_LARAVEL8, "Laravel version \nSupported: php73-laravel8, php81-laravel9")
	useRouteParam = flag.Bool("use-route-param", false, "Use route parameter")
	removeRouteParam = flag.Bool("remove-route-param", false, "Remove route parameter")
	projectLocation = flag.String("location", "", "Project location")

	flag.Parse()

	if *useRouteParam && *removeRouteParam {
		log.Fatal("cannot use both use-route-param and remove-route-param")
	}
}

func collectionFrom(postmanCollection *postman.Collection) {
	switch *fromFlag {
	case laravel.PHP73_LARAVEL8:
		laravel.Version = laravel.PHP73_LARAVEL8
	case laravel.PHP81_LARAVEL9:
		laravel.Version = laravel.PHP81_LARAVEL9
	}

	laravel.Location = *projectLocation
	*postmanCollection = laravel.MakeCollection(useRouteParam, removeRouteParam)
}

func flagModifier() {
	if !strings.Contains(*baseUrlFlag, "http://") && !strings.Contains(*baseUrlFlag, "https://") {
		baseUrl := PROTOCOL + "://" + *baseUrlFlag
		*baseUrlFlag = baseUrl
	}
}

func printCollection(postmanCollection postman.Collection) {
	for _, item := range postmanCollection.Items {
		fmt.Println(item.Name)
		for _, subItem := range item.Items {
			fmt.Println("--", "[", subItem.Request.Method, "]", subItem.Name, ":", subItem.Request.Url.Raw)
		}
	}
}
