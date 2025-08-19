package models

type ISSCrewResponse struct {
	People  []Astronaut `json:"people"`
	Message string      `json:"message"`
}

type Astronaut struct {
	Name  string `json:"name"`
	Craft string `json:"craft"`
}

type ISSCrewWithPhotosResponse struct {
	People  []AstronautWithPhoto `json:"people"`
	Message string               `json:"message"`
}

type AstronautWithPhoto struct {
	Name     string `json:"name"`
	ImageUrl string `json:"imgUrl"`
}

type NASAImagesResponse struct {
	Collection struct {
		Items []struct {
			Data []struct {
				Title       string   `json:"title"`
				Description string   `json:"description"`
				Keywords    []string `json:"keywords"`
			} `json:"data"`
			Links []struct {
				Href string `json:"href"`
				Rel  string `json:"rel"`
			} `json:"links"`
		} `json:"items"`
	} `json:"collection"`
}
