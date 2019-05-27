package local

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

	provider := Connect()
	localChannel, err := provider.ChannelFor("channel name")
	if err != nil {
		t.Fatal(err)
	}
	localChannel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &mockNetMessage{}
	})

	handlerFiredChan := make(chan struct{})
	handler := net.HandleMessageFunc{
		Type: "rambo",
		Handler: func(msg net.Message) error {
			handlerFiredChan <- struct{}{}
			return nil
		},
	}

	localChannel.Recv(handler)

	localChannel.Send(&mockNetMessage{})

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

			provider := Connect()
			localChannel, err := provider.ChannelFor("channel name")
			if err != nil {
				t.Fatal(err)
			}
			localChannel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
				return &mockNetMessage{}
			})

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

				localChannel.Recv(handler)
			}

			// Unregister specified handlers.
			for _, handlerType := range test.handlersUnregistered {
				localChannel.UnregisterRecv(handlerType)
			}

			// Send a message, all handlers should be called.
			localChannel.Send(&mockNetMessage{})

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

type mockNetMessage struct{}

func (mm *mockNetMessage) Type() string {
	return "mock_message"
}

func (mm *mockNetMessage) Marshal() ([]byte, error) {
	return []byte("some mocked bytes"), nil
}

func (mm *mockNetMessage) Unmarshal(bytes []byte) error {
	return nil
}
