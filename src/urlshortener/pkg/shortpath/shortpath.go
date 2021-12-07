package shortpath

import (
	"math/rand"
	"strings"
)

var chars = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

type Generator interface {
	Generate(len int) string
}

type generator struct{}

func New(seed int64) Generator {
	rand.Seed(seed)

	return &generator{}
}

func (g *generator) Generate(length int) string {
	var sb strings.Builder

	for i := 0; i < length; i++ {
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}

	return sb.String()
}
