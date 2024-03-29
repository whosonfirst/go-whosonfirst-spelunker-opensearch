package opensearch

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
	_ "net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/aaronland/go-pagination/cursor"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	opensearchapi "github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-opensearch/client"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/go-reader"
	_ "github.com/whosonfirst/go-reader-http"
)

const scroll_duration time.Duration = 5 * time.Minute
const scroll_trigger int64 = 10000

type OpenSearchSpelunker struct {
	spelunker.Spelunker
	client *opensearch.Client
	index  string
	reader reader.Reader
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

	if q.Has("reader-uri"){
		
		reader_uri := q.Get("reader_uri")
	
		r, err := reader.NewReader(ctx, reader_uri)
		
		if err != nil {
			return nil, fmt.Errorf("Failed to create reader, %w", err)
		}

		s.reader = r
	}
		
	return s, nil
}

func (s *OpenSearchSpelunker) GetRecordForId(ctx context.Context, id int64, ) ([]byte, error) {

	q := fmt.Sprintf(`{"query": { "ids": { "values": [ %d ] } } }`, id)

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
	}

	body, err := s.searchWithIndex(ctx, req)

	if err != nil {
		slog.Error("Get by ID query failed", "q", q)
		return nil, fmt.Errorf("Failed to retrieve %d, %w", id, err)
	}

	r := gjson.GetBytes(body, "hits.hits.0._source")

	if !r.Exists() {
		return nil, fmt.Errorf("First hit missing")
	}

	return []byte(r.String()), nil
}

func (s *OpenSearchSpelunker) GetFeatureForId(ctx context.Context, id int64, uri_args *uri.URIArgs) ([]byte, error) {

	rel_path, err := uri.Id2RelPath(id, uri_args)

	if err != nil {
		return nil, err
	}

	f_reader := s.reader
	
	if f_reader == nil {

		record, err := s.GetRecordForId(ctx, id)
		
		if err != nil {
			return nil, err
		}
		
		repo_name := gjson.GetBytes(record, "wof:repo")
		reader_uri := fmt.Sprintf("https://raw.githubusercontent.com/whosonfirst-data/%s/master/data", repo_name)
		
		r, err := reader.NewReader(ctx, reader_uri)

		if err != nil {
			return nil, err
		}

		f_reader = r
	}
	
	r, err := f_reader.Read(ctx, rel_path)
	
	if err != nil {
		return nil, err
	}

	defer r.Close()
	return io.ReadAll(r)
}

func (s *OpenSearchSpelunker) facet(ctx context.Context, req *opensearchapi.SearchRequest, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	body, err := s.searchWithIndex(ctx, req)

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

// searchPaginated wraps all the logic for determining whether to do a cursor-based or plain-vanilla-paginated query
func (s *OpenSearchSpelunker) searchPaginated(ctx context.Context, pg_opts pagination.Options, q string) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	scroll_id := ""
	pre_count := false
	use_scroll := false

	if pg_opts.Method() == pagination.Cursor {
		scroll_id = pg_opts.Pointer().(string)
	}

	if scroll_id == "" {
		pre_count = true
	}

	if pre_count {

		count, err := s.countForQuery(ctx, q)

		if err != nil {
			return nil, nil, err
		}

		if count >= scroll_trigger {
			use_scroll = true
		}
	}

	var body []byte
	var err error

	if scroll_id != "" {

		// See this? Neither of these things are documented anywhere.
		// Good times...
		scroll_id = strings.TrimLeft(scroll_id, "after-")
		q = fmt.Sprintf(`{"scroll_id": "%s"}`, scroll_id)

		req := &opensearchapi.ScrollRequest{
			Body:     strings.NewReader(q),
			ScrollID: scroll_id,
			Scroll:   scroll_duration,
		}

		body, err = s.searchWithScroll(ctx, req)

	} else {

		sz := int(pg_opts.PerPage())

		from := int(pg_opts.PerPage() * (pg_opts.Pointer().(int64) - 1))

		req := &opensearchapi.SearchRequest{
			Body: strings.NewReader(q),
			Size: &sz,
			From: &from,
		}

		if use_scroll {
			req.Scroll = scroll_duration
		}

		body, err = s.searchWithIndex(ctx, req)
	}

	// To do: Check for expired scroll
	
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to execute search, %w", err)
	}

	return s.searchResultsToSPR(ctx, pg_opts, body)
}

// https://pkg.go.dev/github.com/opensearch-project/opensearch-go/v2@v2.3.0/opensearchapi#SearchRequest

func (s *OpenSearchSpelunker) searchWithIndex(ctx context.Context, req *opensearchapi.SearchRequest) ([]byte, error) {

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

		body, _ := io.ReadAll(rsp.Body)
		slog.Error(string(body))

		slog.Error("Query failed", "status", rsp.StatusCode)
		return nil, fmt.Errorf("Invalid status")
	}

	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		return nil, fmt.Errorf("Failed to read response, %w", err)
	}

	return body, nil
}

// https://pkg.go.dev/github.com/opensearch-project/opensearch-go/v2/opensearchapi#ScrollRequest

func (s *OpenSearchSpelunker) searchWithScroll(ctx context.Context, req *opensearchapi.ScrollRequest) ([]byte, error) {

	rsp, err := req.Do(ctx, s.client)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute search, %w", err)
	}

	defer rsp.Body.Close()

	// To do: Check for expired cursor...
	
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

func (s *OpenSearchSpelunker) propsToGeoJSON(props []byte) []byte {

	// See this? It's a derived geometry. Still working through
	// how to signal and how to fetch full geometry...

	lat_rsp := gjson.GetBytes(props, "geom:latitude")
	lon_rsp := gjson.GetBytes(props, "geom:longitude")

	lat := lat_rsp.String()
	lon := lon_rsp.String()

	return []byte(`{"type": "Feature", "properties": ` + string(props) + `, "geometry": { "type": "Point", "coordinates": [` + lon + `,` + lat + `] } }`)
}

func (s *OpenSearchSpelunker) countForQuery(ctx context.Context, q string) (int64, error) {

	sz := 0

	req := &opensearchapi.SearchRequest{
		Body: strings.NewReader(q),
		Size: &sz,
	}

	body, err := s.searchWithIndex(ctx, req)

	if err != nil {
		return 0, fmt.Errorf("Failed to determine count for query, %w", err)
	}

	r := gjson.GetBytes(body, "hits.total.value")
	count := r.Int()

	return count, nil
}

func (s *OpenSearchSpelunker) searchResultsToSPR(ctx context.Context, pg_opts pagination.Options, body []byte) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	scroll_rsp := gjson.GetBytes(body, "_scroll_id")
	scroll_id := scroll_rsp.String()

	total_rsp := gjson.GetBytes(body, "hits.total.value")
	count := total_rsp.Int()

	var pg_results pagination.Results
	var pg_err error

	if scroll_id != "" {

		page_count := math.Ceil(float64(count) / float64(pg_opts.PerPage()))

		c_results := new(cursor.CursorResults)
		c_results.TotalCount = count
		c_results.PerPageCount = pg_opts.PerPage()
		c_results.CursorNext = scroll_id
		c_results.PageCount = int64(page_count)

		pg_results = c_results

	} else {

		if pg_opts != nil {
			pg_results, pg_err = countable.NewResultsFromCountWithOptions(pg_opts, count)
		} else {
			pg_results, pg_err = countable.NewResultsFromCount(count)
		}
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
