package main

import (
	"apitodocs/helpers"
	"apitodocs/laravel"
	"apitodocs/postman"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const LARAVEL = 1

func main() {
	err := godotenv.Load(helpers.GetCurrentPath() + ".env")
	if err != nil {
		log.Fatal(err.Error())
	}

	from := 1

	postmanCollection := postman.Collection{}
	switch from {
	case LARAVEL:
		postmanCollection = laravel.MakeCollection()
	}

	for _, item := range postmanCollection.Items {
		fmt.Println(item.Name)
		for _, subItem := range item.Items {
			fmt.Println("  -[", subItem.Request.Method, "]", subItem.Name)
		}
	}

	exportToFile(postmanCollection)
}

func exportToFile(postmanCollection postman.Collection) {
	result, err := json.MarshalIndent(postmanCollection, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(helpers.GetCurrentPath()+"collection.json", result, 0644)
}
