FROM golang:alpine AS builder

ARG CMD_PATH

COPY go.mod /src/

RUN cd /src && go env -w GOPROXY=https://goproxy.cn,direct && go mod download

COPY . /src

WORKDIR /src/${CMD_PATH}

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

FROM debian:stable-slim

ARG CMD_PATH

COPY --from=builder /src/${CMD_PATH}/main /app/main

WORKDIR /app

EXPOSE 8080

CMD ["./main"]