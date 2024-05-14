FROM golang:1.21.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o app cmd/app/main.go

FROM golang:1.21.5-alpine
WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/.env .

ENTRYPOINT ["./app"]