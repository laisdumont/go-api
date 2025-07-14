FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o app

FROM golang:latest

WORKDIR /root/

COPY --from=builder /app/app .
COPY .env .

EXPOSE 8080

CMD ["./app"]
