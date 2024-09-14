package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/mail.v2"
)

type Smailer struct {
	SMTPHost      string
	SMTPPort      int
	SMTPUser      string
	SMTPPassword  string
	TemplatesPath string
}

type SendRequestBody struct {
	From     string      `json:"from"`
	To       string      `json:"to"`
	Subject  string      `json:"subject"`
	Template string      `json:"template"`
	Values   interface{} `json:"values"`
}

func (s *Smailer) SendMail(w http.ResponseWriter, r *http.Request) {
	request := SendRequestBody{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Failed to decode request json body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t, err := template.ParseFiles(filepath.Join(s.TemplatesPath, request.Template+".html"))
	if err != nil {
		log.Printf("Failed to parse template file: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	htmlBuffer := bytes.Buffer{}
	err = t.Execute(&htmlBuffer, request.Values)
	if err != nil {
		log.Printf("Failed to replace html template values: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m := mail.NewMessage()

	m.SetHeader("From", request.From)
	m.SetHeader("To", request.To)
	m.SetHeader("Subject", request.Subject)
	m.SetBody("text/html", htmlBuffer.String())
	d := mail.NewDialer(s.SMTPHost, s.SMTPPort, s.SMTPUser, s.SMTPPassword)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to call SMTP relay: %v", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

func lookupEnvPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("%s env variable is not set", key)
	}
	return value
}

func main() {
	smtpPort, err := strconv.Atoi(lookupEnvPanic("SMTP_PORT"))
	if err != nil {
		log.Fatalf("SMTP_PORT must be a valid int: %v", err)
	}
	templatesPath := lookupEnvPanic("TEMPLATES")
	if stat, err := os.Stat(templatesPath); err != nil {
		log.Fatalf("Failed to open directory at %s path: %v", templatesPath, err)
	} else {
		if !stat.IsDir() {
			log.Fatalf("TEMPLATES is not dir path")
		}
	}

	smailer := &Smailer{
		SMTPHost:      lookupEnvPanic("SMTP_HOST"),
		SMTPPort:      smtpPort,
		SMTPUser:      lookupEnvPanic("SMTP_USER"),
		SMTPPassword:  lookupEnvPanic("SMTP_PASSWORD"),
		TemplatesPath: templatesPath,
	}

	http.HandleFunc("/send", smailer.SendMail)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
