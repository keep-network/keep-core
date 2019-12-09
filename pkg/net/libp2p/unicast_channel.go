package libp2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/libp2p/go-libp2p-core/network"
)

type unicastChannel struct {
	clientIdentity *identity

	closeFn context.CancelFunc

	streamMutex sync.Mutex
	stream      network.Stream

	messageHandlersMutex sync.Mutex
	messageHandlers      []net.HandleMessageFunc

	unmarshalersMutex  sync.Mutex
	unmarshalersByType map[string]func() net.TaggedUnmarshaler
}

func newUnicastChannel(
	ctx context.Context,
	clientIdentity *identity,
	stream network.Stream,
) *unicastChannel {
	_, cancelFn := context.WithCancel(ctx)

	return &unicastChannel{
		clientIdentity:     clientIdentity,
		closeFn:            cancelFn,
		stream:             stream,
		messageHandlers:    make([]net.HandleMessageFunc, 0),
		unmarshalersByType: make(map[string]func() net.TaggedUnmarshaler),
	}
}

func (uc *unicastChannel) Close() error {
	logger.Warning("Closing channel")
	uc.closeFn()
	return uc.stream.Close()
}

func (uc *unicastChannel) Send(message net.TaggedMarshaler) error {
	uc.streamMutex.Lock()
	defer uc.streamMutex.Unlock()

	// Transform net.TaggedMarshaler to a protobuf message
	messageProto, err := uc.messageProto(message)
	if err != nil {
		return err
	}

	messageBytes, err := messageProto.Marshal()
	if err != nil {
		return err
	}

	logger.Warningf("[%v] sending [%v] message bytes to [%v]", uc.clientIdentity.id, len(messageBytes), uc.stream.Conn().RemotePeer())
	return uc.send(messageBytes)
}

func (uc *unicastChannel) send(message []byte) error {
	n, err := uc.stream.Write(message)
	logger.Warningf("[%v] wrote [%v] message bytes", uc.clientIdentity, n)
	//uc.stream.Close()
	return err
}

//FIXME: duplication with channel.go
func (uc *unicastChannel) Recv(handler net.HandleMessageFunc) error {
	uc.messageHandlersMutex.Lock()
	uc.messageHandlers = append(uc.messageHandlers, handler)
	uc.messageHandlersMutex.Unlock()

	return nil
}

func (uc *unicastChannel) handleMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Warning("Context cancelled")
			return
		default:
			logger.Warningf("[%v] Reading...", uc.clientIdentity.id)
			// TODO: dumb implementation, we need something smarter.
			// See streaming multiple messages:
			// https://developers.google.com/protocol-buffers/docs/techniques
			buf := make([]byte, 108)
			n, err := uc.stream.Read(buf)
			if err != nil {
				logger.Error(err)
			}
			logger.Warningf("Read [%v] bytes...", n)
		}
	}
}

//FIXME: duplication with channel.go
func (uc *unicastChannel) UnregisterRecv(handlerType string) error {
	uc.messageHandlersMutex.Lock()
	defer uc.messageHandlersMutex.Unlock()

	removedCount := 0

	// updated slice shares the same backing array and capacity as the original,
	// so the storage is reused for the filtered slice.
	updated := uc.messageHandlers[:0]

	for _, mh := range uc.messageHandlers {
		if mh.Type != handlerType {
			updated = append(updated, mh)
		} else {
			removedCount++
		}
	}

	uc.messageHandlers = updated[:len(uc.messageHandlers)-removedCount]

	return nil
}

//FIXME: duplication with channel.go
func (uc *unicastChannel) RegisterUnmarshaler(unmarshaler func() net.TaggedUnmarshaler) error {
	tpe := unmarshaler().Type()

	uc.unmarshalersMutex.Lock()
	defer uc.unmarshalersMutex.Unlock()

	if _, exists := uc.unmarshalersByType[tpe]; exists {
		return fmt.Errorf("type %s already has an associated unmarshaler", tpe)
	}

	uc.unmarshalersByType[tpe] = unmarshaler
	return nil
}

//FIXME: duplication with channel.go
func (uc *unicastChannel) messageProto(
	message net.TaggedMarshaler,
) (*pb.NetworkMessage, error) {
	payloadBytes, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	senderIdentityBytes, err := uc.clientIdentity.Marshal()
	if err != nil {
		return nil, err
	}

	return &pb.NetworkMessage{
		Payload: payloadBytes,
		Sender:  senderIdentityBytes,
		Type:    []byte(message.Type()),
	}, nil
}
