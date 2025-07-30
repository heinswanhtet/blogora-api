FROM golang:1.24.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl make

RUN curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh

RUN alias goose='/usr/local/bin/goose'

WORKDIR /root/

COPY --from=builder /app/Makefile .
COPY --from=builder /app/migrations/ ./migrations
COPY --from=builder /app/main .


EXPOSE 8080

CMD ["./main"]