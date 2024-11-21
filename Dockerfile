FROM golang:1.23.3-alpine AS build

WORKDIR /app

RUN apk --no-cache add gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=1

RUN go build -o TODO_list ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/TODO_list .

ENV APP_ADDRESS="0.0.0.0:8080"
ENV DB_DRIVER="sqlite3"
ENV DB_NAME="/app/db/todo_list.db"

VOLUME /app/db

CMD ["./TODO_list"]