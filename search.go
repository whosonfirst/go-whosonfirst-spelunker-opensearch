package opensearch

import (
	"context"
	_ "fmt"
	"strings"
	
	"github.com/aaronland/go-pagination"
	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"	
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func (s *OpenSearchSpelunker) Search(ctx context.Context, pg_opts pagination.Options, search_opts *spelunker.SearchOptions, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.searchQuery(search_opts, filters)
	return s.searchPaginated(ctx, pg_opts, q)
}

func (s *OpenSearchSpelunker) SearchFaceted(ctx context.Context, search_opts *spelunker.SearchOptions, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q := s.searchFacetedQuery(search_opts, filters, facets)
	sz := 0

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
		Size: &sz,
	}

	return s.facet(ctx, req, facets)	
}
