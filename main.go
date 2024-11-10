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

var collectionNameFlag *string
var fromFlag *string
var baseUrlFlag *string
var projectLocation *string
var debug *bool

func main() {
	helpers.LoadEnv()
	flagParser()

	postmanCollection := postman.Collection{}

	collectionFrom(&postmanCollection)

	if *debug {
		printResult(postmanCollection)
	}

	helpers.ExportToFileAsJson(postmanCollection, "collection.json")
}

func flagParser() {
	collectionNameFlag = flag.String("name", "collection API", "Custom postman collection name")
	baseUrlFlag = flag.String("base-url", BASE_URL, "Custom base url")
	fromFlag = flag.String("from", laravel.PHP73_LARAVEL8, "Laravel version \nSupported: php73-laravel8, php81-laravel9, php81-laravel10")
	useRouteParam := flag.Bool("use-route-param", false, "Use route parameter")
	removeRouteParam := flag.Bool("remove-route-param", false, "Remove route parameter")
	projectLocation = flag.String("location", "", "Project location")
	debug = flag.Bool("debug", false, "Show collection item in console")

	flag.Parse()

	baseUrlSanitation()

	if *useRouteParam && *removeRouteParam {
		log.Fatal("cannot use both use-route-param and remove-route-param")
	}
	laravel.PathSetting = laravel.IGNORE_ROUTE
	if *removeRouteParam {
		laravel.PathSetting = laravel.REMOVE_ROUTE
	}
	if *useRouteParam {
		laravel.PathSetting = laravel.USE_ROUTE
	}
}

func baseUrlSanitation() {
	full := strings.Split(*baseUrlFlag, "://")
	if len(full) > 1 {
		protocol := full[0]
		base := full[1]

		if protocol != "http" && protocol != "https" {
			*baseUrlFlag = PROTOCOL + "://" + base
		}
	} else {
		*baseUrlFlag = PROTOCOL + "://" + *baseUrlFlag
	}
}

func collectionFrom(postmanCollection *postman.Collection) {

	for _, supportedVersion := range []string{
		laravel.PHP73_LARAVEL8,
		laravel.PHP81_LARAVEL9,
		laravel.PHP81_LARAVEL10,
	} {
		if *fromFlag == supportedVersion {
			laravel.Version = supportedVersion
			laravel.Location = *projectLocation
			laravel.BaseUrl = *baseUrlFlag
			laravel.CollectionName = *collectionNameFlag
			*postmanCollection = laravel.MakeCollection()
			return
		}
	}

	log.Fatal("Unsupported php - framework version")
}

func printResult(postmanCollection postman.Collection) {
	fmt.Println("base-url:", *baseUrlFlag, "\n", "--------------------------------")

	for _, item := range postmanCollection.Items {
		fmt.Println(item.Name)
		for _, subItem := range item.Items {
			fmt.Println("--", "[", subItem.Request.Method, "]", subItem.Name, ":", subItem.Request.Url.Raw)
		}
	}
}
