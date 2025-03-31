package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	url := "http://127.0.0.1:8080/"

	for i := 1; i <= 50; i++ {
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Request %2d: Status %d\n", i, resp.StatusCode)
		time.Sleep(100 * time.Millisecond)
	}
}
