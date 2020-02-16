package main

import (
	"fmt"
	"github.com/farukterzioglu/ChronicleRMQ/consumerserver"
	"log"
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

	connectedChn := make(chan interface{})
	disconnectedChn := make(chan interface{})
	blockChn := make(chan interface{})
	server.AddHandler("connected", connectedChn)
	server.AddHandler("disconnected", disconnectedChn)
	server.AddHandler("block", blockChn)

	go func() {
		for {
			select {
			case <-connectedChn:
				// TODO: Process event
				log.Printf("Connection established with Chronicle...")
			case <-disconnectedChn:
				// TODO: Process event
				log.Printf("Chronicle connection is closed from remote.")
			case p := <-blockChn:
				log.Printf("Got new block %s...", p.([]byte))
			}
		}
	}()

	server.Start()
}
