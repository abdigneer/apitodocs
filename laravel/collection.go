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
const FOLDER_NAMING = 2
const CHILD_PARAM_NAMING = 3

func MakeCollection(useRouteParam *bool, removeRouteParam *bool) postman.Collection {
	routes := []route{}
	err := json.Unmarshal(execute(), &routes)
	if err != nil {
		log.Fatal(err)
	}

	helpers.ExportToFileAsJson(routes, "routes.json")

	return postman.Collection{
		Info:  makeInfo(),
		Items: makeItems(&routes, useRouteParam, removeRouteParam),
	}
}

func makeInfo() postman.CollectionInfo {
	return postman.CollectionInfo{
		PostmanId:  uuid.New().String(),
		Name:       "Api to collection",
		Schema:     "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		ExporterId: strconv.Itoa(int(time.Now().Unix())),
	}
}

func makeItems(routes *[]route, useRouteParam *bool, removeRouteParam *bool) []postman.CollectionItem {
	collectionItems := []postman.CollectionItem{}

	for _, route := range *routes {
		if route.Uri == "" && route.Method == "" {
			continue
		}

		if route.Method == "GET|HEAD" {
			route.Method = "GET"
		}

		if *useRouteParam {
			route.Uri = strings.Replace(route.Uri, "{", "{{", -1)
			route.Uri = strings.Replace(route.Uri, "}", "}}", -1)
		}

		pathSlice := strings.Split(route.Uri, "/")

		itemPath := makeItemPath(pathSlice, route)
		itemPathString := strings.Join(itemPath[:len(itemPath)-1], "/")

		newItem := postman.CollectionItem{
			Name: removingRouteParam(makeItemName(itemPath, ITEM_NAMING), removeRouteParam),
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
					Raw:  "http://localhost:8000/" + removingRouteParam(route.Uri, removeRouteParam),
					Host: []string{"localhost:8000"},
					Path: pathSlice,
				},
			},
		}

		collectionFolderName := makeItemName(itemPath, FOLDER_NAMING)

		if len(collectionItems) == 0 {
			collectionItems = append(collectionItems, postman.CollectionItem{
				Path:  itemPathString,
				Name:  removingRouteParam(collectionFolderName, removeRouteParam),
				Items: []postman.CollectionItem{newItem},
			})
		} else {
			collectionItem := collectionItems[len(collectionItems)-1]

			if collectionItem.Path == itemPathString {
				collectionItems[len(collectionItems)-1].Items = append(collectionItems[len(collectionItems)-1].Items, newItem)
			} else {
				index := findSameItemByName(collectionItems, collectionFolderName)
				if index > -1 {
					collectionItems[index].Items = append(collectionItems[index].Items, newItem)
				} else {
					index := -1
					if len(itemPath) >= 3 {
						routeSub := itemPath[len(itemPath)-2]
						if strings.Contains(routeSub, "{") && strings.Contains(routeSub, "}") {
							routeSub = strings.Replace(routeSub, "{", "", -1)
							routeSub = strings.Replace(routeSub, "}", "", -1)
						}

						if routeSub == itemPath[len(itemPath)-3] || routeSub == "id" {
							index = findSameItemByName(
								collectionItems,
								makeItemName(itemPath, CHILD_PARAM_NAMING))
						}
					}

					if index > -1 {
						collectionItems[index].Items = append(collectionItems[index].Items, newItem)
					} else {
						collectionItems = append(collectionItems, postman.CollectionItem{
							Path:  itemPathString,
							Name:  removingRouteParam(collectionFolderName, removeRouteParam),
							Items: []postman.CollectionItem{newItem},
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

	return cases.Title(language.Tag{}).String(strings.Replace(path[targetIndex], "-", " ", -1))
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
