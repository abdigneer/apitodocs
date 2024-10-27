package laravel

import (
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

func MakeCollection() postman.Collection {
	routes := []route{}
	err := json.Unmarshal(execute(), &routes)
	if err != nil {
		log.Fatal(err)
	}

	return postman.Collection{
		Info:  makeInfo(),
		Items: makeItems(&routes),
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

func makeItems(routes *[]route) []postman.CollectionItem {

	collectionItems := []postman.CollectionItem{}

	for _, route := range *routes {
		if route.Uri == "" && route.Method == "" {
			continue
		}

		if route.Method == "GET|HEAD" {
			route.Method = "GET"
		}

		pathSlice := strings.Split(route.Uri, "/")

		formatedPathSlice := pathSliceModifier(pathSlice, route)
		formatedPath := strings.Join(formatedPathSlice[:len(formatedPathSlice)-1], "/")

		newItem := postman.CollectionItem{
			Name: collectionItemNaming(formatedPathSlice, false),
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
					Raw:  "http://localhost:8000/" + route.Uri,
					Host: []string{"localhost:8000"},
					Path: pathSlice,
				},
			},
		}

		collectionFolderName := collectionItemNaming(formatedPathSlice, true)

		if len(collectionItems) == 0 {
			collectionItems = append(collectionItems, postman.CollectionItem{
				FormatPath: formatedPath,
				Name:       collectionFolderName,
				Items:      []postman.CollectionItem{newItem},
			})
		} else {
			collectionItem := collectionItems[len(collectionItems)-1]

			if collectionItem.FormatPath == formatedPath {
				if collectionItems[len(collectionItems)-1].FormatPath == collectionItem.FormatPath {
					collectionItems[len(collectionItems)-1].Items = append(collectionItems[len(collectionItems)-1].Items, newItem)
				} else {
					collectionItem.Items = append(collectionItem.Items, newItem)
				}
			} else {
				existsItemIndex := -1
				for i, collectionItem := range collectionItems {
					if collectionItem.Name == collectionFolderName {
						existsItemIndex = i
					}
				}

				if existsItemIndex >= 0 {
					collectionItems[existsItemIndex].Items = append(collectionItems[existsItemIndex].Items, newItem)
				} else {
					collectionItems = append(collectionItems, postman.CollectionItem{
						FormatPath: formatedPath,
						Name:       collectionFolderName,
						Items:      []postman.CollectionItem{newItem},
					})
				}
			}
		}
	}

	return collectionItems
}
