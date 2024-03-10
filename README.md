# go-whosonfirst-spelunker-opensearch

Go package implementing the `whosonfirst/go-whosonfirst-spelunker.Spelunker` interface for use with OpenSearch databases.

## Documentation

Documentation is incompete at this time. For starters consult the (also incomplete) documentation in the [whosonfirst/go-whosonfirst-spelunker](https://github.com/whosonfirst/go-whosonfirst-spelunker) package.

## Important

This is work in progress and you should expect things to change, break or simply not work yet.

## Examples

### Indexing

Using the `wof-opensearch-index` tool from the [whosonfirst/go-whosonfirst-opensearch](https://github.com/whosonfirst/go-whosonfirst-opensearch) package:

```
$> bin/wof-opensearch-index \
	-writer-uri 'constant://?val=opensearch2%3A%2F%2Flocalhost%3A9200%2Fspelunker%3Fusername%3Dadmin%26password%3...%26insecure%3Dtrue%26require-tls%3Dtrue' \
	/usr/local/data/whosonfirst-data-admin-ca/
```

Note the unfortunate need to URL escape the `-writer-uri=constant://?val=` parameter which unescaped is the actual `go-whosonfirst-opensearch/writer.OpensearchV2Writer` URI that takes the form of:

```
opensearch2://localhost:9200/spelunker?require-tls=true&insecure=true&debug=false&username=admin&password=s33kret
```

The `wof-opensearch-index` application however expects a [gocloud.dev/runtimevar](https://gocloud.dev/howto/runtimevar/) URI so that you don't need to deply production configuration values with sensitive values (like OpenSearch admin passwords) exposed in them. Under the hood the `wof-opensearch-index` application is using the [sfomuseum/runtimevar](https://github.com/sfomuseum/runtimevar) package to manage the details and this needs to be updated to allow plain (non-runtimevar) strings. Or maybe the `wof-opensearch-index` application needs to be updated. Either way something needs to be updated to avoid the hassle of always needing to URL-escape things.

## See also

* https://github.com/whosonfirst/go-whosonfirst-spelunker
* https://github.com/whosonfirst/go-whosonfirst-spelunker-httpd
* https://github.com/whosonfirst/whosonfirst-opensearch
* https://github.com/whosonfirst/go-whosonfirst-opensearch