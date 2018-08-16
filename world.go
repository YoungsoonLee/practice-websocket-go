package main

import (
	"fmt"
	"time"
)

type World struct {
	ClinetMap map[*Client]bool
	ChanEnter chan *Client
	ChanLeave chan *Client
	BroadCast chan []byte
}

func newWorld() *World {
	return &World{
		ClinetMap: make(map[*Client]bool, 5),
		BroadCast: make(chan []byte),
	}
}

func (w *World) run() {
	w.ChanEnter = make(chan *Client)
	w.ChanLeave = make(chan *Client)

	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case client := <-w.ChanEnter:
			fmt.Println("enter a client")
			w.ClinetMap[client] = true
		case client := <-w.ChanLeave:
			if _, ok := w.ClinetMap[client]; ok {
				delete(w.ClinetMap, client)
				fmt.Println("leave the client")
				close(client.send)
			}
		case message := <-w.BroadCast:
			for client := range w.ClinetMap {
				client.send <- message
			}
		case tick := <-ticker.C:
			for client := range w.ClinetMap {
				client.send <- []byte(tick.String())
			}
			fmt.Println(tick.String())
		}
	}
}
