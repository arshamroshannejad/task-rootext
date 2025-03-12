FROM golang:1.24-alpine AS builder
WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o /go/bin/main ./cmd/server/

FROM alpine:3.21 AS prod
WORKDIR /production
COPY --from=builder /go/bin/main .
ENTRYPOINT ["./main"]