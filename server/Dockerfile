FROM golang:1.22.2 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o chat-app ./server.go

EXPOSE 8080

CMD ["./chat-app"]