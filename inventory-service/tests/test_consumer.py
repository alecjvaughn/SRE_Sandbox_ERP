import pytest
import json
from unittest.mock import MagicMock, patch
from consumer import consume_orders, inventory_events_total

class MockKafkaMessage:
    def __init__(self, value, topic='orders'):
        self.value = value
        self.topic = topic

@patch('consumer.KafkaConsumer')
def test_consume_orders_success(mock_kafka_consumer_cls):
    # Setup mock consumer
    mock_consumer = MagicMock()
    mock_kafka_consumer_cls.return_value = mock_consumer
    
    # Create a mock order event
    order_event = {
        "order_id": "test-order-123",
        "item": "Widget",
        "qty": 5,
        "timestamp": "2026-06-20T12:00:00Z"
    }
    
    mock_message = MockKafkaMessage(json.dumps(order_event).encode('utf-8'))
    
    # Simulate an infinite loop that yields our message, then raises StopIteration
    mock_consumer.__iter__.return_value = iter([mock_message])
    
    # Get initial metric value
    initial_val = inventory_events_total.labels(status='success')._value.get()
    
    # Run the consumer logic
    consume_orders("test-broker:9092", retry=False)
    
    # Verify metric incremented
    final_val = inventory_events_total.labels(status='success')._value.get()
    assert final_val == initial_val + 1

@patch('consumer.KafkaConsumer')
def test_consume_orders_invalid_json(mock_kafka_consumer_cls):
    mock_consumer = MagicMock()
    mock_kafka_consumer_cls.return_value = mock_consumer
    
    mock_message = MockKafkaMessage(b"invalid-json")
    mock_consumer.__iter__.return_value = iter([mock_message])
    
    initial_val = inventory_events_total.labels(status='failed')._value.get()
    
    consume_orders("test-broker:9092", retry=False)
    
    final_val = inventory_events_total.labels(status='failed')._value.get()
    assert final_val == initial_val + 1

@patch('consumer.KafkaConsumer')
def test_consume_orders_missing_fields(mock_kafka_consumer_cls):
    mock_consumer = MagicMock()
    mock_kafka_consumer_cls.return_value = mock_consumer
    
    order_event = {
        "order_id": "test-order-123"
    }
    
    mock_message = MockKafkaMessage(json.dumps(order_event).encode('utf-8'))
    mock_consumer.__iter__.return_value = iter([mock_message])
    
    initial_val = inventory_events_total.labels(status='failed')._value.get()
    
    consume_orders(None, retry=False)
    
    final_val = inventory_events_total.labels(status='failed')._value.get()
    assert final_val == initial_val + 1
