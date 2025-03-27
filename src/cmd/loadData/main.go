package main

import (
	"context"
	"encoding/csv"
	"es/internal/db"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
)

func main() {

	var countSuccessful uint64

	esClient, err := db.NewClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	mapping, err := os.ReadFile("configs/schema.json")
	if err != nil {
		log.Fatalf("Error reading json file: %s", err)
	}

	err = db.SetupIndex(esClient, "places", mapping)
	if err != nil {
		log.Fatalf("Index setup failed: %s", err)
	}

	bi, err := db.NewBulkIndexer(esClient, "places")
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	places := loadPlaces()

	fmt.Printf("Initialized %d places\n", len(places))

	var wg sync.WaitGroup

	for _, place := range places {
		wg.Add(1)
		go func(pl []string) {
			defer wg.Done()
			err = db.AddPlaceToIndex(bi, pl, &countSuccessful)
			if err != nil {
				log.Printf("AddPlaceToIndex: %s", err)
			}
		}(place)
	}

	wg.Wait()

	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	fmt.Printf("Successfully indexed %d/%d places\n", atomic.LoadUint64(&countSuccessful), len(places))

}

func loadPlaces() [][]string {
	filePath := "../materials/data.csv"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	defer file.Close()

	csvr := csv.NewReader(file)
	csvr.Comma = '\t'
	places, err := csvr.ReadAll()
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	return places[1:]
}
