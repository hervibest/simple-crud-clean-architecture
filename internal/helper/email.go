package helper

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

type GomailSender struct {
	dialer *gomail.Dialer
	app    *Application
	from   string
}

type Application struct {
	backend_url  string
	frontend_url string
}

func NewGomailSender(config *viper.Viper, log *logrus.Logger) *GomailSender {
	emailPort, _ := strconv.Atoi(config.GetString("gomail.port"))

	dialer := gomail.NewDialer(config.GetString("gomail.host"),
		emailPort, config.GetString("gomail.username"),
		config.GetString("gomail.password"))

	app := &Application{
		backend_url:  config.GetString("app.backend_url"),
		frontend_url: config.GetString("app.frontend_url"),
	}

	return &GomailSender{
		dialer: dialer,
		app:    app,
		from:   config.GetString("gomail.from"),
	}
}

func (s *GomailSender) SendEmail(email string, token string, category string) error {
	emailData := struct {
		Email       string
		Token       string
		FrontendUrl string
		AppUrl      string
	}{
		Email:       email,
		Token:       token,
		FrontendUrl: s.app.frontend_url,
		AppUrl:      s.app.backend_url}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}

	subject := ""
	filePath := ""

	switch {
	case category == "reset password":
		subject = "Email Reset Password"
		filePath = "request_reset_password.html"

	case category == "new email verification":
		subject = "User Email Verification"
		filePath = "registration.html"
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "hervipro@gmail.com")

	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", subject)

	templatePath := filepath.Join(cwd, "internal/templates/"+filePath)

	html, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("error parsing html template: %w", err)
	}

	var b bytes.Buffer
	if err := html.Execute(&b, emailData); err != nil {
		return fmt.Errorf("error parse html")
	}

	mailer.SetBody("text/html", b.String())

	go func() {
		errd := s.dialer.DialAndSend(mailer)
		if errd != nil {
			fmt.Println("error goroutine dialer:", errd)
		}
	}()

	return nil
}
