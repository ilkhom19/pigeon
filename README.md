# Pigeon | [Visit API Documentation](https://pigeon-cretto.koyeb.app/docs)

###Pigeon is a lightweight and efficient email-sending service built using the Go programming language. The project embraces simplicity and ease of use, making it an ideal choice for developers who need a hassle-free solution for integrating email functionality into their applications.

The project uses OPEN-API 3 code gen

https://github.com/deepmap/oapi-codegen

## Getting started

Create a `.env` file in the project directory and configure it with your Yandex SMTP server credentials. Example:

   ```env
   SMTP_USERNAME=your_smtp_username
   SMTP_PASSWORD=your_smtp_password
   ```

Run the project
```shell
make run
```

## Or


Install oapi-codegen
```shell
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
```

Generate the server code
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
