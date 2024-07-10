package model

import (
	"html/template"
	"strings"

	"github.com/gre-ory/games-go/internal/util/loc"
)

// //////////////////////////////////////////////////
// player

type Player interface {
	HasId() bool
	Id() PlayerId
	Avatar() UserAvatar
	SetAvatar(avatar UserAvatar)
	Name() UserName
	SetName(name UserName)
	Language() UserLanguage
	SetLanguage(language UserLanguage)
	SetCookie(cookie *Cookie)

	IsPlaying() bool
	Status() PlayerStatus
	SetStatus(status PlayerStatus)

	HasGameId() bool
	GameId() GameId
	SetGameId(gameId GameId)
	UnsetGameId()

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

	YourMessage(localizer loc.Localizer) template.HTML
	Message(localizer loc.Localizer) template.HTML

	LabelSlice() []string
	Labels() string
}

// //////////////////////////////////////////////////
// base player

type player struct {
	id       PlayerId
	avatar   UserAvatar
	name     UserName
	language UserLanguage
	status   PlayerStatus
	gameId   GameId
	hasScore bool
	score    PlayerScore
	rank     PlayerRank
	result   PlayerResult
}

func NewPlayer(id PlayerId) Player {
	return &player{
		id: id,
	}
}

func NewPlayerFromCookie(cookie *Cookie) Player {
	return &player{
		id:       cookie.PlayerId(),
		avatar:   cookie.Avatar,
		name:     cookie.Name,
		language: cookie.Language,
	}
}

func (p *player) HasId() bool {
	return p.id != ""
}

func (p *player) Id() PlayerId {
	return p.id
}

func (p *player) Avatar() UserAvatar {
	return p.avatar
}

func (p *player) SetAvatar(avatar UserAvatar) {
	p.avatar = avatar
}

func (p *player) Name() UserName {
	return p.name
}

func (p *player) SetName(name UserName) {
	p.name = name
}

func (p *player) Language() UserLanguage {
	return p.language
}

func (p *player) LocLanguage() loc.Language {
	return loc.Language(p.language)
}

func (p *player) SetLanguage(language UserLanguage) {
	p.language = language
}

func (p *player) SetCookie(cookie *Cookie) {
	if p.id != cookie.PlayerId() {
		panic(ErrInvalidCookie)
	}
	p.avatar = cookie.Avatar
	p.name = cookie.Name
	p.language = cookie.Language
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

func (p *player) HasGameId() bool {
	return p.gameId != ""
}

func (p *player) GameId() GameId {
	return p.gameId
}

func (p *player) SetGameId(gameId GameId) {
	p.gameId = gameId
}

func (p *player) UnsetGameId() {
	p.gameId = ""
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

func (p *player) YourMessage(localizer loc.Localizer) template.HTML {
	switch p.Status() {
	case PlayerStatus_WaitingToJoin:
		return localizer.Loc("YouWaitingToJoin")
	case PlayerStatus_WaitingToStart:
		return localizer.Loc("YouWaitingToStart")
	case PlayerStatus_WaitingToPlay:
		return localizer.Loc("YouWaitingToPlay")
	case PlayerStatus_Playing:
		return localizer.Loc("YouPlaying")
	}
	return ""
}

func (p *player) Message(localizer loc.Localizer) template.HTML {
	switch p.Status() {
	case PlayerStatus_WaitingToJoin:
		return localizer.Loc("PlayerWaitingToJoin")
	case PlayerStatus_WaitingToStart:
		return localizer.Loc("PlayerWaitingToStart")
	case PlayerStatus_WaitingToPlay:
		return localizer.Loc("PlayerWaitingToPlay")
	case PlayerStatus_Playing:
		return localizer.Loc("PlayerPlaying")
	}
	return ""
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
