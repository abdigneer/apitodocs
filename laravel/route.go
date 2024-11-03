package laravel

import (
	"strings"
)

const CLOSURE string = "Closure"

type route struct {
	Name       string   `json:"name"`
	Method     string   `json:"method"`
	Uri        string   `json:"uri"`
	Action     string   `json:"action"`
	Middleware []string `json:"middleware"`
}

func pathSliceModifier(path []string, route route) []string {
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
