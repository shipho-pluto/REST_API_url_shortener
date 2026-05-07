package random

import (
	"math/rand"
	"strings"
	"time"
)

func NewRandiomString(length int) string {
	res := strings.Builder{}

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range length {
		res.WriteString(string(chars[rnd.Intn(len(chars))]))
	}

	return res.String()
}
