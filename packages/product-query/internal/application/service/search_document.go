package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type EsRawHit[T any] struct {
	Source    T
	InnerHits json.RawMessage
}

type EsRawResult[T any] struct {
	Hits  []EsRawHit[T]
	Total int64
}

func SearchDocuments[T any](ctx context.Context, esClient *elasticsearch.TypedClient, index string, queryBody map[string]any) (*EsRawResult[T], error) {
	rawJson, err := json.Marshal(queryBody)
	if err != nil {
		return nil, fmt.Errorf("marshal es query failed: %w", err)
	}

	res, err := esClient.Search().Index(index).Raw(bytes.NewReader(rawJson)).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch transport crash: %w", err)
	}

	results := make([]EsRawHit[T], 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		var source T
		if err := json.Unmarshal(hit.Source_, &source); err != nil {
			return nil, fmt.Errorf("unmarshal es source to struct failed: %w", err)
		}
		innerHits, err := json.Marshal(hit.InnerHits)
		if err != nil {
			return nil, fmt.Errorf("marshal es inner hits failed: %w", err)
		}

		results = append(results, EsRawHit[T]{
			Source:    source,
			InnerHits: innerHits,
		})
	}

	return &EsRawResult[T]{
		Hits:  results,
		Total: totalHits(res.Hits),
	}, nil
}

func DecodeInnerHits[K any](rawInnerHits json.RawMessage, path string) ([]K, error) {
	if len(rawInnerHits) == 0 || string(rawInnerHits) == "null" {
		return nil, nil
	}

	var innerHitsMap map[string]struct {
		Hits struct {
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.Unmarshal(rawInnerHits, &innerHitsMap); err != nil {
		return nil, fmt.Errorf("failed to parse raw inner_hits structure: %w", err)
	}

	targetPathData, exists := innerHitsMap[path]
	if !exists {
		return nil, nil
	}

	esHits := targetPathData.Hits.Hits
	items := make([]K, 0, len(esHits))
	for _, hit := range esHits {
		var item K
		if err := json.Unmarshal(hit.Source, &item); err != nil {
			return nil, fmt.Errorf("failed to unmarshal inner_hit _source to target type: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}

func totalHits(hits types.HitsMetadata) int64 {
	if hits.Total == nil {
		return 0
	}
	return hits.Total.Value
}
