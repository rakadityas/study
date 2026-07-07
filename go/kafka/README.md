# Kafka Example

## What this demonstrates

Apache Kafka is a distributed, partitioned, replicated commit log. Producers
append records to topics; consumers read those records independently, at
their own pace, by tracking an offset. Unlike a traditional message queue,
reading a message does not delete it -- records live for the topic's
retention period (or forever, if compacted) so many independent consumers
can replay the same log.

This example builds a tiny "order events" pipeline and uses it to make the
core Kafka concepts concrete:

- **Topics & partitions**: a topic (`orders.v1`) is split into partitions,
  each an append-only, strictly-ordered log. Ordering is only guaranteed
  *within* a partition, never across partitions.
- **Keys & ordering**: producers key each message by `customer_id`. Kafka
  hashes the key to deterministically pick a partition, so every event for
  a given customer lands in the same partition and is read back in the
  order it was written.
- **Consumer groups & rebalancing**: consumers join a named group
  (`KAFKA_GROUP_ID`). The broker divides the topic's partitions among the
  group's members and automatically rebalances when a member joins or
  leaves. Two consumers in the *same* group split the work (queue-like);
  two consumers in *different* groups each get a full copy of every
  message (pub/sub-like).
- **Offsets & delivery semantics**: consumers commit offsets after
  processing (`CommitInterval: 0` + manual `CommitMessages`), giving
  at-least-once delivery -- a crash before committing replays the message
  on restart instead of losing it.
- **Replication & acks**: the topic's replication factor controls how many
  broker copies of each partition exist; `RequiredAcks` on the producer
  controls how many replicas must confirm a write before it's considered
  successful.

## Run

Start Kafka (KRaft mode -- no ZooKeeper needed) and Kafka UI:

```bash
cd go/kafka
docker compose up -d
```

Kafka UI is at http://localhost:8089 -- use it to browse topics, partitions,
consumer groups, and individual messages as you run the examples below.

Create the topic with 3 partitions:

```bash
cd go/kafka
KAFKA_BROKER='localhost:9092' KAFKA_TOPIC='orders.v1' KAFKA_PARTITIONS=3 \
go run ./cmd/admin
```

### Produce events

```bash
cd go/kafka
KAFKA_BROKER='localhost:9092' KAFKA_TOPIC='orders.v1' PRODUCE_COUNT=30 \
go run ./cmd/producer
```

### Consume events

```bash
cd go/kafka
KAFKA_BROKER='localhost:9092' KAFKA_TOPIC='orders.v1' \
KAFKA_GROUP_ID='order-processors' CONSUMER_NAME='consumer-1' \
go run ./cmd/consumer
```

## Exercises

These are the parts worth actually running, not just reading about:

1. **Per-key ordering**: run the producer once, then the consumer. Notice
   that all `cust-1` events arrive in creation order relative to each
   other, even though `cust-2`/`cust-3` events are interleaved with them
   across partitions.

2. **Consumer group load-splitting**: with the topic at 3 partitions,
   open two more terminals and run the consumer again with the *same*
   `KAFKA_GROUP_ID` but different `CONSUMER_NAME` (`consumer-2`,
   `consumer-3`). Run the producer again. Each consumer only sees a subset
   of partitions -- the group divided the work. Kill one consumer mid-run
   and watch Kafka UI's "Consumer Groups" tab show a rebalance as its
   partitions get reassigned to the survivors.

3. **Broadcast vs. queue**: run two consumers with *different*
   `KAFKA_GROUP_ID` values and produce again. Both consumers now see
   *every* message -- independent groups each get their own copy of the
   log.

4. **More partitions than consumers vs. fewer**: recreate the topic with
   `KAFKA_PARTITIONS=1` and re-run the multi-consumer exercise. Only one
   consumer in the group ever receives messages -- you can't have more
   active consumers *doing work* than partitions in a group, extra
   consumers just sit idle. This is why partition count is your ceiling
   on consumer-group parallelism.

5. **At-least-once replay**: start a consumer, let it process a few
   messages, then `Ctrl+C` it *before* it would have committed the next
   offset (e.g. add a `time.Sleep` before the commit in
   [main.go](./cmd/consumer/main.go) to make this easy to hit). Restart
   it with the same group ID -- it resumes from the last committed
   offset and reprocesses anything in flight, demonstrating why consumers
   must be idempotent under at-least-once delivery.

## Why Kafka over a plain queue

- **Replay**: consumers can re-read history (new consumer groups, bug
  fixes, backfills) because reading doesn't delete records.
- **Fan-out**: many independent consumer groups can read the same topic
  without competing for messages.
- **Ordering with scale**: partitioning gives you both parallelism *and*
  ordering guarantees, as long as you pick keys carefully.
- **Throughput**: Kafka is built around sequential disk I/O and batching,
  so it sustains very high write/read throughput compared to most
  traditional brokers.

The tradeoffs: no per-message priority, no easy "peek and skip", and
consumers -- not the broker -- are responsible for tracking what they've
already processed.
