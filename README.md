# SMailer
a simple (smtp) mailer.   
i didn't find such simple docker ready project, so here it is.

## Configuration
Set the next environment variables:
```
  SMTP_HOST - smtp host
  SMTP_PORT - smtp port
  SMTP_USER - smtp user
  SMTP_PASSWORD - smtp password
  TEMPLATES - html templates directory
```

## API
`/send` handle accepts json body:
```
type SendRequestBody struct {
  From     string      `json:"from"`
  To       string      `json:"to"`
  Subject  string      `json:"subject"`
  Template string      `json:"template"`
  Values   interface{} `json:"values"`
}
```

template is a filename without extension in templates directory   
values is a json map of keys/values, or whatever more complex that can be consumed by the go html/template module   

## Run
`go run main.go`   
or   
`docker run -e ... mintyleaf/smailer -e`   

## Credits 
* [Email template based on this](https://github.com/leemunroe/responsive-html-email-template)
* [Go mail module](https://github.com/go-gomail/gomail)
