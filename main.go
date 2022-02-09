package main

//import --> const --> var --> init()

// curl -X POST -d "{\"MessageType\": \"Alert\", \"CameraId\": 42 }" http://localhost:9000/json
// curl -X POST -d "{\"MessageType\": \"Alert\", \"CameraId\": 42, \"FileName\": \"docs.tar.gz\" }" http://localhost:9000/json

// Multithread test:
// curl -X POST -d "{\"MessageType\": \"Alert1\"}" http://localhost:9000/json | curl -X POST -d "{\"MessageType\": \"Alert2\"}" http://localhost:9000/json | curl -X POST -d "{\"MessageType\": \"Alert3\"}" http://localhost:9000/json  |  curl -X POST -d "{\"MessageType\": \"Alert4\"}" http://localhost:9000/json |  curl -X POST -d "{\"MessageType\": \"Alert5\"}" http://localhost:9000/json

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"webserv/untargz"
	//"github.com/zagoranov/webserv/untargz"
)

var settings struct {
	ServerMode bool   `json:"serverMode"`
	SourceDir  string `json:"sourceDir"`
	TargetDir  string `json:"targetDir"`
}

func readConfig() {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Println("opening config file", err.Error())
		panic(err)
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&settings); err != nil {
		log.Println("parsing config file", err.Error())
		panic(err)
	}
	//fmt.Printf("%v %s %s", settings.ServerMode, settings.SourceDir, settings.TargetDir)
	defer configFile.Close()
}

type incomingPacket struct {
	CameraId    int
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
	packet := cacheList.Back()
	if packet != nil {
		log.Println("Proceeding message: '", packet.Value.(incomingPacket).MessageType, "' , from camera:", packet.Value.(incomingPacket).CameraId, "List size:", cacheList.Len())

		//time.Sleep(time.Millisecond * 300)    //test!

		if packet.Value.(incomingPacket).FileName != "" {
			log.Println("Extracting file:", packet.Value.(incomingPacket).FileName)
			r, err := os.Open(packet.Value.(incomingPacket).FileName)
			if err != nil {
				log.Println("Error opening file", err)
			} else {
				untargz.ExtractTarGz(r)
			}
			defer r.Close()
		}

		cacheList.Remove(packet)
	}
	defer mu.Unlock()
}

func jsonHandler(rw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var packet incomingPacket
	err := decoder.Decode(&packet)
	if err != nil {
		log.Println("jsonHandler: ", err)
		fmt.Fprintf(rw, "Bad data received from you, no comprendo")
		//panic(err)
	} else {
		cacheList.PushBack(packet)
		fmt.Fprintf(rw, "Ok")
		go writeData()
	}
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
	readConfig()
	log.Printf("Server started")
	http.HandleFunc("/", HomeRouterHandler)
	http.HandleFunc("/json", jsonHandler)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
