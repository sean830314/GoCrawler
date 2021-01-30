
Run Cassandra Container
```
docker run --name some-cassandra -p 9042:9042 -d cassandra:latest
```
Cassandra exec sql
```
docker exec -it some-cassandra bash
root@6405f1f27115:/# cqlsh
```
Cassandra list keyspace
```
cqlsh> desc keyspaces;
```
Cassandra use keyspace
```
cqlsh> use ptt_keyspace;
```
Cassandra list table
```
cqlsh:ptt_keyspace> desc tables;
```
Cassandra list table schemas
```
cqlsh:ptt_keyspace> desc table article;
```
Cassandra list table records
```
cqlsh:ptt_keyspace> select * from article;
```
Cassandra count table records
```
cqlsh:ptt_keyspace> select COUNT(*) from article;
```
Run RabbitMQ Container
```
docker run --name some-rabbitmq -p 5672:5672 -d rabbitmq:latest
```
RabbitMQ list jobs
```
rabbitmqctl list_queues
```
