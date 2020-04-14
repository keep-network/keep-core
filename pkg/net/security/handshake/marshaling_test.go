package handshake

import (
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/pbutils"
)

func TestAct1MessageRoundTrip(t *testing.T) {
	message := &Act1Message{
		nonce1: 100,
	}

	unmarshaler := &Act1Message{}

	err := pbutils.RoundTrip(message, unmarshaler)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(message, unmarshaler) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzAct1MessageUnmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&Act1Message{})
}

func TestAct2MessageRoundTrip(t *testing.T) {
	challenge := [32]byte{
		100, 101, 102, 103, 104, 105, 106, 107,
		100, 101, 102, 103, 104, 105, 106, 107,
		100, 101, 102, 103, 104, 105, 106, 107,
		100, 101, 102, 103, 104, 105, 106, 107,
	}

	message := &Act2Message{
		nonce2:    100,
		challenge: challenge,
	}

	unmarshaler := &Act2Message{}

	err := pbutils.RoundTrip(message, unmarshaler)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(message, unmarshaler) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzAct2MessageUnmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&Act2Message{})
}

func TestAct3MessageRoundTrip(t *testing.T) {
	challenge := [32]byte{
		100, 101, 102, 103, 104, 105, 106, 107,
		100, 101, 102, 103, 104, 105, 106, 107,
		100, 101, 102, 103, 104, 105, 106, 107,
		100, 101, 102, 103, 104, 105, 106, 107,
	}

	message := &Act3Message{
		challenge: challenge,
	}

	unmarshaler := &Act3Message{}

	err := pbutils.RoundTrip(message, unmarshaler)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(message, unmarshaler) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzAct3MessageUnmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&Act3Message{})
}
