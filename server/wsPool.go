package server

import (
	"fmt"
	"github.com/SArtemJ/wstest/messages"
	"golang.org/x/net/websocket"
)

type WsPool struct {
	Clients        map[string]*websocket.Conn
	NewClients     chan *websocket.Conn
	RemoveClients  chan *websocket.Conn
	StreamMessages chan messages.Message
}

func NewWsPool() *WsPool {
	clients := make(map[string]*websocket.Conn)
	news := make(chan *websocket.Conn)
	rm := make(chan *websocket.Conn)
	stream := make(chan messages.Message)

	return &WsPool{
		Clients:        clients,
		NewClients:     news,
		RemoveClients:  rm,
		StreamMessages: stream,
	}
}

func (wsp *WsPool) Start() {
	for {
		select {
		case newClient := <-wsp.NewClients:
			wsp.ConnectClient(newClient)
		case exitClient := <-wsp.RemoveClients:
			wsp.DisconnectClient(exitClient)
		case newMessage := <-wsp.StreamMessages:
			wsp.StreamMsg(newMessage)
		}
	}
}

func (wsp *WsPool) ConnectClient(newCl *websocket.Conn) {
	wsp.Clients[newCl.RemoteAddr().String()] = newCl
}

func (wsp *WsPool) DisconnectClient(client *websocket.Conn) {
	delete(wsp.Clients, client.RemoteAddr().String())
}

func (wsp *WsPool) StreamMsg(msg messages.Message) {
	var removeFromSenders []string
	for _, client := range wsp.Clients {
		if isSender := messages.SendersPool.Load(client.RemoteAddr().String()); !isSender {
			err := websocket.JSON.Send(client, msg)
			if err != nil {
				fmt.Println("Error broadcasting message: ", err)
				return
			}
		} else {
			removeFromSenders = append(removeFromSenders, client.RemoteAddr().String())
		}
	}
	for _, item := range removeFromSenders {
		messages.SendersPool.Store(item, false)
	}
}
