FROM golang:1.26-alpine

COPY . /app/

WORKDIR /app

RUN go mod download

CMD ["go", "run", "cmd/api/main.go"]
