package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type Settings struct {
	Deck        []uint `json:"deck"`
	ShowCards   bool   `json:"show_cards"`
	RestartGame bool   `json:"restart_game"`
	MaxPlayers  int    `json:"max_players"`
}

type Move struct {
	Username          string `json:"username"`
	SelectedCard      int    `json:"selected_card"`
	SelectedCardIndex int    `json:"selected_card_index"`
	RemovePlayer      bool   `json:"remove_player"`
}

type NewMove struct {
	Settings `json:"settings"`
	Move     `json:"move"`
	Player   string `json:"player"`
}

type RoomData struct {
	Settings `json:"settings"`
	Moves    []Move   `json:"moves"`
	Players  []string `json:"players"`
}

type Room struct {
	mu       sync.Mutex
	RoomData RoomData
}

var rooms = make(map[string]*Room)
var clients = make(map[*websocket.Conn]string)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *Room) addRoomData(m NewMove) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(m.Settings.Deck) != 0 {
		r.RoomData.Settings.Deck = m.Settings.Deck
	}

	r.RoomData.Settings.ShowCards = m.Settings.ShowCards
	r.RoomData.Settings.RestartGame = m.Settings.RestartGame
	r.RoomData.Settings.MaxPlayers = m.Settings.MaxPlayers

	if m.Move.Username != "" {
		for i, move := range r.RoomData.Moves {
			if move.Username == m.Username {
				r.RoomData.Moves[i] = m.Move
				return
			}
		}

		r.RoomData.Moves = append(r.RoomData.Moves, m.Move)
	}

	if m.Player != "" {
		for i, player := range r.RoomData.Players {
			if player == m.Player {
				r.RoomData.Players[i] = m.Player
				return
			}
		}

		r.RoomData.Players = append(r.RoomData.Players, m.Player)
	}
}

func (r *Room) getRoomData() RoomData {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.RoomData
}

func (r *Room) sendRoomData(roomId string) {
	for client, roomID := range clients {
		if roomID == roomId {
			prev, _ := json.Marshal(r.getRoomData())

			var oldRoomData RoomData
			json.Unmarshal([]byte(prev), &oldRoomData)

			if err := client.WriteMessage(websocket.TextMessage, []byte(prev)); err != nil {
				log.Println("Error sending previous room data: ", err)
				return
			}
		}
	}
}

func (r *Room) restartMoves() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.RoomData.Settings.ShowCards = false
	r.RoomData.Settings.RestartGame = false

	for i := range r.RoomData.Moves {
		r.RoomData.Moves[i].SelectedCard = -1
		r.RoomData.Moves[i].SelectedCardIndex = -1
	}
}

func (r *Room) removePlayer(player string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, p := range r.RoomData.Players {
		if p == player {
			r.RoomData.Players = append(r.RoomData.Players[:i], r.RoomData.Players[i+1:]...)
		}
	}

	for i, m := range r.RoomData.Moves {
		if m.Username == player {
			r.RoomData.Moves = append(r.RoomData.Moves[:i], r.RoomData.Moves[i+1:]...)
			return
		}
	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error upgrading connection:", err)
	}
	defer conn.Close()

	roomId := r.URL.Path[len("/ws/"):]
	clients[conn] = roomId

	room, ok := rooms[roomId]
	if !ok {
		room = &Room{}
		rooms[roomId] = room
	}

	room.sendRoomData(roomId)

	for {
		_, newMove, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading ws message:", err)
			delete(clients, conn)
			return
		}

		var parsedNewMove NewMove
		json.Unmarshal([]byte(newMove), &parsedNewMove)

		removePlayer := parsedNewMove.RemovePlayer
		restartGame := parsedNewMove.RestartGame

		if removePlayer {
			room.removePlayer(parsedNewMove.Move.Username)
			delete(clients, conn)
		} else if restartGame {
			room.restartMoves()
		} else {
			if parsedNewMove.Move.Username != "" {
				log.Println(parsedNewMove.Move.Username + " selected card " + strconv.Itoa(parsedNewMove.Move.SelectedCard))
			}

			room.addRoomData(parsedNewMove)
		}

		room.sendRoomData(roomId)
	}
}
