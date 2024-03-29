package model

type UserLanguage string

var (
	Fr_UserLanguage UserLanguage = "fr"
	En_UserLanguage UserLanguage = "en"
)

func (l UserLanguage) Validate() error {
	if l == "" {
		return ErrMissingUserLanguage
	}
	switch l {
	case Fr_UserLanguage:
		return nil
	case En_UserLanguage:
		return nil
	}
	return ErrInvalidUserLanguage
}
