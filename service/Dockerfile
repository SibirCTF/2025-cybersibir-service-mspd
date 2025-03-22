FROM golang:alpine AS builder

RUN apk update --no-cache

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/MSPD service/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/MSPD .

EXPOSE 1015

CMD ["/app/MSPD"]