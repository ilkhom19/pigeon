all: install generate run

install:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v2.0.0

generate:
	oapi-codegen -config _oas/server.cfg.yaml _oas/openapi3.yaml

run:
	go run main.go

visit:
	open http://localhost:8000/docs