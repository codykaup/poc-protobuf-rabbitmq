package main

import (
	"context"
	"log"
	"sync"

	"github.com/codykaup/poc-protobuf-rabbitmq/generated"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5673/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"test-queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-local
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	messageRead := make(chan bool)

	var sync sync.WaitGroup
	sync.Add(1)
	go consumeMessages(ch, q.Name, &sync, messageRead)

	sendMessage(ch, q.Name)
	<-messageRead

	err = ch.Close()
	if err != nil {
		log.Fatalf("Failed to close channel: %v", err)
	}

	sync.Wait()
	log.Printf("Done!")
}

func consumeMessages(ch *amqp.Channel, queueName string, sync *sync.WaitGroup, messageRead chan bool) {
	defer sync.Done()

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for d := range msgs {
		message := generated.Sample{}
		err = proto.Unmarshal(d.Body, &message)
		if err != nil {
			log.Fatalf("Failed to unmarshal message: %v", err)
		}

		log.Printf("<-- Received a message: %s", message.GetMessage())

		messageRead <- true
		close(messageRead)
	}
}

func sendMessage(ch *amqp.Channel, queueName string) {
	message := generated.Sample{
		Id:      1,
		Message: "From protobufs!",
	}

	log.Printf("--> Sending message: %v", message.GetMessage())

	body, err := proto.Marshal(&message)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	err = ch.PublishWithContext(context.Background(), "", queueName, false, false, amqp.Publishing{
		Body: body,
	})
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
}
