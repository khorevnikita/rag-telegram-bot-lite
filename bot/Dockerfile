# Stage 1: Build
FROM golang:1.23 AS bot_base

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

FROM bot_base as bot_builder
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o start_bot ./cmd/bot/main.go

FROM bot_base as api_builder
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o start_api ./cmd/api/main.go

# Stage 2: Run
FROM alpine:latest as bot

# Устанавливаем переменные среды
ENV GIN_MODE=release

# Устанавливаем ffmpeg
RUN apk add --no-cache ffmpeg

# Создаем непривилегированного пользователя для повышения безопасности
RUN adduser -D -g '' appuser

# Устанавливаем рабочую директорию
WORKDIR /home/appuser

# Копируем собранный бинарный файл из предыдущего этапа
COPY --from=bot_builder /app/start_bot .
COPY --from=bot_base /app/bot.yaml .

# Устанавливаем права на выполнение
RUN chmod +x start_bot

# Переключаемся на непривилегированного пользователя
USER appuser

RUN mkdir "storage"

# Запуск приложения
CMD ["./start_bot"]

# Stage 2: Run
FROM alpine:latest as api

# Устанавливаем переменные среды
ENV GIN_MODE=release

# Создаем непривилегированного пользователя для повышения безопасности
RUN adduser -D -g '' appuser

# Устанавливаем рабочую директорию
WORKDIR /home/appuser

COPY --from=api_builder /app/start_api .
COPY --from=bot_base /app/templates ./templates
COPY --from=bot_base /app/bot.yaml .

# Устанавливаем права на выполнение
RUN chmod +x start_api

# Переключаемся на непривилегированного пользователя
USER appuser

# Запуск приложения
CMD ["./start_api"]
