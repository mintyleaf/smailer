package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Smailer struct {
	SMTPHost      string
	SMTPPort      int
	SMTPUser      string
	SMTPPassword  string
	TemplatesPath string
	Token         string
}

type SendRequestBody struct {
	From     string      `json:"from"`
	To       string      `json:"to"`
	Subject  string      `json:"subject"`
	Template string      `json:"template"`
	Values   interface{} `json:"values"`
}

type Address struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Mail struct {
	From    Address   `json:"from"`
	To      []Address `json:"to"`
	Subject string    `json:"subject"`
	HTML    string    `json:"html"`
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

	url := "https://send.api.mailtrap.io/api/send"
	method := "POST"
	mail := Mail{
		From: Address{
			Email: request.From,
		},
		To: []Address{
			Address{
				Email: request.To,
			},
		},
		Subject: request.Subject,
		HTML:    htmlBuffer.String(),
	}

	payload, err := json.Marshal(mail)
	if err != nil {
		log.Printf("Failed to marshal payload buffer %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))

	if err != nil {
		log.Printf("Failed to create mailtrap api request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Header.Add("Authorization", "Bearer "+s.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to do mailtrap api request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read mailtrap response body: %s, %v", body, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// m := mail.NewMessage()

	// m.SetHeader("From", request.From)
	// m.SetHeader("To", request.To)
	// m.SetHeader("Subject", request.Subject)
	// m.SetBody("text/html", htmlBuffer.String())
	// d := mail.NewDialer(s.SMTPHost, s.SMTPPort, s.SMTPUser, s.SMTPPassword)
	// if err := d.DialAndSend(m); err != nil {
	// 	log.Printf("Failed to call SMTP relay: %v", err)
	// 	http.Error(w, err.Error(), http.StatusServiceUnavailable)
	// 	return
	// }
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
		Token:         lookupEnvPanic("TOKEN"),
	}

	http.HandleFunc("/send", smailer.SendMail)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
