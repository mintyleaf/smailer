FROM golang:latest AS builder
LABEL maintainer="mintyleaf <mintyleafdev@gmail.com>"

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o smailer

FROM alpine
WORKDIR /

COPY --from=builder /build/smailer /smailer

EXPOSE 8080
CMD ["/smailer"]
