package main

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	chimiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/runtime/middleware"
	"pigeon/api"
)

func init() {
	_username = os.Getenv("SMTP_USERNAME")
	_password = os.Getenv("SMTP_PASSWORD")
	_salt = os.Getenv("SALT")
}

func main() {
	smtpConfig := &SmtpConfig{
		SMTPServer:   "smtp.yandex.com",
		SMTPPort:     "587",
		SMTPUsername: _username,
		SMTPPassword: _password,
	}
	service := NewEmailService(smtpConfig, &_salt)
	s := NewServer(service)

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatal(err)
	}
	swagger.Servers = nil

	router := chi.NewRouter()

	router.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(swagger)
	})
	router.Handle("/swagger/", middleware.SwaggerUI(middleware.SwaggerUIOpts{
		Path:    "/swagger/",
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
		api.NewStrictHandler(s, nil),
		api.ChiServerOptions{
			BaseURL:    "/api/v1",
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

type SmtpConfig struct {
	SMTPServer   string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
}

var (
	_username string
	_password string
	_salt     string
)

func NewServer(emailService *EmailService) api.StrictServerInterface {
	return &server{
		EmailService: emailService,
	}
}

type server struct {
	*EmailService
}

type EmailService struct {
	config *SmtpConfig
	salt   *string
}

func NewEmailService(config *SmtpConfig, salt *string) *EmailService {
	return &EmailService{
		config: config,
		salt:   salt,
	}
}

func (s *EmailService) PostSendMail(ctx context.Context, request api.PostSendMailRequestObject) (api.PostSendMailResponseObject, error) {
	data := request.Body

	from := s.config.SMTPUsername
	to := []string{string(data.Receiver)}
	msg := "From: " + from + "\n" +
		"To: " + string(data.Receiver) + "\n" +
		"Subject: " + data.Subject + "\n\n" +
		data.Body

	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPServer)
	err := smtp.SendMail(s.config.SMTPServer+":"+s.config.SMTPPort, auth, from, to, []byte(msg))

	if err != nil {
		log.Println("Failed to send email:", err)
		response := api.PostSendMail400JSONResponse{
			Status:  "error",
			Message: "Failed to send the email: " + err.Error(),
		}
		return response, nil
	}

	response := api.PostSendMail200JSONResponse{
		Status:  "success",
		Message: "Email sent successfully",
	}

	return response, nil
}

func (s *EmailService) PostBookaroomVerify(ctx context.Context, request api.PostBookaroomVerifyRequestObject) (api.PostBookaroomVerifyResponseObject, error) {
	receiver := request.Body.Receiver

	rand.Seed(time.Now().UnixNano())
	codeStr := fmt.Sprintf("%d", rand.Intn(9000)+1000)
	salt := *s.salt

	hasher := sha512.New()
	hasher.Write([]byte(codeStr + salt))
	secretCode := hex.EncodeToString(hasher.Sum(nil))

	from := s.config.SMTPUsername
	to := []string{string(receiver)}
	msg := "From: " + from + "\n" +
		"To: " + string(receiver) + "\n" +
		"Subject: " + "Bookaroom Email Verification" + "\n\n" +
		"Please enter the following secret code: " + codeStr

	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPServer)
	err := smtp.SendMail(s.config.SMTPServer+":"+s.config.SMTPPort, auth, from, to, []byte(msg))

	if err != nil {
		log.Println("Failed to send email:", err)
		response := api.PostBookaroomVerify400JSONResponse{
			Status:  "error",
			Message: "Failed to send the email: " + err.Error(),
		}
		return response, nil
	}

	response := api.PostBookaroomVerify200JSONResponse{
		Status:  "success",
		Message: "Email sent successfully",
		Hash:    secretCode,
	}

	return response, nil
}
