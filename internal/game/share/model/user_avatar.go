package model

type UserAvatar int

func (a UserAvatar) Validate() error {
	if a == 0 {
		return ErrMissingUserAvatar
	}
	if 1 <= a && a <= 20 {
		return nil
	}
	return ErrInvalidUserAvatar
}

func GetAvailableAvatars() [][]UserAvatar {
	result := make([][]UserAvatar, 0, 5)
	for i := 0; i < 4; i++ {
		result = append(result, []UserAvatar{
			UserAvatar((5 * i) + 1),
			UserAvatar((5 * i) + 2),
			UserAvatar((5 * i) + 3),
			UserAvatar((5 * i) + 4),
			UserAvatar((5 * i) + 5),
		})
	}
	return result
}
