FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /service /app/cmd/service

EXPOSE 8000

CMD ["/service"]
