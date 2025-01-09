# Telegram Bot для работы с RAG системой

Этот проект представляет собой Telegram-бота, предназначенного для взаимодействия пользователей с RAG системой. Бот
обрабатывает вопросы с помощью AI и сохраняет логи сообщений для генерации отчетов.

## Основные возможности

- Обработка вопросов от пользователей и ответов с помощью AI.
- Хранение данных о пользователях, сообщениях и логах.
- Статистика для администраторов через команду `/report`.
- Анкетирование пользователей и составление контекстов запроса
- Биллинг-система (по подписке) с интеграцией через YooKassa и CloudPayments
- Возможность гибко настраивать текста, меню бота, различные команды, статические запросы или подготовленные запросы к ИИ

## Основные команды

| Команда         | Описание                                                                                      |
|-----------------|-----------------------------------------------------------------------------------------------|
| `/start`        | Запускает бота и регистрирует нового пользователя.                                            | |
| `/menu`         | Открывает главное меню, если функционал включен                                               |
| `/form`         | Переходит в режим анкетирования пользователя, если функционал включен.                        |
| `/subscription` | Возвращает информацию о подписке пользователя, если включен биллинг.                          |
| `/report`       | Доступно только администраторам. Выводит отчет по подписчикам и сообщениям за разные периоды. |
| `/*`            | Возвращает статичное сообщение согласно конфигурации бота.                                    |

## Требования

- Docker
- Docker Compose
- PostgreSQL (версия 14 или выше)
- Golang (для разработки и создания миграций)

## Настройка

1. Клонируйте репозиторий:

    ```bash
    git clone https://github.com/khonikdev/gorag-tg-bot.git
    ```

2. Создайте `.env` файл в корне проекта и заполните его переменными окружения:

    ```plaintext
    # Telegram Bot
    APP_NAME=my_bot  # название проекта
    
    # данные для локальной БД
    DB_USER=my_bot_user
    DB_PASSWORD=my_bot_password
    DB_NAME=my_bot_db
    DB_PORT=5432
    EXPOSE_DB_PORT=5431
   
    # Порт для работы API
    EXPOSE_API_PORT=3000
   
    # Параметры docker compose 
    COMPOSE_PROJECT_NAME=rag
    NETWORK_NAME=rag-net
    ```
3. Создайте `bot/bot.yaml` файл и сконфигурируйте бота

