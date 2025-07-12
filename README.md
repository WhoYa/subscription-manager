# Subscription Manager 2.0

**Самохостингный сервис для прозрачного управления групповыми подписками**

Subscription Manager 2.0 - это Go-based REST API + Telegram бот система, которая помогает владельцам подписок прозрачно собирать деньги с друзей/семьи за совместные подписки типа Netflix, Spotify, YouTube Premium и др.

## 🌟 Ключевые возможности

- **Управление пользователями и подписками** с гибкими настройками ценообразования
- **Система ручного управления курсами валют** админом
- **Автоматический расчет платежей** с учетом курсов валют и наценок
- **Подробная аналитика прибыли** для администраторов
- **Логирование всех платежей** с полной историей
- **REST API** для интеграции с внешними системами
- **Готовность к интеграции с Telegram ботом**

## 🏗️ Архитектура

### Технологический стек
- **Backend**: Go 1.24+ с Fiber web framework
- **Database**: PostgreSQL 17 с GORM ORM
- **Containerization**: Docker + Docker Compose
- **API Style**: RESTful with JSON

### Структура проекта
```
subscription-manager/
├── cmd/api/                    # Main application entry point
├── internal/                   # Private application code
│   ├── handlers/              # HTTP handlers (controllers)
│   ├── repository/            # Data layer (GORM implementations)
│   ├── service/               # Business logic layer
│   ├── app/                   # Application bootstrap
│   └── util/                  # Utilities and helpers
├── pkg/                       # Public packages
│   └── db/                    # Database models, enums, migrations
├── docker-compose.yml         # Development environment
├── dockerfile                 # Production container
└── Subscription_Manager.postman_collection.json  # API documentation
```

### Архитектурные принципы
- **Layered Architecture**: Handlers → Services → Repositories
- **Interface-first approach**: Repository pattern с интерфейсами
- **Environment-based configuration**: Все настройки через переменные окружения
- **Database-first design**: PostgreSQL с строгими ограничениями

## 🚀 Быстрый старт

### Предварительные требования
- Docker и Docker Compose
- Git

### Установка и запуск

1. **Клонируйте репозиторий**
```bash
git clone https://github.com/WhoYa/subscription-manager
cd subscription-manager
```

2. **Настройте переменные окружения**
```bash
cp .env.example .env
# Отредактируйте .env под ваши нужды
```

3. **Запустите сервисы**
```bash
make st
```

Это команда выполнит:
- Сборку Docker образа API
- Запуск PostgreSQL и API контейнеров
- Автоматическую миграцию базы данных

### Доступные команды Make

```bash
make build    # Собрать Docker образ
make up       # Запустить сервисы
make ps       # Показать статус контейнеров  
make rm       # Удалить контейнеры и volumes
make st       # Полная пересборка и запуск
```

## 📋 Базовый workflow

### 1. Создание администратора
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "tg_id": 123456789,
    "username": "admin",
    "fullname": "Admin User", 
    "is_admin": true
  }'
```

### 2. Установка курсов валют (админ)
```bash
curl -X POST http://localhost:8080/api/admin/{admin_user_id}/currency/set \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "USD",
    "rate": 95.50
  }'
```

### 3. Создание подписки
```bash
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Netflix",
    "base_price": 15.99,
    "base_currency": "USD",
    "period_days": 30
  }'
```

### 4. Подписка пользователя на сервис
```bash
curl -X POST http://localhost:8080/api/users/{user_id}/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "subscription_id": "{subscription_id}",
    "pricing_mode": "percent",
    "markup_percent": 10.0
  }'
```

### 5. Расчет платежа
```bash
curl "http://localhost:8080/api/calculate/{user_id}/{subscription_id}?due_date=2024-12-15"
```

## 💰 Система ценообразования

Система поддерживает гибкие модели ценообразования:

### Режимы наценки (PricingMode)
- **none**: Без наценки, чистая цена
- **percent**: Процентная наценка
- **fixed**: Фиксированная доплата в рублях

### Приоритет настроек
1. Индивидуальные настройки пользователя
2. Глобальные настройки системы
3. Базовая цена без наценки

### Пример расчета
- Netflix: $15.99 (базовая цена)
- Курс USD: 95.50₽
- Наценка пользователя: 10%
- **Итого: (15.99 × 95.50) × 1.10 = 1,681.89₽**

## 🔧 API Endpoints

### Основные группы endpoints:

- **`/api/users`** - Управление пользователями
- **`/api/subscriptions`** - Управление подписками  
- **`/api/currency_rates`** - Работа с курсами валют
- **`/api/payments`** - Логирование платежей
- **`/api/settings`** - Глобальные настройки
- **`/api/calculate`** - Расчет платежей

### Административные endpoints:

- **`/api/admin/{admin_id}/currency`** - Управление курсами валют
- **`/api/admin/{admin_id}/profit`** - Аналитика прибыли

Полная документация API доступна в Postman коллекции: `Subscription_Manager.postman_collection.json`

## 🎯 Roadmap

### ✅ Completed (95-100%)
- **Data Layer**: Модели, миграции, репозитории
- **Core API**: CRUD операции для всех сущностей
- **Payment Calculation**: Система расчета платежей  
- **Currency Management**: Ручное управление курсами админом
- **Profit Analytics**: Детальная аналитика прибыли
- **Admin Panel**: API endpoints для администрирования

### 🚧 In Progress (10%)
- **Service Layer**: Бизнес-логика и валидации
- **Error Handling**: Улучшенная обработка ошибок
- **Input Validation**: Комплексная валидация входных данных

### 📋 Planned (0%)
- **Telegram Bot**: 
  - Уведомления о предстоящих платежах
  - Интерфейс для пользователей через Telegram
  - Запрос курсов валют у админа
- **Scheduler Service**:
  - Автоматические уведомления по расписанию
  - Напоминания о просроченных платежах
- **Enhanced Features**:
  - Групповые подписки (несколько пользователей на одну подписку)
  - Automatic payment reminders
  - Statistics dashboard
- **DevOps & Production**:
  - CI/CD pipeline
  - Production-ready deployment
  - Monitoring and logging
  - Backup strategies

## 🔒 Security Features

- **Admin Access Control**: Middleware проверки прав администратора
- **Input Validation**: Валидация всех входящих данных
- **SQL Injection Protection**: Использование GORM с параметризованными запросами
- **Environment Isolation**: Секреты только через переменные окружения

## 🧪 Development

### Структура базы данных

**Основные таблицы:**
- `users` - Пользователи с Telegram интеграцией
- `subscriptions` - Сервисы подписок
- `user_subscriptions` - Связь пользователей с подписками + настройки ценообразования
- `payment_logs` - История всех платежей
- `currency_rates` - Курсы валют (ручное управление)
- `global_settings` - Глобальные настройки системы

### Поддерживаемые валюты
- **USD** - Доллар США
- **EUR** - Евро  
- **RUB** - Российский рубль (базовая валюта для расчетов)

### Health Check
```bash
curl http://localhost:8080/api/healthz
```

## 📝 Contributing

1. Форкните репозиторий
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Зафиксируйте изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📄 License

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для подробностей.

## 🤝 Support

При возникновении вопросов или проблем:
1. Проверьте [Issues](https://github.com/WhoYa/subscription-manager/issues)
2. Создайте новый Issue с подробным описанием
3. Используйте Postman коллекцию для тестирования API

---

**Subscription Manager 2.0** - делает групповые подписки прозрачными и удобными! 🎯
