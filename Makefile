all: install generate run

install:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v2.0.0

generate:
	oapi-codegen -config config/server.cfg.yaml docs/openapi.yaml

run:
	go run cmd/main.go

visit:
	open http://localhost:8000/docs