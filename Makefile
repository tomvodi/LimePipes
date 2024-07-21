
.PHONY: test test-cover lint cover-html server

cli:
	go build -o ./limepipes-cli limepipes/cmd/limepipes-cli

mocks:
	mockery

test:
	go test ./...

test-cover:
	go test ./... -coverprofile cover.out

lint:
	golangci-lint run

cover-html: test-cover
	go tool cover -html=cover.out

server:
	./scripts/generate_server.sh

create_test_certificates:
	mkdir -p build && \
	cd build && \
	pwd && \
	openssl req -new -subj "/C=US/ST=Utah/CN=localhost" -newkey rsa:2048 -nodes -keyout localhost.key -out localhost.csr && \
	openssl x509 -req -days 365 -in localhost.csr -signkey localhost.key -out localhost.crt


