package model

// //////////////////////////////////////////////////
// player

type Player interface {
	HasId() bool
	Id() PlayerId
	Status() PlayerStatus
	SetStatus(status PlayerStatus)
	Avatar() UserAvatar
	SetAvatar(avatar UserAvatar)
	Name() UserName
	SetName(name UserName)
	Language() UserLanguage
	SetLanguage(language UserLanguage)
	SetCookie(cookie *Cookie)
	HasGameId() bool
	GameId() GameId
	SetGameId(gameId GameId)
	UnsetGameId()
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
