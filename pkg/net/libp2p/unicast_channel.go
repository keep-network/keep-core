package libp2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/keep-network/keep-core/pkg/net/retransmission"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/libp2p/go-libp2p-core/network"
)

type unicastChannel struct {
	clientIdentity *identity

	streamMutex sync.Mutex
	stream      network.Stream

	messageHandlersMutex sync.Mutex
	messageHandlers      []*messageHandler

	unmarshalersMutex  sync.Mutex
	unmarshalersByType map[string]func() net.TaggedUnmarshaler
}

func (uc *unicastChannel) RemotePeerID() string {
	return uc.stream.Conn().RemotePeer().String()
}

// TODO: GC of unicast channel instances. Is Close() handled by the channel manager?
func (uc *unicastChannel) Close() error {
	logger.Debugf("closing unicast channel with peer [%v]", uc.RemotePeerID())
	return uc.stream.Close()
}

func (uc *unicastChannel) Send(ctx context.Context, message net.TaggedMarshaler) error {
	uc.streamMutex.Lock()
	defer uc.streamMutex.Unlock()

	messageProto, err := uc.messageProto(message)
	if err != nil {
		return err
	}

	messageBytes, err := messageProto.Marshal()
	if err != nil {
		return err
	}

	logger.Debugf(
		"[peer:%v] sending [%v] message bytes to peer [%v]",
		uc.clientIdentity.id,
		len(messageBytes),
		uc.RemotePeerID(),
	)

	return uc.send(messageBytes)
}

func (uc *unicastChannel) send(message []byte) error {
	n, err := uc.stream.Write(message)
	logger.Debugf("[peer:%v] wrote [%v] message bytes", uc.clientIdentity.id, n)
	return err
}

//FIXME: duplication with channel.go
func (uc *unicastChannel) Recv(ctx context.Context, handler func(m net.Message)) {
	messageHandler := &messageHandler{
		ctx:     ctx,
		channel: make(chan retransmission.NetworkMessage),
	}

	uc.messageHandlersMutex.Lock()
	uc.messageHandlers = append(uc.messageHandlers, messageHandler)
	uc.messageHandlersMutex.Unlock()

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Debug("context is done, removing handler")
				uc.removeHandler(messageHandler)
				return

			case msg := <-messageHandler.channel:
				// Go language specification says that if one or more of the
				// communications in the select statement can proceed, a single
				// one that will proceed is chosen via a uniform pseudo-random
				// selection.
				// Thus, it can happen this communication is called when ctx is
				// already done. Since we guarantee in the network channel API
				// that handler is not called after ctx is done (client code
				// could e.g. perform come cleanup), we need to double-check
				// the context state here.
				if messageHandler.ctx.Err() != nil {
					continue
				}

				handler(msg)
			}
		}
	}()
}

func (uc *unicastChannel) handleMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			logger.Debugf(
				"[peer:%v] processing message from peer [%v]",
				uc.clientIdentity.id,
				uc.RemotePeerID(),
			)

			// TODO: dumb implementation, we need something smarter.
			// See streaming multiple messages:
			// https://developers.google.com/protocol-buffers/docs/techniques
			buf := make([]byte, 108)
			n, err := uc.stream.Read(buf)
			if err != nil {
				logger.Error(err)
			}

			// TODO: pass message to handlers

			logger.Debugf(
				"[peer:%v] read message of [%v] bytes from peer [%v]",
				uc.clientIdentity.id,
				n,
				uc.RemotePeerID(),
			)
		}
	}
}

func (uc *unicastChannel) removeHandler(handler *messageHandler) {
	uc.messageHandlersMutex.Lock()
	defer uc.messageHandlersMutex.Unlock()

	for i, h := range uc.messageHandlers {
		if h.channel == handler.channel {
			uc.messageHandlers[i] = uc.messageHandlers[len(uc.messageHandlers)-1]
			uc.messageHandlers = uc.messageHandlers[:len(uc.messageHandlers)-1]
		}
	}
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
