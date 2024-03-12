package opensearch

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-opensearch/client"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type OpenSearchSpelunker struct {
	spelunker.Spelunker
	client *opensearch.Client
	index  string
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
		index:  "spelunker",
	}

	return s, nil
}

func (s *OpenSearchSpelunker) GetById(ctx context.Context, id int64) ([]byte, error) {

	q := fmt.Sprintf(`{"query": { "ids": { "values": [ %d ] } } }`, id)

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
	}

	body, err := s.search(ctx, req)

	if err != nil {
		slog.Error("Get by ID query failed", "q", q)
		return nil, fmt.Errorf("Failed to retrieve %d, %w", id, err)
	}

	r := gjson.GetBytes(body, "hits.hits.0._source")

	if !r.Exists() {
		return nil, fmt.Errorf("First hit missing")
	}

	return s.propsToGeoJSON(r.String()), nil
}

func (s *OpenSearchSpelunker) GetAlternateGeometryById(ctx context.Context, id int64, alt_geom *uri.AltGeom) ([]byte, error) {

	return nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) Search(ctx context.Context, pg_opts pagination.Options, search_opts *spelunker.SearchOptions) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.searchQuery(search_opts)

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),

		// pagination...
	}

	body, err := s.search(ctx, req)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to execute search, %w", err)
	}

	return s.searchResultsToSPR(ctx, pg_opts, body)
}

func (s *OpenSearchSpelunker) GetRecent(ctx context.Context, pg_opts pagination.Options, d time.Duration, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	return nil, nil, spelunker.ErrNotImplemented
}

func (s *OpenSearchSpelunker) facet(ctx context.Context, req *opensearchapi.SearchRequest, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	body, err := s.search(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("Failed to query facets, %w", err)
	}

	aggs_rsp := gjson.GetBytes(body, "aggregations")

	if !aggs_rsp.Exists() {
		return nil, fmt.Errorf("Failed to derive facets, missing")
	}

	facetings := make([]*spelunker.Faceting, 0)

	for k, rsp := range aggs_rsp.Map() {

		facet_results := make([]*spelunker.FacetCount, 0)

		buckets_rsp := rsp.Get("buckets")

		for _, b := range buckets_rsp.Array() {

			k_rsp := b.Get("key")
			v_rsp := b.Get("doc_count")

			fc := &spelunker.FacetCount{
				Key:   k_rsp.String(),
				Count: v_rsp.Int(),
			}

			facet_results = append(facet_results, fc)
		}

		var bucket_facet *spelunker.Facet

		for _, f := range facets {
			if f.String() == k {
				bucket_facet = f
				break
			}
		}

		if bucket_facet == nil {
			return nil, fmt.Errorf("Failed to determine facet for bucket key '%s'", k)
		}

		faceting := &spelunker.Faceting{
			Facet:   bucket_facet,
			Results: facet_results,
		}

		facetings = append(facetings, faceting)
	}

	return facetings, nil
}

func (s *OpenSearchSpelunker) search(ctx context.Context, req *opensearchapi.SearchRequest) ([]byte, error) {

	// https://pkg.go.dev/github.com/opensearch-project/opensearch-go/v2@v2.3.0/opensearchapi#SearchRequest

	if len(req.Index) == 0 {
		req.Index = []string{
			s.index,
		}
	}

	rsp, err := req.Do(ctx, s.client)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute search, %w", err)
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {

		// body, _ := io.ReadAll(rsp.Body)
		// slog.Error(string(body))

		slog.Error("Query failed", "status", rsp.StatusCode)
		return nil, fmt.Errorf("Invalid status")
	}

	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed to read response, %w", err)
	}

	return body, nil
}

func (s *OpenSearchSpelunker) propsToGeoJSON(props string) []byte {

	// See this? It's a derived geometry. Still working through
	// how to signal and how to fetch full geometry...

	lat_rsp := gjson.Get(props, "geom:latitude")
	lon_rsp := gjson.Get(props, "geom:longitude")

	lat := lat_rsp.String()
	lon := lon_rsp.String()

	return []byte(`{"type": "Feature", "properties": ` + props + `, "geometry": { "type": "Point", "coordinates": [` + lon + `,` + lat + `] } }`)
}

func (s *OpenSearchSpelunker) searchResultsToSPR(ctx context.Context, pg_opts pagination.Options, body []byte) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	r := gjson.GetBytes(body, "hits.total.value")
	count := r.Int()

	var pg_results pagination.Results
	var pg_err error

	if pg_opts != nil {
		pg_results, pg_err = countable.NewResultsFromCountWithOptions(pg_opts, count)
	} else {
		pg_results, pg_err = countable.NewResultsFromCount(count)
	}

	if pg_err != nil {
		return nil, nil, pg_err
	}

	// START OF put me in a(nother) function

	hits_r := gjson.GetBytes(body, "hits.hits")
	count_hits := len(hits_r.Array())

	results := make([]wof_spr.StandardPlacesResult, count_hits)

	for idx, r := range hits_r.Array() {

		src := r.Get("_source")
		sp_spr, err := NewSpelunkerRecordSPR([]byte(src.String()))

		if err != nil {
			slog.Error("Failed to derive SPR from result", "index", idx, "error", err)
			return nil, nil, fmt.Errorf("Failed to derive SPR from result, %w", err)
		}

		results[idx] = sp_spr
	}

	spr_results := NewSpelunkerStandardPlacesResults(results)

	return spr_results, pg_results, nil
}
