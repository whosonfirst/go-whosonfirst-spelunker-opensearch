package opensearch

import (
	"context"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func (s *OpenSearchSpelunker) GetConcordances(ctx context.Context) (*spelunker.Faceting, error) {

	return nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) HasConcordance(ctx context.Context, pg_opts pagination.Options, namespace string, predicate string, value string, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	return nil, nil, spelunker.ErrNotImplemented
}
