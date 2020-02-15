module github.com/farukterzioglu/ChronicleRMQ/chronicle-consumer-rmq

go 1.13

require (
	github.com/farukterzioglu/ChronicleRMQ/consumerserver v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.4.1
)

replace github.com/farukterzioglu/ChronicleRMQ/consumerserver => ../consumer-server
