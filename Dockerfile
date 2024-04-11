FROM golang:1.22 AS builder

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/Groq2API .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/Groq2API /app/Groq2API

CMD ["/app/Groq2API"]
