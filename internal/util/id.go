package util

import (
	"github.com/jaevor/go-nanoid"
)

var tokenAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
var GenerateUserToken = Must(nanoid.CustomASCII(tokenAlphabet, 10))
