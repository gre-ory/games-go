package model

import (
	"html/template"
	"strings"
	"time"

	"github.com/gre-ory/games-go/internal/util/dict"
	"github.com/gre-ory/games-go/internal/util/list"
	"github.com/gre-ory/games-go/internal/util/loc"
)

// //////////////////////////////////////////////////
// game

type Game[PlayerT Player] interface {
	Id() GameId

	IsStarted() bool
	IsStopped() bool
	WasStarted() bool
	Status() GameStatus
	SetStatus(status GameStatus)
	SetStarted()
	SetStopped()
	IsMarkedForDeletion() bool
	MarkForDeletion()

	CreatedAt() time.Time

	CanJoin() bool
	CanStart() bool
	UpdateJoinStatus()

	Round() int
	FirstRound()
	NextRound()
	Order() [][]PlayerId
	SetOrder(order [][]PlayerId)
	SetRandomOrder()
	OrderedPlayerIds(index int) []PlayerId
	OrderedPlayers(index int) []PlayerT
	OrderedPlayerId(index int) PlayerId
	OrderedPlayer(index int) PlayerT
	RoundPlayerIds() []PlayerId
	RoundPlayerId() PlayerId
	RoundPlayers() []PlayerT
	RoundPlayer() PlayerT
	SetPlayingRoundPlayers()
	SetPlayingRoundPlayer()

	HasPlayers() bool
	NbPlayer() int
	Players() []PlayerT
	FilterPlayers(filterFn func(player PlayerT) bool) []PlayerT
	AttachPlayer(player PlayerT)
	DetachPlayer(player PlayerT)
	HasUser(userId UserId) bool
	HasPlayer(playerId PlayerId) bool
	Player(id PlayerId) (PlayerT, bool)
	MustPlayer(id PlayerId) PlayerT
	PlayerLabelSlice(id PlayerId) []string
	PlayerLabels(id PlayerId) string

	IsPlayingPlayer(playerId PlayerId) bool
	PlayingPlayer() (PlayerT, bool)
	PlayingPlayers() []PlayerT
	NonPlayingPlayers() []PlayerT
	SetPlayingPlayer(playerIds ...PlayerId)

	RankIdFn(leftId, rightId PlayerId) RankResult
	RankPlayerFn(left, right PlayerT) RankResult
	RankPlayers()

	UpdatePlayerScoreFn(player PlayerT)
	UpdateScores()

	SetLoosers(looserIds ...PlayerId)
	SetWinners(winnerIds ...PlayerId)
	SetTie()

	YourPlayerMessage(localizer loc.Localizer, playerId PlayerId) template.HTML
	PlayerMessage(localizer loc.Localizer, playerId PlayerId) template.HTML
	PlayerStatusIcon(playerId PlayerId) string

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

func NewGame[PlayerT Player](minNbPlayer, maxNbPlayer int) Game[PlayerT] {
	return &game[PlayerT]{
		id:          GenerateGameId(),
		status:      GameStatus_JoinableNotStartable,
		createdAt:   time.Now(),
		minNbPlayer: minNbPlayer,
		maxNbPlayer: maxNbPlayer,
		players:     make(map[PlayerId]PlayerT),
		order:       make([][]PlayerId, 0),
		ranks:       make([][]PlayerId, 0),
		round:       0,
	}
}

type game[PlayerT Player] struct {
	id          GameId
	status      GameStatus
	createdAt   time.Time
	minNbPlayer int
	maxNbPlayer int
	players     map[PlayerId]PlayerT
	round       int
	order       [][]PlayerId
	ranks       [][]PlayerId
}

func (g *game[PlayerT]) Id() GameId {
	return g.id
}

func (g *game[PlayerT]) IsStarted() bool {
	return g.Status().IsStarted()
}

func (g *game[PlayerT]) IsStopped() bool {
	return g.Status().IsStopped()
}

func (g *game[PlayerT]) WasStarted() bool {
	return g.IsStarted() || g.IsStopped()
}

func (g *game[PlayerT]) Status() GameStatus {
	return g.status
}

func (g *game[PlayerT]) SetStatus(status GameStatus) {
	g.status = status
}

func (g *game[PlayerT]) SetStarted() {
	g.status = GameStatus_Started
}

func (g *game[PlayerT]) SetStopped() {
	g.status = GameStatus_Stopped
}

func (g *game[PlayerT]) IsMarkedForDeletion() bool {
	return g.status == GameStatus_MarkedForDeletion
}

func (g *game[PlayerT]) MarkForDeletion() {
	g.SetStatus(GameStatus_MarkedForDeletion)
}

func (g *game[PlayerT]) CreatedAt() time.Time {
	return g.createdAt
}

func (g *game[PlayerT]) CanJoin() bool {
	return g.maxNbPlayer == 0 || len(g.players) < g.maxNbPlayer
}

func (g *game[PlayerT]) CanStart() bool {
	return len(g.players) >= g.minNbPlayer
}

func (g *game[PlayerT]) UpdateJoinStatus() {
	if g.WasStarted() || g.IsMarkedForDeletion() {
		return
	}
	if !g.CanJoin() {
		g.SetStatus(GameStatus_NotJoinableAndStartable)
		for _, player := range g.Players() {
			player.SetStatus(PlayerStatus_WaitingToStart)
		}
	} else if g.CanStart() {
		g.SetStatus(GameStatus_JoinableAndStartable)
		for _, player := range g.Players() {
			player.SetStatus(PlayerStatus_WaitingToStart)
		}
	} else {
		g.SetStatus(GameStatus_JoinableNotStartable)
		for _, player := range g.Players() {
			player.SetStatus(PlayerStatus_WaitingToJoin)
		}
	}
}

func (g *game[PlayerT]) Start() {
	g.SetStatus(GameStatus_Started)

}

func (g *game[PlayerT]) Round() int {
	return g.round
}

func (g *game[PlayerT]) FirstRound() {
	g.round = 1
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

func (g *game[PlayerT]) OrderedPlayerIds(index int) []PlayerId {
	return g.order[index%len(g.order)]
}

func (g *game[PlayerT]) OrderedPlayers(index int) []PlayerT {
	playerIds := g.OrderedPlayerIds(index)
	return list.Convert(playerIds, func(playerId PlayerId) PlayerT {
		return g.MustPlayer(playerId)
	})
}

func (g *game[PlayerT]) OrderedPlayerId(index int) PlayerId {
	playerIds := g.OrderedPlayerIds(index)
	if len(playerIds) == 0 {
		panic(ErrPlayerNotFound)
	}
	return playerIds[0]
}

func (g *game[PlayerT]) OrderedPlayer(index int) PlayerT {
	playerId := g.OrderedPlayerId(index)
	return g.MustPlayer(playerId)
}

func (g *game[PlayerT]) RoundPlayerIds() []PlayerId {
	orderIndex := g.Round() % g.NbPlayer()
	return g.order[orderIndex]
}

func (g *game[PlayerT]) RoundPlayerId() PlayerId {
	playerIds := g.RoundPlayerIds()
	if len(playerIds) == 0 {
		panic(ErrPlayerNotFound)
	}
	return playerIds[0]
}

func (g *game[PlayerT]) RoundPlayers() []PlayerT {
	playerIds := g.RoundPlayerIds()
	return list.Convert(playerIds, func(playerId PlayerId) PlayerT {
		return g.MustPlayer(playerId)
	})
}

func (g *game[PlayerT]) RoundPlayer() PlayerT {
	playerId := g.RoundPlayerId()
	return g.MustPlayer(playerId)
}

func (g *game[PlayerT]) SetPlayingRoundPlayers() {
	g.SetPlayingPlayer(g.RoundPlayerIds()...)
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

func (g *game[PlayerT]) Players() []PlayerT {
	return dict.Values(g.players)
}

func (g *game[PlayerT]) FilterPlayers(filterFn func(player PlayerT) bool) []PlayerT {
	return dict.Filter(g.players, filterFn)
}

func (g *game[PlayerT]) AttachPlayer(player PlayerT) {
	if player.Id().GameId() != g.Id() {
		panic(ErrWrongPlayer)
	}
	g.players[player.Id()] = player
}

func (g *game[PlayerT]) DetachPlayer(player PlayerT) {
	if player.Id().GameId() != g.Id() {
		panic(ErrWrongPlayer)
	}
	delete(g.players, player.Id())
	player.SetStatus(PlayerStatus_WaitingToJoin)
}

func (g *game[PlayerT]) HasUser(userId UserId) bool {
	for playerId := range g.players {
		if playerId.MatchUser(userId) {
			return true
		}
	}
	return false
}

func (g *game[PlayerT]) HasPlayer(playerId PlayerId) bool {
	return dict.ContainsKey(g.players, playerId)
}

func (g *game[PlayerT]) Player(id PlayerId) (PlayerT, bool) {
	return dict.Get(g.players, id)
}

func (g *game[PlayerT]) MustPlayer(id PlayerId) PlayerT {
	player, found := g.Player(id)
	if !found {
		panic(ErrPlayerNotFound)
	}
	return player
}

func (g *game[PlayerT]) PlayerLabelSlice(id PlayerId) []string {
	player, found := g.Player(id)
	if !found {
		return []string{"error"}
	}
	return player.LabelSlice()
}

func (g *game[PlayerT]) PlayerLabels(id PlayerId) string {
	return strings.Join(g.PlayerLabelSlice(id), " ")
}

func (g *game[PlayerT]) IsPlayingPlayer(playerId PlayerId) bool {
	player, found := g.Player(playerId)
	if !found {
		return false
	}
	return player.IsPlaying()
}

func (g *game[PlayerT]) PlayingPlayer() (PlayerT, bool) {
	return dict.First(g.players, func(player PlayerT) bool {
		return player.Status().IsPlaying()
	})
}

func (g *game[PlayerT]) PlayingPlayers() []PlayerT {
	return dict.Filter(g.players, func(player PlayerT) bool {
		return player.Status().IsPlaying()
	})
}

func (g *game[PlayerT]) NonPlayingPlayers() []PlayerT {
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
	left, leftFound := g.Player(leftId)
	right, rightFound := g.Player(rightId)
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
			player, found := g.Player(playerId)
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

func (g *game[PlayerT]) YourPlayerMessage(localizer loc.Localizer, playerId PlayerId) template.HTML {
	player, found := g.Player(playerId)
	if !found {
		return localizer.Loc("Error", ErrPlayerNotFound.Error())
	}
	if player.HasResult() {
		result := player.Result()
		switch {
		case result.IsWin():
			return localizer.Loc("YouWin")
		case result.IsTie():
			return localizer.Loc("YouTie")
		case result.IsLoose():
			return localizer.Loc("YouLoose")
		}
	}
	return player.Status().YourMessage(localizer)
}

func (g *game[PlayerT]) PlayerMessage(localizer loc.Localizer, playerId PlayerId) template.HTML {
	player, found := g.Player(playerId)
	if !found {
		return localizer.Loc("Error", ErrPlayerNotFound.Error())
	}
	if player.HasResult() {
		result := player.Result()
		switch {
		case result.IsWin():
			return localizer.Loc("PlayerWin")
		case result.IsTie():
			return localizer.Loc("PlayerTie")
		case result.IsLoose():
			return localizer.Loc("PlayerLoose")
		}
	}
	return player.Status().Message(localizer)
}

func (g *game[PlayerT]) PlayerStatusIcon(playerId PlayerId) string {
	player, found := g.Player(playerId)
	if !found {
		return ""
	}
	if player.HasResult() {
		return player.Result().Icon()
	}
	return player.Status().Icon()
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
