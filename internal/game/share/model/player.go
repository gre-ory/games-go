package model

import (
	"strings"
)

// //////////////////////////////////////////////////
// player

type Player interface {
	User() User

	HasId() bool
	Id() PlayerId

	GameId() GameId

	IsPlaying() bool
	Status() PlayerStatus
	SetStatus(status PlayerStatus)

	HasScore() bool
	Score() PlayerScore
	SetScore(score PlayerScore)
	AddScore(score PlayerScore)
	RemoveScore(score PlayerScore)
	UnsetScore()

	HasRank() bool
	Rank() PlayerRank
	SetRank(rank PlayerRank)
	UnsetRank()

	HasResult() bool
	Result() PlayerResult
	SetResult(result PlayerResult)
	SetWin()
	SetTie()
	SetLoose()
	UnsetResult()

	LabelSlice() []string
	Labels() string
}

// //////////////////////////////////////////////////
// base player

type player struct {
	user     User
	id       PlayerId
	status   PlayerStatus
	gameId   GameId
	hasScore bool
	score    PlayerScore
	rank     PlayerRank
	result   PlayerResult
}

func NewPlayer(gameId GameId, userId UserId) Player {
	return &player{
		user:   NewUser(userId),
		id:     NewPlayerId(gameId, userId),
		gameId: gameId,
	}
}

func NewPlayerFromUser(gameId GameId, user User) Player {
	return &player{
		user:   NewUserFromUser(user),
		id:     NewPlayerId(gameId, user.Id()),
		gameId: gameId,
	}
}

func (p *player) User() User {
	return p.user
}

func (p *player) HasId() bool {
	return p.id != ""
}

func (p *player) Id() PlayerId {
	return p.id
}

func (p *player) IsPlaying() bool {
	return p.status.IsPlaying()
}

func (p *player) Status() PlayerStatus {
	return p.status
}

func (p *player) SetStatus(status PlayerStatus) {
	p.status = status
}

func (p *player) GameId() GameId {
	return p.gameId
}

func (p *player) HasScore() bool {
	return p.hasScore
}

func (p *player) Score() PlayerScore {
	return p.score
}

func (p *player) SetScore(score PlayerScore) {
	p.hasScore = true
	p.score = score
}

func (p *player) AddScore(score PlayerScore) {
	p.hasScore = true
	p.score += score
}

func (p *player) RemoveScore(score PlayerScore) {
	p.hasScore = true
	p.score -= score
}

func (p *player) UnsetScore() {
	p.hasScore = false
	p.score = 0
}

func (p *player) HasRank() bool {
	return p.rank != 0
}

func (p *player) Rank() PlayerRank {
	return p.rank
}

func (p *player) SetRank(rank PlayerRank) {
	p.rank = rank
}

func (p *player) UnsetRank() {
	p.rank = 0
}

func (p *player) HasResult() bool {
	return p.result != 0
}

func (p *player) Result() PlayerResult {
	return p.result
}

func (p *player) SetResult(result PlayerResult) {
	p.result = result
}

func (p *player) SetWin() {
	p.result = PlayerResult_Win
}

func (p *player) SetTie() {
	p.result = PlayerResult_Tie
}

func (p *player) SetLoose() {
	p.result = PlayerResult_Loose
}

func (p *player) UnsetResult() {
	p.result = PlayerResult_Unknown
}

func (p *player) LabelSlice() []string {
	labels := make([]string, 0)
	labels = append(labels, "player")
	labels = append(labels, p.Status().LabelSlice()...)
	if p.HasRank() {
		labels = append(labels, p.Rank().LabelSlice()...)
	}
	if p.HasResult() {
		labels = append(labels, p.Result().LabelSlice()...)
		if p.HasRank() {
			labels = append(labels, p.Rank().MedalLabelSlice()...)
		}
	}
	return labels
}

func (p *player) Labels() string {
	return strings.Join(p.LabelSlice(), " ")
}
