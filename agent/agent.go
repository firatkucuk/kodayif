
package main



// ---------------------------------------------------------------------------------------------------------------------

import "bytes"
import "encoding/json"
import "fmt"
import "github.com/streadway/amqp"
import "io/ioutil"
import "log"
import "net"
import "net/http"
import "os"
import "strings"



// ---------------------------------------------------------------------------------------------------------------------

const CONFIGURATION_FILE string = "config.json"
var   ipAddress          string



// ---------------------------------------------------------------------------------------------------------------------

func failOnError(err error, msg string) {

  if err != nil {
    log.Fatalf("%s: %s", msg, err)
    panic(fmt.Sprintf("%s: %s", msg, err))
  }
}



// ---------------------------------------------------------------------------------------------------------------------

type Configuration struct {
  Controller       string
  MqConnString     string
}



// ---------------------------------------------------------------------------------------------------------------------

type OperationRequest struct {
  Operation        uint8
  FilePath         string
  Term             string
}



// ---------------------------------------------------------------------------------------------------------------------

type OperationMessage struct {
  Request          OperationRequest
  Uuid             string
}



// ---------------------------------------------------------------------------------------------------------------------

type OperationReply struct {
  State            bool
  Uuid             string
  IpAddress        string
}



// ---------------------------------------------------------------------------------------------------------------------

func reply(controller string, uuid string, state bool) {

  operationReply    := OperationReply {
    State     : state,
    Uuid      : uuid,
    IpAddress : ipAddress,
  }
  messageBytes, err := json.Marshal(operationReply)

  if err != nil {
    log.Println("Seialization error")
    return
  }

  resp, err := http.Post("http://" + controller + "/reply", "application/json", bytes.NewBuffer(messageBytes))

  if err != nil {
    log.Println("Cannot send reply")
  }

  if (resp.StatusCode != 200) {
    log.Println("Reply send failed.")
  }
}



// ---------------------------------------------------------------------------------------------------------------------

func processMessage(controller string, message []byte) {

  operationMessage := OperationMessage{}
  err              := json.Unmarshal(message, &operationMessage)

  if err != nil {
    log.Println("JSON Decoding error")
  }

  operationRequest := operationMessage.Request
  uuid             := operationMessage.Uuid

  if operationRequest.Operation == 1 {
    if _, err := os.Stat(operationRequest.FilePath); err == nil {
      reply(controller, uuid, true)
    } else {
      reply(controller, uuid, false)
    }
  } else if (operationRequest.Operation == 2) {
    if _, err := os.Stat(operationRequest.FilePath); err == nil {
      content, err := ioutil.ReadFile(operationRequest.FilePath)

      if err != nil {
        log.Println("Cannot read file")
        reply(controller, uuid, false)
      }

      if strings.Contains(string(content), operationRequest.Term) {
        reply(controller, uuid, true)
      } else {
        reply(controller, uuid, false)
      }
    } else {
      reply(controller, uuid, false)
    }
  }
}



// ---------------------------------------------------------------------------------------------------------------------

func connectToMq(controller string, mqConnString string) {

  connection, err := amqp.Dial(mqConnString)
  failOnError(err, "Failed to connect to RabbitMQ")
  defer connection.Close()

  channel, err := connection.Channel()
  failOnError(err, "Failed to open a channel")
  defer channel.Close()

  err = channel.ExchangeDeclare(
    "kodayif", // name
    "fanout",  // type
    true ,     // durable
    false,     // auto-deleted
    false,     // internal
    false,     // no-wait
    nil,       // arguments
  )
  failOnError(err, "Failed to declare an exchange")

  queue, err := channel.QueueDeclare(
    "",    // name
    false, // durable
    false, // delete when usused
    true,  // exclusive
    false, // no-wait
    nil,   // arguments
  )
  failOnError(err, "Failed to declare a queue")

  err = channel.QueueBind(
    queue.Name, // queue name
    "",         // routing key
    "kodayif",  // exchange
    false,
    nil,
  )
  failOnError(err, "Failed to bind a queue")

  msgs, err := channel.Consume(
    queue.Name, // queue
    "",         // consumer
    true,       // auto-ack
    false,      // exclusive
    false,      // no-local
    false,      // no-wait
    nil,        // args
  )
  failOnError(err, "Failed to register a consumer")

  forever := make(chan bool)

  go func() {
    for msg := range msgs {
      processMessage(controller, msg.Body)
    }
  }()

  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
  <-forever
}



// ---------------------------------------------------------------------------------------------------------------------

func parseConfiguration() Configuration {

  file, err := os.Open(CONFIGURATION_FILE)

  if err != nil {
    log.Fatal("Couldn't read config file!")
  }

  decoder       := json.NewDecoder(file)
  configuration := Configuration{}
  err            = decoder.Decode(&configuration)

  if err != nil {
    log.Fatal("Syntax error in configuration file.")
  }

  return configuration
}



// ----------------------------------------------------------------------------------------------------------------xx-

func ipLookUp() string {

  interfaces, err := net.Interfaces()

  if err != nil {
    log.Println("Couln't fetch interface list")
  }

  for _, i := range interfaces {
    addrs, _ := i.Addrs()

    for _, addr := range addrs {
      address := addr.String()

      if !strings.HasPrefix(address, "127") && !strings.HasPrefix(address, ":") {
        return strings.Split(address, "/")[0]
      }
    }
  }

  return ""
}



// ----------------------------------------------------------------------------------------------------------------xx-

func main() {

  ipAddress = ipLookUp()

  if (ipAddress == "") {
    log.Fatal("Couldn't find a valid interface address")
  }

  configuration := parseConfiguration()
  connectToMq(configuration.Controller, configuration.MqConnString)
}
