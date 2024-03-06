package opensearch

import (
	"fmt"
)

func (s OpenSearchSpelunker) idQuery(id int64) string {
	return fmt.Sprintf(`{"query": { "ids": { "values": [ %d ] } } }`, id)
}

func (s OpenSearchSpelunker) descendantsQuery(id int64) string {
	return fmt.Sprintf(`{"query": { "term": { "wof:belongsto":  %d  } } }`, id)
}
