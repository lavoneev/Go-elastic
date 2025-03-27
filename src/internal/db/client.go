package db

import (
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
)

func NewClient() (*elasticsearch.Client, error) {
	retryBackoff := backoff.NewExponentialBackOff()

	cfg := elasticsearch.Config{
		Addresses:     []string{"http://localhost:9200"},
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(attempt int) time.Duration {
			if attempt == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
	}

	return elasticsearch.NewClient(cfg)
}
