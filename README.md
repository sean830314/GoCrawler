
# GoCrawler
## Cassandra Command
### Cassandra exec sql
```
docker exec -it some-cassandra bash
root@6405f1f27115:/# cqlsh
```
### Cassandra list keyspace
```
cqlsh> desc keyspaces;
```
### Cassandra use keyspace
```
cqlsh> use ptt_keyspace;
```
### Cassandra list table
```
cqlsh:ptt_keyspace> desc tables;
```
### Cassandra list table schemas
```
cqlsh:ptt_keyspace> desc table article;
```
### Cassandra list table records
```
cqlsh:ptt_keyspace> select * from article;
```
### Cassandra count table records
```
cqlsh:ptt_keyspace> select COUNT(*) from article;
```
## Rabbitmq Command
### Run RabbitMQ Container
```
docker run --name some-rabbitmq -p 5672:5672 -d rabbitmq:latest
```
### RabbitMQ list jobs
```
rabbitmqctl list_queues
```
## Build images
### Build go-crawler-service images
```
GOCRAWLER_IMAGE_TAGS=latest make build-gocrawler-service-image
```

### Build ptt-crawler-consumer images
```
PTT_CONSUMER_IMAGE_TAGS=latest make build-ptt-consumer-image
```
## Deploy service with docker-compose
### deploy service with docker-compose
```
cd deploy
docker-compose up -d
```
