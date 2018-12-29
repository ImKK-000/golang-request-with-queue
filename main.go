package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type fakeAPIHandler struct{}

func (fakeAPIHandler fakeAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// save request time
	requestTime := time.Now().Format(time.RFC3339Nano)

	// wait for implement
	time.Sleep(time.Second)

	// generate response message to json format
	responseBodyStruct := struct {
		RequestTime  string `json:"request_time"`
		ResponseTime string `json:"response_time"`
	}{
		RequestTime:  requestTime,
		ResponseTime: time.Now().Format(time.RFC3339Nano),
	}
	responseBody, _ := json.MarshalIndent(responseBodyStruct, "", "  ")

	// send response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

// use simple queue with slice
var queue []string
var chanQueue = make(chan string)

type apiHandler struct{}

func (apiHandler apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(""))
		return
	}

	// send response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))

	rawRequestBody, _ := ioutil.ReadAll(r.Body)
	var requestBody struct {
		URL string `json:"url"`
	}
	json.Unmarshal(rawRequestBody, &requestBody)

	chanQueue <- requestBody.URL
}

func requestAPI(url string) string {
	rawResponse, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	response, _ := ioutil.ReadAll(rawResponse.Body)
	return string(response)
}

func main() {
	// fake api
	go http.ListenAndServe(":9999", &fakeAPIHandler{})
	go http.ListenAndServe(":9980", &apiHandler{})

	for {
		preQueue := <-chanQueue
		queue = append(queue, preQueue)

		// trigger queue to send request
		if len(queue) == 5 {
			queueBuffer := queue[:5]
			for _, qBuffer := range queueBuffer {
				fmt.Println(requestAPI(qBuffer))
			}
			queue = queue[5:]
		}
	}
}
