FROM golang:1.25.5-alpine3.22 AS builder


WORKDIR /app
COPY go.mod ./
#RUN go mod download
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o siwakasen main.go

FROM alpine:latest

WORKDIR /home/app
#RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

COPY --from=builder /app/siwakasen /home/app/siwakasen
#COPY --from=builder /app/migration /home/app/migration

EXPOSE 80

# Run the binary
CMD ["./siwakasen"]