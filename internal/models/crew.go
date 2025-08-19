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
