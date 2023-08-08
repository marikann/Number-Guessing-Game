package main

import (
	"MatchMaking/models"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func createPlayer(nickname string) (*models.Player, string) {
	id := generateUUID()
	player := &models.Player{ID: id, Nickname: nickname}
	return player, id
}

var mu sync.Mutex

// Yardımcı fonksiyonlar

// UUID üretme fonksiyonu
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Oda ID'si üretme fonksiyonu
func generateRoomID() int {
	return rand.Intn(1000000) + 1
}

// Rastgele sayı üretme fonksiyonu
func generateRandomNumber() int {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(100) + 1

	return randomNumber
}

// Oyuncu sonuçlarını hesaplayan fonksiyon
func calculateResults(players []*models.Player, secret int) []*models.PlayerResult {
	results := make([]*models.PlayerResult, 0)

	trophys := []int{
		30, 10, 0,
	}

	for _, player := range players {
		player.GuessDiffWithSecret = abs(player.Guess - secret)
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].GuessDiffWithSecret < players[j].GuessDiffWithSecret
	})

	rank := 0

	// Tahminleri değerlendir
	for i, p := range players {
		if p.Guess != 0 {
			rank = i + 1
			mu.Lock()
			p.Trophy += trophys[i]
			mu.Unlock()
			results = append(results, &models.PlayerResult{Rank: rank, Player: p.ID, Guess: p.Guess, DeltaTrophy: trophys[i]})
		} else {
			mu.Lock()
			p.Trophy += 0
			mu.Unlock()
			results = append(results, &models.PlayerResult{Rank: -1, Player: p.ID})
		}

		fmt.Println("player", i, "    ", p.Trophy)
	}

	return results
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func findWebSocketConnection(id string) (*websocket.Conn, error) {
	for conn := range connections {
		if conn.id == id {
			return conn.conn, nil
		}
	}
	return nil, fmt.Errorf("Kullanıcı için WebSocket bağlantısı bulunamadı")
}
