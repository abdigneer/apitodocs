package laravel

import "strings"

const CLOSURE string = "Closure"

type route struct {
	Name       string   `json:"name"`
	Method     string   `json:"method"`
	Uri        string   `json:"uri"`
	Action     string   `json:"action"`
	Middleware []string `json:"middleware"`
}

func removingRouteParam(routeUri string, removeRouteParam *bool) string {
	if *removeRouteParam {
		routeUri = strings.Replace(routeUri, "{", "", -1)
		routeUri = strings.Replace(routeUri, "}", "", -1)
	}
	return routeUri
}

func usingRouteParam(routeUri *string, useRouteParam *bool) {
	if *useRouteParam {
		*routeUri = strings.Replace(*routeUri, "{", "{{", -1)
		*routeUri = strings.Replace(*routeUri, "}", "}}", -1)
	}
}
