package main

import (
	"TP_801_GO/global"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

func main() {
	http.HandleFunc("/", global.Handler)
	fmt.Println("Server listening on :8081")

	//ts := NewSpace("ts")

	wg.Add(1)

	global.SendToConn("Initilisation system !")

	go log.Fatal(http.ListenAndServe(":8081", nil))

	wg.Wait()
}
