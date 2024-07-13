
API_GEN_DIR=./internal/api_gen

.PHONY: test test-cover lint cover-html server

cli:
	go build -o ./limepipes-cli limepipes/cmd/limepipes-cli

mocks:
	go generate mockgen ./...

test:
	go test ./...

test-cover:
	go test ./... -coverprofile cover.out

lint:
	golangci-lint run

cover-html: test-cover
	go tool cover -html=cover.out

# TODO: Add openAPI spec from external repo limepipes-api
server:
	openapi-generator-cli generate \
		-i ./api/openapi-spec/openapi.yaml \
		-g go-gin-server \
		-o ${API_GEN_DIR} \
		--additional-properties=packageName=api_gen
