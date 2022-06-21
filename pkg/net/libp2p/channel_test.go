package libp2p

import (
	"context"
	"encoding/hex"
	"github.com/keep-network/keep-core/pkg/operator"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pubsubpb "github.com/libp2p/go-libp2p-pubsub/pb"
)

func TestRegisterAndFireHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	channel := &channel{}

	handlerFiredChan := make(chan struct{})
	channel.Recv(ctx, func(msg net.Message) {
		handlerFiredChan <- struct{}{}
	})

	channel.deliver(&mockNetMessage{})

	select {
	case <-handlerFiredChan:
		return
	case <-ctx.Done():
		t.Errorf("Expected handler not called")
	}
}

func TestUnregisterHandler(t *testing.T) {
	tests := map[string]struct {
		handlersRegistered   []string
		handlersUnregistered []string
		handlersFired        []string
	}{
		"unregister the first registered handler": {
			handlersRegistered:   []string{"a", "b", "c"},
			handlersUnregistered: []string{"a"},
			handlersFired:        []string{"b", "c"},
		},
		"unregister the last registered handler": {
			handlersRegistered:   []string{"a", "b", "c"},
			handlersUnregistered: []string{"c"},
			handlersFired:        []string{"a", "b"},
		},
		"unregister handler registered in the middle": {
			handlersRegistered:   []string{"a", "b", "c"},
			handlersUnregistered: []string{"b"},
			handlersFired:        []string{"a", "c"},
		},
		"unregister various handlers": {
			handlersRegistered:   []string{"a", "b", "c", "d", "e", "f", "g"},
			handlersUnregistered: []string{"a", "c", "f", "g"},
			handlersFired:        []string{"b", "d", "e"},
		},
		"unregister all handlers": {
			handlersRegistered:   []string{"a", "b", "c"},
			handlersUnregistered: []string{"a", "b", "c"},
			handlersFired:        []string{},
		},
	}

	for testName, test := range tests {
		test := test
		t.Run(testName, func(t *testing.T) {
			channel := &channel{}

			handlersFiredMutex := &sync.Mutex{}
			handlersFired := []string{}

			handlerCancellations := map[string]context.CancelFunc{}

			// Register all handlers. If the handler is called, append its
			// type to `handlersFired` slice.
			for _, handlerName := range test.handlersRegistered {
				handlerType := handlerName

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				handlerCancellations[handlerName] = cancel

				channel.Recv(ctx, func(msg net.Message) {
					handlersFiredMutex.Lock()
					handlersFired = append(handlersFired, handlerType)
					handlersFiredMutex.Unlock()
				})
			}

			// Cancel the specified handlers
			for _, handlerName := range test.handlersUnregistered {
				handlerCancellations[handlerName]()
			}

			// Deliver message, all handlers should be called
			channel.deliver(&mockNetMessage{})

			// Handlers are fired asynchronously; wait for them
			time.Sleep(500 * time.Millisecond)

			sort.Strings(handlersFired)
			if !reflect.DeepEqual(test.handlersFired, handlersFired) {
				t.Errorf(
					"Unexpected handlers fired\nExpected: %v\nActual:   %v\n",
					test.handlersFired,
					handlersFired,
				)
			}
		})
	}
}

func TestUnregisterWhenHandling(t *testing.T) {
	channel := &channel{}

	ctx, cancel := context.WithCancel(context.Background())

	receivedCount := 0
	stopAt := 90

	channel.Recv(ctx, func(msg net.Message) {
		receivedCount++

		if receivedCount == stopAt {
			cancel()
		}
	})

	go func() {
		for i := 0; i < 300; i++ {
			channel.deliver(&mockNetMessage{seqno: uint64(i)})
		}
	}()

	time.Sleep(500 * time.Millisecond)

	if receivedCount != stopAt {
		t.Fatalf("unexpected number of received messages: [%v]", receivedCount)
	}
}

func TestCreateTopicValidator(t *testing.T) {
	operatorPublicKeys := make([]*operator.PublicKey, 5)
	for i := range operatorPublicKeys {
		_, operatorPublicKey, _ := operator.GenerateKeyPair(DefaultCurve)
		operatorPublicKeys[i] = operatorPublicKey
	}

	authorizations := map[string]bool{
		toEncodedBytes(t, operatorPublicKeys[0]): true,
		toEncodedBytes(t, operatorPublicKeys[3]): true,
	}

	filter := func(publicKey *operator.PublicKey) bool {
		_, isAuthorized := authorizations[toEncodedBytes(t, publicKey)]
		return isAuthorized
	}

	validator := createTopicValidator(filter)

	expectedResults := []bool{true, false, false, true, false}
	for i, operatorPublicKey := range operatorPublicKeys {
		libp2pPublicKey, err := OperatorPublicKeyToLibp2pPublicKey(operatorPublicKey)
		if err != nil {
			t.Fatal(err)
		}

		authorID, _ := peer.IDFromPublicKey(libp2pPublicKey)
		authorIDBytes, _ := authorID.Marshal()
		message := &pubsubpb.Message{From: authorIDBytes}

		actualResult := validator(nil, peer.ID(rune(i)), &pubsub.Message{Message: message})

		if expectedResults[i] != actualResult {
			t.Errorf(
				"Unexpected result for public key of index [%v]\n"+
					"Expected: %v\nActual:   %v\n",
				i,
				expectedResults[i],
				actualResult,
			)
		}
	}
}

func toEncodedBytes(t *testing.T, publicKey *operator.PublicKey) string {
	publicKeyBytes, err := operator.MarshalUncompressed(publicKey)
	if err != nil {
		t.Fatal(err)
	}

	return hex.EncodeToString(publicKeyBytes)
}

type mockNetMessage struct {
	seqno uint64
}

func (mnm *mockNetMessage) TransportSenderID() net.TransportIdentifier {
	return &mockTransportIdentifier{"donald duck"}
}

func (mnm *mockNetMessage) Payload() interface{} {
	panic("not implemented in mock")
}

func (mnm *mockNetMessage) Type() string {
	panic("not implemented in mock")
}

func (mnm *mockNetMessage) SenderPublicKey() []byte {
	panic("not implemented in mock")
}

func (mnm *mockNetMessage) Seqno() uint64 {
	return mnm.seqno
}

type mockTransportIdentifier struct {
	transportID string
}

func (mti *mockTransportIdentifier) String() string {
	return mti.transportID
}
