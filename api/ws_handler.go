package api

import "sync"

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
