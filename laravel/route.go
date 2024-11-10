package laravel

import (
	"strings"
)

const CLOSURE string = "Closure"

const IGNORE_ROUTE = 0
const USE_ROUTE = 1
const REMOVE_ROUTE = 2

var PathSetting int

type route struct {
	Name       string   `json:"name"`
	Method     string   `json:"method"`
	Uri        string   `json:"uri"`
	Action     string   `json:"action"`
	Middleware []string `json:"middleware"`
}

func uriModifier(routeUri string) string {
	if PathSetting == REMOVE_ROUTE {
		routeUri = removeRouteParamSyntax(routeUri)
	}
	if PathSetting == USE_ROUTE {
		routeUri = strings.Replace(routeUri, "{", "{{", -1)
		routeUri = strings.Replace(routeUri, "}", "}}", -1)
	}

	return routeUri
}

func removeRouteParamSyntax(text string) string {
	if strings.Contains(text, "{") && strings.Contains(text, "}") {
		text = strings.Replace(text, "{", "", -1)
		text = strings.Replace(text, "}", "", -1)
	}

	return text
}
