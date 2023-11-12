# Pigeon

The project uses OPEN-API 3 code gen 

https://github.com/deepmap/oapi-codegen

## Getting started

Install oapi-codegen
```shell
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
```

Generate the server code
```shell
oapi-codegen -config _oas/server.cfg.yaml _oas/openapi3.yaml
```

Run server
```shell
go run main.go
```

Visit Swagger
```shell
open http://localhost:8000/swagger/
```
