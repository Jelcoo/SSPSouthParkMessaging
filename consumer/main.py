import json
import os

import pika

RABBITMQ_URL = os.getenv("RABBITMQ_URL")
QUEUE = os.getenv("RABBITMQ_QUEUE")


def on_message(channel, method, properties, body):
    try:
        message = json.loads(body)
        print(f"[{message.get('id')}] {message.get('author')}: {message.get('body')}", flush=True)
    except json.JSONDecodeError:
        print(f"Received non-JSON message: {body!r}", flush=True)
    channel.basic_ack(delivery_tag=method.delivery_tag)


def main():
    connection = pika.BlockingConnection(pika.URLParameters(RABBITMQ_URL))
    channel = connection.channel()
    channel.queue_declare(queue=QUEUE, durable=True)
    channel.basic_consume(queue=QUEUE, on_message_callback=on_message)

    print(f"Waiting for messages on '{QUEUE}'...", flush=True)
    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        channel.stop_consuming()
    finally:
        connection.close()


if __name__ == "__main__":
    main()
