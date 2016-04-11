
package main



// ---------------------------------------------------------------------------------------------------------------------

import "encoding/json"
import "errors"
import "github.com/nu7hatch/gouuid"
import "log"
import "net/http"



// ---------------------------------------------------------------------------------------------------------------------

type OperationRequest struct {
  Operation        uint8
  FilePath         string
  Term             string
}



// ---------------------------------------------------------------------------------------------------------------------

type OperationResponse struct {
  Uuid             string
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
}



// ---------------------------------------------------------------------------------------------------------------------

func validate(request *http.Request) (OperationRequest, error) {

  decoder          := json.NewDecoder(request.Body)
  operationRequest := OperationRequest{}
  err              := decoder.Decode(&operationRequest)

  if err != nil {
    log.Println("JSON Decoding error")

    return operationRequest, err
  }

  if operationRequest.Operation == 1 {
    if operationRequest.FilePath == "" {
      return operationRequest, errors.New("filePath cannot be empty")
    }
  } else if operationRequest.Operation == 2 {
    if operationRequest.FilePath == "" || operationRequest.Term == "" {
      return operationRequest, errors.New("filePath and term cannot be empty")
    }
  } else {
    return operationRequest, errors.New("Invalid operation number")
  }

  return operationRequest, nil
}



// ---------------------------------------------------------------------------------------------------------------------

func sendActionHandler(responseWriter http.ResponseWriter, request *http.Request) {

  operationRequest, err := validate(request)

  if err != nil {
    responseWriter.WriteHeader(http.StatusBadRequest)
    return
  }

  opUuid, err      := uuid.NewV4()
  operationMessage := OperationMessage {
    Request : operationRequest,
    Uuid    : opUuid.String(),
  }

  err = sendToMq(operationMessage)

  if err != nil {
    log.Println(err)
    responseWriter.WriteHeader(http.StatusInternalServerError)
    return
  }

  header := responseWriter.Header()
  header.Add("Content-Type", "application/json")

  encoder           := json.NewEncoder(responseWriter)
  operationResponse := OperationResponse{Uuid : operationMessage.Uuid}
  encoder.Encode(operationResponse)
}



// ---------------------------------------------------------------------------------------------------------------------

func replyActionHandler(responseWriter http.ResponseWriter, request *http.Request) {

  decoder        := json.NewDecoder(request.Body)
  operationReply := OperationReply{}
  err            := decoder.Decode(&operationReply)

  if err != nil {
    responseWriter.WriteHeader(http.StatusBadRequest)
    return
  }

  log.Println(operationReply)

  header := responseWriter.Header()
  header.Add("Content-Type", "application/json")
}



// ---------------------------------------------------------------------------------------------------------------------

func startWebServer(listenAddress string) {

  var fileServer = http.FileServer(http.Dir("static"))
  http.Handle("/", fileServer)

  http.HandleFunc("/send", sendActionHandler)
  http.HandleFunc("/reply", replyActionHandler)
  http.ListenAndServe(listenAddress, nil)
}
