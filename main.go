package main

import (
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	gameServer := NewGameServer()
	go gameServer.MatchMake()

	http.HandleFunc("/register", RegisterHandler(gameServer))
	http.HandleFunc("/stats", StatsHandler(gameServer))
	go gameServer.StartWebSocketServer()

	// Set up CORS headers
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	go func() {
		//açık goroutine lerime bakmak için kulladım
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Wrap the default ServeMux with the CORS middleware
	handler := corsMiddleware(http.DefaultServeMux)

	err := http.ListenAndServe(":1234", handler)
	if err != nil {
		log.Fatal("Web sunucu hatasi:", err)
	}
}

func RegisterHandler(gs *GameServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nickname := r.FormValue("nickname")
		if nickname == "" {
			http.Error(w, "Nickname eksik", http.StatusBadRequest)
			return
		}

		player := gs.RegisterPlayer(nickname)
		response := map[string]string{"id": player.ID}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func StatsHandler(gs *GameServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := gs.GetStats()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
