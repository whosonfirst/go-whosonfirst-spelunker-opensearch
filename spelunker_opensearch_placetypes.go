package opensearch

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aaronland/go-pagination"
	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func (s *OpenSearchSpelunker) HasPlacetype(ctx context.Context, pg_opts pagination.Options, pt *placetypes.WOFPlacetype, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.hasPlacetypesQuery(pt.Name, filters)
	sz := int(pg_opts.PerPage())

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
		Size: &sz,
		// pagination offset, scroll wah-wah here...
	}

	body, err := s.search(ctx, req)

	if err != nil {
		slog.Error("Placetypes query failed", "q", q)
		return nil, nil, fmt.Errorf("Failed to retrieve placetypes, %w", err)
	}

	return s.searchResultsToSPR(ctx, pg_opts, body)
}

func (s *OpenSearchSpelunker) GetPlacetypes(ctx context.Context) (*spelunker.Faceting, error) {

	pt_facet := spelunker.NewFacet("placetype")

	facets := []*spelunker.Facet{
		pt_facet,
	}

	q := s.matchAllFacetedQuery(facets)
	sz := 0

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
		Size: &sz,
	}

	f, err := s.facet(ctx, req, facets)

	if err != nil {
		return nil, fmt.Errorf("Failed to facet placetypes, %w", err)
	}

	return f[0], nil
}
