CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAG=-s -w

# This probably won't work yet...
SPELUNKER_URI=opensearch://?dsn=https://localhost:9200

OSPSWD=dkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj

server:
	go run -mod $(GOMOD) cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri $(SPELUNKER_URI)

# https://opensearch.org/docs/latest/install-and-configure/install-opensearch/docker/
# --rm \

os:
	docker run \
		-it \
		-p 9200:9200 \
		-p 9600:9600 \
		-e "discovery.type=single-node" \
		-e "OPENSEARCH_INITIAL_ADMIN_PASSWORD=$(OSPSWD)" \
		-v opensearch-data1:/usr/local/data/opensearch \
		opensearchproject/opensearch:latest
