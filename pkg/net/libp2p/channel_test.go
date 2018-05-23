package libp2p

import (
	"fmt"
	"testing"

	"github.com/keep-network/keep-core/pkg/net"
	peer "github.com/libp2p/go-libp2p-peer"
)

func TestRegisterIdentifier(t *testing.T) {
	t.Parallel()

	var (
		peerID    = &identity{id: peer.ID("")}
		testProto = testProtocolIdentifier(struct{}{})
	)

	tests := map[string]struct {
		transportIdentifier net.TransportIdentifier
		protocolIdentifier  net.ProtocolIdentifier
		transportMap        map[net.TransportIdentifier]net.ProtocolIdentifier
		protocolMap         map[net.ProtocolIdentifier]net.TransportIdentifier
		expectedError       func(string) string
	}{
		"invalid transport identifier": {
			transportIdentifier: &testTransportIdentifier{},
			protocolIdentifier:  nil,
			transportMap:        make(map[net.TransportIdentifier]net.ProtocolIdentifier),
			protocolMap:         make(map[net.ProtocolIdentifier]net.TransportIdentifier),
			expectedError: func(name string) string {
				return fmt.Sprintf(
					"incorrect type for transportIdentifier: [%v] in channel [%s]",
					&testTransportIdentifier{}, name,
				)
			},
		},
		"protocol identifier already exists": {
			transportIdentifier: peerID,
			protocolIdentifier:  testProto,
			transportMap: map[net.TransportIdentifier]net.ProtocolIdentifier{
				&testTransportIdentifier{}: testProto,
			},
			protocolMap: make(map[net.ProtocolIdentifier]net.TransportIdentifier),
			expectedError: func(name string) string {
				return fmt.Sprintf(
					"protocol identifier in channel [%s] already associated with [%v]",
					name, peerID,
				)
			},
		},
		"transport identifier already exists": {
			transportIdentifier: peerID,
			protocolIdentifier:  testProto,
			transportMap:        make(map[net.TransportIdentifier]net.ProtocolIdentifier),
			protocolMap: map[net.ProtocolIdentifier]net.TransportIdentifier{
				testProto: peerID,
			},
			expectedError: func(name string) string {
				return fmt.Sprintf(
					"transport identifier in channel [%s] already associated with [%v]",
					name, testProto,
				)
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			testChannel := &channel{
				name: "test",
				transportToProtoIdentifiers: test.transportMap,
				protoToTransportIdentifiers: test.protocolMap,
			}

			err := testChannel.RegisterIdentifier(test.transportIdentifier, test.protocolIdentifier)
			if err != nil && test.expectedError(testChannel.name) != err.Error() {
				t.Errorf("\ngot error: %v\nwant error: %v", err, test.expectedError(testChannel.name))
			}
		})
	}
}

type testProtocolIdentifier struct{}
type testTransportIdentifier struct{}

var _ net.TransportIdentifier = (*testTransportIdentifier)(nil)
var _ net.ProtocolIdentifier = (*testProtocolIdentifier)(nil)

func (t *testTransportIdentifier) ProviderName() string { return "test" }
