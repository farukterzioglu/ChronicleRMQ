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

	connectedChn := server.AddHandler(consumerserver.CONNECTED)
	disconnectedChn := server.AddHandler(consumerserver.DISCONNECTED)
	blockChn := server.AddHandler(consumerserver.BLOCK)
	blockLoggerChn := server.AddHandler(consumerserver.BLOCK)

	go func() {
		for {
			select {
			case <-connectedChn:
				// TODO: Process event
				log.Printf("Connection established with Chronicle...")
			case <-disconnectedChn:
				// TODO: Process event
				log.Printf("Chronicle connection is closed from remote.")
			case block := <-blockChn:
				log.Printf("Processing block %s...", block.([]byte))
			case p := <-blockLoggerChn:
				log.Printf("Got new block %s...", p.([]byte))
			case tx := <-server.AddHandler(consumerserver.TX):
				log.Printf("Got new tx %s...", tx.([]byte))
			}
		}
	}()

	server.Start()
}
