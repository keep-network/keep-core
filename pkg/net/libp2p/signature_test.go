package libp2p

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
	libp2pcrypto "github.com/libp2p/go-libp2p-crypto"
)

func TestVerifyMessageSignature(t *testing.T) {
	identity, err := newTestIdentity()

	ch := &channel{
		clientIdentity: identity,
	}

	msg := []byte("It's not much of a tail, but I'm sort of attached to it.")

	signature, err := ch.sign(msg)
	if err != nil {
		t.Fatal(err)
	}

	if err := ch.verify(identity.id, msg, signature); err != nil {
		t.Fatal(err)
	}
}

// Check if a signature created with a key other than the expected
// is considered as incorrect.
func TestDetectUnexpectedMessageSignature(t *testing.T) {
	identity, err := newTestIdentity()

	ch := &channel{
		clientIdentity: identity,
	}

	msg := []byte("It's not much of a tail, but I'm sort of attached to it.")

	signature, err := ch.sign(msg)
	if err != nil {
		t.Fatal(err)
	}

	anotherIdentity, err := newTestIdentity()
	if err != nil {
		t.Fatal(err)
	}

	err = ch.verify(anotherIdentity.id, msg, signature)
	if err == nil {
		t.Fatal("signature validation should fail")
	}

	if !strings.HasPrefix(err.Error(), "invalid signature") {
		t.Fatalf("error other than expected: %v", err)
	}
}

// Check if a malformed signature is considered incorrect
func TestDetectMalformedMessageSignature(t *testing.T) {
	identity, err := newTestIdentity()

	ch := &channel{
		clientIdentity: identity,
	}

	msg := []byte("It's not much of a tail, but I'm sort of attached to it.")

	signature := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	err = ch.verify(identity.id, msg, signature)
	if err == nil {
		t.Fatal("signature validation should fail")
	}

	if !strings.HasPrefix(err.Error(), "failed to verify signature") {
		t.Fatalf("error other than expected: %v", err)
	}
}

// Integration test simulating malicious adversary tampering the network message
// put into the channel. There are two messages sent:
// - one with a valid signature evaluated with sender's key
// - one with a valid signature evaluated with key other than sender's key
// The first message should be properly delivered, the second message should get
// rejected.
func TestRejectMessageWithUnexpectedSignature(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	staticKey, err := key.GenerateStaticNetworkKey(crand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	ch, err := createTestChannel(ctx, staticKey)
	if err != nil {
		t.Fatal(err)
	}

	if err := registerIdentity(ch, staticKey, "honest"); err != nil {
		t.Fatal(err)
	}

	honestPayload := "I did know once, only I've sort of forgotten."
	maliciousPayload := "You never can tell with bees."

	// Create and publish message with a correct signature...
	envelope, err := ch.sealEnvelope(nil, &testMessage{Payload: honestPayload})
	if err != nil {
		t.Fatal(err)
	}

	envelopeBytes, err := envelope.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	ch.pubsub.Publish(ch.name, envelopeBytes)

	// Create and publish message with a signature created with other key than
	// sender's...
	adversaryKey, err := key.GenerateStaticNetworkKey(crand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	if err := registerIdentity(ch, adversaryKey, "malicious"); err != nil {
		t.Fatal(err)
	}

	envelope, err = ch.sealEnvelope(nil, &testMessage{Payload: maliciousPayload})
	if err != nil {
		t.Fatal(err)
	}

	adversarySignature, err := adversaryKey.Sign(envelope.Message)
	if err != nil {
		t.Fatal(err)
	}
	envelope.Signature = adversarySignature

	envelopeBytes, err = envelope.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	ch.pubsub.Publish(ch.name, envelopeBytes)

	var (
		honestRecvChan    = make(chan net.Message)
		maliciousRecvChan = make(chan net.Message)
	)

	filterBySender := func(msg net.Message) error {
		protocolIdentifier, ok := msg.ProtocolSenderID().(*protocolIdentifier)
		if !ok {
			return fmt.Errorf(
				"expected: type *protocolIdentifier\ngot:   type [%v]",
				protocolIdentifier,
			)
		}
		if protocolIdentifier.id == "malicious" {
			maliciousRecvChan <- msg
		}
		if protocolIdentifier.id == "honest" {
			honestRecvChan <- msg
		}
		return nil
	}

	ensureNonMaliciousMessage := func(t *testing.T, msg net.Message) error {
		testPayload, ok := msg.Payload().(*testMessage)
		if !ok {
			return fmt.Errorf(
				"expected: payload type string\ngot:   payload type [%v]",
				testPayload,
			)
		}

		if honestPayload != testPayload.Payload {
			return fmt.Errorf(
				"expected: message payload [%s]\ngot:   payload [%s]",
				honestPayload,
				testPayload.Payload,
			)
		}
		return nil
	}

	// Check if the message with correct signature has been properly delivered
	// and if the message with incorrect signature has been dropped...
	if err := ch.Recv(net.HandleMessageFunc{
		Type: "test",
		Handler: func(msg net.Message) error {
			return filterBySender(msg)
		},
	}); err != nil {
		t.Fatal(err)
	}

	honestMessageDelivered := false

	for {
		select {
		case msg := <-maliciousRecvChan:
			if err := ensureNonMaliciousMessage(t, msg); err != nil {
				t.Fatal(err)
			}
		case msg := <-honestRecvChan:
			if err := ensureNonMaliciousMessage(t, msg); err != nil {
				t.Fatal(err)
			}

			honestMessageDelivered = true

			// Ensure all messages are flushed before exiting
			time.Sleep(500 * time.Millisecond)
			return
		case <-ctx.Done():
			if !honestMessageDelivered {
				t.Fatal("expected message not delivered")
			}

			return
		}
	}
}

func registerIdentity(
	ch *channel,
	staticKey libp2pcrypto.PrivKey,
	protocolID string,
) error {
	identity, err := createIdentity(staticKey)
	if err != nil {
		return err
	}

	if err := ch.RegisterIdentifier(
		networkIdentity(identity.id),
		&protocolIdentifier{id: protocolID},
	); err != nil {
		return err
	}

	return nil
}

// createTestChannel creates and initializes `BroadcastChannel` with all
// underlying libp2p setup steps. Created instance is then casted to
// `lib2p.channel` type so the private interface is available and can be
// tested.
func createTestChannel(
	ctx context.Context,
	staticKey *key.StaticNetworkKey,
) (*channel, error) {
	networkConfig := Config{Port: 8080}

	provider, err := Connect(ctx, networkConfig, staticKey)
	if err != nil {
		return nil, err
	}

	broadcastChannel, err := provider.ChannelFor("testchannel")
	if err != nil {
		return nil, err
	}

	if err := broadcastChannel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &testMessage{} },
	); err != nil {
		return nil, err
	}

	ch, ok := broadcastChannel.(*channel)
	if !ok {
		return nil, fmt.Errorf("unexpected channel type")
	}

	return ch, nil
}
