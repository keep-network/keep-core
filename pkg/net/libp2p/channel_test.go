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
		transportMap        map[net.TransportIdentifier]net.ProtocolIdentifier
		protocolMap         map[net.ProtocolIdentifier]net.TransportIdentifier
		expectedError       string
	}{
		"invalid transport identifier": {
			transportIdentifier: &testTransportIdentifier{},
			protocolIdentifier:  nil,
			transportMap:        make(map[net.TransportIdentifier]net.ProtocolIdentifier),
			protocolMap:         make(map[net.ProtocolIdentifier]net.TransportIdentifier),
			expectedError:       fmt.Sprintf("incorrect type for transportIdentifier: [%v]", &testTransportIdentifier{}),
		},
		"protocol identifier already exists": {
			transportIdentifier: peerID,
			protocolIdentifier:  testProto,
			transportMap: map[net.TransportIdentifier]net.ProtocolIdentifier{
				&testTransportIdentifier{}: testProto,
			},
			protocolMap:   make(map[net.ProtocolIdentifier]net.TransportIdentifier),
			expectedError: fmt.Sprintf("already have a protocol identifier in channel [%s] associated with [%v]", ch.name, peerID),
		},
		"transport identifier already exists": {
			transportIdentifier: peerID,
			protocolIdentifier:  testProto,
			transportMap:        make(map[net.TransportIdentifier]net.ProtocolIdentifier),
			protocolMap: map[net.ProtocolIdentifier]net.TransportIdentifier{
				testProto: peerID,
			},
			expectedError: fmt.Sprintf("already have a transport identifier in channel [%s] associated with [%v]", ch.name, testProto),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			ch.transportToProtoIdentifiers = test.transportMap
			ch.protoToTransportIdentifiers = test.protocolMap

			err := ch.RegisterIdentifier(test.transportIdentifier, test.protocolIdentifier)
			if err != nil && test.expectedError != err.Error() {
				t.Errorf("\ngot error: %v\nwant error: %v", err, test.expectedError)
			}
		})
	}
}

type testProtocolIdentifier struct{}
type testTransportIdentifier struct{}

var _ net.TransportIdentifier = (*testTransportIdentifier)(nil)
var _ net.ProtocolIdentifier = (*testProtocolIdentifier)(nil)

func (t *testTransportIdentifier) ProviderName() string { return "" }
