package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	ma "github.com/multiformats/go-multiaddr"
)

// GetIPv4Address returns this node's IPv4 IP Address
// If more than one IP address found, call GetPreferredOutboundIP
// 127.0.0.1 will be returned if no other IPv4 addresses are found;
// otherwise, the non 127.0.0.1 address will be returned
// Assumes node has at least one interface (and the 127.0.0.1 address)
func GetIPv4Address(ips []ma.Multiaddr) string {
	myIPAddress := "127.0.0.1"
	var ipv4s []string
	for _, ip := range ips {
		ipAddr := ip.String()
		if strings.Contains(ipAddr, ".") && !strings.Contains(ipAddr, "127.0.0.1") {
			// Ex: ip = "/ip4/192.168.10.103/tcp/27001"
			ipv4s = append(ipv4s, strings.Split(ipAddr, "/")[2])
		}
	}
	if len(ipv4s) == 1 {
		myIPAddress = ipv4s[0]
	} else if len(ipv4s) > 1 {
		preferredIPAddress, err := GetPreferredOutboundIP()
		if err != nil {
			myIPAddress = preferredIPAddress
		}
	}
	return myIPAddress
}

// GetPreferredOutboundIP gets the preferred outbound ip address
func GetPreferredOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "9.9.9.9:9999")
	if err != nil {
		return "", err
	}
	defer closeConn(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

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

// Closable wraps Close() method
type Closable interface {
	Close() error
}

func closeConn(conn Closable) {
	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
