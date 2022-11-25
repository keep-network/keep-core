package electrum

import "strings"

// Protocol is a type used for electrum protocol enumeration.
type Protocol int

// Electrum protocols enumeration.
const (
	Unknown Protocol = iota
	TCP
	SSL
)

var (
	protocolMap = map[string]Protocol{
		"tcp": TCP,
		"ssl": SSL,
	}
)

func (n Protocol) String() string {
	return []string{"unknown", "tcp", "ssl"}[n]
}

// ParseProtocol parses string.
func ParseProtocol(str string) (Protocol, bool) {
	c, ok := protocolMap[strings.ToLower(str)]
	return c, ok
}
