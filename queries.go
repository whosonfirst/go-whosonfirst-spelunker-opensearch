package opensearch

import (
	"fmt"
	"strings"
	"log/slog"
	"time"
	
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
	
	return fmt.Sprintf(`{"query": %s, "aggs": { %s } }`, q, str_aggs)
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
	
	return fmt.Sprintf(`{"query": %s, "aggs": { %s } }`, q, str_aggs)
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

func (s OpenSearchSpelunker) hasConcordanceQuery(namespace string, predicate string, value any, filters []spelunker.Filter) string {

	q := s.hasConcordanceQueryCriteria(namespace, predicate, value, filters)
	return fmt.Sprintf(`{"query": %s }`, q)		
}

func (s OpenSearchSpelunker) hasConcordanceFacetedQuery(namespace string, predicate string, value any, filters []spelunker.Filter, facets []*spelunker.Facet) string {

	q := s.hasConcordanceQueryCriteria(namespace, predicate, value, filters)
	str_aggs := s.facetsToAggregations(facets)
	
	return fmt.Sprintf(`{"query": %s, "aggs": { %s } }`, q, str_aggs)
}

func (s OpenSearchSpelunker) hasConcordanceQueryCriteria(namespace string, predicate string, value any, filters []spelunker.Filter) string {

	var q string

	str_value := fmt.Sprintf("%v", value)

	// Basically we need to index "magic 8"s...
	
	switch {
	case namespace != "" && predicate != "" && str_value != "":
		q = fmt.Sprintf(`{ "term": { "wof:concordances.%s:%s":  { "value": "%s", "case_insensitive": true } } }`, namespace, predicate, str_value)
	case namespace != "" && predicate != "":
		q = fmt.Sprintf(`{ "wildcard": { "wof:concordances_machinetags.keyword":  { "value": "%s:%s=*", "case_insensitive": true }  } }`, namespace, predicate)
	case namespace != "":
		q = fmt.Sprintf(`{ "prefix": { "wof:concordances_sources":  { "value": "%s", "case_insensitive": true }  } }`, namespace)		
	case predicate != "" && str_value != "":
		q = fmt.Sprintf(`{ "wildcard": { "wof:concordances_machinetags":  { "value": "*:%s=%s", "case_insensitive": true }  } }`, predicate, value)
	case predicate != "":
		q = fmt.Sprintf(`{ "wildcard": { "wof:concordances_machinetags":  { "value": "*:%s", "case_insensitive": true }  } }`, predicate)
	case namespace != "" && str_value != "":
		q = fmt.Sprintf(`{ "wildcard": { "wof:concordances_machinetags":  { "value": "%s:*=%s", "case_insensitive": true }  } }`, namespace, value)		
	case value != nil:
		q = fmt.Sprintf(`{ "wildcard": { "wof:concordances_machinetags":  { "value": "*:*=%s", "case_insensitive": true }  } }`, value)		
	default:
		
	}

	slog.Info("Concordance", "namespace", namespace, "predicate", predicate, "value", value)
	slog.Info(q)
	
	if len(filters) == 0 {
		return q
	}

	must := []string{
		q,
	}

	return s.mustQueryWithFiltersCriteria(must, filters)	
}


func (s OpenSearchSpelunker) getRecentQuery(d time.Duration, filters []spelunker.Filter) string {
	
	q := s.getRecentQueryCriteria(d, filters)
	return fmt.Sprintf(`{"query": %s }`, q )
}

func (s OpenSearchSpelunker) getRecentFacetedQuery(d time.Duration, filters []spelunker.Filter, facets []*spelunker.Facet) string {
	
	q := s.getRecentQueryCriteria(d, filters)
	str_aggs := s.facetsToAggregations(facets)
	
	return fmt.Sprintf(`{"query": %s, "aggs": { %s } }`, q, str_aggs)
}

func (s OpenSearchSpelunker) getRecentQueryCriteria(d time.Duration, filters []spelunker.Filter) string {

	now := time.Now()
	ts := now.Unix()

	then := ts - int64(d.Seconds())
	
	q := fmt.Sprintf(`{ "range": { "wof:lastmodified": { "gte": %d  } } }`, then)

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
