package model

type Cookie struct {
	Id       UserId
	Name     UserName
	Avatar   UserAvatar
	Language UserLanguage
}

func NewCookie() *Cookie {
	id := NewUserId()
	return &Cookie{
		Id:       id,
		Name:     UserName(id),
		Avatar:   1,
		Language: UserLanguage_Fr,
	}
}

func (c *Cookie) PlayerId() PlayerId {
	return PlayerId(c.Id)
}

func (c *Cookie) Sanitize() {
	if c.Id != "" {
		if err := c.Id.Validate(); err != nil {
			c.Id = ""
		}
	}
	if c.Name != "" {
		if err := c.Name.Validate(); err != nil {
			c.Name = ""
		}
	}
	if c.Avatar != 0 {
		if err := c.Avatar.Validate(); err != nil {
			c.Avatar = 0
		}
	}
	if c.Language != "" {
		if err := c.Language.Validate(); err != nil {
			c.Language = ""
		}
	}
}

func (c *Cookie) Validate() error {
	if c.Id == "" {
		return ErrMissingUserId
	}
	if c.Name != "" {
		if err := c.Name.Validate(); err != nil {
			return err
		}
	}
	if c.Avatar != 0 {
		if err := c.Avatar.Validate(); err != nil {
			return err
		}
	}
	if c.Language != "" {
		if err := c.Language.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Cookie) Data() map[string]any {
	if c != nil && c.Id != "" {
		return map[string]any{
			"user": c,
		}
	}
	return nil
}
