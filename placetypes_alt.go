package opensearch

import (
	"context"
	"fmt"
	"strings"

	"github.com/aaronland/go-pagination"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func (s *OpenSearchSpelunker) GetAlternatePlacetypes(ctx context.Context) (*spelunker.Faceting, error) {

	pt_facet := spelunker.NewFacet("placetypealt")

	facets := []*spelunker.Facet{
		pt_facet,
	}

	q := s.matchAllFacetedQuery(facets)
	sz := 0

	req := &opensearchapi.SearchReq{
		Indices: []string{
			s.index,
		},
		Body: strings.NewReader(q),
		Params: opensearchapi.SearchParams{
			Size: &sz,
		},
	}

	f, err := s.facet(ctx, req, facets)

	if err != nil {
		return nil, fmt.Errorf("Failed to facet alternate placetypes, %w", err)
	}

	return f[0], nil
}

func (s *OpenSearchSpelunker) HasAlternatePlacetype(ctx context.Context, pg_opts pagination.Options, pt string, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.hasAlternatePlacetypeQuery(pt, filters)
	return s.searchPaginated(ctx, pg_opts, q)
}

func (s *OpenSearchSpelunker) HasAlternatePlacetypeFaceted(ctx context.Context, pt string, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q := s.hasAlternatePlacetypeFacetedQuery(pt, filters, facets)
	sz := 0

	req := &opensearchapi.SearchReq{
		Indices: []string{
			s.index,
		},
		Body: strings.NewReader(q),
		Params: opensearchapi.SearchParams{
			Size: &sz,
		},
	}

	return s.facet(ctx, req, facets)
}
