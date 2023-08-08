package main

import (
	"MatchMaking/models"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type GameServer struct {
	Players     map[string]*models.Player
	WaitingList []*models.Player
	Rooms       map[int]*Room
	Lock        sync.Mutex
}

func NewGameServer() *GameServer {
	return &GameServer{
		Players: make(map[string]*models.Player),
		Rooms:   make(map[int]*Room),
		Lock:    sync.Mutex{},
	}
}

// Oyuncu istatistikleri alma fonksiyonu
func (gs *GameServer) GetStats() map[string]interface{} {
	gs.Lock.Lock()
	defer gs.Lock.Unlock()

	activeRooms := gs.getActiveRooms()

	return map[string]interface{}{
		"registeredPlayers": len(gs.Players),
		"activeRooms":       activeRooms,
	}
}

func (gs *GameServer) getActiveRooms() []map[string]interface{} {
	activeRooms := make([]map[string]interface{}, 0)

	for _, room := range gs.Rooms {
		if !room.Finished {
			activeRooms = append(activeRooms, map[string]interface{}{
				"id":     room.ID,
				"secret": room.Secret,
			})
		}
	}

	return activeRooms
}

// Websocket sunucu başlatma fonksiyonu
func (gs *GameServer) StartWebSocketServer() {
	http.HandleFunc("/ws", gs.handleWebSocket)
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		log.Fatal("Websocket sunucu hatası:", err)
	}
}

// Websocket bağlantı işleme fonksiyonu
func (gs *GameServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade hatası:", err)
		return
	}
	defer conn.Close()

	for {
		var command models.WebSocketCommand
		err := conn.ReadJSON(&command)
		if err != nil {
			fmt.Println("Websocket okuma hatası:", err)
			break
		}

		gs.handleWebSocketCommand(conn, &command)
	}
}

func (gs *GameServer) handleWebSocketCommand(conn *websocket.Conn, command *models.WebSocketCommand) {
	switch command.Cmd {
	case "join":
		gs.handleJoinCommand(conn, command)
	case "guess":
		gs.handleGuessCommand(conn, command)
	default:
		fmt.Println("Geçersiz komut:", command.Cmd)
	}
}

func (gs *GameServer) MatchMake() {
	interval := 30 * time.Second
	ticker := time.Tick(interval)

	for range ticker {
		groups := gs.divideIntoGroupsOfThree()

		gs.createRoomsForValidGroups(groups)

		gs.updateWaitingList(groups)
	}
}

func (gs *GameServer) divideIntoGroupsOfThree() [][]*models.Player {
	gs.Lock.Lock()
	defer gs.Lock.Unlock()
	groups := make([][]*models.Player, 0)

	// Divide players into groups of 3
	for i := 0; i < len(gs.WaitingList); i += 3 {
		end := i + 3
		if end > len(gs.WaitingList) {
			end = len(gs.WaitingList)
		}
		groups = append(groups, gs.WaitingList[i:end])
	}

	return groups
}

func (gs *GameServer) updateWaitingList(groups [][]*models.Player) {
	newGroups := make([][]*models.Player, 0)
	for _, group := range groups {
		if len(group) < 3 {
			// If the group size is less than 3, exclude it from the waiting list
			newGroups = append(newGroups, group)
		}
	}
	gs.Lock.Lock()
	gs.WaitingList = gs.getFirstGroupIfAny(newGroups)
	gs.Lock.Unlock()
}

func (gs *GameServer) getFirstGroupIfAny(groups [][]*models.Player) []*models.Player {

	if len(groups) > 0 {
		return groups[0]
	}

	return nil
}

func (gs *GameServer) createRoomsForValidGroups(groups [][]*models.Player) {
	for _, group := range groups {
		gs.Lock.Lock()
		if len(group) == 3 {
			roomID := generateRoomID()
			secret := generateRandomNumber()
			room := NewRoom(roomID, secret, group, gs)
			gs.Rooms[roomID] = room
		}
		gs.Lock.Unlock()
	}
}

func sendWebSocketReply(conn *websocket.Conn, cmd, err, reply string) {
	replyData := models.WebSocketReply{Cmd: cmd, Error: err, Reply: reply}
	conn.WriteJSON(replyData)
}

// "join" komutunu işleme fonksiyonu
func (gs *GameServer) handleJoinCommand(conn *websocket.Conn, command *models.WebSocketCommand) {

	player, ok := gs.Players[command.ID]
	if !ok {
		sendWebSocketReply(conn, "join", "notRegistered", "")
		return
	}
	gs.Lock.Lock()
	gs.WaitingList = append(gs.WaitingList, player)
	gs.Lock.Unlock()
	sendWebSocketReply(conn, "join", "waiting", "")

	newConnection := &WebSocketConnection{conn: conn, id: player.ID}
	gs.Lock.Lock()
	connections[newConnection] = true
	gs.Lock.Unlock()

}

// "guess" komutunu işleme fonksiyonu
func (gs *GameServer) handleGuessCommand(conn *websocket.Conn, command *models.WebSocketCommand) {
	gs.Lock.Lock()
	defer gs.Lock.Unlock()

	player, ok := gs.Players[command.ID]
	if !ok {
		sendWebSocketReply(conn, "guess", "notRegistered", "")
		return
	}

	room, ok := gs.Rooms[command.Room]
	if !ok {
		sendWebSocketReply(conn, "guess", "notInRoom", "")
		return
	}

	for _, p := range room.Players {
		if p.ID == player.ID {
			p.Guess = command.Data
			break
		}
	}

	sendWebSocketReply(conn, "guess", "guessReceived", "")
	room.GuessChan <- command.Data
}

// Yeni oyuncu kaydetme fonksiyonu
func (gs *GameServer) RegisterPlayer(nickname string) *models.Player {

	player, id := createPlayer(nickname)
	gs.Lock.Lock()
	gs.Players[id] = player
	gs.Lock.Unlock()
	return player
}
