package main

import (
	"fmt"
	"github.com/farukterzioglu/ChronicleRMQ/consumerserver"
)

func main() {
	fmt.Println("Chronicle consumer is listening...")

	var server consumerserver.IConsumerServer
	server = consumerserver.NewConsumerServer(consumerserver.Options{
		AckEvery:    100,
		Port:        "8800",
		Async:       false,
		Interactive: false,
	})

	connChn := make(chan interface{})
	server.AddHandler("connected", connChn)

	go func() {
		select {
		case evt := <-connChn:
			if evt == nil {
				return
			}

			// TODO: Process event
			fmt.Println("Connected...")
		}
	}()

	server.Start()
}
