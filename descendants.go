package opensearch

import (
	"context"
	"strings"

	"github.com/aaronland/go-pagination"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func (s *OpenSearchSpelunker) GetDescendants(ctx context.Context, pg_opts pagination.Options, id int64, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.descendantsQuery(id, filters)
	return s.searchPaginated(ctx, pg_opts, q)
}

func (s *OpenSearchSpelunker) GetDescendantsFaceted(ctx context.Context, id int64, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q := s.descendantsFacetedQuery(id, filters, facets)
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

func (s *OpenSearchSpelunker) CountDescendants(ctx context.Context, id int64) (int64, error) {

	filters := make([]spelunker.Filter, 0)

	q := s.descendantsQuery(id, filters)
	return s.countForQuery(ctx, q)
}
