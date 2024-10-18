# build bin
FROM golang:1.22.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /chat-box

# run bin
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /chat-box /app/chat-box

EXPOSE 8000

CMD ["./chat-box"]