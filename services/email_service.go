package services

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"pigeon/api"
	"pigeon/config"
	"time"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(config *config.Config) *EmailService {
	return &EmailService{
		config: config,
	}
}

func (s *EmailService) PostSendMail(_ context.Context, request api.PostSendMailRequestObject) (api.PostSendMailResponseObject, error) {
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

func (s *EmailService) PostBookaroomVerify(_ context.Context, request api.PostBookaroomVerifyRequestObject) (api.PostBookaroomVerifyResponseObject, error) {
	receiver := request.Body.Email

	rand.Seed(time.Now().UnixNano())
	codeStr := fmt.Sprintf("%d", rand.Intn(9000)+1000)
	salt := s.config.VerificationSalt

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
