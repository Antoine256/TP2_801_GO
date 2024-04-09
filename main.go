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

	wg.Add(1)

	go log.Fatal(http.ListenAndServe(":8081", nil))

	wg.Wait()
}
