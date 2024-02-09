package util

import (
	"github.com/jaevor/go-nanoid"
)

var idAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var GenerateId = Must(nanoid.CustomASCII(idAlphabet, 4))
