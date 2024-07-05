package model

import (
	"strings"
	"time"

	"github.com/gre-ory/games-go/internal/util/dict"
	"github.com/gre-ory/games-go/internal/util/list"
)

// //////////////////////////////////////////////////
// game

type Game[PlayerT Player] interface {
	Id() GameId

	Status() GameStatus
	SetStatus(status GameStatus)
	MarkForDeletion()

	CreatedAt() time.Time

	Start()
	Round() int
	NextRound()
	Order() []PlayerId
	SetOrder(order []PlayerId)
	SetRandomOrder()
	GetOrderedPlayerId(index int) PlayerId
	GetOrderedPlayer(index int) PlayerT
	GetRoundPlayerId() PlayerId
	SetPlayingRoundPlayer()
	Stop()

	HasPlayers() bool
	NbPlayer() int
	GetPlayers() []PlayerT
	FilterPlayers(filterFn func(player PlayerT) bool) []PlayerT
	AttachPlayer(player PlayerT)
	DetachPlayer(player PlayerT)
	HasPlayer(playerId PlayerId) bool
	GetPlayer(id PlayerId) (PlayerT, bool)

	GetPlayingPlayer() (PlayerT, bool)
	GetPlayingPlayers() []PlayerT
	GetNonPlayingPlayers() []PlayerT
	SetPlayingPlayer(playerIds ...PlayerId)

	PlayerResult(playerId PlayerId) PlayerResult
	SetWinnerOthers(playerIds ...PlayerId)
	SetWinners(winnerIds ...PlayerId)
	SetTie()
	SetWinRank(playerId PlayerId, rank int)
	SetTieRank(playerId PlayerId, rank int)
	SetLooseRank(playerId PlayerId, rank int)

	LabelSlice() []string
	Labels() string
}

// //////////////////////////////////////////////////
// game

func NewGame[PlayerT Player]() Game[PlayerT] {
	return &game[PlayerT]{
		id:            GenerateGameId(),
		status:        GameStatus_JoinableNotStartable,
		createdAt:     time.Now(),
		players:       make(map[PlayerId]PlayerT),
		order:         make([]PlayerId, 0),
		playerResults: make(map[PlayerId]PlayerResult),
		round:         0,
	}
}

type game[PlayerT Player] struct {
	id            GameId
	status        GameStatus
	createdAt     time.Time
	players       map[PlayerId]PlayerT
	round         int
	order         []PlayerId
	playerResults map[PlayerId]PlayerResult
}

func (g *game[PlayerT]) Id() GameId {
	return g.id
}

func (g *game[PlayerT]) Status() GameStatus {
	return g.status
}

func (g *game[PlayerT]) SetStatus(status GameStatus) {
	g.status = status
}

func (g *game[PlayerT]) MarkForDeletion() {
	g.SetStatus(GameStatus_MarkedForDeletion)
}

func (g *game[PlayerT]) CreatedAt() time.Time {
	return g.createdAt
}

func (g *game[PlayerT]) Start() {
	g.SetStatus(GameStatus_Started)
	g.round = 1
}

func (g *game[PlayerT]) Round() int {
	return g.round
}

func (g *game[PlayerT]) NextRound() {
	g.round++
}

func (g *game[PlayerT]) Order() []PlayerId {
	return g.order
}

func (g *game[PlayerT]) SetOrder(order []PlayerId) {
	for _, playerId := range order {
		if !g.HasPlayer(playerId) {
			panic(ErrPlayerNotFound)
		}
	}
	g.order = order
}

func (g *game[PlayerT]) SetRandomOrder() {
	ids := dict.ConvertToList(g.players, dict.Key)
	list.Shuffle(ids)
	g.SetOrder(ids)
}

func (g *game[PlayerT]) GetOrderedPlayerId(index int) PlayerId {
	if index < 0 || index >= len(g.order) {
		panic(ErrPlayerNotFound)
	}
	return g.order[index]
}

func (g *game[PlayerT]) GetOrderedPlayer(index int) PlayerT {
	player, found := g.GetPlayer(g.GetOrderedPlayerId(index))
	if !found {
		panic(ErrPlayerNotFound)
	}
	return player
}

func (g *game[PlayerT]) GetRoundPlayerId() PlayerId {
	orderIndex := g.Round() % g.NbPlayer()
	return g.order[orderIndex]
}

