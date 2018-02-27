package relayclient

import (
	"context"
	"fmt"
	"log"
	"sync"

	floodsub "github.com/libp2p/go-floodsub"
	ci "github.com/libp2p/go-libp2p-crypto"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
)

// Group is the concrete type implementing the broadcast.Channel
// interface as well as containing the client's identity, a list of
// Group members, and a pubsub subscription along with corresponding
// buffer channels for messages going in and out of our Subscription
type Group struct {
	// implements the broadcast.Channel interface
	ctx  context.Context
	name string
	id   *Identity

	sub              *floodsub.Subscription
	incomingMessages chan *Message
	outgoingMessages chan *Message

	members []*Member  // value from on-chain
	mu      sync.Mutex // guards members

}

// Members are other staked clients that are in our group gossipsub mesh
type Member struct {
	ID     peer.ID   // libp2p concept
	PubKey ci.PubKey // on-chain identifying information
}

// Message is the information we send over the wire,
// with the raw data and signed message (Member's private key)
type Message struct {
	Sender    *Member
	Data      string
	Signature string

	seqno int
}

// Name returns the name of the group, also referenced to as the
// floodsub topic, and the hashed concatenation of all public keys
// listed in the on-chain group registry
func (g *Group) Name() string {
	return g.name
}

// TODO: get and increment a seq, use the Protocol, sign entire message
// Send handles signing a message and ensuring that it's put on the queue
// to be flushed and broadcasted to all members of the group.
func (g *Group) Send(message *Message) bool {
	msg, err := signBroadcastMessage(g.id.privKey, message)
	if err != nil {
		log.Println("Failed signing message with err: %s", err)
		return false
	}
	g.outgoingMessages <- msg

	return true
}

// RecvChan returns a processed, inbound message from a group
func (g *Group) RecvChan() <-chan *Message {
	return g.incomingMessages
}

func signBroadcastMessage(pk ci.PrivKey, message *Message) (*Message, error) {
	signed, err := pk.Sign([]byte(message.Data))
	if err != nil {
		return nil, err
	}
	message.Signature = string(signed)
	return message, nil
}

func (g *Group) handleGroupMessages(ctx context.Context, r *dht.IpfsDHT) {
	defer g.sub.Cancel()
	// TODO: obey ctx.Done()
	for {
		msg, err := g.sub.Next(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		if err := g.handleMessage(ctx, msg, r); err != nil {
			log.Println(err)
			return
		}
	}
}

func (g *Group) flushMessages(ctx context.Context, fs *floodsub.PubSub) {
	for {
		select {
		case msg := <-g.outgoingMessages:
			// TODO: send whole message, not just signature
			if err := fs.Publish(g.name, []byte(msg.Signature)); err != nil {
				log.Println("Error publishing message %#v to group %s", msg, g.name)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (g *Group) handleMessage(ctx context.Context, msg *floodsub.Message, r *dht.IpfsDHT) error {
	// TODO:
	// Step one, given the message, see who the from is
	sender := msg.GetFrom()
	fmt.Printf("SENDER: %s\n", sender)

	// step two, look up that peer in the dht
	pub, err := r.GetPublicKey(ctx, sender)
	if err != nil {
		return err
	}
	fmt.Printf("WE HAVE PUBKEY: %s\n", pub)

	// TODO: step 3, verify that the message is coming from a valid group member
	//  - now check that the peer's public key is part of the group
	// data := isSignedByGroupMember(pub, raw)

	// TODO: step 4, don't know them? add the peer
	// we added them in our magic GetPublicKey function above
	// of note, if they fail step 3, I guess we should remove them from the peerstore?

	// TODO: do this as well? n.Sub.AddPeer(peerid, floodsub.GossipSubID)
	// TODO: do we need to construct a new connection via swarm.NewStreamToPeer
	// TODO: How can I measure Peer grafting?

	// per the bradfield class, if these messages have a ttl, we need to check that?
	// where am I storing these messages? Are messages ordered? Might I have recv this before?

	// TODO: step 6, slap these messages onto our event loop or...?
	log.Printf("GOT: %+v", msg)
	log.Printf("GOT FROM: %+v", msg.GetFrom())
	log.Printf("GOT Data: %s", msg.GetData())
	log.Printf("GOT Seqno: %d", msg.GetSeqno())
	log.Printf("GOT TopicIDs: %d", msg.GetTopicIDs())

	m := &Message{}
	g.incomingMessages <- m

	return nil
}

// func isSignedByGroupMember(pub ci.PubKey, msg []byte) string {
// 	dst := make([]byte, hex.DecodedLen(len(msg)))
// 	n, err := hex.Decode(dst, msg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	pieces := strings.Split(fmt.Sprintf("%s", dst[:n]), "||")
// 	ok, err := pub.Verify([]byte(pieces[1]), []byte(pieces[0]))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if !ok {
// 		fmt.Errorf("Failed to validate signature\n")
// 	}
// 	return string(pieces[1])
// }

func (g *Group) assertGroupMembership(pub ci.PubKey) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, m := range g.members {
		if pub.Equals(m.PubKey) {
			return true
		}
	}
	return false
}
