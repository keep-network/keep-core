package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
)

const HomeIPv4Address = "127.0.0.1"

// GetMyIPv4Address returns this node's IPv4 IP Address
// If myPreferredOutboundIP (from config file) is valid, return that
// If more than one IP address found, call GetPreferredOutboundIP
// 127.0.0.1 will be returned if no other IPv4 addresses are found
// Assumes node has at least one interface (and the 127.0.0.1 address)
func GetMyIPv4Address(myPreferredOutboundIP string) string {
	myIPAddress := HomeIPv4Address
	ifaces, err := net.Interfaces()
	if err != nil {
		return HomeIPv4Address
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
					if myIPAddress == myPreferredOutboundIP {
						return myIPAddress // myPreferredOutboundIP is valid
					}
					if myIPAddress != HomeIPv4Address {
						myIPs = append(myIPs, myIPAddress)
					}
				}
			}
		}
	}
	if len(myIPs) > 1 {
		myIPAddress = GetPreferredOutboundIP()
	}
	if len(myPreferredOutboundIP) > 0 {
		fmt.Printf("preferred-ip-address (%s) not valid - using %s instead\n", myPreferredOutboundIP, myIPAddress)
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
