package opensearch

import (
	"context"
	"fmt"
	_ "log/slog"
	"net/url"
	"time"

	"github.com/aaronland/go-pagination"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/whosonfirst/go-whosonfirst-opensearch/client"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type OpenSearchSpelunker struct {
	spelunker.Spelunker
	client *opensearch.Client
}

func init() {
	ctx := context.Background()
	spelunker.RegisterSpelunker(ctx, "opensearch", NewOpenSearchSpelunker)
}

func NewOpenSearchSpelunker(ctx context.Context, uri string) (spelunker.Spelunker, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	dsn := q.Get("dsn")

	if dsn == "" {
		return nil, fmt.Errorf("Missing ?dsn= parameter")
	}

	cl, err := client.NewClient(ctx, dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to create opensearch client, %w", err)
	}

	s := &OpenSearchSpelunker{
		client: cl,
	}

	return s, nil
}

func (s *OpenSearchSpelunker) GetById(ctx context.Context, id int64) ([]byte, error) {

	return nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) GetAlternateGeometryById(ctx context.Context, id int64, alt_geom *uri.AltGeom) ([]byte, error) {

	return nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) GetDescendants(ctx context.Context, pg_opts pagination.Options, id int64, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	return nil, nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) GetDescendantsFaceted(ctx context.Context, id int64, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	return nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) CountDescendants(ctx context.Context, id int64) (int64, error) {

	return 0, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) HasPlacetype(ctx context.Context, pg_opts pagination.Options, pt *placetypes.WOFPlacetype, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	return nil, nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) Search(ctx context.Context, pg_opts pagination.Options, search_opts *spelunker.SearchOptions) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	return nil, nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) GetRecent(ctx context.Context, pg_opts pagination.Options, d time.Duration, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	return nil, nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) GetPlacetypes(ctx context.Context) (*spelunker.Faceting, error) {

	return nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) GetConcordances(ctx context.Context) (*spelunker.Faceting, error) {

	return nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) HasConcordance(ctx context.Context, pg_opts pagination.Options, namespace string, predicate string, value string, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	return nil, nil, spelunker.ErrNotImplemented
}
