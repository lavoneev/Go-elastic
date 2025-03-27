package models

type Place struct {
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location GeoPoint `json:"location"`
}

type GeoPoint struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type EsSearch struct {
	Hits struct {
		Total struct {
			Value int `json:"value,omitempty"`
		} `json:"total,omitempty"`
		Hits []struct {
			Source Place `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
