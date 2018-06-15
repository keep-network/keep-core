package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	ci "github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
)

func header(header string) {
	dashes := strings.Repeat("-", len(header))
	fmt.Printf("\n%s\n%s\n%s\n", dashes, header, dashes)
}

type libp2pMessage struct {
	//ID      *libp2p.Identity
	Payload string
}

// Type of this message
func (m *libp2pMessage) Type() string {
	return "test/unmarshaler"
}

// Marshal this message
func (m *libp2pMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal this message
func (m *libp2pMessage) Unmarshal(bytes []byte) error {
	var message libp2pMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		fmt.Println("hit this error")
		return err
	}
	m.Payload = message.Payload

	return nil
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

// PeerNetParams is a struct to bundle together the four things
// you need to run a connection with a peer: id, 2keys, and addr.
type PeerNetParams struct {
	ID      peer.ID
	PrivKey ci.PrivKey
	PubKey  ci.PubKey
	Addr    ma.Multiaddr
}

// ZeroLocalTCPAddress is the "zero" tcp local multiaddr. This means: /ip4/127.0.0.1/tcp/0
var ZeroLocalTCPAddress ma.Multiaddr

func init() {
	// initialize ZeroLocalTCPAddress
	maddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	if err != nil {
		panic(err)
	}
	ZeroLocalTCPAddress = maddr
}

// GetOutboundIP gets the preferred outbound ip
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "9.9.9.9:9999")
	if err != nil {
		log.Fatal(err)
	}
	defer closeConn(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
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
