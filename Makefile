CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-spelunker cmd/wof-spelunker/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/wof-spelunker-httpd cmd/wof-spelunker-httpd/main.go

server-linux:
	GOOS=linux GOARCH=amd64	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o work/wof-spelunker-httpd cmd/wof-spelunker-httpd/main.go

# Targets for running the Spelunker locally

# https://github.com/whosonfirst/go-whosonfirst-database
OS_INDEX=/usr/local/whosonfirst/go-whosonfirst-database/bin/wof-opensearch-index

# https://github.com/whosonfirst/whosonfirst-database
WHOSONFIRST_OPENSEARCH=/usr/local/whosonfirst/go-whosonfirst-database/opensearch

# https://github.com/aaronland/go-tools
URLESCAPE=$(shell which urlescape)

# Opensearch server

# This is for debugging. Do not change this at your own risk.
# (That means you should change this.)
OS_PSWD=dkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj

# CACHE_URI=gocache://
CACHE_URI=ristretto://
ENC_CACHE_URI=$(shell $(URLESCAPE) $(CACHE_URI))

READER_URI=https://data.whosonfirs.org
ENC_READER_URI=$(shell $(URLESCAPE) $(READER_URI))

# opensearch2%3A%2F%2Flocalhost%3A9200%2Fspelunker%3Frequire-tls%3Dtrue%26insecure%3Dtrue%26debug%3Dfalse%26username%3Dadmin%26password%3D$(OS_PSWD)

CLIENT_URI="https://localhost:9200/spelunker?username=admin&password=$(OS_PSWD)&insecure=true&require-tls=true"
ENC_CLIENT_URI=$(shell $(URLESCAPE) $(CLIENT_URI))

SPELUNKER_URI=opensearch://?client-uri=$(ENC_CLIENT_URI)&cache-uri=$(ENC_CACHE_URI)&reader-uri=$(ENC_READER_URI)

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
	cat $(WHOSONFIRST_OPENSEARCH)/schema/2.x/settings.spelunker.json | \
	curl -k \
		-H 'Content-type:application/json' \
		-XPUT \
		https://admin:$(OS_PSWD)@localhost:9200/spelunker/_settings \
		-d @-

# Opensearch indexing

index-local:
	$(OS_INDEX) \
		-writer-uri 'constant://?val=$(ENC_CLIENT_URI)' \
		$(REPO)

# Spelunker server

server-local:
	@make cli
	./bin/wof-spelunker-httpd \
		-server-uri http://localhost:8080 \
		-spelunker-uri '$(SPELUNKER_URI)' \
		-protomaps-api-key '$(APIKEY)'


# Lambda stuff

lambda:
	@make lambda-server

lambda-server:
	if test -f bootstrap; then rm -f bootstrap; fi
	if test -f server.zip; then rm -f server.zip; fi
	GOARCH=arm64 GOOS=linux go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -tags lambda.norpc -o bootstrap cmd/wof-spelunker-httpd/main.go
	zip server.zip bootstrap
	rm -f bootstrap
