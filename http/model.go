package http

type RequestType string

const (
	GET                               = "GET"
	POST                              = "POST"
	PUT                               = "PUT"
	DELETE                            = "DELETE"
	PATCH                             = "PATCH"
	TypeJSON              RequestType = "json"
	TypeXML               RequestType = "xml"
	TypeUrlencoded        RequestType = "urlencoded"
	TypeForm              RequestType = "form"
	TypeFormData          RequestType = "form-data"
	TypeMultipartFormData RequestType = "multipart-form-data"
)

var types = map[RequestType]string{
	TypeJSON:              "application/json",
	TypeXML:               "application/xml",
	TypeUrlencoded:        "application/x-www-form-urlencoded",
	TypeForm:              "application/x-www-form-urlencoded",
	TypeFormData:          "application/x-www-form-urlencoded",
	TypeMultipartFormData: "multipart/form-data",
}

type File struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}
