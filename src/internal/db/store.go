package db

import (
	"bytes"
	"encoding/json"
	"es/internal/models"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]models.Place, int, error)
	GetClosetPlaces(limit int, lat, lon float64) ([]models.Place, error)
}

type ElasticsearchStore struct {
	client *elasticsearch.Client
}

func NewElasticsearchStore(client *elasticsearch.Client) Store {
	return &ElasticsearchStore{client: client}
}

func (es *ElasticsearchStore) GetPlaces(limit int, offset int) ([]models.Place, int, error) {
	query := map[string]any{
		"from": offset,
		"size": limit,
		"query": map[string]any{
			"match_all": map[string]any{},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, fmt.Errorf("error encoding query: %s", err)
	}

	res, err := es.client.Search(
		es.client.Search.WithIndex("places"),
		es.client.Search.WithBody(&buf),
		es.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("error searching through es: %s", err)
	}
	defer res.Body.Close()

	var EsSearch models.EsSearch

	if err := json.NewDecoder(res.Body).Decode(&EsSearch); err != nil {
		return nil, 0, fmt.Errorf("error parsing the response body: %s", err)
	}

	places := make([]models.Place, len(EsSearch.Hits.Hits))
	for i, hit := range EsSearch.Hits.Hits {
		places[i] = hit.Source
	}

	return places, EsSearch.Hits.Total.Value, nil
}

func (es *ElasticsearchStore) GetClosetPlaces(limit int, lat, lon float64) ([]models.Place, error) {
	query := map[string]any{
		"size": limit,
		"sort": []map[string]any{
			{
				"_geo_distance": map[string]any{
					"location":        map[string]float64{"lat": lat, "lon": lon},
					"order":           "asc",
					"unit":            "km",
					"distance_type":   "arc",
					"ignore_unmapped": true,
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %s", err)
	}

	res, err := es.client.Search(
		es.client.Search.WithIndex("places"),
		es.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("error searching through es: %s", err)
	}
	defer res.Body.Close()

	var EsSearch models.EsSearch

	if err := json.NewDecoder(res.Body).Decode(&EsSearch); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	places := make([]models.Place, len(EsSearch.Hits.Hits))
	for i, hit := range EsSearch.Hits.Hits {
		places[i] = hit.Source
	}

	return places, nil
}
