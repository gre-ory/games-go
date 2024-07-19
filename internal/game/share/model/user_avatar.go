package model

import (
	"fmt"
	"html/template"
)

type UserAvatar int

func (a UserAvatar) XS() template.HTML {
	return a.ExtraSmallHtml()
}

func (a UserAvatar) ExtraSmallHtml() template.HTML {
	return a.html("xs")
}

func (a UserAvatar) S() template.HTML {
	return a.SmallHtml()
}

func (a UserAvatar) SmallHtml() template.HTML {
	return a.html("s")
}

func (a UserAvatar) M() template.HTML {
	return a.MediumHtml()
}

func (a UserAvatar) MediumHtml() template.HTML {
	return a.html("m")
}

func (a UserAvatar) L() template.HTML {
	return a.LargeHtml()
}

func (a UserAvatar) LargeHtml() template.HTML {
	return a.html("s")
}

func (a UserAvatar) XL() template.HTML {
	return a.ExtraLargeHtml()
}

func (a UserAvatar) ExtraLargeHtml() template.HTML {
	return a.html("xl")
}

func (a UserAvatar) Html() template.HTML {
	return a.MediumHtml()
}

func (a UserAvatar) html(size string) template.HTML {
	if a != 0 {
		return template.HTML(fmt.Sprintf("<div class=\"avatar-%d %s\"></div>", a, size))
	}
	return ""
}

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
