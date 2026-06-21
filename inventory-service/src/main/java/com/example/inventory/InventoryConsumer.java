package com.example.inventory;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.micrometer.core.instrument.MeterRegistry;
import io.quarkus.runtime.Startup;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.ConsumerRecords;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.jboss.logging.Logger;

import java.time.Duration;
import java.util.Collections;
import java.util.Properties;

@ApplicationScoped
@Startup
public class InventoryConsumer {

    private static final Logger LOG = Logger.getLogger(InventoryConsumer.class);

    @ConfigProperty(name = "kafka.bootstrap.servers", defaultValue = "localhost:9092")
    String bootstrapServers;

    @ConfigProperty(name = "kafka.topic", defaultValue = "orders")
    String topic;

    @Inject
    ObjectMapper objectMapper;

    @Inject
    MeterRegistry meterRegistry;

    private KafkaConsumer<String, String> consumer;
    private volatile boolean running = true;

    @PostConstruct
    void init() {
        Properties props = new Properties();
        props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
        props.put(ConsumerConfig.GROUP_ID_CONFIG, "inventory-group");
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class.getName());
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        consumer = new KafkaConsumer<>(props);
        consumer.subscribe(Collections.singletonList(topic));

        new Thread(this::pollLoop, "kafka-consumer-thread").start();
        LOG.infov("Kafka consumer started listening on {0} for topic {1}", bootstrapServers, topic);
    }

    private void pollLoop() {
        try {
            while (running) {
                ConsumerRecords<String, String> records = consumer.poll(Duration.ofMillis(100));
                for (ConsumerRecord<String, String> record : records) {
                    processRecord(record);
                }
            }
        } catch (Exception e) {
            if (running) {
                LOG.error("Error in Kafka poll loop", e);
            }
        } finally {
            consumer.close();
        }
    }

    private void processRecord(ConsumerRecord<String, String> record) {
        try {
            JsonNode order = objectMapper.readTree(record.value());
            String item = order.has("item") ? order.get("item").asText() : "Unknown";
            LOG.infov("Processing stock reduction for item: {0}", item);
            meterRegistry.counter("inventory_events_total", "type", "order_received").increment();
        } catch (Exception e) {
            meterRegistry.counter("inventory_errors_total", "type", "payload_parse_error").increment();
            LOG.error("Failed to parse order payload", e);
        }
    }

    @PreDestroy
    void shutdown() {
        running = false;
        consumer.wakeup(); // Break out of poll loop immediately
    }
}
