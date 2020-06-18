package constant

//Article is exported struct
type Article struct {
	Title  string `json:"title"`
	Desc   string `json:"description"`
	Author string `json:"author"`
}

// Articles is exported struct
type Articles []Article

// Sucess is exporeted struct
type Sucess struct {
	Message string `json:"message"`
}
