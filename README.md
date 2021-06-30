
# GoCrawler
This service is a crawler for social platform data.
* * *
## Run service locally

### Deploy infra services
```
cd deploy && docker-compose -f infra.yml
```
### Initialize postgresql databases
create admin databases
```
docker exec -it postgres bash
psql -U postgres
CREATE DATABASE backend_admin;
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
```
* * *
### Deploy traefik as reverse proxy
```
cd deploy
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
* [jaeger & gin example](https://www.yisu.com/zixun/372364.html)
* [opentracing-jaeger](https://jeremyxu2010.github.io/2018/07/%E7%A0%94%E7%A9%B6%E8%B0%83%E7%94%A8%E9%93%BE%E8%B7%9F%E8%B8%AA%E6%8A%80%E6%9C%AF%E4%B9%8Bjaeger/)
* [gin+gorm+router 快速搭建 crud restful API 接口](https://learnku.com/articles/23548/gingormrouter-quickly-build-crud-restful-api-interface)
* [go-gin-example](https://github.com/eddycjy/go-gin-example)
