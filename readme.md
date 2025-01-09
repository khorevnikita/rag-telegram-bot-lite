# Telegram Bot for RAG System Integration

This project is a Telegram bot designed to facilitate user interaction with a RAG (Retrieval-Augmented Generation)
system. The bot processes user queries using AI and logs messages for generating reports.

## Key Features

- Process user queries and provide AI-powered responses.
- Store data about users, messages, and logs.
- Provide administrators with statistics via the `/report` command.
- Conduct user surveys and create contextual requests.
- Subscription-based billing system with integration through YooKassa and CloudPayments.
- Flexible configuration for bot texts, menus, various commands, static or prepared AI queries.

## Main Commands

| Command         | Description                                                                                            |
|-----------------|--------------------------------------------------------------------------------------------------------|
| `/start`        | Starts the bot and registers a new user.                                                               |
| `/menu`         | Opens the main menu, if enabled.                                                                       |
| `/form`         | Switches to user survey mode, if enabled.                                                              |
| `/subscription` | Provides subscription information, if billing is enabled.                                              |
| `/report`       | Available to administrators only. Displays a report on subscribers and messages for different periods. |
| `/*`            | Returns a static message according to the bot configuration.                                           |

## Requirements

- Docker
- Docker Compose
- PostgreSQL (version 14 or higher)
- Golang (for development and migration creation)

## Setup

1. Clone the repository:

    ```bash
    git clone https://github.com/khorevnikita/rag-telegram-bot-lite
    ```

2. Create a `.env` file in the project root and fill it with the following environment variables:

    ```plaintext
    # Telegram Bot
    APP_NAME=my_bot  # project name

    # Local database credentials
    DB_USER=my_bot_user
    DB_PASSWORD=my_bot_password
    DB_NAME=my_bot_db
    DB_PORT=5432
    EXPOSE_DB_PORT=5431

    # API port
    EXPOSE_API_PORT=3000

    # Docker Compose parameters
    COMPOSE_PROJECT_NAME=rag
    NETWORK_NAME=rag-net
    ```

3. Copy and customize the `bot.yaml` file from the provided example:

    ```bash
    cp bot/bot.yaml.example bot/bot.yaml
    ```

   The configuration file contains the following top-level keys:

    ```yaml
    app_name: string
    app_url: string
    bot_token: string
    rag_api_endpoint: string
    rag_api_token: string

    database:
      host: string
      port: string
      user: string
      password: string
      name: string

    ai_mode_message: string
    temporary_message: string

    commands: []
    menu: {}
    modules: {}
    integrations: {}
    ```

## Deployment with Docker Compose

1. Build and start the containers:

    ```bash
    docker-compose up -d
    ```

2. Apply database migrations to create necessary tables and indexes:

    ```bash
    make migrate-up
    ```

3. To create a new migration (optional):

    ```bash
    make migration
    ```

4. To rollback migrations (optional):

    ```bash
    make migrate-down
    ```

5. Seed the database according to the configuration file:

    ```bash
    make seed-form
    make seed-contexts
    make seed-payment-providers
    ```

## Project Structure

- `bot/cmd/bot`: Main file for bot execution.
- `bot/cmd/api`: Main file for API execution.
- `bot/cmd/seeder/seeds`: Database seeders.
- `bot/controllers`: Controllers for API operations.
- `bot/core/bot`: Global Telegram Bot entity.
- `bot/database/migrations`: SQL files for database migrations.
- `bot/models`: Models for database operations.
- `bot/handlers`: Logic for bot command processing.
- `bot/services`: Logic for interacting with various services.

## Support

For help or support, contact us via email: `khonikdev@gmail.com`.

