
package main



// ---------------------------------------------------------------------------------------------------------------------

import "github.com/mediocregopher/radix.v2/redis"
import "log"



// ---------------------------------------------------------------------------------------------------------------------

var client *redis.Client



// ---------------------------------------------------------------------------------------------------------------------

func getAllKeys(uuid string) map[string]string {

  allKeys, _ := client.Cmd("HGETALL", uuid).Map()

  return allKeys
}



// ---------------------------------------------------------------------------------------------------------------------

func createHashKey(operationReply OperationReply) {

  client.Cmd("HSET", operationReply.Uuid, operationReply.IpAddress, operationReply.State)
}



// ---------------------------------------------------------------------------------------------------------------------

func connectToRedis(redisAddress string) {

  var err error

  client, err = redis.Dial("tcp", redisAddress)

  if err != nil {
  	log.Fatal(err)
  }
}



// ---------------------------------------------------------------------------------------------------------------------

func disconnectRedis() {

  client.Close()
}
