package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/victordantasdev/scrum-poker-api/api"
)

func TestWsHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(api.WsHandler))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]
	ws, _, err := websocket.DefaultDialer.Dial(wsURL+"/ws/room1", nil)
	if err != nil {
		t.Fatalf("Error starting websocket connection: %v", err)
	}
	defer ws.Close()

	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Error reading server message: %v", err)
	}

	expectedMessage := `{"settings":{"deck":null,"show_cards":false,"restart_game":false,"max_players":0},"moves":null,"players":null}`
	if string(message) != expectedMessage {
		t.Errorf("Incorrect initial message received. Expected: %s, Received: %s", expectedMessage, string(message))
	}

	initialSettings := `{"settings":{"deck":[1,2,3],"show_cards":true,"restart_game":false,"max_players":3},"moves":null,"players":null}`
	err = ws.WriteMessage(websocket.TextMessage, []byte(initialSettings))
	if err != nil {
		t.Fatalf("Error sending message to server: %v", err)
	}

	_, response, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Error reading server response: %v", err)
	}

	expectedResponse := `{"settings":{"deck":[1,2,3],"show_cards":true,"restart_game":false,"max_players":3},"moves":null,"players":null}`
	if string(response) != expectedResponse {
		t.Errorf("Incorrect server response. Expected: %s, Received: %s", expectedResponse, string(response))
	}

	initialMove := `{"move":{"username":"test","selected_card":3,"selected_card_index":2}}`
	err = ws.WriteMessage(websocket.TextMessage, []byte(initialMove))
	if err != nil {
		t.Fatalf("Error sending message to server: %v", err)
	}

	_, response, err = ws.ReadMessage()
	if err != nil {
		t.Fatalf("Error reading server response: %v", err)
	}

	expectedResponse = `{"settings":{"deck":[1,2,3],"show_cards":false,"restart_game":false,"max_players":0},"moves":[{"username":"test","selected_card":3,"selected_card_index":2,"remove_player":false}],"players":null}`
	if string(response) != expectedResponse {
		t.Errorf("Incorrect server response. Expected: %s, Received: %s", expectedResponse, string(response))
	}
}
