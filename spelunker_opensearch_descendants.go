package opensearch

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aaronland/go-pagination"
	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func (s *OpenSearchSpelunker) GetDescendants(ctx context.Context, pg_opts pagination.Options, id int64, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.descendantsQuery(id, filters)
	sz := int(pg_opts.PerPage())

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
		Size: &sz,
		// pagination offset, scroll wah-wah here...
	}

	body, err := s.search(ctx, req)

	if err != nil {
		slog.Error("Count descendants query failed", "q", q)
		return nil, nil, fmt.Errorf("Failed to retrieve %d, %w", id, err)
	}

	return s.searchResultsToSPR(ctx, pg_opts, body)
}

func (s *OpenSearchSpelunker) GetDescendantsFaceted(ctx context.Context, id int64, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q := s.descendantsFacetedQuery(id, filters, facets)
	sz := 0

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
		Size: &sz,
	}

	return s.facet(ctx, req, facets)
}

func (s *OpenSearchSpelunker) CountDescendants(ctx context.Context, id int64) (int64, error) {

	filters := make([]spelunker.Filter, 0)

	q := s.descendantsQuery(id, filters)
	sz := 0

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
		Size: &sz,
	}

	body, err := s.search(ctx, req)

	if err != nil {
		slog.Error("Count descendants query failed", "q", q)
		return 0, fmt.Errorf("Failed to retrieve %d, %w", id, err)
	}

	r := gjson.GetBytes(body, "hits.total.value")
	return r.Int(), nil
}
