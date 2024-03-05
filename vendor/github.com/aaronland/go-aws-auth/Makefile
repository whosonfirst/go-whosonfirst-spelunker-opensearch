GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/aws-mfa-session cmd/aws-mfa-session/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/aws-get-credentials cmd/aws-get-credentials/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/aws-set-env cmd/aws-set-env/main.go
