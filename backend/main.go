package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/papaaannn/8ballpool/engine"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type GameServer struct {
    mu    sync.Mutex
    game  *engine.Game
    conns map[*websocket.Conn]bool
}

func newGameServer() *GameServer {
    return &GameServer{
        game:  engine.NewGame(),
        conns: make(map[*websocket.Conn]bool),
    }
}

func (gs *GameServer) handleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Failed to upgrade to WebSocket:", err)
        return
    }
    defer func() {
        gs.mu.Lock()
        delete(gs.conns, ws)
        gs.mu.Unlock()
        ws.Close()
    }()

    gs.mu.Lock()
    gs.conns[ws] = true
    gs.mu.Unlock()

    for {
        var msg map[string]interface{}
        err := ws.ReadJSON(&msg)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Unexpected WebSocket closure: %v", err)
            } else {
                log.Printf("WebSocket closed: %v", err)
            }
            break
        }

        if msg["type"] == "shoot" {
            angle := msg["angle"].(float64)
            power := msg["power"].(float64)
            gs.mu.Lock()
            gs.game.ShootBall(angle, power)
            gs.mu.Unlock()
        } else if msg["type"] == "restart" {
            gs.mu.Lock()
            gs.game = engine.NewGame()
            gs.mu.Unlock()
        }
    }
}

func (gs *GameServer) run() {
    ticker := time.NewTicker(time.Second / 60) // 60 FPS update rate
    defer ticker.Stop()

    for {
        <-ticker.C
        gs.mu.Lock()
        gs.game.Update()
        gameState, _ := json.Marshal(gs.game)
        for conn := range gs.conns {
            err := conn.WriteMessage(websocket.TextMessage, gameState)
            if err != nil {
                log.Printf("Error writing message: %v", err)
                conn.Close()
                delete(gs.conns, conn)
            }
        }
        gs.mu.Unlock()
    }
}

func main() {
    gameServer := newGameServer()
    http.HandleFunc("/ws", gameServer.handleConnections)
    go gameServer.run()

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}