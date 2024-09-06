FROM golang:1.22.4-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./chat-box

EXPOSE 8000

CMD ["./chat-box"]