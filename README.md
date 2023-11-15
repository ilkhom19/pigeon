# Pigeon

The project uses OPEN-API 3 code gen 

https://github.com/deepmap/oapi-codegen

## Getting started

Install oapi-codegen
```shell
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
```

Generate
```shell
oapi-codegen -config docs/server.cfg.yaml docs/openapi3.yaml
```

Run server
```shell
go run main.go
```

Visit Swagger
```shell
open http://localhost:8000/docs
```
