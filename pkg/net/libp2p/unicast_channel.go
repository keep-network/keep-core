package libp2p

import (
	"bufio"
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/operator"
	"io"
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/net/internal"

	"github.com/libp2p/go-libp2p-core/peer"

	protoio "github.com/gogo/protobuf/io"
	"github.com/gogo/protobuf/proto"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	"github.com/libp2p/go-libp2p-core/network"
)

const (
	readerMaxSize = 1 << 20
	sendTimeout   = 10 * time.Second
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

func (uc *unicastChannel) Send(message net.TaggedMarshaler) error {
	ctx, cancel := context.WithTimeout(context.Background(), sendTimeout)
	defer cancel()

	logger.Debugf(
		"[%v] sending message to peer [%v]",
		uc.clientIdentity.id,
		uc.remotePeerID,
	)

	streamSuccess := make(chan network.Stream)
	streamError := make(chan error)

	go func() {
		stream, err := uc.streamFactory(ctx, uc.remotePeerID)
		if err != nil {
			streamError <- err
		}
		streamSuccess <- stream
	}()

	select {
	case stream := <-streamSuccess:
		messageProto, err := uc.messageProto(message)
		if err != nil {
			return err
		}

		err = signMessage(messageProto, uc.clientIdentity.privKey)
		if err != nil {
			return err
		}

		return uc.send(stream, messageProto)
	case err := <-streamError:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
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

	defer func() {
		if err := stream.Close(); err != nil {
			logger.Errorf("could not close stream: [%v]", err)
		}
	}()

	err := writeMsg(message)
	if err != nil {
		resetErr := stream.Reset()
		if resetErr != nil {
			logger.Errorf("could not reset stream: [%v]", resetErr)
		}
		return err
	}

	return err
}

func (uc *unicastChannel) messageProto(
	message net.TaggedMarshaler,
) (*pb.UnicastNetworkMessage, error) {
	payloadBytes, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	senderIdentityBytes, err := uc.clientIdentity.Marshal()
	if err != nil {
		return nil, err
	}

	return &pb.UnicastNetworkMessage{
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
			break
		}
	}
}

func (uc *unicastChannel) SetUnmarshaler(unmarshaler func() net.TaggedUnmarshaler) {
	tpe := unmarshaler().Type()

	uc.unmarshalersMutex.Lock()
	defer uc.unmarshalersMutex.Unlock()

	uc.unmarshalersByType[tpe] = unmarshaler
}

func (uc *unicastChannel) handleStream(stream network.Stream) {
	if stream.Conn().RemotePeer() != uc.remotePeerID {
		// A unicast channel is created for a specific remote peer.
		// All streams incoming from other peers should be dropped.
		return
	}

	go func() {
		reader := protoio.NewDelimitedReader(stream, readerMaxSize)

		for {
			messageProto := new(pb.UnicastNetworkMessage)
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
				"[%v] received message from peer [%v]",
				uc.clientIdentity.id,
				uc.remotePeerID,
			)

			// Every message should be independent from any other message.
			go func(message *pb.UnicastNetworkMessage) {
				if err := uc.processMessage(message); err != nil {
					logger.Error(err)
					return
				}
			}(messageProto)
		}
	}()
}

func (uc *unicastChannel) processMessage(message *pb.UnicastNetworkMessage) error {
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
			"messages from peer [%v] are not supported by the "+
				"unicast channel for peer [%v]",
			senderIdentifier.id,
			uc.remotePeerID,
		)
	}

	if err := verifyMessageSignature(message, senderIdentifier.pubKey); err != nil {
		return err
	}

	operatorPublicKey, err := networkPublicKeyToOperatorPublicKey(senderIdentifier.pubKey)
	if err != nil {
		return fmt.Errorf(
			"sender [%v] with key [%v] is not of correct type",
			senderIdentifier.id,
			senderIdentifier.pubKey,
		)
	}

	operatorPublicKeyBytes, err := operator.MarshalUncompressed(operatorPublicKey)
	if err != nil {
		return fmt.Errorf(
			"cannot marshal sender [%v] key [%v]: [%v]",
			senderIdentifier.id,
			senderIdentifier.pubKey,
			err,
		)
	}

	uc.deliver(internal.BasicMessage(
		senderIdentifier.id,
		unmarshaled,
		string(message.Type),
		operatorPublicKeyBytes,
		uint64(0),
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
			"couldn't find unmarshaler for type [%s]",
			messageType,
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
