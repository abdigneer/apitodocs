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

func collectionItemNaming(path []string, isFolder bool) string {
	targetIndex := len(path) - 1
	if isFolder {
		targetIndex = len(path) - 2
	}

	return cases.Title(language.Tag{}).String(strings.Replace(path[targetIndex], "-", " ", -1))
}

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

func removingRouteParam(routeUri string, removeRouteParam *bool) string {
	if *removeRouteParam {
		routeUri = strings.Replace(routeUri, "{", "", -1)
		routeUri = strings.Replace(routeUri, "}", "", -1)
	}
	return routeUri
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

		formatedPathSlice := pathSliceModifier(pathSlice, route)
		formatedPath := strings.Join(formatedPathSlice[:len(formatedPathSlice)-1], "/")

		newItem := postman.CollectionItem{
			Name: removingRouteParam(collectionItemNaming(formatedPathSlice, false), removeRouteParam),
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

		collectionFolderName := collectionItemNaming(formatedPathSlice, true)

		if len(collectionItems) == 0 {
			collectionItems = append(collectionItems, postman.CollectionItem{
				FormatPath: formatedPath,
				Name:       removingRouteParam(collectionFolderName, removeRouteParam),
				Items:      []postman.CollectionItem{newItem},
			})
		} else {
			collectionItem := collectionItems[len(collectionItems)-1]

			if collectionItem.FormatPath == formatedPath {
				collectionItems[len(collectionItems)-1].Items = append(collectionItems[len(collectionItems)-1].Items, newItem)
			} else {
				index := findSameItemByFolderName(collectionItems, collectionFolderName)
				if index > -1 {
					collectionItems[index].Items = append(collectionItems[index].Items, newItem)
				} else {
					index := -1
					if len(formatedPathSlice) >= 3 {
						routeSub := formatedPathSlice[len(formatedPathSlice)-2]
						if strings.Contains(routeSub, "{") && strings.Contains(routeSub, "}") {
							routeSub = strings.Replace(routeSub, "{", "", -1)
							routeSub = strings.Replace(routeSub, "}", "", -1)
						}

						if routeSub == formatedPathSlice[len(formatedPathSlice)-3] || routeSub == "id" {
							index = findSameItemByFolderName(
								collectionItems,
								cases.Title(language.Tag{}).String(strings.Replace(formatedPathSlice[len(formatedPathSlice)-3], "-", " ", -1)))
						}
					}

					if index > -1 {
						collectionItems[index].Items = append(collectionItems[index].Items, newItem)
					} else {
						collectionItems = append(collectionItems, postman.CollectionItem{
							FormatPath: formatedPath,
							Name:       removingRouteParam(collectionFolderName, removeRouteParam),
							Items:      []postman.CollectionItem{newItem},
						})
					}
				}
			}
		}
	}

	return collectionItems
}

func findSameItemByFolderName(collectionItems []postman.CollectionItem, collectionFolderName string) int {
	sameItemIndex := -1
	for i, collectionItem := range collectionItems {
		if collectionItem.Name == collectionFolderName {
			sameItemIndex = i
		}
	}

	return sameItemIndex
}
