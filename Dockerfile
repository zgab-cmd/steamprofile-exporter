FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/steamprofile-exporter .

FROM scratch
COPY --from=builder /app/steamprofile-exporter /steamprofile-exporter
EXPOSE 9805
ENTRYPOINT ["/steamprofile-exporter"]
