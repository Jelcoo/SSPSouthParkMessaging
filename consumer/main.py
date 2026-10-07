import json
import logging
import os

import pika

RABBITMQ_URL = os.getenv("RABBITMQ_URL")
QUEUE = os.getenv("RABBITMQ_QUEUE")

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(message)s", datefmt="%Y-%m-%d %H:%M:%S")
logging.getLogger("pika").setLevel(logging.WARNING)
log = logging.getLogger(__name__)


def on_message(channel, method, properties, body):
    try:
        message = json.loads(body)
        log.info(f"[{message.get('id')}] {message.get('author')}: {message.get('body')}")
    except json.JSONDecodeError:
        log.info(f"Received non-JSON message: {body!r}")
    channel.basic_ack(delivery_tag=method.delivery_tag)


def main():
    connection = pika.BlockingConnection(pika.URLParameters(RABBITMQ_URL))
    channel = connection.channel()
    channel.queue_declare(queue=QUEUE, durable=True)
    channel.basic_consume(queue=QUEUE, on_message_callback=on_message)

    log.info(f"Waiting for messages on '{QUEUE}'...")
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        channel.stop_consuming()
    finally:
        connection.close()


if __name__ == "__main__":
    main()
