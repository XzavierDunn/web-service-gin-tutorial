package models

type AlbumRecord struct {
	PK     string  `json:"pk"`
	SK     string  `json:"sk"`
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

type Response struct {
	Message string `json:"message"`
}