func (g *game[PlayerT]) GetRoundPlayer() PlayerT {
	player, found := g.GetPlayer(g.GetRoundPlayerId())
	if !found {
		panic(ErrPlayerNotFound)
	}
	return player
}

func (g *game[PlayerT]) SetPlayingRoundPlayer() {
	g.SetPlayingPlayer(g.GetRoundPlayerId())
}

func (g *game[PlayerT]) Stop() {
	g.SetStatus(GameStatus_Stopped)
	g.round = 0
}

func (g *game[PlayerT]) HasPlayers() bool {
	return len(g.players) > 0
}

func (g *game[PlayerT]) NbPlayer() int {
	return len(g.players)
}

func (g *game[PlayerT]) GetPlayers() []PlayerT {
	return dict.Values(g.players)
}

func (g *game[PlayerT]) FilterPlayers(filterFn func(player PlayerT) bool) []PlayerT {
	return dict.Filter(g.players, filterFn)
}

func (g *game[PlayerT]) AttachPlayer(player PlayerT) {
	g.players[player.Id()] = player
	g.playerResults[player.Id()] = PlayerResult_Unknown
	player.SetGameId(g.id)
}

func (g *game[PlayerT]) DetachPlayer(player PlayerT) {
	delete(g.players, player.Id())
	delete(g.playerResults, player.Id())
	player.UnsetGameId()
}

func (g *game[PlayerT]) HasPlayer(playerId PlayerId) bool {
	return dict.ContainsKey(g.players, playerId)
}

func (g *game[PlayerT]) GetPlayer(id PlayerId) (PlayerT, bool) {
	return dict.Get(g.players, id)
}

func (g *game[PlayerT]) GetPlayingPlayer() (PlayerT, bool) {
	return dict.First(g.players, func(player PlayerT) bool {
		return player.Status().IsPlaying()
	})
}

func (g *game[PlayerT]) GetPlayingPlayers() []PlayerT {
	return dict.Filter(g.players, func(player PlayerT) bool {
		return player.Status().IsPlaying()
	})
}

func (g *game[PlayerT]) GetNonPlayingPlayers() []PlayerT {
	return dict.Filter(g.players, func(player PlayerT) bool {
		return player.Status().IsWaitingToPlay()
	})
}

func (g *game[PlayerT]) SetPlayingPlayer(playerIds ...PlayerId) {
	count := 0
	for _, player := range g.players {
		if list.Contains(playerIds, player.Id()) {
			count++
			player.SetStatus(PlayerStatus_Playing)
		} else {
			player.SetStatus(PlayerStatus_WaitingToPlay)
		}
	}
	if len(playerIds) != count {
		panic(ErrPlayerNotFound)
	}
}

func (g *game[PlayerT]) PlayerResult(playerId PlayerId) PlayerResult {
	return dict.MustGet(g.playerResults, playerId)
}

func (g *game[PlayerT]) SetWinnerOthers(playerIds ...PlayerId) {
	otherIds := list.Filter(g.order, func(otherId PlayerId) bool {
		return !list.Contains(playerIds, otherId)
	})
	g.SetWinners(otherIds...)
}

func (g *game[PlayerT]) SetWinners(winnerIds ...PlayerId) {
	for playerId := range g.playerResults {
		if list.Contains(winnerIds, playerId) {
			g.playerResults[playerId] = NewWinResult()
		} else {
			g.playerResults[playerId] = NewLooseResult()
		}
	}
}

func (g *game[PlayerT]) SetTie() {
	for playerId := range g.playerResults {
		g.playerResults[playerId] = NewTieResult()
	}
}

func (g *game[PlayerT]) SetWinRank(playerId PlayerId, rank int) {
	if g.HasPlayer(playerId) {
		g.playerResults[playerId] = NewWinRankResult(rank)
	}
}

func (g *game[PlayerT]) SetTieRank(playerId PlayerId, rank int) {
	if g.HasPlayer(playerId) {
		g.playerResults[playerId] = NewTieRankResult(rank)
	}
}

func (g *game[PlayerT]) SetLooseRank(playerId PlayerId, rank int) {
	if g.HasPlayer(playerId) {
		g.playerResults[playerId] = NewLooseRankResult(rank)
	}
}

func (g *game[PlayerT]) LabelSlice() []string {
	labels := make([]string, 0)
	labels = append(labels, "game")
	labels = append(labels, g.status.Labels()...)
	return labels
}

func (g *game[PlayerT]) Labels() string {
	return strings.Join(g.LabelSlice(), " ")
}
