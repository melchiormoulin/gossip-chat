package chat

import (
	"golang.org/x/net/websocket"
	"log"
)

type Chat struct {
	clients          []*websocket.Conn
	addClientChan    chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	Messages         *chan string
}

func InitChat(messages *chan string)*Chat {
	return &Chat{
		clients:          make([]*websocket.Conn,0),
		addClientChan:    make(chan *websocket.Conn),
		removeClientChan: make(chan *websocket.Conn),
		Messages:         messages,
	}
}

func (chat *Chat) Loop(ws *websocket.Conn) {
	chat.addClient(ws)
	go chat.run()
	go chat.receive(ws)
	chat.send()
}

func (chat *Chat) receive(ws *websocket.Conn) {
	for {
		var message string
		if err := websocket.Message.Receive(ws, &message); err != nil {
			log.Println("Can't receive from ws", err)
			chat.removeClient(ws)
			break
		}
		log.Println("received from ws :" + message)
		broadcasts.QueueBroadcast(simpleBroadcast(message))
		*chat.Messages <- message

	}
}

func (chat *Chat) send() {
	for {
		select {
		case messageReceived := <-*chat.Messages:
			for _, ws := range chat.clients {
				log.Println("sending to ws :" + messageReceived)
				if err := websocket.Message.Send(ws, messageReceived); err != nil {
					log.Println("Can't send to ws", err)
					break
				}
				log.Println("sent to ws :" + messageReceived)
			}
		}
	}
}

// run receives from the hub channels and calls the appropriate hub method
func (chat *Chat) run() {
	for {
		select {
		case conn := <-chat.addClientChan:
			chat.addClient(conn)
		case conn := <-chat.removeClientChan:
			chat.removeClient(conn)
		}
	}
}

// removeClient removes a conn from the pool
func (chat *Chat) removeClient(conn *websocket.Conn) {
	for i, tmpCon := range chat.clients {
		if tmpCon ==conn {
			chat.clients[i] = chat.clients [len(chat.clients )-1]
			chat.clients = chat.clients[:len(chat.clients)-1]
		}
	}
}

// addClient adds a conn to the pool
func (chat *Chat) addClient(conn *websocket.Conn) {
	chat.clients= append(chat.clients, conn)
}
