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
	global.InitTupleSpace()
	global.AddDetectorOnAllBatiments()

	http.HandleFunc("/", global.Handler)
	fmt.Println("Server listening on :8081")

	wg.Add(2)

	go global.HandleWrites()
	go log.Fatal(http.ListenAndServe(":8081", nil))

	wg.Wait()
}
