package cmd

import (
	"encoding/json"
	"fmt"
	ci "github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"strconv"
)


func header(header string) {
	dashes := strings.Repeat("-", len(header))
	fmt.Printf("\n%s\n%s\n%s\n", dashes, header, dashes)
}


type libp2pMessage struct {
	//ID      *libp2p.Identity
	Payload string
}

func (m *libp2pMessage) Type() string {
	return "test/unmarshaler"
}

func (m *libp2pMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *libp2pMessage) Unmarshal(bytes []byte) error {
	var message libp2pMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		fmt.Println("hit this error")
		return err
	}
	//m.ID = message.ID
	m.Payload = message.Payload

	return nil
}

//func generateDeterministicNetworkConfig() (*libp2p.Config, error) {
//	p, err := randPeerNetParamsOrFatal()
//	if err != nil {
//		return &libp2p.Config{}, err
//	}
//	pi := &libp2p.Identity{ID: libp2p.NetworkIdentity(p.ID), PrivKey: p.PrivKey, PubKey: p.PubKey}
//	return &libp2p.Config{Port: 8080, ListenAddrs: []ma.Multiaddr{p.Addr}, Identity: pi}, err
//}


func configFromPortsAndPeers(port int, peers []string)  *libp2p.Config {
	return &libp2p.Config{Port: port, Peers: peers}
}



// PeerNetParams is a struct to bundle together the four things
// you need to run a connection with a peer: id, 2keys, and addr.
type PeerNetParams struct {
	ID      peer.ID
	PrivKey ci.PrivKey
	PubKey  ci.PubKey
	Addr    ma.Multiaddr
}

//X: newPeer()
func randPeerNetParamsOrFatal() (PeerNetParams, error) {
	p, err := RandPeerNetParams()
	if err != nil {
		return PeerNetParams{}, err
	}
	return *p, nil
}

// ZeroLocalTCPAddress is the "zero" tcp local multiaddr. This means:
//   /ip4/127.0.0.1/tcp/0
var ZeroLocalTCPAddress ma.Multiaddr

func init() {
	// initialize ZeroLocalTCPAddress
	maddr, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	if err != nil {
		panic(err)
	}
	ZeroLocalTCPAddress = maddr
}

//X: KeysAndID()
func RandPeerNetParams() (*PeerNetParams, error) {
	var p PeerNetParams
	var err error
	p.Addr = ZeroLocalTCPAddress
	p.PrivKey, p.PubKey, err = RandTestKeyPair(1024)
	if err != nil {
		return nil, err
	}
	p.ID, err = peer.IDFromPublicKey(p.PubKey)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

var generatedPairs int64 = 0

func RandTestKeyPair(bits int) (ci.PrivKey, ci.PubKey, error) {
	seed := time.Now().UnixNano()

	// workaround for low time resolution
	seed += atomic.AddInt64(&generatedPairs, 1) << 32

	r := rand.New(rand.NewSource(seed))
	return ci.GenerateKeyPairWithReader(ci.RSA, bits, r)
}

type testMessage struct {
	Payload string
}

func (m *testMessage) Type() string {
	return "test/unmarshaler"
}

func (m *testMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *testMessage) Unmarshal(bytes []byte) error {
	var message testMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		fmt.Println("hit this error")
		return err
	}
	m.Payload = message.Payload

	return nil
}

func portFromMa(url string) (int, error) {
	s := strings.Split(url, "/")
	return strconv.Atoi(s[4])
}