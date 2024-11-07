FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /usr/src/app

RUN apk add --no-cache ca-certificates upx

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -a -ldflags '-w -extldflags "-static"' -o bin/main cmd/stone-api/main.go
RUN upx bin/main

FROM scratch

COPY --from=builder /usr/src/app/bin/main /main
COPY --from=builder /usr/src/app/config.toml /config.toml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080
CMD ["/main"]