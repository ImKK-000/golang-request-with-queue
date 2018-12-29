package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type apiHandler struct{}

func (apiHandler apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// save request time
	requestTime := time.Now().Format(time.RFC3339)

	// wait for implement
	time.Sleep(time.Second)

	// generate response message to json format
	responseBodyStruct := struct {
		RequestTime  string `json:"request_time"`
		ResponseTime string `json:"response_time"`
	}{
		RequestTime:  requestTime,
		ResponseTime: time.Now().Format(time.RFC3339),
	}
	responseBody, _ := json.MarshalIndent(responseBodyStruct, "", "  ")

	// send response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

const serveURL = "http://127.0.0.1:9999"

func main() {
	// fake api
	go http.ListenAndServe(":9999", &apiHandler{})

	rawResponse, err := http.Get(serveURL)
	if err != nil {
		log.Fatal(err)
	}
	response, err := ioutil.ReadAll(rawResponse.Body)
	fmt.Println(string(response))
}
