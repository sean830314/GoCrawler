package queue

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

type RabbitmqConfig struct {
	Host     string
	Port     int
	Account  string
	Password string
}

func (rc RabbitmqConfig) Producing(message []byte) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", rc.Account, rc.Password, rc.Host, rc.Port))
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
	body := message
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
	if err != nil {
		logrus.Error(fmt.Sprintf("Failed to publish a message: %v", err))
	}
	logrus.Info(fmt.Sprintf(" [x] Sent %s", body))
}
