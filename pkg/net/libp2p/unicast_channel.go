package libp2p

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"sync"

	"github.com/keep-network/keep-core/pkg/net/internal"
	"github.com/keep-network/keep-core/pkg/net/key"

	"github.com/libp2p/go-libp2p-core/helpers"
	"github.com/libp2p/go-libp2p-core/peer"

	protoio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/libp2p/go-libp2p-core/network"
)

type streamFactory func(ctx context.Context, peerID peer.ID) (network.Stream, error)

type unicastChannel struct {
	clientIdentity *identity

	remotePeerID peer.ID

	streamFactory streamFactory

	messageHandlersMutex sync.Mutex
	messageHandlers      []*unicastMessageHandler

	unmarshalersMutex  sync.Mutex
	unmarshalersByType map[string]func() net.TaggedUnmarshaler
}

type unicastMessageHandler struct {
	ctx     context.Context
	channel chan net.Message
}

func (uc *unicastChannel) RemotePeerID() string {
	return uc.remotePeerID.String()
}

func (uc *unicastChannel) Send(ctx context.Context, message net.TaggedMarshaler) error {
	messageProto, err := uc.messageProto(message)
	if err != nil {
		return err
	}

	logger.Debugf(
		"[%v] sending message [%v] to peer [%v]",
		uc.clientIdentity.id,
		hex.EncodeToString(messageProto.Payload),
		uc.RemotePeerID(),
	)

	stream, err := uc.streamFactory(ctx, uc.remotePeerID)
	if err != nil {
		//TODO: consider retrying on error
		logger.Errorf("[%v] could not create stream: [%v]", uc.clientIdentity.id, err)
		return err
	}

	return uc.send(stream, messageProto)
}

func (uc *unicastChannel) send(stream network.Stream, message proto.Message) error {
	writer := bufio.NewWriter(stream)
	protoWriter := protoio.NewDelimitedWriter(writer)

	writeMsg := func(msg proto.Message) error {
		err := protoWriter.WriteMsg(msg)
		if err != nil {
			return err
		}

		return writer.Flush()
	}

	defer helpers.FullClose(stream)

	err := writeMsg(message)
	if err != nil {
		_ = stream.Reset()
		return err
	}

	return err
}

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

func (uc *unicastChannel) Recv(ctx context.Context, handler func(m net.Message)) {
	messageHandler := &unicastMessageHandler{
		ctx:     ctx,
		channel: make(chan net.Message),
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

func (uc *unicastChannel) removeHandler(handler *unicastMessageHandler) {
	uc.messageHandlersMutex.Lock()
	defer uc.messageHandlersMutex.Unlock()

	for i, h := range uc.messageHandlers {
		if h.channel == handler.channel {
			uc.messageHandlers[i] = uc.messageHandlers[len(uc.messageHandlers)-1]
			uc.messageHandlers = uc.messageHandlers[:len(uc.messageHandlers)-1]
		}
	}
}

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

func (uc *unicastChannel) handleStream(stream network.Stream) {
	if stream.Conn().RemotePeer() != uc.remotePeerID {
		logger.Warningf(
			"[%v] stream [%v] from peer [%v] is not supported "+
				"by the unicast channel of peer [%v]",
			uc.clientIdentity.id,
			stream.Protocol(),
			stream.Conn().RemotePeer(),
			uc.remotePeerID,
		)
		return
	}

	go func() {
		reader := protoio.NewDelimitedReader(stream, 1<<20)

		for {
			messageProto := new(pb.NetworkMessage)
			err := reader.ReadMsg(messageProto)
			if err != nil {
				if err != io.EOF {
					_ = stream.Reset()
				} else {
					_ = stream.Close()
				}
				return
			}

			logger.Debugf(
				"[%v] read message [%v] from peer [%v]",
				uc.clientIdentity.id,
				hex.EncodeToString(messageProto.Payload),
				uc.remotePeerID,
			)

			// Every message should be independent from any other message.
			go func(message *pb.NetworkMessage) {
				if err := uc.processMessage(message); err != nil {
					logger.Error(err)
					return
				}
			}(messageProto)
		}
	}()
}

func (uc *unicastChannel) processMessage(message *pb.NetworkMessage) error {
	// The protocol type is on the envelope; let's pull that type
	// from our map of unmarshallers.
	unmarshaled, err := uc.getUnmarshalingContainerByType(string(message.Type))
	if err != nil {
		return err
	}

	if err := unmarshaled.Unmarshal(message.GetPayload()); err != nil {
		return err
	}

	// Construct an identifier from the sender.
	senderIdentifier := &identity{}
	if err := senderIdentifier.Unmarshal(message.Sender); err != nil {
		return err
	}

	if senderIdentifier.id != uc.remotePeerID {
		return fmt.Errorf(
			"messages from peer [%v] is not supported by the "+
				"unicast channel of peer [%v]",
			senderIdentifier.id,
			uc.remotePeerID,
		)
	}

	networkKey := key.Libp2pKeyToNetworkKey(senderIdentifier.pubKey)
	if networkKey == nil {
		return fmt.Errorf(
			"sender [%v] with key [%v] is not of correct type",
			senderIdentifier.id,
			senderIdentifier.pubKey,
		)
	}

	uc.deliver(internal.BasicMessage(
		senderIdentifier.id,
		unmarshaled,
		string(message.Type),
		key.Marshal(networkKey),
	))

	return err
}

func (uc *unicastChannel) getUnmarshalingContainerByType(messageType string) (
	net.TaggedUnmarshaler,
	error,
) {
	uc.unmarshalersMutex.Lock()
	defer uc.unmarshalersMutex.Unlock()

	unmarshaler, found := uc.unmarshalersByType[messageType]
	if !found {
		return nil, fmt.Errorf(
			"couldn't find unmarshaler for type %s", messageType,
		)
	}

	return unmarshaler(), nil
}

func (uc *unicastChannel) deliver(message net.Message) {
	uc.messageHandlersMutex.Lock()
	snapshot := make([]*unicastMessageHandler, len(uc.messageHandlers))
	copy(snapshot, uc.messageHandlers)
	uc.messageHandlersMutex.Unlock()

	for _, handler := range snapshot {
		go func(message net.Message, handler *unicastMessageHandler) {
			select {
			case handler.channel <- message:
			// Nothing to do here; we block until the message is handled
			// or until the context gets closed.
			// This way we don't lose any message but also don't stay
			// with any dangling goroutines if there is no longer anyone
			// to receive messages.
			case <-handler.ctx.Done():
				return
			}
		}(message, handler)
	}
}
