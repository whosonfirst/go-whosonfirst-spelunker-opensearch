package funcs

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	
	"github.com/whosonfirst/go-whosonfirst-sources"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func NameForSource(source string) string {

	nspred := strings.Split(source, ":")
	prefix := nspred[0]
	
	src, err := sources.GetSourceByPrefix(prefix)

	if err != nil {
		return prefix
	}

	return src.Fullname
}

func FormatNumber(i int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", i)
}

func AppendPagination(uri string, k string, v any) string {

	u, err := url.Parse(uri)

	if err != nil {
		slog.Error("Failed to parse URI to append pagination", "uri", uri, "error", err)
		return "#"
	}

	q := u.Query()
	q.Set(k, fmt.Sprintf("%v", v))

	u.RawQuery = q.Encode()
	return u.String()
}
