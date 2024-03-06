# go-whosonfirst-spelunker-opensearch

Go package implementing the `whosonfirst/go-whosonfirst-spelunker.Spelunker` interface for use with OpenSearch databases.

## Documentation

Documentation is incompete at this time. For starters consult the (also incomplete) documentation in the [whosonfirst/go-whosonfirst-spelunker](https://github.com/whosonfirst/go-whosonfirst-spelunker) package.

## Important

This is work in progress and you should expect things to change, break or simply not work yet.

Also, nothing (at all) works yet.

## Examples

### Indexing

Using the `wof-opensearch-index` tool from the [whosonfirst/go-whosonfirst-opensearch](https://github.com/whosonfirst/go-whosonfirst-opensearch) package:

```
$> bin/wof-opensearch-index \
	-writer-uri 'constant://?val=opensearch2%3A%2F%2Flocalhost%3A9200%2Fspelunker%3Fusername%3Dadmin%26password%3...%26insecure%3Dtrue%26require-tls%3Dtrue' \
	/usr/local/data/whosonfirst-data-admin-ca/
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-spelunker
* https://github.com/whosonfirst/go-whosonfirst-spelunker-httpd
* https://github.com/whosonfirst/whosonfirst-opensearch
* https://github.com/whosonfirst/go-whosonfirst-opensearch