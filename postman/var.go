package postman

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

type CollectionInfo struct {
	PostmanId  string `json:"_postman_id"`
	Name       string `json:"name"`
	Schema     string `json:"schema"`
	ExporterId string `json:"_exporter_id"`
}

type CollectionItem struct {
	FormatPath string           `json:"-"`
	Name       string           `json:"name"`
	Items      []CollectionItem `json:"item"`
	Request    ItemRequest      `json:"request"`
	Response   []struct{}       `json:"response"`
}

type Collection struct {
	Info  CollectionInfo   `json:"info"`
	Items []CollectionItem `json:"item"`
}