```yaml
app_name: ""
app_url: ""
bot_token: ""
rag_api_endpoint: "https://ai.medichain.ai"
rag_api_token: ""

database:
  host: "db"
  port: "5432"
  user: ""
  password: ""
  name: ""

ai_mode_message: "Я готов ответить на ваши вопросы."
temporary_message: "Пожалуйста, подождите, я обрабатываю ваш запрос..."

commands:
  - name: start
    label: "Запустить бота"
    enabled: true
    show_in_menu: true
    message: "Здравствуйте! Чем могу помочь?"
  - name: about
    label: "О программе"
    enabled: true
    show_in_menu: true
    message: "О программе: Этот бот создан для помощи пользователям. Разработан командой энтузиастов."
  - name: help
    label: "Помощь и техподдержка"
    enabled: true
    show_in_menu: true
    message: "Если у вас возникли вопросы или нужна помощь, пожалуйста, свяжитесь с нашей службой поддержки."
  - name: report
    label: ""
    enabled: true
    show_in_menu: false
    message: ""
  - name: form
    label: "Анкета"
    enabled: true
    show_in_menu: true
    message: ""
  - name: subscription
    label: "Подписка"
    enabled: true
    show_in_menu: true
    message: ""
  - name: prompt
    label: "Задать вопрос"
    enabled: true
    show_in_menu: true
    message: "Выберите наиболее подходящую команду боту"
    replies:
      - Подготовь мне рацион питания на неделю
      - Какие продукты не содержат глицерин?
      - Какой у меня индекс массы тела?
      - Какая калорийность у подсолнечного масла?
  - name: themas
    label: "План питания"
    enabled: true
    show_in_menu: true
    message: "Выберите необходимое действие"
    actions:
      - label: "Научные основы потери веса и почему так трудно его удержать"
        prompt: "Объясни научные основы потери веса и причины, по которым сложно удержать сниженный вес."
        act_data: "1"
        act_unique: "command_prompt" # DEFAULT VALUE
      - label: "План тренировок"
        prompt: "Подготовь мне план тренировок исходя из моих реалий и целей"
        act_data: "2"
        act_unique: "command_prompt" # DEFAULT VALUE

menu:
  enabled: true
  label: "Главное меню"
  items:
    - key: "document"
      enabled: true
      button_label: "Документ"
      message: ""
      context: ""
      actions: [ ]
    - key: "training"
      enabled: true
      button_label: "ИИ тренер"
      message: ""
      context: ""
      actions:
        - label: "Начать"
          prompt: ""
          act_data: "training"
          act_unique: "command_prompt"


modules:
  likes:
    enabled: true
    like_label: "👍 Нравится"
    like_response: "Спасибо за вашу оценку!"
    dislike_label: "👎 Не нравится"
    dislike_response: "Спасибо за ваш отзыв, мы постараемся улучшить работу."
  files:
    enabled: true
  audio:
    enabled: true
  form:
    enabled: true
    show_on_start: true
    can_skip: false
    disclaimer_label: "Чтобы улучшить работу бота, пожалуйста, заполните небольшую анкету."
    start_label: "Заполнить анкету"
    later_label: "Вернуться к анкете позже"
    later_message: "Чтобы вернуться к анкете позже, выберите соответствующий пункт в меню или введите команду /form."
    allow_edit: true
    edit_label: "Отредактировать анкету"
    select_question_message: "Выберите вопрос"
    view_label: "Посмотреть анкету"
    completed_message: "Спасибо за заполнение анкеты!"
    wrong_option_message: "Пожалуйста, выберите один из доступных вариантов"
    custom_option_message: "Введите свой вариант"
    next_question_label: "Перейти к следующему вопросу"
    more_option_label: "Выбрать еще вариант"
    option_saved_message: "Ваш ответ сохранён"
    add_option_message: "Выберите еще вариант?"
    context_prefix: "Информация о пользователе-враче:"
    questions:
      - text: "Выберите ваш пол"
        order: 0
        is_required: true
        type: select
        selectable_options_count: 1
        options:
          - text: "Мужской"
          - text: "Женский"
      - text: "Укажите возраст"
        order: 1
        is_required: true
        type: select
        selectable_options_count: 1
        options:
          - text: "до 18"
          - text: "18-25"
          - text: "26-35"
          - text: "36-45"
          - text: "46-55"
          - text: "56-65"
          - text: "старше 65"
      - text: "Высота (см)"
        order: 2
        is_required: true
        type: number
        hint: "Введите высоту в см, например, 170"
      - text: "Вес"
        order: 3
        is_required: true
        type: number
        hint: "Введите вес в кг, например, 70"
      - text: "У вас есть аллергия?"
        order: 4
        is_required: true
        type: select
        selectable_options_count: 1
        hint: "Если есть, укажите"
        options:
          - text: "Да"
            require_additional_text: true
          - text: "Нет"
      - text: "Какие ваши цели?"
        order: 5
        is_required: true
        type: select
        selectable_options_count: 9
        options:
          - text: "Improve overall health"
          - text: "Lose weight"
          - text: "Gain weight"
          - text: "Increase energy levels"
          - text: "Enhance athletic performance"
          - text: "Improve digestion"
          - text: "Strengthen the immune system"
          - text: "Manage a health condition"
          - text: "Other"
            require_additional_text: true
      - text: "Расскажите немного о себе"
        order: 6
        is_required: true
        type: text


  billing:
    enabled: true
    subscription_alert: "Подписка на сервис отсутствует или истекла"
    subscription_message: "Подписка действительна до: {expires_at}"
    subscribe_btn: "Перейдите к оплате"
    unsubscribe_btn: "Отменить подписку"
    unsubscribe_confirmation_message: "Вы точно уверены, что хотите отменить подписку?\n
    Без подписки вы потеряете доступ к функционалу DIMA, который экономит ваше время и помогает в работе\\.\n
    👇 Подтвердите действие или вернитесь назад\\."
    unsubscribe_cancel: "Экономить время"
    unsubscribe_confirm: "Отменить подписку"
    bye_message: "Подписка отменена. Вы сможете пользоваться ботом до конца оплаченного периода. Следующего списания не произойдет."
    not_enough_amount_notification: "Некорректная сумма оплаты, попробуйте еще раз"
    subscription_granted_notification: "Подписка успешно продлена. Спасибо, что Вы с нами!"
    providers:
      cloud_payments:
        enabled: false
        public_key: ""
        secret_key: ""
        period: 1
        period_unit: "Month" # [ Day, Week, Month ]
        price: 1000
      yoo_kassa:
        enabled: true
        token: ""
        period: 1
        period_unit: "Month"
        price: 99000
integrations: { }
```
## Развертывание с помощью Docker Compose

1. Постройте и запустите контейнеры:

    ```bash
    docker-compose up -d
    ```

2. Примените миграции для базы данных, чтобы создать необходимые таблицы и индексы:

    ```bash
    make migrate-up
    ```

3. Для создания новой миграции (опционально):

    ```bash
    make migration
    ```

4. Для отката миграций (опционально):

    ```bash
    make migrate-down
    ```
5. Заполнение таблиц согласно конфигурационному файлу:

    ```bash
    make seed-form
    make seed-contexts
    make seed-payment-providers
    ```

## Структура проекта

- `bot/cmd/bot`: Основной файл запуска бота.
- `bot/cmd/api`: Основной файл запуска API.
- `bot/cmd/seeder/seeds`: Сидеры БД
- `bot/controllers`: Контроллеры для работы с API.
- `bot/core/bot`: Глобальная сущность Telegram Bot.
- `bot/database/migrations`: SQL-файлы для миграций базы данных.
- `bot/models`: Модели для работы с базой данных.
- `bot/handlers`: Логика обработки команд бота.
- `bot/services`: Логика взаимодействия с различными сервисами

## Логирование и отладка

- Логирование запросов и ответов AI производится в консоль. Для включения более подробного логирования установите
  переменную окружения `DEBUG=true`.
- Все сообщения пользователей и ответы AI сохраняются в таблице `messages` для последующего анализа.

## Поддержка

Для получения помощи или поддержки, свяжитесь с нами по электронной почте: `khonikdev@gmail.com`.