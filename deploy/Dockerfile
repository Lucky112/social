
FROM golang:1.22-alpine3.18 AS builder

WORKDIR /app

COPY ../ .

RUN go mod tidy && go build -C cmd/app -o /app/social


FROM alpine:3.18

COPY --from=builder /app/social /app/social

EXPOSE 15000

# config along the path is provided in docker-compose
CMD ["/app/social", "-config", "/app/config/config.yaml"]

