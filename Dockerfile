FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o auth-service main.go

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/auth-service ./auth-service
EXPOSE 8081

ENTRYPOINT ["./auth-service"]