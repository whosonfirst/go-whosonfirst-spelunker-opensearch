package opensearch

import (
	"fmt"
	"strings"
	"log/slog"
	
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

	q := s.descendantsQueryCriteria(id, filters)
	return fmt.Sprintf(`{"query": %s }`, q)
}

func (s OpenSearchSpelunker) descendantsFacetedQuery(id int64, filters []spelunker.Filter, facets []*spelunker.Facet) string {

	q := s.descendantsQueryCriteria(id, filters)
	str_aggs := s.facetsToAggregations(facets)
	
	return fmt.Sprintf(`{"query": %s }, "aggs": { %s } }`, q, str_aggs)
}

func (s OpenSearchSpelunker) descendantsQueryCriteria(id int64, filters []spelunker.Filter) string {

	q := fmt.Sprintf(`{ "term": { "wof:belongsto":  %d  } }`, id)
	
	if len(filters) == 0 {
		return q
	}

	must := []string{
		q,
	}

	return s.mustQueryWithFiltersCriteria(must, filters)
	return fmt.Sprintf(`{"query": %s }`, q)	
}

func (s OpenSearchSpelunker) hasPlacetypeQuery(pt string, filters []spelunker.Filter) string {

	q := s.hasPlacetypeQueryCriteria(pt, filters)
	return fmt.Sprintf(`{"query": %s }`, q)		
}

func (s OpenSearchSpelunker) hasPlacetypeFacetedQuery(pt string, filters []spelunker.Filter, facets []*spelunker.Facet) string {

	q := s.hasPlacetypeQueryCriteria(pt, filters)
	str_aggs := s.facetsToAggregations(facets)
	
	return fmt.Sprintf(`{"query": %s }, "aggs": { %s } }`, q, str_aggs)
}

func (s OpenSearchSpelunker) hasPlacetypeQueryCriteria(pt string, filters []spelunker.Filter) string {

	q := fmt.Sprintf(`{ "term": { "wof:placetype":  "%s"  } }`, pt)
	
	if len(filters) == 0 {
		return q
	}

	must := []string{
		q,
	}

	return s.mustQueryWithFiltersCriteria(must, filters)	
}

func (s OpenSearchSpelunker) matchAllFacetedQuery(facets []*spelunker.Facet) string {

	str_aggs := s.facetsToAggregations(facets)
	return fmt.Sprintf(`{"query": { "match_all": {} }, "aggs": { %s } }`, str_aggs)
}


// https://opensearch.org/docs/latest/aggregations/
// https://opensearch.org/docs/latest/aggregations/bucket/terms/

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

func (s OpenSearchSpelunker) mustQueryWithFiltersCriteria(must []string, filters []spelunker.Filter) string {

	for _, f := range filters {

		switch f.Scheme() {
		case "placetype":
			must = append(must, fmt.Sprintf(`{ "term": { "wof:placetype": "%s" } }`, f.Value()))
		case "country":
			must = append(must, fmt.Sprintf(`{ "term": { "wof:country": "%s" } }`, f.Value()))			
		default:
			slog.Warn("Unsupported filter scheme", "scheme", f.Scheme())
		}
	}
	
	str_must := strings.Join(must, ",")
	
	return fmt.Sprintf(`{ "bool": { "must": [ %s ] } }`, str_must)	
}
