CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAG=-s -w

# This is for debugging. Do not change this at your own risk.
# (That means you should change this.)
OSPSWD=dkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj

# This probably won't work yet...
SPELUNKER_URI=opensearch://?dsn=https://admin:$(OSPSWD)@localhost:9200

server:
	go run -mod $(GOMOD) cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri '$(SPELUNKER_URI)'

# https://opensearch.org/docs/latest/install-and-configure/install-opensearch/docker/
#
# And then:
# curl -v -k https://admin:$(OSPSWD)@localhost:9200/

os:
	docker run \
		-it \
		-p 9200:9200 \
		-p 9600:9600 \
		-e "discovery.type=single-node" \
		-e "OPENSEARCH_INITIAL_ADMIN_PASSWORD=$(OSPSWD)" \
		-v opensearch-data1:/usr/local/data/opensearch \
		opensearchproject/opensearch:latest
