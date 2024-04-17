FROM golang:1.21.5 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o app cmd/app/main.go



FROM scratch

COPY --from=builder /app/app /bin/app

ENTRYPOINT ["/bin/app"]