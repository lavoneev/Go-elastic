package main

import (
	"encoding/json"
	"es/internal/db"
	"es/internal/models"
	"net/http"
	"strconv"
	"text/template"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		http.Error(w, "Invalid 'page' value: 'foo'.", http.StatusBadRequest)
		return
	}

	esClient, err := db.NewClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	store := db.NewElasticsearchStore(esClient)

	limit := 10
	offset := (page - 1) * limit
	places, total, err := store.GetPlaces(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl, err := template.ParseFiles("templates/places.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	LastPage := float64(total) / float64(limit)
	if float64(int(LastPage)) < LastPage {
		LastPage += 1
	}

	data := struct {
		Places   []models.Place
		Total    int
		Page     int
		NextPage int
		PrevPage int
		LastPage int
	}{
		Places:   places,
		Total:    total,
		Page:     page,
		NextPage: page + 1,
		PrevPage: page - 1,
		LastPage: int(LastPage),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ApiPlacesHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		http.Error(w, "Invalid 'page' value: 'foo'.", http.StatusBadRequest)
		return
	}

	esClient, err := db.NewClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	store := db.NewElasticsearchStore(esClient)

	limit := 10
	offset := (page - 1) * limit
	places, total, err := store.GetPlaces(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	LastPage := float64(total) / float64(limit)
	if float64(int(LastPage)) < LastPage {
		LastPage += 1
	}

	data := struct {
		Places   []models.Place `json:"places"`
		Total    int            `json:"-"`
		Page     int            `json:"-"`
		PrevPage int            `json:"prev_page"`
		NextPage int            `json:"next_page"`
		LastPage int            `json:"last_page"`
	}{
		Places:   places,
		Total:    total,
		Page:     page,
		NextPage: page + 1,
		PrevPage: page - 1,
		LastPage: int(LastPage),
	}

	w.Header().Set("Content-type", "application/json")

	JSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if _, err := w.Write(JSON); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ApiRecommendHandler(w http.ResponseWriter, r *http.Request) {
	latStr, lonStr := r.URL.Query().Get("lat"), r.URL.Query().Get("lon")
	lat, latErr := strconv.ParseFloat(latStr, 64)
	lon, lonErr := strconv.ParseFloat(lonStr, 64)
	if latErr != nil || lonErr != nil {
		http.Error(w, latErr.Error()+lonErr.Error(), http.StatusInternalServerError)
	}

	esClient, err := db.NewClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	store := db.NewElasticsearchStore(esClient)

	limit := 3
	places, err := store.GetClosetPlaces(limit, lat, lon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")

	JSON, err := json.MarshalIndent(places, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if _, err := w.Write(JSON); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func main() {
	http.HandleFunc("/", RootHandler)

	http.HandleFunc("/api/places", ApiPlacesHandler)

	http.HandleFunc("/api/recommend", ApiRecommendHandler)

	http.ListenAndServe(":8888", nil)

}
