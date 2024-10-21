package main

import (
	"apitodocs/helpers"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const CLOSURE string = "Closure"

type RequestUrlQuery struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RequestUrl struct {
	Raw   string            `json:"raw"`
	Host  []string          `json:"host"`
	Path  []string          `json:"path"`
	Query []RequestUrlQuery `json:"query"`
}

type RequestBody struct {
	Mode    string `json:"mode"`
	Raw     string `json:"raw"`
	Options struct {
		Raw struct {
			Language string `json:"language"`
		} `json:"raw"`
	} `json:"options"`
}

type RequestHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type ItemRequest struct {
	Method  string          `json:"method"`
	Headers []RequestHeader `json:"header"`
	Url     RequestUrl      `json:"url"`
	Body    RequestBody     `json:"body"`
}

type PostmanCollectionInfo struct {
	PostmanId  string `json:"_postman_id"`
	Name       string `json:"name"`
	Schema     string `json:"schema"`
	ExporterId string `json:"_exporter_id"`
}

type PostmanCollectionItem struct {
	FormatPath string                  `json:"-"`
	Name       string                  `json:"name"`
	Items      []PostmanCollectionItem `json:"item"`
	Request    ItemRequest             `json:"request"`
	Response   []struct{}              `json:"response"`
}

type PostmanCollection struct {
	Info  PostmanCollectionInfo   `json:"info"`
	Items []PostmanCollectionItem `json:"item"`
}

func pathSliceModifier(path []string, route Route) []string {
	if route.Action != CLOSURE {
		if strings.Split(route.Action, "@")[1] == "index" {
			newName := "index"
			path = append(path, newName)
		}

		if strings.Contains(path[len(path)-1], "{") {
			path[len(path)-1] = strings.Split(route.Action, "@")[1]
		}
	}

	return path
}

type Route struct {
	Name       string   `json:"name"`
	Method     string   `json:"method"`
	Uri        string   `json:"uri"`
	Action     string   `json:"action"`
	Middleware []string `json:"middleware"`
}

func collectionItemNaming(path []string, isFolder bool) string {
	targetIndex := len(path) - 1
	if isFolder {
		targetIndex = len(path) - 2
	}

	return cases.Title(language.Tag{}).String(strings.Replace(path[targetIndex], "-", " ", -1))
}

func makeCollection(routes *[]Route) []PostmanCollectionItem {
	collectionItems := []PostmanCollectionItem{}

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

		newItem := PostmanCollectionItem{
			Name: collectionItemNaming(formatedPathSlice, false),
			Request: ItemRequest{
				Method: route.Method,
				Headers: []RequestHeader{
					{
						Key:   "Content-Type",
						Value: "application/json",
						Type:  "text",
					},
				},
				Url: RequestUrl{
					Raw:  "http://localhost:8000/" + route.Uri,
					Host: []string{"localhost:8000"},
					Path: pathSlice,
				},
			},
		}

		collectionFolderName := collectionItemNaming(formatedPathSlice, true)

		if len(collectionItems) == 0 {
			collectionItems = append(collectionItems, PostmanCollectionItem{
				FormatPath: formatedPath,
				Name:       collectionFolderName,
				Items:      []PostmanCollectionItem{newItem},
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
					collectionItems = append(collectionItems, PostmanCollectionItem{
						FormatPath: formatedPath,
						Name:       collectionFolderName,
						Items:      []PostmanCollectionItem{newItem},
					})
				}
			}
		}
	}

	return collectionItems
}

func main() {
	err := godotenv.Load(helpers.GetCurrentPath() + ".env")
	if err != nil {
		log.Fatal(err.Error())
	}

	executable := os.Getenv("EXECUTABLE")
	executables := append(
		strings.Split(executable, " "),
		"route:list", "--columns=name,method,uri,middleware,action", "--json", "--sort=uri",
	)

	output, err := exec.Command(executables[0], executables[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	routes := []Route{}
	err = json.Unmarshal(output, &routes)
	if err != nil {
		log.Fatal(err)
	}

	newPostmanCollection := PostmanCollection{
		Info: PostmanCollectionInfo{
			PostmanId:  uuid.New().String(),
			Name:       "Api to collection",
			Schema:     "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
			ExporterId: strconv.Itoa(int(time.Now().Unix())),
		},
		Items: makeCollection(&routes),
	}

	for _, item := range newPostmanCollection.Items {
		fmt.Println(item.Name)
		for _, subItem := range item.Items {
			fmt.Println("  -[", subItem.Request.Method, "]", subItem.Name)
		}
	}

	result, err := json.MarshalIndent(newPostmanCollection, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(helpers.GetCurrentPath()+"collection.json", result, 0644)
}
