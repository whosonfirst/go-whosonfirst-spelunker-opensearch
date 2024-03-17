package opensearch

import (
	"fmt"
	"strings"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
)

// Something something something do all of this with templates...

func (s OpenSearchSpelunker) matchAllQuery() string {
	return `{"query": { "match_all": {} }}`
}

func (s OpenSearchSpelunker) idQuery(id int64) string {
	return fmt.Sprintf(`{"query": { "ids": { "values": [ %d ] } } }`, id)
}

func (s OpenSearchSpelunker) descendantsQuery(id int64, filters []spelunker.Filter) string {
	return fmt.Sprintf(`{"query": { "term": { "wof:belongsto":  %d  } } }`, id)
}

func (s OpenSearchSpelunker) hasPlacetypesQuery(pt string, filters []spelunker.Filter) string {
	return fmt.Sprintf(`{"query": { "term": { "wof:placetype":  "%s"  } } }`, pt)
}

func (s OpenSearchSpelunker) matchAllFacetedQuery(facets []*spelunker.Facet) string {

	str_aggs := s.facetsToAggregations(facets)
	return fmt.Sprintf(`{"query": { "match_all": {} }, "aggs": { %s } }`, str_aggs)
}

func (s OpenSearchSpelunker) hasPlacetypeFacetedQuery(pt string, filters []spelunker.Filter, facets []*spelunker.Facet) string {

	str_aggs := s.facetsToAggregations(facets)	
	return fmt.Sprintf(`{"query": { "term": { "wof:placetype":  "%s"  } }, "aggs": { %s } }`, pt, str_aggs)
}

// https://opensearch.org/docs/latest/aggregations/
// https://opensearch.org/docs/latest/aggregations/bucket/terms/

func (s OpenSearchSpelunker) descendantsFacetedQuery(id int64, filters []spelunker.Filter, facets []*spelunker.Facet) string {
	
	str_aggs := s.facetsToAggregations(facets)
	return fmt.Sprintf(`{"query": { "term": { "wof:belongsto":  %d  } }, "aggs": { %s } }`, id, str_aggs)
}

func (s OpenSearchSpelunker) searchQuery(search_opts *spelunker.SearchOptions) string {

	return fmt.Sprintf(`{"query": { "term": { "names_all": "%s" } } }`, search_opts.Query)
}

func (s OpenSearchSpelunker) facetsToAggregations(facets []*spelunker.Facet) string {

	count_facets := len(facets)
	aggs := make([]string, count_facets)

	for i, f := range facets {

		// This will probably need to be done in a switch statement eventually...
		facet_field := fmt.Sprintf("wof:%s", f)

		aggs[i] = fmt.Sprintf(`"%s": { "terms": { "field": "%s", "size": 1000 } }`, f, facet_field)
	}

	return strings.Join(aggs, ",")	
}
