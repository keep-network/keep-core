package local

import "math/rand"

var letterRunes = [52]rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y',
	'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
	'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

type localIdentifier string

func (li localIdentifier) String() string {
	return string(li)
}

func randomIdentifier() localIdentifier {
	runes := make([]rune, 32)
	for i := range runes {
		runes[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return localIdentifier(runes)
}
