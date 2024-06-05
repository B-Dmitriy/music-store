FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && go mod verify
RUN apk add gcc
RUN apk add musl-dev
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o musicstore ./cmd/main.go


FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/musicstore .
COPY --from=builder /app/data/musicstore.db ./data/musicstore.db
COPY --from=builder /app/config/main.yaml ./config/main.yaml
RUN apt-get update && apt-get install -y gcc && apt-get install -y musl-dev
EXPOSE 5050
CMD ["./musicstore"]