FROM golang:1.25.5-alpine3.22 AS builder


WORKDIR /app
COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o siwakasen main.go

FROM alpine:latest

WORKDIR /home/app

COPY --from=builder /app/siwakasen /home/app/siwakasen

EXPOSE 80

CMD ["./siwakasen"]
