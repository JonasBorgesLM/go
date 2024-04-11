FROM golang:1.21.5
WORKDIR /app

COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o app src/app/main.go

CMD ["./app"]