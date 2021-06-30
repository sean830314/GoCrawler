# Infra service command
## MongoDB Command
### Mongo exec mongo command
```
docker exec -it mongo bash
root@6405f1f27115:/# mongo
```
### Mongo list databases
```
show databases;
```
### Mongo use database
```
use log_db;
```
### Mongo list tables
```
show tables;
```
### Mongo list table records
```
db.crawler_service.find();
```
* * *
## Cassandra Command
### Cassandra exec sql
```
docker exec -it some-cassandra bash
root@6405f1f27115:/# cqlsh
```
### Cassandra list keyspaces
```
cqlsh> desc keyspaces;
```
### Cassandra use keyspace
```
cqlsh> use social_data;
```
### Cassandra list table
```
cqlsh:social_data> desc tables;
```
### Cassandra list table schemas
```
cqlsh:social_data> desc table ptt_article;
```
### Cassandra list table records
```
cqlsh:social_data> select * from ptt_article;
```
### Cassandra count table records
```
cqlsh:social_data> select COUNT(*) from ptt_article;
```
* * *
## Jaeger
### Deploy Jaeger on local
```
docker run -d --name jaeger  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp  -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one:latest
```
* * *
## Rabbitmq Command
### Run RabbitMQ Container
```
docker run --name some-rabbitmq -p 5672:5672 -d rabbitmq:latest
```
### RabbitMQ list jobs
```
rabbitmqctl list_queues
```
* * *
