package opensearch

import (
	"context"
	"strings"
	"time"
	
	"github.com/aaronland/go-pagination"
	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"	
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func (s *OpenSearchSpelunker) GetRecent(ctx context.Context, pg_opts pagination.Options, d time.Duration, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.getRecentQuery(d, filters)
	return s.searchPaginated(ctx, pg_opts, q)
}

func (s *OpenSearchSpelunker) GetRecentFaceted(ctx context.Context, d time.Duration, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q := s.getRecentFacetedQuery(d, filters, facets)
	sz := 0

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
		Size: &sz,
	}

	return s.facet(ctx, req, facets)	
}
