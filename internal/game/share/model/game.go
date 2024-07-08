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
	Order() [][]PlayerId
	SetOrder(order [][]PlayerId)
	SetRandomOrder()
	GetOrderedPlayerIds(index int) []PlayerId
	GetOrderedPlayers(index int) []PlayerT
	GetOrderedPlayerId(index int) PlayerId
	GetOrderedPlayer(index int) PlayerT
	GetRoundPlayerIds() []PlayerId
	GetRoundPlayerId() PlayerId
	GetRoundPlayers() []PlayerT
	GetRoundPlayer() PlayerT
	SetPlayingRoundPlayers()
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
	MustGetPlayer(id PlayerId) PlayerT

	GetPlayingPlayer() (PlayerT, bool)
	GetPlayingPlayers() []PlayerT
	GetNonPlayingPlayers() []PlayerT
	SetPlayingPlayer(playerIds ...PlayerId)

	RankIdFn(leftId, rightId PlayerId) RankResult
	RankPlayerFn(left, right PlayerT) RankResult
	RankPlayers()

	UpdatePlayerScoreFn(player PlayerT)
	UpdateScores()

	SetLoosers(looserIds ...PlayerId)
	SetWinners(winnerIds ...PlayerId)
	SetTie()

	LabelSlice() []string
	Labels() string
}

type RankResult int

const (
	RankResult_Left RankResult = iota
	RankResult_Equal
	RankResult_Right
)

// //////////////////////////////////////////////////
// game

func NewGame[PlayerT Player]() Game[PlayerT] {
	return &game[PlayerT]{
		id:        GenerateGameId(),
		status:    GameStatus_JoinableNotStartable,
		createdAt: time.Now(),
		players:   make(map[PlayerId]PlayerT),
		order:     make([][]PlayerId, 0),
		ranks:     make([][]PlayerId, 0),
		round:     0,
	}
}

