package libp2p

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"sync"

	"github.com/libp2p/go-libp2p-core/helpers"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/keep-network/keep-core/pkg/net/retransmission"

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
	messageHandlers      []*messageHandler

	unmarshalersMutex  sync.Mutex
	unmarshalersByType map[string]func() net.TaggedUnmarshaler
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
		}
	}()
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
