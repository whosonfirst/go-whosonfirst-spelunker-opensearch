package opensearch

import (
	"fmt"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
)

func (s OpenSearchSpelunker) idQuery(id int64) string {
	return fmt.Sprintf(`{"query": { "ids": { "values": [ %d ] } } }`, id)
}

func (s OpenSearchSpelunker) descendantsQuery(id int64) string {
	return fmt.Sprintf(`{"query": { "term": { "wof:belongsto":  %d  } } }`, id)
}

func (s OpenSearchSpelunker) searchQuery(search_opts *spelunker.SearchOptions) string {

	return fmt.Sprintf(`{"query": { "term": { "names_all": "%s" } } }`, search_opts.Query)
}
