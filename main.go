package main

//import --> const --> var --> init()

// curl -X POST -d "{\"MessageType\": \"Alert\"}" http://localhost:9000/json

// Multithread test:
// curl -X POST -d "{\"MessageType\": \"Alert1\"}" http://localhost:9000/json | curl -X POST -d "{\"MessageType\": \"Alert2\"}" http://localhost:9000/json | curl -X POST -d "{\"MessageType\": \"Alert3\"}" http://localhost:9000/json  |  curl -X POST -d "{\"MessageType\": \"Alert4\"}" http://localhost:9000/json |  curl -X POST -d "{\"MessageType\": \"Alert5\"}" http://localhost:9000/json

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type incomingPacket struct {
	DateTime    string
	MessageType string
	MessageText string
	Value       float64
	FileName    string
}

var (
	cacheList list.List
	mu        sync.Mutex
)

func writeData() {
	mu.Lock()
	log.Println("Last message:", cacheList.Back().Value.(incomingPacket).MessageType, ", List size:", cacheList.Len())
	//time.Sleep(time.Millisecond * 300)
	defer mu.Unlock()
}

func jsonHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var packet incomingPacket
	err := decoder.Decode(&packet)
	if err != nil {
		log.Fatal("jsonResp: ", err)
		//panic(err)
	}
	cacheList.PushBack(packet)
	go writeData()
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	log.Println("path", r.URL.Path)
	log.Println("scheme", r.URL.Scheme)
	log.Println(r.Form["url_long"])
	for k, v := range r.Form {
		log.Println("key:", k)
		log.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Not very hospitable wellcome page")
}

func main() {
	log.Printf("Запуск сервера")
	http.HandleFunc("/", HomeRouterHandler)
	http.HandleFunc("/json", jsonHandler)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
