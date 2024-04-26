FROM golang:1.22.2-alpine3.19

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go mod download && go mod verify

RUN go build -o ./build/chat-ws -v cmd/main.go

CMD ["./build/chat-ws"]