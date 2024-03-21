package opensearch

import (
	"context"
	"fmt"
	"strings"
	
	"github.com/aaronland/go-pagination"
	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"	
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func (s *OpenSearchSpelunker) GetConcordances(ctx context.Context) (*spelunker.Faceting, error) {

	c_facet := spelunker.NewFacet("concordances_sources.keyword")

	facets := []*spelunker.Facet{
		c_facet,
	}

	q := s.matchAllFacetedQuery(facets)
	sz := 0

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
		Size: &sz,
	}

	f, err := s.facet(ctx, req, facets)

	if err != nil {
		return nil, fmt.Errorf("Failed to facet concordances, %w", err)
	}

	return f[0], nil
}

func (s *OpenSearchSpelunker) HasConcordance(ctx context.Context, pg_opts pagination.Options, namespace string, predicate string, value any, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.hasConcordanceQuery(namespace, predicate, value, filters)
	return s.searchPaginated(ctx, pg_opts, q)
	
	return nil, nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) HasConcordanceFaceted(ctx context.Context, namespace string, predicate string, value any, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q := s.hasConcordanceFacetedQuery(namespace, predicate, value, filters, facets)
	sz := 0

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
		Size: &sz,
	}

	return s.facet(ctx, req, facets)	
}
