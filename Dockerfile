# syntax=docker/dockerfile:1

FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend-app

COPY internal/db/migrations/*.sql /migrations/

CMD "/backend-app" "--port" $SERVER_PORT "--db_source" $DATABASE_URL "--migrations" "/migrations"


