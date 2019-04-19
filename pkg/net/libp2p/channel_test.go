package libp2p

import (
	"context"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
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
