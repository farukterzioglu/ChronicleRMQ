## chronicle-consumer: a Golang library for Chronicle data processing & an application for streaming to RabbitMq

The Golang library (./consumer-server) allows an application receive real-time notifications or
historical data from an EOSIO blockchain, such as Telos, EOS, WAX,
Lynxchain, Europechain, and many others.

The Golang application (./chronicle-consumer-rmq) uses this library and streams blockchain data to a RabbitMq instance.  

[Chronicle](https://github.com/EOSChronicleProject/eos-chronicle) is a
software package for receiving and decoding the data flow that is
exported by `state_history_plugin` of `nodeos`, the blockchain node
daemon [developed by Block One](https://developers.eos.io/).

## Installing

`go get github.com/farukterzioglu/ChronicleRMQ/consumerserver`  

## Usage  
(More detailed docs and samples is on the way...)  

#### With docker-compose  
`docker-compose -f docker-compose-test.yml run chronicle-consumer`
`docker build -t chronicle-consumer-rmq -f chronicle-consumer-rmq/Dockerfile .`
`docker run -p 8800:8800 chronicle-consumer-rmq`  

Based on;  
https://github.com/EOSChronicleProject/eos-chronicle  
https://github.com/EOSTribe/eos-chronicle-docker  
https://github.com/eostitan/eosio-chronicle-indexer-template  

Inspired from;  
https://github.com/EOSChronicleProject/chronicle-consumer-npm  