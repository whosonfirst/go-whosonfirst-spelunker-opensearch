CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAG=-s -w

# This is for debugging. Do not change this at your own risk.
# (That means you should change this.)
OS_PSWD=dkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj

# https://github.com/whosonfirst/go-whosonfirst-opensearch
OS_INDEX=/usr/local/whosonfirst/go-whosonfirst-opensearch/bin/wof-opensearch-index

# URL escaping, sigh...
SPELUNKER_URI=opensearch://?dsn=https%3A%2F%2Flocalhost%3A9200%2Fspelunker%3Fusername%3Dadmin%26password%3Ddkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj%26insecure%3Dtrue%26require-tls%3Dtrue

index:
	$(OS_INDEX) \
		-writer-uri 'constant://?val=opensearch2%3A%2F%2Flocalhost%3A9200%2Fspelunker%3Frequire-tls%3Dtrue%26insecure%3Dtrue%26debug%3Dfalse%26username%3Dadmin%26password%3D$(OS_PSWD)' \
		$(REPO)

server:
	go run -mod $(GOMOD) cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri '$(SPELUNKER_URI)' \
		-protomaps-api-key '$(APIKEY)'

# https://opensearch.org/docs/latest/install-and-configure/install-opensearch/docker/
#
# And then:
# curl -v -k https://admin:$(OS_PSWD)@localhost:9200/

os:
	docker run \
		-it \
		-p 9200:9200 \
		-p 9600:9600 \
		-e "discovery.type=single-node" \
		-e "OPENSEARCH_INITIAL_ADMIN_PASSWORD=$(OS_PSWD)" \
		-v opensearch-data1:/usr/local/data/opensearch \
		opensearchproject/opensearch:latest
