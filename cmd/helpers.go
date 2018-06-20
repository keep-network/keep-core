package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
)

func header(header string) {
	dashes := strings.Repeat("-", len(header))
	fmt.Printf("\n%s\n%s\n%s\n", dashes, header, dashes)
}

type testMessage struct {
	Payload string
}

// Type of this message
func (m *testMessage) Type() string {
	return "test/unmarshaler"
}

// Marshal this message
func (m *testMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal this message
func (m *testMessage) Unmarshal(bytes []byte) error {
	var message testMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		fmt.Println("hit this error")
		return err
	}
	m.Payload = message.Payload

	return nil
}
