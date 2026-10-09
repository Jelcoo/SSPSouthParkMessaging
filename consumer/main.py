import json
import logging
import os

import pika

RABBITMQ_URL = os.getenv("RABBITMQ_URL")
EXCHANGE = os.getenv("RABBITMQ_EXCHANGE")
BINDING_KEY = os.getenv("RABBITMQ_BINDING_KEY", "messages.#")

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(message)s", datefmt="%Y-%m-%d %H:%M:%S")
logging.getLogger("pika").setLevel(logging.WARNING)
log = logging.getLogger(__name__)


def on_message(channel, method, properties, body):
    try:
        message = json.loads(body)
        log.info(f"[{method.routing_key}] [{message.get('id')}] {message.get('author')}: {message.get('body')}")
    except json.JSONDecodeError:
        log.info(f"[{method.routing_key}] Received non-JSON message: {body!r}")
    channel.basic_ack(delivery_tag=method.delivery_tag)


def main():
    connection = pika.BlockingConnection(pika.URLParameters(RABBITMQ_URL))
    channel = connection.channel()
    channel.exchange_declare(exchange=EXCHANGE, exchange_type="topic", durable=True)
    queue = channel.queue_declare(queue="", exclusive=True).method.queue
    channel.queue_bind(queue=queue, exchange=EXCHANGE, routing_key=BINDING_KEY)
    channel.basic_consume(queue=queue, on_message_callback=on_message)

    log.info(f"Listening on exchange '{EXCHANGE}' with '{BINDING_KEY}'...")
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        channel.stop_consuming()
    finally:
        connection.close()


if __name__ == "__main__":
    main()
