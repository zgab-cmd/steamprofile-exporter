FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN apk add --no-cache ca-certificates \
    && update-ca-certificates \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/steamprofile-exporter .

FROM scratch
LABEL org.opencontainers.image.source="https://github.com/zgab-cmd/steamprofile-exporter" \
      org.opencontainers.image.licenses="MIT"
COPY --from=builder /app/steamprofile-exporter /steamprofile-exporter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 9805
ENTRYPOINT ["/steamprofile-exporter"]
