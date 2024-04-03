package opensearch

import (
	"context"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func (s *OpenSearchSpelunker) GetTags(ctx context.Context) (*spelunker.Faceting, error) {
	return nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) HasTag(ctx context.Context, pg_opts pagination.Options, tag string, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) HasTagFaceted(ctx context.Context, tag string, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {
	return nil, spelunker.ErrNotImplemented
}
