package libp2p

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-core/pkg/net/security/handshake"

	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	msmux "github.com/multiformats/go-multistream"

	protoio "github.com/gogo/protobuf/io"
)

func (c *channel) InitiateRequestForNonceHandler(
	ctx context.Context,
	peerID peer.ID,
) {
	// Check to see if we have a nonce for that Peer, but that happens out of here
	fmt.Println("peerID: ", peerID)
	fmt.Println("ID: ", NonceHandshakeID)
	fmt.Println("host: ", c.p2phost)
	stream, err := c.p2phost.NewStream(ctx, peerID, NonceHandshakeID)
	if err != nil {
		fmt.Printf("failed to open stream: [%v]\n", err)
		return
	}

	fmt.Println("setting protocol")
	stream.SetProtocol(NonceHandshakeID)

	if err := msmux.SelectProtoOrFail(NonceHandshakeID, stream); err != nil {
		fmt.Printf("failed to select protocol: [%v]\n", err)
		return
	}

	// initiator station

	initiatorConnectionReader := protoio.NewDelimitedReader(stream, maxFrameSize)
	initiatorConnectionWriter := protoio.NewDelimitedWriter(stream)

	//
	// Act 1
	//

	initiatorAct1, err := handshake.InitiateHandshake()
	if err != nil {
		fmt.Printf("initiator failed initializing handshake: [%v]\n", err)
		return
	}

	act1WireMessage, err := initiatorAct1.Message().Marshal()
	if err != nil {
		fmt.Printf("initiator failed marshaling: [%v]\n", err)
		return
	}

	if err := initiatorSendAct1(
		act1WireMessage,
		initiatorConnectionWriter,
		c.clientIdentity.privKey,
		c.clientIdentity.id,
	); err != nil {
		fmt.Printf("initiator failed sending act 1: [%v]\n", err)
		return
	}

	initiatorAct2 := initiatorAct1.Next()

	//
	// Act 2
	//

	act2Message, err := initiatorReceiveAct2(initiatorConnectionReader, stream.Conn().RemotePeer())
	if err != nil {
		fmt.Printf("initiator failed receving act 2: [%v]\n", err)
		return
	}

	initiatorAct3, err := initiatorAct2.Next(act2Message)
	if err != nil {
		fmt.Printf("initiator failed changing state: [%v]\n", err)
		return
	}

	//
	// Act 3
	// Modified version of the complete handshake protocol: we just care
	// about the agreed upon challenge (which will become our nonce). We set
	// the calculated challenge as our starting nonce (with this peer), and
	// exit the protocol (no finalizing or sending off an Act3 response).

	act3Message := initiatorAct3.Message()
	if err := c.setInitiatorNonce(peerID, act3Message); err != nil {
		fmt.Printf("[%v]\n", err)
		return
	}

	go inet.FullClose(stream)
	return
}

func (c *channel) setInitiatorNonce(
	peerID peer.ID,
	act3Message *handshake.Act3Message,
) error {
	if act3Message == nil {
		return fmt.Errorf("failed to provide valid act3Message\n")
	}

	c.messageCache.nonceServiceLock.Lock()
	ns := c.messageCache.nonceService[peerID]

	ns.initial = act3Message.Nonce()
	ns.latest = act3Message.Nonce()
	ns.max = act3Message.Nonce()
	ns.used[act3Message.Nonce()] = true

	fmt.Println("Got a goddamn nonce %+v\n", act3Message.Nonce())
	c.messageCache.nonceServiceLock.Unlock()
	return nil
}

func (c *channel) setResponderNonce(
	peerID peer.ID,
	act2Message *handshake.Act2Message,
) error {
	if act2Message == nil {
		return fmt.Errorf("failed to provide valid act2Message")
	}

	c.messageCache.nonceServiceLock.Lock()
	ns := c.messageCache.nonceService[peerID]

	ns.initial = act2Message.Nonce()
	ns.latest = act2Message.Nonce()
	ns.max = act2Message.Nonce()
	ns.used[act2Message.Nonce()] = true

	c.messageCache.nonceServiceLock.Unlock()
	return nil
}

func (c *channel) respondToRequestForNonceHandler(stream inet.Stream) {
	defer inet.FullClose(stream)

	// responder station
	responderConnectionReader := protoio.NewDelimitedReader(stream, maxFrameSize)
	responderConnectionWriter := protoio.NewDelimitedWriter(stream)

	//
	// Act 1
	//

	act1Message, _, err := responderReceiveAct1(responderConnectionReader)
	if err != nil {
		fmt.Printf("responder failed receving act 1: [%v]\n", err)
	}

	responderAct2, err := handshake.AnswerHandshake(act1Message)
	if err != nil {
		fmt.Printf("responder failed parsing act 1: [%v]\n", err)
	}

	//
	// Act 2
	// Modified version of the complete handshake protocol: we just care
	// about the agreed upon challenge (which will become our nonce). We set
	// the calculated challenge as our starting nonce (with this peer), send
	// off our nonce, n2, (for the initiator to calculate the challenge), and
	// exit the protocol (no finalizing or waiting for an Act3 response).

	// This nonce is the new starting point for communications with this peer
	act2Message := responderAct2.Message()
	if err := c.setResponderNonce(
		stream.Conn().RemotePeer(),
		act2Message,
	); err != nil {
		fmt.Printf("[%v]\n", err)
	}

	act2WireMessage, err := act2Message.Marshal()
	if err != nil {
		fmt.Printf("responder failed marshaling act 2: [%v]\n", err)
	}

	if err := responderSendAct2(
		act2WireMessage,
		responderConnectionWriter,
		c.clientIdentity.privKey,
		c.clientIdentity.id,
	); err != nil {
		fmt.Printf("responder failed sending act 2: [%v]\n", err)
	}
}
