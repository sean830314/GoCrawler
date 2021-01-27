package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to connect to RabbitMQ: %v", err))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to open a channel: %v", err))
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to declare a queue: %v", err))
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to set QoS: %v", err))
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to register a consumer: %v", err))
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			logrus.Info(fmt.Sprintf("Received a message: %v", string(d.Body)))
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			logrus.Info("Done")
			d.Ack(false)
		}
	}()
	logrus.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
