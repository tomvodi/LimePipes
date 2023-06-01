
API_GEN_DIR=./internal/api_gen

.PHONY: test test-cover lint cover-html server

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

server:
	openapi-generator-cli generate \
		-i ./api/openapi-spec/openapi.yaml \
		-g go-gin-server \
		-o ${API_GEN_DIR} \
		--additional-properties=packageName=api_gen
