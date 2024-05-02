FROM golang:1.21.5 AS builder
WORKDIR /app

COPY . ./

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o app cmd/app/main.go



FROM scratch
WORKDIR /app

COPY --from=builder /app/app /bin/app

ENTRYPOINT ["/bin/app"]