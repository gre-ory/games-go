package model

type UserName string

func (n UserName) Validate() error {
	if n == "" {
		return ErrMissingUserName
	}
	return nil
}
