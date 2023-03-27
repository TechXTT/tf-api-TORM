package email

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	jwt "github.com/hacktues-9/tf-api/pkg/jwt"
)

func parseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("./pkg/email/%s", templateFileName))
	if err != nil {
		return "", err
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}

// GmailService : Gmail client for sending email
var GmailService *gmail.Service

func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     os.Getenv("GMAIL_CLIENT_ID"),
		ClientSecret: os.Getenv("GMAIL_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  os.Getenv("GMAIL_ACCESS_TOKEN"),
		RefreshToken: os.Getenv("GMAIL_REFRESH_TOKEN"),
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		fmt.Println("Email service is initialized")
	}
}

func SendEmailOAUTH2(to string, data interface{}, template string) (bool, error) {

	emailBody, err := parseTemplate(template, data)
	if err != nil {
		fmt.Println("Error parsing template: ", err)
		return false, err
	}

	var message gmail.Message

	emailTo := "To: " + to + "\r\n"
	subject := "Subject: " + "Tues Fest 2023 Vote" + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err = GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		fmt.Println("Error sending email: ", err)
		return false, err
	}
	return true, nil
}

func GenerateVerificationLink(email string, privateKey string, publicKey string, TokenTTL time.Duration) string {
	token, err := jwt.CreateToken(TokenTTL, email, privateKey, publicKey)
	if err != nil {
		fmt.Println("Error creating token: ", err)
		return ""
	}
	return fmt.Sprintf("https://tuesfest.bg/verify/%s", token)
}

func ValidateEmailToken(token string) (string, error) {
	sub, err := jwt.ValidateStringToken(token, os.Getenv("PUBLIC_KEY"))
	if err != nil {
		fmt.Println("Error validating token: ", err)
		return "", err
	}
	return sub, nil
}
