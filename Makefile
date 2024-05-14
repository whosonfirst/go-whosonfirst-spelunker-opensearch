CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAG=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-spelunker-httpd cmd/httpd/main.go

# Targets for running the Spelunker locally

# https://github.com/whosonfirst/go-whosonfirst-opensearch
OS_INDEX=/usr/local/whosonfirst/go-whosonfirst-opensearch/bin/wof-opensearch-index

# https://github.com/whosonfirst/whosonfirst-opensearch
WHOSONFIRST_OPENSEARCH=/usr/local/whosonfirst/whosonfirst-opensearch

# https://github.com/aaronland/go-tools
URLESCAPE=$(shell which urlescape)

# CACHE_URI=gocache://
CACHE_URI=ristretto://
ENC_CACHE_URI=$(shell $(URLESCAPE) $(CACHE_URI))

DSN=https://localhost:9200/spelunker?username=admin&password=dkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj&insecure=true&require-tls=true&cache-uri=$(ENC_CACHE_URI)
ENC_DSN=$(shell $(URLESCAPE) $(DSN))

SPELUNKER_URI=opensearch://?dsn=$(DSN)

# Opensearch server

# This is for debugging. Do not change this at your own risk.
# (That means you should change this.)
OS_PSWD=dkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj

# https://opensearch.org/docs/latest/install-and-configure/install-opensearch/docker/
#
# And then:
# curl -v -k https://admin:$(OS_PSWD)@localhost:9200/

opensearch-local:
	docker run \
		-it \
		-p 9200:9200 \
		-p 9600:9600 \
		-e "discovery.type=single-node" \
		-e "OPENSEARCH_INITIAL_ADMIN_PASSWORD=$(OS_PSWD)" \
		-v opensearch-data1:/usr/local/data/opensearch \
		opensearchproject/opensearch:latest

# Opensearch "spelunker" index

spelunker-local:
	@make spelunker-local-index
	@make spelunker-local-fieldlimit

spelunker-local-index:
	cat $(WHOSONFIRST_OPENSEARCH)/schema/2.x/mappings.spelunker.json | \
		curl -k \
		-H 'Content-Type: application/json' \
		-X PUT \
		https://admin:$(OS_PSWD)@localhost:9200/spelunker \
		-d @-

spelunker-local-fieldlimit:
	curl -k \
		-H 'Content-type:application/json' \
		-XPUT https://admin:$(OS_PSWD)@localhost:9200/spelunker/_settings \
		-d '{"index.mapping.total_fields.limit": $(FIELD_LIMIT)}'

# Opensearch indexing

index-local:
	$(OS_INDEX) \
		-writer-uri 'constant://?val=opensearch2%3A%2F%2Flocalhost%3A9200%2Fspelunker%3Frequire-tls%3Dtrue%26insecure%3Dtrue%26debug%3Dfalse%26username%3Dadmin%26password%3D$(OS_PSWD)' \
		$(REPO)

# Spelunker server

server-local:
	go run -mod $(GOMOD) cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri '$(SPELUNKER_URI)' \
		-protomaps-api-key '$(APIKEY)'


# Lambda stuff

lambda:
	@make lambda-server

lambda-server:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/httpd/main.go
	zip server.zip bootstrap
	rm -f bootstrap
