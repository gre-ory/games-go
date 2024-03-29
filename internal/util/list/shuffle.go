package list

import (
	"math/rand"
)

func Shuffle[T any](items []T) {
	rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })
}
