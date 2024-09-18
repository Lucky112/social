# Используем базовый образ Golang для сборки
FROM golang:1.22-alpine3.18 AS builder

# Создаём рабочую директорию
WORKDIR /app

# Копируем исходный код
COPY ../ .

# Сборка приложения
RUN go mod tidy && go build -C cmd/app -o /app/social

# Используем минимальный образ для запуска
FROM alpine:3.18

# Копируем бинарник из предыдущего шага
COPY --from=builder /app/social /app/social

# Копируем файл конфигурации
COPY --from=builder /app/deploy/config.yaml /app/config.yaml

# Порт, который будет слушать сервер
EXPOSE 15000

# Команда для запуска приложения
CMD ["/app/social", "-config", "/app/config.yaml"]
