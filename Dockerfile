FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin ./cmd

FROM alpine:3.23.5
WORKDIR /app
COPY --from=builder /app/bin .
COPY --from=builder /app/configs ./configs

RUN apk add --no-cache ca-certificates
EXPOSE 8443
ENTRYPOINT ["/app/bin"]
