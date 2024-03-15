package opensearch

import (
	"errors"
)

var ErrCursorIsExpired = errors.New("Query cursor has expired")
