package main

import (
	"fmt"
	"github.com/farukterzioglu/ChronicleRMQ/consumerserver"
)

func main() {
	fmt.Println("Chronicle consumer is listening...")

	server := consumerserver.NewConsumerServer(consumerserver.Options{
		AckEvery:    100,
		Port:        "8800",
		Async:       false,
		Interactive: false,
	})

	server.Start()
}
