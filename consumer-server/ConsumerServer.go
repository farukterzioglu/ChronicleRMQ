package consumerserver

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var typeMap = map[int16]string{
	1001: "fork",
	1002: "block",
	1003: "tx",
	1004: "abi",
	1006: "abiError",
	1005: "abiRemoved",
	1007: "tableRow",
	1008: "encoderError",
	1009: "pause",
	1010: "blockCompleted",
	1011: "permission",
	1012: "permissionLink",
	1013: "accMetadata",
}

const (
	// Time to wait before force close on connection.
	closeGracePeriod = 2 * time.Second
)

type Options struct {
	AckEvery    int8
	Port        string
	Host        string
	Async       bool
	Interactive bool
}

type IConsumerServer interface {
	Start()
	Stop()
}

type ConsumerServer struct {
	confirmedBlock      int64
	unconfirmedBlock    int64
	ackEvery            int8
	wsPort              string
	wsHost              string
	async               bool
	interactive         bool
	chronicleConnection *websocket.Conn
}

var _ IConsumerServer = ConsumerServer{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewConsumerServer(opts Options) ConsumerServer {
	consumerServer := ConsumerServer{
		wsPort:           opts.Port,
		ackEvery:         opts.AckEvery,
		interactive:      opts.Interactive,
		confirmedBlock:   0,
		unconfirmedBlock: 0,
	}

	consumerServer.wsHost = "0.0.0.0"
	if len(opts.Host) > 0 {
		consumerServer.wsHost = opts.Host
	}

	consumerServer.ackEvery = 100
	if opts.AckEvery > 0 {
		consumerServer.ackEvery = opts.AckEvery
	}

	// TODO: Handle async mode
	consumerServer.async = opts.Async

	return consumerServer
}

func (s ConsumerServer) Start() {
	fmt.Printf("Starting Chronicle consumer on %s:%s\n", s.wsHost, s.wsPort)
	fmt.Printf("Acknowledging every %d blocks\n", s.ackEvery)

	http.HandleFunc("/", s.wsEndpoint)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", s.wsHost, s.wsPort), nil))
}

var kConsumerServerClientConnected = false

func (s ConsumerServer) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer socket.Close()

	if kConsumerServerClientConnected {
		fmt.Println("Rejected a new Chronicle connection because one is active already")
		socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(400, ""))
		time.Sleep(closeGracePeriod)
		socket.Close()
		return
	}

	kConsumerServerClientConnected = true
	s.chronicleConnection = socket
	// TODO: Emit connected message

	socket.SetCloseHandler(func(code int, text string) error {
		fmt.Println("Chronicle connection is closed from remote")
		kConsumerServerClientConnected = false
		// TODO: Emit closed message

		var message []byte
		if code != websocket.CloseNoStatusReceived {
			message = websocket.FormatCloseMessage(code, "")
		}
		socket.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	reader(socket)
}

var latestBlock int32 = 0

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// TODO: Get & check message type

		// TODO: Emit message
		fmt.Println(string(p))

		// TODO: Ack block

		latestBlock += 1
		if err := conn.WriteMessage(messageType, []byte(fmt.Sprintf("%d", latestBlock))); err != nil {
			log.Println(err)
			return
		}
	}
}

func (s ConsumerServer) Stop() {
	kConsumerServerClientConnected = false
	s.chronicleConnection.Close()
}
