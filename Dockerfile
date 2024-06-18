
FROM golang:1.21.0

# Set destination for COPY
WORKDIR /app

# Download Go modules
ENV GO111MODULE=on

ENV HOST='35.222.136.171:2775'
ENV USER='smppclient1'
ENV PASS='password'

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


EXPOSE 8080

# Run
CMD ["./smpp-connector"]