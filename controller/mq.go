
package main



// ---------------------------------------------------------------------------------------------------------------------

import "encoding/json"
import "fmt"
import "github.com/streadway/amqp"
import "log"



// ---------------------------------------------------------------------------------------------------------------------

var connection *amqp.Connection
var channel    *amqp.Channel
var queue       amqp.Queue



// ---------------------------------------------------------------------------------------------------------------------

func failOnError(err error, msg string) {

  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}



// ---------------------------------------------------------------------------------------------------------------------

func sendToMq(operationMessage OperationMessage) error {

  messageBytes, err := json.Marshal(operationMessage)

  if err != nil {
    log.Println("Seialization error")
    return err
  }

  err = channel.Publish(
    "",         // exchange
    queue.Name, // routing key
    false,      // mandatory
    false,      // immediate
    amqp.Publishing {
      ContentType : "application/json",
      Body        : messageBytes,
    })

  return err
}



// ---------------------------------------------------------------------------------------------------------------------

func connectToMq(connString string) {

  connection, err := amqp.Dial(connString)

  if err != nil {
    defer connection.Close()
    log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
  }

  channel, err = connection.Channel()

  if err != nil {
    defer channel.Close()
    defer connection.Close()
    log.Fatalf("%s: %s", "Failed to open a channel", err)
  }

  queue, err = channel.QueueDeclare(
    "kodayif", // name
    false,     // durable
    false,     // delete when unused
    false,     // exclusive
    false,     // no-wait
    nil,       // arguments
  )

  if err != nil {
    defer channel.Close()
    defer connection.Close()
    log.Fatalf("%s: %s", "Failed to declare a queue", err)
  }
}



// ---------------------------------------------------------------------------------------------------------------------

func disconnectMq() {

  connection.Close()
  channel.Close()
}
