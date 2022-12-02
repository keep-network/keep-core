package electrum

import (
	"fmt"
	"strings"
)

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

func (p Protocol) String() string {
	return []string{"unknown", "tcp", "ssl"}[p]
}

// ParseProtocol parses string.
func ParseProtocol(str string) (Protocol, bool) {
	c, ok := protocolMap[strings.ToLower(str)]
	return c, ok
}

// UnmarshalText deserializes bytes to a Protocol.
func (p *Protocol) UnmarshalText(text []byte) error {
	protocol, ok := ParseProtocol(string(text))
	if !ok {
		return fmt.Errorf("failed to parse protocol string %s", text)
	}

	*p = protocol

	return nil
}
