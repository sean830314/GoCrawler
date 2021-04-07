
# GoCrawler
This service is a crawler for social platform data.
* * *
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
## Build images
### Build go-crawler-service images
```
GOCRAWLER_IMAGE_TAGS=latest make build-gocrawler-service-image
```

### Build ptt-crawler-consumer images
```
CONSUMER_IMAGE_TAGS=latest make build-consumer-image  
```
## Deploy service with docker-compose

### Add domain name to /etc/hosts
```
echo "127.0.0.1 gocrawler2.microservice.com\n127.0.0.1 gocrawler1.microservice.com" >> /etc/hosts
cd deploy
```
* * *
### Deploy traefik as reverse proxy
```
docker create network traefik_net
docker-compose -f traefik-docker-compose.yml up -d
```
### Deploy app microservices with docker-compose
```
docker-compose up -d
```
### Test 
```
http://gocrawler1.microservice.com/ping
http://gocrawler2.microservice.com/ping
http://gocrawler1.microservice.com/swagger/index.html
http://gocrawler2.microservice.com/swagger/index.html
```
## Reference
* [traefik docs](https://doc.traefik.io/traefik/)
* [traefik github](https://github.com/traefik/traefik)
* [traefik 2.0 Example](https://github.com/DoTheEvo/Traefik-v2-examples)
* [cassandra  docker](https://hub.docker.com/_/cassandra)
* [cassandra cqlsh tutorial](https://www.tutorialspoint.com/cassandra/cassandra_cqlsh.htm)
* [rabbitmq docker](https://hub.docker.com/_/rabbitmq)
* [rabbitmq getstarted](https://www.rabbitmq.com/getstarted.html)
* [mongodb docker](https://hub.docker.com/_/mongo)
* [mongodb tutorial](https://www.tutorialspoint.com/mongodb/index.htm)
* [fluent example](https://github.com/sean830314/service-tool-note/tree/master/fluentd)
* [fluent-logger-golang github](https://github.com/fluent/fluent-logger-golang)
