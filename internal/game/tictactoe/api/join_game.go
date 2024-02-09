package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/game/tictactoe/service"
)

func HandleJoinGame(w http.ResponseWriter, r *http.Request) {
	game, player := service.JoinGame("AAAA", "player 2")
	tpl.ExecuteTemplate(w, "play.tpl", map[string]any{
		"player": player,
		"game":   game,
	})
}
