# SMailer
a simple (smtp) mailer [Docker hub](https://hub.docker.com/r/mintyleaf/smailer).   
i didn't find such simple docker ready project, so here it is.

## Configuration
Set the next environment variables:
```
  TOKEN - mailtrap api token
```

## API
`/send` handle accepts json body:
```
type SendRequestBody struct {
  From     string `json:"from"`
  To       string `json:"to"`
  Subject  string `json:"subject"`
  Data     string `json:"data"`
}
```
Data is the raw HTML data of the mail

## Run
`go run main.go`   
or   
`docker run -e ... mintyleaf/smailer -e`   

## Credits 
* [Email template based on this](https://github.com/leemunroe/responsive-html-email-template)
* [Go mail module](https://github.com/go-gomail/gomail)
