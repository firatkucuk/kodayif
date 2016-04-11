
package main



// ---------------------------------------------------------------------------------------------------------------------

import "encoding/json"
import "log"
import "os"


// ---------------------------------------------------------------------------------------------------------------------

const CONFIGURATION_FILE string = "config.json"


// ---------------------------------------------------------------------------------------------------------------------

type Configuration struct {
  ListenAddress    string
  MqConnString     string
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


// ---------------------------------------------------------------------------------------------------------------------

func main() {

  configuration := parseConfiguration()

  connectToMq(configuration.MqConnString)
  startWebServer(configuration.ListenAddress)
}
