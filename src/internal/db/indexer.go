package db

import (
	"bytes"
	"context"
	"encoding/json"
	"es/internal/models"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

func SetupIndex(es *elasticsearch.Client, indexName string, mapping []byte) error {
	res, err := es.Indices.Delete(
		[]string{indexName},
		es.Indices.Delete.WithIgnoreUnavailable(true),
	)

	if err != nil || res.IsError() {
		return err
	}
	res.Body.Close()

	res, err = es.Indices.Create(
		indexName,
		es.Indices.Create.WithBody(bytes.NewReader(mapping)),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Error response: %s", res.String())
	}
	return nil
}

func NewBulkIndexer(es *elasticsearch.Client, indexName string) (esutil.BulkIndexer, error) {
	biCfg := esutil.BulkIndexerConfig{
		Index:         indexName,
		Client:        es,
		NumWorkers:    8,
		FlushBytes:    5 * 1024 * 1024,
		FlushInterval: 30 * time.Second,
	}

	return esutil.NewBulkIndexer(biCfg)
}

func AddPlaceToIndex(bi esutil.BulkIndexer, place []string, countSuccessful *uint64) error {
	lon, err := strconv.ParseFloat(place[4], 64)
	if err != nil {
		return err
	}
	lat, err := strconv.ParseFloat(place[5], 64)
	if err != nil {
		return err
	}

	pl := models.Place{
		Name:    place[1],
		Address: place[2],
		Phone:   place[3],
		Location: models.GeoPoint{
			Longitude: lon,
			Latitude:  lat,
		},
	}

	data, err := json.Marshal(pl)
	if err != nil {
		return err
	}

	err = bi.Add(
		context.Background(),
		esutil.BulkIndexerItem{
			Action: "index",

			DocumentID: place[0],

			Body: bytes.NewReader(data),

			OnSuccess: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem) {
				atomic.AddUint64(countSuccessful, 1)
			},

			OnFailure: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: %s", err)
				} else {
					log.Printf("ERROR: %s: %s", biri.Error.Type, biri.Error.Reason)
				}
			},
		},
	)

	if err != nil {
		return err
	}

	return nil
}
