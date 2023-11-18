package main

import (
	"context"
	"encoding/json"
	chimiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-openapi/runtime/middleware"
	"log"
	"net/http"
	"pigeon/api"
	"pigeon/config"
	"pigeon/services"
)

func main() {
	config.LoadEnvs()
	appConfig := config.NewConfig()
	service := services.NewEmailService(appConfig)
	appServer := NewServer(service)

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}
	swagger.Servers = nil

	router := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(cors.Handler)

	router.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(swagger)
	})
	router.Handle("/docs", middleware.SwaggerUI(middleware.SwaggerUIOpts{
		Path:    "/docs",
		SpecURL: "/swagger/doc.json",
	}, nil))

	validator := chimiddleware.OapiRequestValidatorWithOptions(
		swagger,
		&chimiddleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
					return nil
				},
			},
		},
	)

	apiServer := api.HandlerWithOptions(
		api.NewStrictHandler(appServer, nil),
		api.ChiServerOptions{
			BaseURL:    "",
			BaseRouter: router,
			Middlewares: []api.MiddlewareFunc{
				validator,
			},
		},
	)

	addr := ":8000"
	httpServer := http.Server{
		Addr:    addr,
		Handler: apiServer,
	}

	log.Println("Server listening on", addr)
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func NewServer(emailService *services.EmailService) api.StrictServerInterface {
	return &server{
		EmailService: emailService,
	}
}

type server struct {
	*services.EmailService
}
