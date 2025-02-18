FROM golang:1.24-alpine AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go/bin/server /go/src/app/main.go

FROM alpine:3 AS certs
RUN apk --no-cache add ca-certificates

FROM gcr.io/distroless/static-debian12 AS app

COPY --from=builder /go/bin/server /server
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/server"]