type game[PlayerT Player] struct {
	id        GameId
	status    GameStatus
	createdAt time.Time
	players   map[PlayerId]PlayerT
	round     int
	order     [][]PlayerId
	ranks     [][]PlayerId
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

func (g *game[PlayerT]) Order() [][]PlayerId {
	return g.order
}

func (g *game[PlayerT]) SetOrder(order [][]PlayerId) {
	for _, round := range order {
		for _, playerId := range round {
			if !g.HasPlayer(playerId) {
				panic(ErrPlayerNotFound)
			}
		}
	}
	g.order = order
}

func (g *game[PlayerT]) SetRandomOrder() {
	ids := dict.ConvertToList(g.players, dict.Key)
	list.Shuffle(ids)
	order := make([][]PlayerId, 0, len(ids))
	for _, id := range ids {
		order = append(order, []PlayerId{id})
	}
	g.SetOrder(order)
}

func (g *game[PlayerT]) GetOrderedPlayerIds(index int) []PlayerId {
	return g.order[index%len(g.order)]
}

func (g *game[PlayerT]) GetOrderedPlayers(index int) []PlayerT {
	playerIds := g.GetOrderedPlayerIds(index)
	return list.Convert(playerIds, func(playerId PlayerId) PlayerT {
		return g.MustGetPlayer(playerId)
	})
}

func (g *game[PlayerT]) GetOrderedPlayerId(index int) PlayerId {
	playerIds := g.GetOrderedPlayerIds(index)
	if len(playerIds) == 0 {
		panic(ErrPlayerNotFound)
	}
	return playerIds[0]
}

func (g *game[PlayerT]) GetOrderedPlayer(index int) PlayerT {
	playerId := g.GetOrderedPlayerId(index)
	return g.MustGetPlayer(playerId)
}

func (g *game[PlayerT]) GetRoundPlayerIds() []PlayerId {
	orderIndex := g.Round() % g.NbPlayer()
	return g.order[orderIndex]
}

func (g *game[PlayerT]) GetRoundPlayerId() PlayerId {
	playerIds := g.GetRoundPlayerIds()
	if len(playerIds) == 0 {
		panic(ErrPlayerNotFound)
	}
	return playerIds[0]
}

func (g *game[PlayerT]) GetRoundPlayers() []PlayerT {
	playerIds := g.GetRoundPlayerIds()
	return list.Convert(playerIds, func(playerId PlayerId) PlayerT {
		return g.MustGetPlayer(playerId)
	})
}

func (g *game[PlayerT]) GetRoundPlayer() PlayerT {
	playerId := g.GetRoundPlayerId()
	return g.MustGetPlayer(playerId)
}

func (g *game[PlayerT]) SetPlayingRoundPlayers() {
	g.SetPlayingPlayer(g.GetRoundPlayerIds()...)
}

func (g *game[PlayerT]) SetPlayingRoundPlayer() {
	g.SetPlayingRoundPlayers()
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
	player.SetGameId(g.id)
}

func (g *game[PlayerT]) DetachPlayer(player PlayerT) {
	delete(g.players, player.Id())
	player.UnsetGameId()
}

func (g *game[PlayerT]) HasPlayer(playerId PlayerId) bool {
	return dict.ContainsKey(g.players, playerId)
}

func (g *game[PlayerT]) GetPlayer(id PlayerId) (PlayerT, bool) {
	return dict.Get(g.players, id)
}

func (g *game[PlayerT]) MustGetPlayer(id PlayerId) PlayerT {
	player, found := g.GetPlayer(id)
	if !found {
		panic(ErrPlayerNotFound)
	}
	return player
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

func (g *game[PlayerT]) RankIdFn(leftId, rightId PlayerId) RankResult {
	left, leftFound := g.GetPlayer(leftId)
	right, rightFound := g.GetPlayer(rightId)
	if !leftFound && !rightFound {
		return RankResult_Equal
	} else if leftFound && !rightFound {
		return RankResult_Left
	} else if !leftFound && rightFound {
		return RankResult_Right
	} else {
		return g.RankPlayerFn(left, right)
	}
}

func (g *game[PlayerT]) RankPlayerFn(left, right PlayerT) RankResult {
	leftScore := left.Score()
	rightScore := right.Score()
	if leftScore > rightScore {
		return RankResult_Left
	} else if rightScore > leftScore {
		return RankResult_Right
	} else {
		return RankResult_Equal
	}
}

func (g *game[PlayerT]) RankPlayers() {

	ranks := make([][]PlayerId, 0, len(g.players))
	for leftId := range g.players {
		added := false
		for i, group := range ranks {
			rightId := group[0]
			switch g.RankIdFn(leftId, rightId) {
			case RankResult_Left:
				if i > 0 {
					ranks = append(ranks[:i], append([][]PlayerId{{leftId}}, ranks[i:]...)...)
				} else {
					ranks = append([][]PlayerId{{leftId}}, ranks...)
				}
				added = true
				continue
			case RankResult_Right:
				continue
			case RankResult_Equal:
				ranks[i] = append(ranks[i], leftId)
				added = true
			}
			if added {
				break
			}
		}
		if !added {
			ranks = append(ranks, []PlayerId{leftId})
		}
	}

	for rankIndex, playerIds := range ranks {
		for _, playerId := range playerIds {
			player, found := g.GetPlayer(playerId)
			if !found {
				panic(ErrPlayerNotFound)
			}
			player.SetRank(PlayerRank(rankIndex + 1))
		}
	}
}

func (g *game[PlayerT]) UpdatePlayerScoreFn(player PlayerT) {
	player.UnsetScore()
}

func (g *game[PlayerT]) UpdateScores() {
	for _, player := range g.players {
		g.UpdatePlayerScoreFn(player)
	}
	g.RankPlayers()
}

func (g *game[PlayerT]) SetLoosers(looserIds ...PlayerId) {
	for playerId, player := range g.players {
		if list.Contains(looserIds, playerId) {
			player.SetLoose()
		} else {
			player.SetWin()
		}
	}
}

func (g *game[PlayerT]) SetWinners(winnerIds ...PlayerId) {
	for playerId, player := range g.players {
		if list.Contains(winnerIds, playerId) {
			player.SetWin()
		} else {
			player.SetLoose()
		}
	}
}

func (g *game[PlayerT]) SetTie() {
	for _, player := range g.players {
		player.SetTie()
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
