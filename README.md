# South Park Messaging

A small message-queue demo built around South Park quotes. A Go API publishes messages to a RabbitMQ topic exchange, routed by author, and a Python consumer subscribes to the routing keys it cares about and logs the messages.

## How it works

```
                 POST /shit-post
  HTTP client ─────────────────────┐
                                   ▼
                          ┌─────────────────┐  messages.<author>  ┌─────────────────────┐  messages.#  ┌──────────────────┐
  cron (random quote) ──▶ │  api (Go, Gin)  │ ──────────────────▶ │ south-park (topic)  │ ───────────▶ │ consumer (Python)│ ──▶ logs
                          └─────────────────┘                     └─────────────────────┘   binding    └──────────────────┘
```

- **api** (`api/`): a Go service with a hexagonal (ports and adapters) layout.
  - It exposes `POST /shit-post`, which accepts a message, gives it a UUID and publishes it to the exchange.
  - It runs a cron job (`CRON_SCHEDULE`) that publishes a random predefined South Park quote on each tick.
  - Messages are published as persistent JSON to a durable topic exchange (`RABBITMQ_EXCHANGE`).
- **consumer** (`consumer/`): a Python script using `pika`. It listens on the exchange for routing keys matching `RABBITMQ_BINDING_KEY`, logs each message as `[routing key] [id] author: body` and acknowledges it.
- **rabbitmq**: the message broker, with the management UI enabled.

### Routing keys

Each message is published with the routing key `messages.<author>`. The author is lowercased and every run of characters other than letters and digits becomes a `-`. For example:

| Author          | Routing key                 |
| --------------- | --------------------------- |
| `Cartman`       | `messages.cartman`          |
| `Mr. Mackey`    | `messages.mr-mackey`        |
| `Scott Tenorman`| `messages.scott-tenorman`   |

Consumers choose which messages they get with the binding key. `*` matches exactly one word and `#` matches zero or more:

| Binding key                       | Receives                          |
| --------------------------------- | --------------------------------- |
| `messages.#` (default)            | every message                     |
| `messages.cartman`                | only Cartman                      |
| `messages.mr-*`                   | nothing: `*` matches a whole word, not part of one |

Every running consumer gets its own copy of each message that matches its binding key. To add a subscriber, start another consumer with whatever `RABBITMQ_BINDING_KEY` it should listen to.

> [!NOTE]
> AMQP can only deliver to queues, so each consumer binds a temporary queue that RabbitMQ names for it and deletes when the consumer disconnects. As a result, a message published while no consumer is listening is dropped. Nothing is buffered.

### Project layout

```
api/
  cmd/main.go                      # wiring: publisher, service, cron, HTTP router
  internal/core/domain/            # Message model
  internal/core/ports/             # MessengerService and MessagePublisher interfaces
  internal/core/services/          # MessengerService: assigns an ID and publishes
  internal/adapters/handler/       # Gin HTTP handler
  internal/adapters/messaging/     # RabbitMQ topic exchange publisher
  internal/adapters/scheduler/     # cron job with predefined quotes
consumer/
  main.py                          # listens on the exchange and logs messages
docker-compose.yml                 # api + consumer + rabbitmq
.env.example                       # configuration template
```

## Prerequisites

To run with Docker (recommended):

- [Docker](https://docs.docker.com/get-docker/) with Docker Compose v2

To run the services locally without Docker:

- Go 1.26+
- Python 3.13+ and `pip`
- A running RabbitMQ instance (for example, just the `rabbitmq` service from the compose file)

## Configuration

Copy the example env file and adjust it as needed:

```sh
cp .env.example .env
```

| Variable                | Used by            | Description                                                         | Example                              |
| ----------------------- | ------------------ | ------------------------------------------------------------------- | ------------------------------------ |
| `API_PORT`              | api                | Port the HTTP server listens on                                     | `3000`                               |
| `RABBITMQ_URL`          | api, consumer      | AMQP connection string                                              | `amqp://admin:admin@rabbitmq:5672/`  |
| `RABBITMQ_DEFAULT_USER` | rabbitmq           | Broker username                                                     | `admin`                              |
| `RABBITMQ_DEFAULT_PASS` | rabbitmq           | Broker password                                                     | `admin`                              |
| `RABBITMQ_EXCHANGE`     | api, consumer      | Name of the topic exchange                                          | `south-park`                         |
| `RABBITMQ_BINDING_KEY`  | consumer           | Routing keys the consumer listens to (defaults to `messages.#`)     | `messages.#`                         |
| `CRON_SCHEDULE`         | api                | How often a random quote is sent ([robfig/cron] syntax)             | `@every 10s`                         |

[robfig/cron]: https://pkg.go.dev/github.com/robfig/cron/v3

> [!NOTE]
> `docker-compose.yml` maps host port `3000` to container port `3000`. If you change `API_PORT`, update the port mapping too.

## Running

### With Docker Compose

```sh
docker compose up --build
```

This starts RabbitMQ, waits for its health check to pass, then starts the API and the consumer. You should see the consumer log a random quote every 10 seconds (with the default `CRON_SCHEDULE`):

```
consumer-1  | 2026-10-09 12:00:10 [messages.cartman] [3f1c...] Cartman: Respect my authoritah!
```

| Service                | URL                      |
| ---------------------- | ------------------------ |
| API                    | http://localhost:3000    |
| RabbitMQ management UI | http://localhost:15672   |

Log in to the management UI with `RABBITMQ_DEFAULT_USER` / `RABBITMQ_DEFAULT_PASS`.

Stop everything with `docker compose down` (add `-v` to also remove the RabbitMQ data volume).

### Locally

Neither service reads `.env` on its own, so export the variables in your shell first. When RabbitMQ runs on your machine, use `localhost` instead of `rabbitmq` as the host in `RABBITMQ_URL`.

```sh
# Start only the broker
docker compose up -d rabbitmq

set -a; source .env; set +a
export RABBITMQ_URL=amqp://admin:admin@localhost:5672/

# API
cd api && go run ./cmd

# Consumer (in another terminal, with the same variables exported)
cd consumer
pip install -r requirements.txt
python main.py
```

## API

### `POST /shit-post`

Publishes a message to the exchange with the routing key `messages.<author>`.

**Request body**

```json
{
  "author": "Butters",
  "body": "Oh, hamburgers."
}
```

Both `author` and `body` are required.

**Example**

```sh
curl -X POST http://localhost:3000/shit-post \
  -H "Content-Type: application/json" \
  -d '{"author": "Butters", "body": "Oh, hamburgers."}'
```

**Responses**

| Status                      | Body                                                                 |
| --------------------------- | -------------------------------------------------------------------- |
| `202 Accepted`              | The published message, including its generated `id`                 |
| `400 Bad Request`           | `{"error": "..."}` when the body is invalid or a field is missing    |
| `500 Internal Server Error` | `{"error": "..."}` when publishing to RabbitMQ fails                 |

A `202` means the exchange accepted the message. It doesn't guarantee that any consumer was listening.

```json
{
  "id": "6b0f4a8e-2c1d-4f3a-9e7b-1a2b3c4d5e6f",
  "author": "Butters",
  "body": "Oh, hamburgers."
}
```
