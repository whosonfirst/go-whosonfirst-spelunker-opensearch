package opensearch

import (
	"context"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func (s *OpenSearchSpelunker) HasPlacetype(ctx context.Context, pg_opts pagination.Options, pt *placetypes.WOFPlacetype, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	return nil, nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) GetPlacetypes(ctx context.Context) (*spelunker.Faceting, error) {

	return nil, spelunker.ErrNotImplemented
}
