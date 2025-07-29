# Docker Compose Files

1. Create folders

```bash
mkdir -p volumes/kafka/broker-1/
mkdir -p volumes/kafka/broker-2/
mkdir -p volumes/kafka/broker-3/
mkdir -p volumes/zookeeper/data
mkdir -p volumes/zookeeper/transactions
```

2. Run docker compose

```bash
docker compose -f common.yml -f zookeeper.yml up
docker compose -f common.yml -f kafka_cluster.yml up
docker compose -f common.yml -f init_kafka.yml up
```