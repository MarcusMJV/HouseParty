package models

type Song struct {
	Id		  string `json:"id"`
	URI	   string `json:"uri"`
	Name	 string `json:"name"`
	Artists []string `json:"artists"`
	Album	 string `json:"album"`
	Image	 Image `json:"image"`
	DurationMs int `json:"duration_ms"`
	Explicit bool `json:"explicit"`
	ExternalURL string `json:"external_url"`
}

type Image struct {
	URL string `json:"url"`
	Width int `json:"width"`
	Height int `json:"height"`
}