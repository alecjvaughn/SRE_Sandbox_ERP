import json
import logging
from kafka import KafkaConsumer
from prometheus_client import Counter
import os

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

inventory_events_total = Counter(
    'inventory_events_total',
    'Total number of inventory events processed',
    ['status']
)

import time

def consume_orders(broker_url=None, retry=True):
    if broker_url is None:
        broker_url = os.getenv("KAFKA_BROKERS", "kafka:9092")
        
    consumer = None
    while True:
        try:
            if not consumer:
                consumer = KafkaConsumer(
                    'orders',
                    bootstrap_servers=[broker_url],
                    group_id='inventory-group',
                    auto_offset_reset='earliest',
                    enable_auto_commit=True,
                    consumer_timeout_ms=1000 # To allow safe testing without blocking forever
                )
                logger.info(f"Connected to Kafka broker at {broker_url}. Listening for 'orders'...")

            for message in consumer:
                try:
                    # Decode message
                    payload = json.loads(message.value.decode('utf-8'))
                    logger.info(f"Received order event: {payload}")

                    # Extract fields
                    order_id = payload.get('order_id')
                    item = payload.get('item')
                    qty = payload.get('qty')

                    if not all([order_id, item, qty is not None]):
                        logger.error(f"Invalid payload format: {payload}")
                        inventory_events_total.labels(status='failed').inc()
                        continue

                    # Simulate stock decrement
                    logger.info(f"Decremented stock for {qty}x {item} (Order ID: {order_id})")

                    # Update metrics
                    inventory_events_total.labels(status='success').inc()

                except json.JSONDecodeError as e:
                    logger.error(f"Failed to decode message JSON: {e}")
                    inventory_events_total.labels(status='failed').inc()
                except Exception as e:
                    logger.error(f"Unexpected error processing message: {e}")
                    inventory_events_total.labels(status='failed').inc()
            
            # If the loop exits naturally, break out unless we are retrying forever
            if not retry:
                break
        except Exception as e:
            logger.error(f"Failed to start or run consumer: {e}. Retrying in 5 seconds...")
            if not retry:
                break
            time.sleep(5)

if __name__ == "__main__":
    consume_orders()
