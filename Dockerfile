FROM golang:1.14-alpine
RUN apk add --no-cache git postgresql-client
WORKDIR /apps/pusher
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o bin/pusher cmd/pusher/pusher.go
EXPOSE 8080

CMD ["bin/pusher"]
