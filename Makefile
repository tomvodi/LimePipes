
API_GEN_DIR=./internal/api_gen

mocks:
	go generate mockgen ./...

test:
	go test ./...

server:
	openapi-generator-cli generate \
		-i ./api/openapi-spec/openapi.yaml \
		-g go-gin-server \
		-o ${API_GEN_DIR} \
		--additional-properties=packageName=api_gen
