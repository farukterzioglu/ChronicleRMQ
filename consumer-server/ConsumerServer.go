package consumerserver

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type EventType string

const (
	BLOCK        EventType = "block"
	CONNECTED    EventType = "connected"
	DISCONNECTED EventType = "disconnected"
)

type MessageType int16

var typeMap = map[MessageType]EventType{
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
	AddHandler(t EventType) chan interface{}
	RemoveHandler(t EventType, ch chan interface{})
}

type consumerServer struct {
	confirmedBlock      int64
	unconfirmedBlock    int64
	ackEvery            int8
	wsPort              string
	wsHost              string
	async               bool
	interactive         bool
	chronicleConnection *websocket.Conn
	messageHandlers     map[EventType][]chan interface{}
}

var _ IConsumerServer = &consumerServer{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewConsumerServer(opts Options) IConsumerServer {
	consumerServer := consumerServer{
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

	return &consumerServer
}

func (s consumerServer) Start() {
	fmt.Printf("Starting Chronicle consumer on %s:%s\n", s.wsHost, s.wsPort)
	fmt.Printf("Acknowledging every %d blocks\n", s.ackEvery)

	http.HandleFunc("/", s.wsEndpoint)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", s.wsHost, s.wsPort), nil))
}

var kConsumerServerClientConnected = false

func (s *consumerServer) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket connection
	// TODO : remove default websocket logs
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
	s.emit(CONNECTED, struct{ remoteAddress interface{} }{remoteAddress: socket.RemoteAddr()})

	socket.SetCloseHandler(func(code int, text string) error {
		kConsumerServerClientConnected = false
		s.emit(DISCONNECTED, struct{ remoteAddress interface{} }{remoteAddress: socket.RemoteAddr()})

		var message []byte
		if code != websocket.CloseNoStatusReceived {
			message = websocket.FormatCloseMessage(code, "")
		}
		socket.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	s.startReader(socket)
}

var latestBlock int32 = 0

func (s *consumerServer) startReader(conn *websocket.Conn) {
	for {
		// read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if err := s.processMessage(p, conn); err != nil {
			fmt.Print(err)
		}
	}
}

func (s *consumerServer) processMessage(p []byte, conn *websocket.Conn) error {
	// TODO: Get message type
	var msgType MessageType = 1002
	eventType, ok := typeMap[msgType]
	if !ok {
		return errors.New("Unknown msgType: " + fmt.Sprint(msgType))
	}

	s.emit(eventType, p)

	// TODO: Ack block

	latestBlock += 1
	if err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", latestBlock))); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *consumerServer) Stop() {
	kConsumerServerClientConnected = false
	s.chronicleConnection.Close()

	// TODO: Close handler channels
}

func (s *consumerServer) AddHandler(t EventType) chan interface{} {
	ch := make(chan interface{})
	if s.messageHandlers == nil {
		s.messageHandlers = make(map[EventType][]chan interface{})
	}
	if _, ok := s.messageHandlers[t]; ok {
		s.messageHandlers[t] = append(s.messageHandlers[t], ch)
	} else {
		s.messageHandlers[t] = []chan interface{}{ch}
	}
	return ch
}

// RemoveSitter removes an event listener from the Dog struct instance
func (s *consumerServer) RemoveHandler(t EventType, ch chan interface{}) {
	if _, ok := s.messageHandlers[t]; ok {
		for i := range s.messageHandlers[t] {
			if s.messageHandlers[t][i] == ch {
				s.messageHandlers[t] = append(s.messageHandlers[t][:i], s.messageHandlers[t][i+1:]...)
				break
			}
		}
	}
}

func (s *consumerServer) emit(t EventType, response interface{}) {
	if _, ok := s.messageHandlers[t]; ok {
		for _, handler := range s.messageHandlers[t] {
			go func(handler chan interface{}) {
				handler <- response
			}(handler)
		}
	}
}
