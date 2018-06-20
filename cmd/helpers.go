package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

// GetIPv4Address returns this node's IPv4 IP Address
// If more than one IP address found, call GetPreferredOutboundIP
// 127.0.0.1 will be returned if no other IPv4 addresses are found;
// otherwise, the non 127.0.0.1 address will be returned
// Assumes node has at least one interface (and the 127.0.0.1 address)
func GetIPv4Address() string {
	myIPAddress := "127.0.0.1"
	ifaces, err := net.Interfaces()
	if err != nil {
		return "127.0.0.1"
	}
	var myIPs []string
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			addrString := addr.String()
			ip, _, err := net.ParseCIDR(addrString)
			if err == nil {
				myIPBytes := ip.To4()
				if myIPBytes != nil {
					myIPAddress = myIPBytes.String()
					if myIPAddress != "127.0.0.1" {
						myIPs = append(myIPs, myIPAddress)
					}
				}
			}
		}
	}
	if len(myIPs) > 1 {
		myIPAddress = GetPreferredOutboundIP()
	}
	return myIPAddress
}

// GetPreferredOutboundIP gets the preferred outbound ip address
func GetPreferredOutboundIP() string {
	conn, err := net.Dial("udp", "9.9.9.9:9999")
	if err != nil {
		log.Fatal(err)
	}
	defer closeConn(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
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
