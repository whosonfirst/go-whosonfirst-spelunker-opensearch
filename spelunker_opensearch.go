package opensearch

import (
	"context"
	"fmt"
	"io"
	_ "log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/aaronland/go-pagination"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/tidwall/gjson"
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

	q := fmt.Sprintf(`{"query": { "ids": { "values": [ %d ] } } }`, id)

	// v2/opensearchapi

	rsp, err := s.client.Search(
		s.client.Search.WithIndex("spelunker"),
		s.client.Search.WithBody(strings.NewReader(q)),
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute search, %w", err)
	}

	// v3/opensearchapi

	/*
		req := &opensearchapi.SearchRequest{
		        Body: strings.NewReader(q),
		}

		rsp, err := s.client.Search(ctx, req)

		if err != nil {
			return nil, fmt.Errorf("Failed to execute search, %w", err)
		}
	*/

	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid status")
	}

	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed to read response, %w", err)
	}

	r := gjson.GetBytes(body, "hits.hits.0._source")

	if !r.Exists() {
		return nil, fmt.Errorf("First hit missing")
	}

	// Need to turn this in to GeoJSON...

	return []byte(r.String()), nil
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

	return 0, nil // spelunker.ErrNotImplemented
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
