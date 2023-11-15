# Pigeon

The project uses OPEN-API 3 code gen
https://github.com/deepmap/oapi-codegen

## Getting started

Run the project
```shell
make run
```

## Or


Install oapi-codegen
```shell
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
```

Generate
```shell
oapi-codegen -config config/server.cfg.yaml docs/openapi.yaml
```

Run server
```shell
go run cmd/main.go
```

Visit Swagger
```shell
http://localhost:8000/docs
```
