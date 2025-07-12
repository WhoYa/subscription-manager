# Docker Development Guide

## Архитектура

Проект состоит из 3 независимых контейнеров:

- **db** - PostgreSQL 17 база данных
- **api** - REST API сервер (Fiber/Go)  
- **bot** - Telegram бот (административный интерфейс)

## Быстрый старт

1. Скопируйте пример конфигурации:
```bash
cp .env.example .env
```

2. Отредактируйте `.env` файл:
   - `TOKEN` - токен Telegram бота от @BotFather
   - `ADMINS` - список Telegram ID администраторов (через запятую)
   - Остальные настройки можно оставить по умолчанию

3. Запустите все сервисы:
```bash
make up
```

## Makefile команды

### Основные команды
- `make build` - собрать все Docker образы
- `make up` - запустить все сервисы
- `make down` - остановить все сервисы
- `make restart` - перезапустить все сервисы
- `make ps` - показать статус сервисов

### Логи
- `make logs` - логи всех сервисов
- `make logs-api` - логи только API
- `make logs-bot` - логи только бота
- `make logs-db` - логи только базы данных

### Развертывание
- `make rebuild` - полная пересборка и запуск
- `make clean` - полная очистка (контейнеры, сети, volumes)

### Частичный запуск
- `make up-api` - запустить только API и DB
- `make up-bot` - запустить только бота

## Порты

- **API**: `http://localhost:8080`
- **PostgreSQL**: `localhost:5432`
- **Telegram Bot**: работает через Telegram API

## Переменные окружения

### Обязательные
- `TOKEN` - токен Telegram бота
- `ADMINS` - Telegram ID администраторов (например: `123456789,987654321`)

### Опциональные
- `PORT=8080` - порт API сервера
- `DB_*` - настройки PostgreSQL (по умолчанию подходят для Docker)
- `API_BASE_URL` - устанавливается автоматически в docker-compose

## Примеры использования

### Разработка только API
```bash
make up-api
# API доступен на http://localhost:8080
```

### Разработка только бота
```bash
make up-api    # сначала запустить API
make up-bot    # затем запустить бота
```

### Полный стек
```bash
make up
# Все сервисы запущены
```

### Отладка
```bash
make logs-bot   # смотреть логи бота
make logs-api   # смотреть логи API
```

## Troubleshooting

### Бот не запускается
1. Проверьте TOKEN в `.env`
2. Убедитесь что API сервис запущен и здоров
3. Проверьте логи: `make logs-bot`

### API недоступен
1. Проверьте что PostgreSQL запущен: `make logs-db`
2. Проверьте логи API: `make logs-api`
3. Убедитесь что порт 8080 свободен

### База данных
1. Данные сохраняются в Docker volume `db_data`
2. Для полной очистки: `make clean` 