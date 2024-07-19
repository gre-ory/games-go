package model

type User interface {
	HasId() bool
	Id() UserId
	SetId(id UserId)
	IsUser(userId UserId) bool
	IsNotUser(userId UserId) bool

	HasName() bool
	Name() UserName
	SetName(name UserName)

	HasAvatar() bool
	Avatar() UserAvatar
	SetAvatar(avatar UserAvatar)

	HasLanguage() bool
	Language() UserLanguage
	SetLanguage(language UserLanguage)

	SetCookie(cookie *Cookie)
}

func NewUser(id UserId) User {
	return &user{
		id:       id,
		name:     UserName(id),
		avatar:   1,
		language: UserLanguage_Fr,
	}
}

func NewUserFromCookie(cookie *Cookie) User {
	return &user{
		id:       cookie.Id,
		name:     cookie.Name,
		avatar:   cookie.Avatar,
		language: cookie.Language,
	}
}

func NewUserFromUser(other User) User {
	return &user{
		id:       other.Id(),
		name:     other.Name(),
		avatar:   other.Avatar(),
		language: other.Language(),
	}
}

type user struct {
	id       UserId
	name     UserName
	avatar   UserAvatar
	language UserLanguage
}

func (u *user) HasId() bool {
	return u.id != ""
}

func (u *user) Id() UserId {
	return u.id
}

func (u *user) SetId(userId UserId) {
	u.id = userId
}

func (u *user) IsUser(userId UserId) bool {
	return u.id == userId
}

func (u *user) IsNotUser(userId UserId) bool {
	return u.id != userId
}

func (u *user) HasName() bool {
	return u.name != ""
}

func (u *user) Name() UserName {
	return u.name
}

func (u *user) SetName(name UserName) {
	u.name = name
}

func (u *user) HasAvatar() bool {
	return u.avatar != 0
}

func (u *user) Avatar() UserAvatar {
	return u.avatar
}

func (u *user) SetAvatar(avatar UserAvatar) {
	u.avatar = avatar
}

func (u *user) HasLanguage() bool {
	return u.language != ""
}

func (u *user) Language() UserLanguage {
	return u.language
}

func (u *user) SetLanguage(language UserLanguage) {
	u.language = language
}

func (u *user) SetCookie(cookie *Cookie) {
	if u.id != cookie.Id {
		panic(ErrInvalidCookie)
	}
	u.name = cookie.Name
	u.avatar = cookie.Avatar
	u.language = cookie.Language
}
