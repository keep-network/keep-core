package libp2p

import (
	"fmt"
	"os"
	"testing"

	"github.com/keep-network/keep-core/pkg/net"
	peer "github.com/libp2p/go-libp2p-peer"
)

func TestMain(m *testing.M) {
	// TODO: Invoke code that builds up state
	code := m.Run()
	// TODO: close open conns
	os.Exit(code)
}

func TestRegisterIdentifier(t *testing.T) {
	t.Parallel()

	var (
		ch        = &channel{name: "test"}
		peerID    = &peerIdentifier{id: peer.ID("")}
		testProto = testProtocolIdentifier(struct{}{})
	)

	tests := map[string]struct {
		transportIdentifier net.TransportIdentifier
		protocolIdentifier  net.ProtocolIdentifier
		tMap                map[net.TransportIdentifier]net.ProtocolIdentifier
		pMap                map[net.ProtocolIdentifier]net.TransportIdentifier
		errorString         string
	}{
		"invalid transport identifier": {
			transportIdentifier: &testTransportIdentifier{},
			protocolIdentifier:  nil,
			tMap:                nil,
			pMap:                nil,
			errorString:         fmt.Sprintf("incorrect type for transportIdentifier: [%v]", &testTransportIdentifier{}),
		},
		"protocol identifier already exists": {
			transportIdentifier: peerID,
			protocolIdentifier:  testProto,
			tMap: map[net.TransportIdentifier]net.ProtocolIdentifier{
				&testTransportIdentifier{}: testProto,
			},
			pMap:        nil,
			errorString: fmt.Sprintf("already have a protocol identifier in channel [%s] associated with [%v]", ch.name, peerID),
		},
		"transport identifier already exists": {
			transportIdentifier: peerID,
			protocolIdentifier:  testProto,
			tMap:                nil,
			pMap: map[net.ProtocolIdentifier]net.TransportIdentifier{
				testProto: peerID,
			},
			errorString: fmt.Sprintf("already have a transport identifier in channel [%s] associated with [%v]", ch.name, testProto),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.tMap != nil {
				ch.transportToProtoIdentifiers = tt.tMap
			} else {
				ch.transportToProtoIdentifiers = make(map[net.TransportIdentifier]net.ProtocolIdentifier)
			}
			if tt.pMap != nil {
				ch.protoToTransportIdentifiers = tt.pMap
			} else {
				ch.protoToTransportIdentifiers = make(map[net.ProtocolIdentifier]net.TransportIdentifier)
			}
			err := ch.RegisterIdentifier(tt.transportIdentifier, tt.protocolIdentifier)
			if err != nil && tt.errorString != err.Error() {
				t.Errorf("\ngot: %v\nwant: %v", err, tt.errorString)
			}
		})
	}
}

type testProtocolIdentifier struct{}
type testTransportIdentifier struct{}

var _ net.TransportIdentifier = (*testTransportIdentifier)(nil)
var _ net.ProtocolIdentifier = (*testProtocolIdentifier)(nil)

func (t *testTransportIdentifier) ProviderName() string { return "" }
