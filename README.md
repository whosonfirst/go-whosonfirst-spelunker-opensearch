# go-whosonfirst-spelunker-opensearch

Go package implementing the `whosonfirst/go-whosonfirst-spelunker.Spelunker` interface for use with OpenSearch databases.

## Documentation

Documentation is incompete at this time. For starters consult the (also incomplete) documentation in the [whosonfirst/go-whosonfirst-spelunker](https://github.com/whosonfirst/go-whosonfirst-spelunker) package.

## Examples

Note: All the examples assume a "local" setup meaning there is local instance of OpenSearch running on port 9200.

* For an example of how to run a local OpenSearch instance from a Docker container [consult the `os` Makefile target](https://github.com/whosonfirst/go-whosonfirst-spelunker-opensearch/blob/main/Makefile#L28-L36) in this package.

* For an example of how to create a "Spelunker" index and mappings in an OpenSearch index [consult the `spelunker-local` target](https://github.com/whosonfirst/whosonfirst-opensearch/blob/main/Makefile#L5-L15) in the `whosonfirst/whosonfirst-opensearch` package.

### Indexing

Using the `wof-opensearch-index` tool from the [whosonfirst/go-whosonfirst-database](https://github.com/whosonfirst/go-whosonfirst-database) package:

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

## Tools

### server

```
$> make server-local
go run -mod vendor cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri 'opensearch://?client-uri=https%3A%2F%2Flocalhost%3A9200%2Fspelunker%3Fusername%3Dadmin%26password%3Ddkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj%26insecure%3Dtrue%26require-tls%3Dtrue&cache-uri=ristretto%3A%2F%2F&reader-uri=https%3A%2F%2Fdata.whosonfirst.org'
		
2024/03/11 09:06:51 INFO Listening for requests address=http://localhost:8080
```

See all the URL escaped gibberish in the `-spelunker-uri` flag? It's the same issues described in the docs for the `wof-opensearch-index` tool above.

## See also

* https://github.com/whosonfirst/go-whosonfirst-spelunker
* https://github.com/whosonfirst/go-whosonfirst-spelunker-httpd
* https://github.com/whosonfirst/go-whosonfirst-database
