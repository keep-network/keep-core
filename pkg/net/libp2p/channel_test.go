package libp2p

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/keep-network/keep-core/pkg/net"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pubsubpb "github.com/libp2p/go-libp2p-pubsub/pb"
)

func TestRegisterAndFireHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	channel := &channel{}

	handlerFiredChan := make(chan struct{})
	handler := net.HandleMessageFunc{
		Type: "theboar",
		Handler: func(msg net.Message) error {
			handlerFiredChan <- struct{}{}
			return nil
		},
	}
	channel.Recv(handler)

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
		"unregister two first registered handlers with the same type": {
			handlersRegistered:   []string{"a", "a", "b", "c", "d"},
			handlersUnregistered: []string{"a"},
			handlersFired:        []string{"b", "c", "d"},
		},
		"unregister two last registered handlers with the same type": {
			handlersRegistered:   []string{"a", "b", "c", "d", "d"},
			handlersUnregistered: []string{"d"},
			handlersFired:        []string{"a", "b", "c"},
		},
		"unregister various handlers with the same type": {
			handlersRegistered:   []string{"a", "f", "b", "e", "c", "f", "e"},
			handlersUnregistered: []string{"e", "f"},
			handlersFired:        []string{"a", "b", "c"},
		},
		"unregister handler not previously registered": {
			handlersRegistered:   []string{"a", "b", "c"},
			handlersUnregistered: []string{"z"},
			handlersFired:        []string{"a", "b", "c"},
		},
	}

	for testName, test := range tests {
		test := test
		t.Run(testName, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			channel := &channel{}

			handlersFiredMutex := &sync.Mutex{}
			handlersFired := []string{}

			// Register all handlers. If the handler is called, append its
			// type to `handlersFired` slice.
			for _, handlerType := range test.handlersRegistered {
				handlerType := handlerType
				handler := net.HandleMessageFunc{
					Type: handlerType,
					Handler: func(msg net.Message) error {
						handlersFiredMutex.Lock()
						handlersFired = append(handlersFired, handlerType)
						handlersFiredMutex.Unlock()
						return nil
					},
				}

				channel.Recv(handler)
			}

			// Unregister specified handlers.
			for _, handlerType := range test.handlersUnregistered {
				channel.UnregisterRecv(handlerType)
			}

			// Deliver message, all handlers should be called.
			channel.deliver(&mockNetMessage{})

			// Handlers are fired asynchronously; wait for them.
			<-ctx.Done()

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

func TestCreateTopicValidator(t *testing.T) {
	publicKeys := make([]crypto.PubKey, 5)
	for i := range publicKeys {
		_, publicKey, _ := crypto.GenerateSecp256k1Key(rand.Reader)
		publicKeys[i] = publicKey
	}

	authorizations := map[string]bool{
		toEncodedBytes(toEcdsaPublicKey(publicKeys[0])): true,
		toEncodedBytes(toEcdsaPublicKey(publicKeys[3])): true,
	}

	filter := func(publicKey *ecdsa.PublicKey) bool {
		_, isAuthorized := authorizations[toEncodedBytes(publicKey)]
		return isAuthorized
	}

	validator := createTopicValidator(filter)

	expectedResults := []bool{true, false, false, true, false}
	for i, publicKey := range publicKeys {
		authorID, _ := peer.IDFromPublicKey(publicKey)
		authorIDBytes, _ := authorID.Marshal()
		message := &pubsubpb.Message{From: authorIDBytes}

		actualResult := validator(nil, peer.ID(i), &pubsub.Message{Message: message})

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

func toEcdsaPublicKey(publicKey crypto.PubKey) *ecdsa.PublicKey {
	secp256k1PublicKey, _ := publicKey.(*crypto.Secp256k1PublicKey)
	return (*btcec.PublicKey)(secp256k1PublicKey).ToECDSA()
}

func toEncodedBytes(publicKey *ecdsa.PublicKey) string {
	return hex.EncodeToString(
		elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y),
	)
}

type mockNetMessage struct {
}

func (mnm *mockNetMessage) TransportSenderID() net.TransportIdentifier {
	panic("not implemented in mock")
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
