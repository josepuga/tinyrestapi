package main

// Handler is a "special" struct used by the Server

type Handler struct {
	Paths        []string          // All paths used for requests. ie: ["/books", "/libros", "/llibres"]
	Methods      map[string]bool   // Methods. Only true are avail. ie: ["GET"]true
	ReadOnlyMode bool              // To prevent call to a not allowed methods
	jsonData     JSONData          // The JSONData struct than contains all data
}

func NewHandler() *Handler {
	result := &Handler{
        jsonData: *NewJSONData(),
	}
	// By default all methods are disabled.
	result.Methods = map[string]bool{
		"POST":   false,
		"GET":    false,
		"PUT":    false,
		"PATCH":  false,
		"DELETE": false,
        "HEAD": false,
	}
	return result
}

// SetData sets the data to Items struct after replace aliases
func (h *Handler) SetData(data []byte) error {
	return h.jsonData.SetData(data)
}

// AddPath: Add a Path an its "variants" to the Handler. ie: "books" ==> "/books", "/books/:id"
func (h *Handler) AddPath(pathName string) {
	h.Paths = append(h.Paths, "/"+pathName)
	h.Paths = append(h.Paths, "/"+pathName+"/:id")
}

func (h *Handler) AddAlias(alias, internal string) {
    h.jsonData.AddAlias(alias, internal)
}