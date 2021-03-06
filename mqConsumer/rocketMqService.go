package mqConsumer

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
)

//   还是写死了      它只能消费创建订单的信息  这个应该将其弄成common的吧

//一个common的消费者
//	queueName :="orderMessage"
func ReceiveMessageNormalConsumer(queueName string) string {

	//TODO  要将这个改成一个函数吧

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		errResult := fmt.Sprintf("connect to the rabbitMq failed %s", err)
		fmt.Println(errResult)
		err = errors.New(errResult)
		//return err
	}

	//failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		errResult := fmt.Sprintf("Failed to open a channel %s", err)
		fmt.Println(errResult)

		err = errors.New(errResult)
		//return err
	}

	//failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		errResult := fmt.Sprintf("Failed to declare a queue %s", err)
		fmt.Println(errResult)

		err = errors.New(errResult)
		//return err
	}
	//failOnError(err, "Failed to declare a queue")

	//Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args Table) (<-chan Delivery, error)
	// When we return fro

	msgs, err := ch.Consume(
		q.Name, // queue

		"",    // consumer true,   // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		false,
		nil,
	)
	d := <-msgs

	return string(d.Body)

}




func MqConsumerCommon(queueName string) interface{} {
	//TODO  要将这个改成一个函数吧
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		errResult := fmt.Sprintf("connect to the rabbitMq failed %s", err)
		fmt.Println(errResult)
		err = errors.New(errResult)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		errResult := fmt.Sprintf("Failed to open a channel %s", err)
		fmt.Println(errResult)
		err = errors.New(errResult)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		errResult := fmt.Sprintf("Failed to declare a queue %s", err)
		fmt.Println(errResult)

		err = errors.New(errResult)
		//return err
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return msgs
}

