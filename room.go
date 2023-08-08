package main

import (
	"MatchMaking/models"
	"fmt"
	"sync"
	"time"
)

type Room struct {
	ID         int
	Secret     int
	Players    []*models.Player
	GuessChan  chan int
	GuessCount int
	Finished   bool
	Lock       sync.Mutex
}

func NewRoom(ID int, secret int, players []*models.Player, gameServer *GameServer) *Room {
	room := &Room{
		ID:        ID,
		Secret:    secret,
		Players:   players,
		GuessChan: make(chan int),
		Finished:  false,
	}

	go room.waitForGuessesAndHandleGameOver(gameServer)

	// Notify players that they have joined the room
	room.notifyPlayersJoined()

	return room
}

func (room *Room) waitForGuessesAndHandleGameOver(gameServer *GameServer) {
	for {
		select {
		case <-room.GuessChan:
			room.Lock.Lock()
			fmt.Println("Girdi")
			room.GuessCount++
			room.Lock.Unlock()
			if room.GuessCount == 3 {
				room.handleGameOver()
				return
			}
			continue
			/*
					if !allPlayersGuessed(room.Players) {
						room.Lock.Unlock()
						continue
					}
					room.Lock.Unlock()
				return
			*/

		// 20 seconds have passed
		case <-time.After(20 * time.Second):
			room.handleGameOver()
			return
		}
	}

}

func (room *Room) handleGameOver() {
	room.Finished = true
	result := calculateResults(room.Players, room.Secret)
	event := models.WebSocketEvent{Event: "gameOver", Secret: room.Secret, Rankings: result}

	for _, p := range room.Players {
		room.Lock.Lock()
		p.Guess = 0
		conn, err := findWebSocketConnection(p.ID)
		if err == nil {
			conn.WriteJSON(event)
		}
		room.Lock.Unlock()

	}

}

func (room *Room) notifyPlayersJoined() {
	event := models.WebSocketEvent{Event: "joinedRoom", Room: room.ID, Secret: room.Secret}

	for _, p := range room.Players {
		conn, err := findWebSocketConnection(p.ID)
		if err == nil {
			conn.WriteJSON(event)
		}
	}
}
