package tictactoe

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/game/tictactoe/api"
)

func Register(router *http.ServeMux) {
	router.HandleFunc("/tictactoe", api.HandleIndex)
	router.HandleFunc("/tictactoe/new", api.HandleNewGame)
	router.HandleFunc("/tictactoe/join", api.HandleJoinGame)
}
