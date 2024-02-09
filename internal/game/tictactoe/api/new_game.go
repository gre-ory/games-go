package api

import (
	"net/http"

	"github.com/gre-ory/games-go/internal/game/tictactoe/service"
)

func HandleNewGame(w http.ResponseWriter, r *http.Request) {
	game, player := service.NewGame("player 1")
	tpl.ExecuteTemplate(w, "play.tpl", map[string]any{
		"player": player,
		"game":   game,
	})
}
