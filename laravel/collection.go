package laravel

import (
	"apitodocs/helpers"
	"apitodocs/postman"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const ITEM_NAMING = 1
const PARENT_NAMING = 2
const CHILD_NAMING = 3

func MakeCollection() postman.Collection {
	routes := []route{}
	err := json.Unmarshal(execute(), &routes)
	if err != nil {
		log.Fatal(err)
	}

	helpers.ExportToFileAsJson(routes, "routes.json")

	return postman.Collection{
		Info:  makeInfo(),
		Items: makeItems(&routes),
	}
}

func makeInfo() postman.CollectionInfo {
	return postman.CollectionInfo{
		PostmanId:  uuid.New().String(),
		Name:       CollectionName,
		Schema:     "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		ExporterId: strconv.Itoa(int(time.Now().Unix())),
	}
}

func makeItems(routes *[]route) []postman.CollectionItem {
	collectionItems := []postman.CollectionItem{}

	for _, route := range *routes {
		if route.Uri == "" && route.Method == "" {
			continue
		}

		if route.Method == "GET|HEAD" {
			route.Method = "GET"
		}

		if route.Method == "PUT|PATCH" {
			route.Method = "PUT"
		}

		itemPathSlice := makeItemPath(strings.Split(route.Uri, "/"), route)
		itemUrl := uriModifier(route.Uri)
		item := postman.CollectionItem{
			Name: makeItemName(itemPathSlice, ITEM_NAMING),
			Request: postman.ItemRequest{
				Method: route.Method,
				Headers: []postman.RequestHeader{
					{
						Key:   "Content-Type",
						Value: "application/json",
						Type:  "text",
					},
				},
				Url: postman.RequestUrl{
					Raw:  BaseUrl + "/" + itemUrl,
					Host: []string{BaseUrl},
					Path: strings.Split(itemUrl, "/"),
				},
			},
		}

		parentItemPathString := strings.Join(itemPathSlice[:len(itemPathSlice)-1], "/")
		parentItemName := makeItemName(itemPathSlice, PARENT_NAMING)
		if len(collectionItems) == 0 {
			collectionItems = append(collectionItems, postman.CollectionItem{
				Path:  parentItemPathString,
				Name:  parentItemName,
				Items: []postman.CollectionItem{item},
			})
		} else {
			collectionItem := collectionItems[len(collectionItems)-1]

			if collectionItem.Path == parentItemPathString {
				collectionItems[len(collectionItems)-1].Items = append(collectionItems[len(collectionItems)-1].Items, item)
			} else {
				index := findSameItemByName(collectionItems, parentItemName)
				if index > -1 {
					collectionItems[index].Items = append(collectionItems[index].Items, item)
				} else {
					index := -1
					if len(itemPathSlice) >= 3 {
						segment := removeRouteParamSyntax(itemPathSlice[len(itemPathSlice)-2])
						// is route param?
						if segment == itemPathSlice[len(itemPathSlice)-3] || segment == "id" {
							index = findSameItemByName(collectionItems, makeItemName(itemPathSlice, CHILD_NAMING))
						}
					}

					if index > -1 {
						collectionItems[index].Items = append(collectionItems[index].Items, item)
					} else {
						collectionItems = append(collectionItems, postman.CollectionItem{
							Path:  parentItemPathString,
							Name:  parentItemName,
							Items: []postman.CollectionItem{item},
						})
					}
				}
			}
		}
	}

	return collectionItems
}

func makeItemPath(path []string, route route) []string {
	if route.Action != CLOSURE {
		if len(strings.Split(route.Action, "@")) > 1 {
			if strings.Split(route.Action, "@")[1] == "index" {
				newName := "index"
				path = append(path, newName)
			} else if strings.Split(route.Action, "@")[1] == "store" && route.Method == "POST" {
				newName := "store"
				path = append(path, newName)
			}
		}

		if strings.Contains(path[len(path)-1], "{") {
			path[len(path)-1] = strings.Split(route.Action, "@")[1]
		}
	}

	return path
}

func makeItemName(path []string, fromLast int) string {
	targetIndex := len(path) - fromLast
	return cases.Title(language.Tag{}).String(
		removeRouteParamSyntax(
			strings.Replace(path[targetIndex], "-", " ", -1),
		),
	)
}

func findSameItemByName(collectionItems []postman.CollectionItem, collectionItemName string) int {
	sameItemIndex := -1
	for i, collectionItem := range collectionItems {
		if collectionItem.Name == collectionItemName {
			sameItemIndex = i
		}
	}

	return sameItemIndex
}
