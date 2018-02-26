package node

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/dfinity/go-dfinity-crypto/bls"
	floodsub "github.com/libp2p/go-floodsub"
	ci "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
)

// Only to be used in the core process loop
type GroupManager struct {
	Groups map[string]*Group
	mu     sync.Mutex // guards Group

	pubsub *floodsub.PubSub // pub/sub for group communication
	dht    *dht.IpfsDHT
	host   host.Host
}

type Group struct {
	name string

	members []*Member  // value from on-chain
	mu      sync.Mutex // guards members

	sub              *floodsub.Subscription
	incomingMessages chan *Message
}

type Member struct {
	ID  peer.ID   // libp2p concept
	PK  ci.PubKey // on-chain identifying information
	BLS bls.ID
}

type Message struct {
	From     *Member
	Receiver *Member
	Data     string

	seqno int
}

func NewGroupManager(fs *floodsub.PubSub, h host.Host, d *dht.IpfsDHT) *GroupManager {
	return &GroupManager{Groups: make(map[string]*Group), pubsub: fs, host: h, dht: d}
}

// TODO:
func (gm *GroupManager) SyncActiveGroups() []*Group {
	// FIXME: make a network call to ethereum, get group registry, unmarshal it,
	// make our updates to gm.Groups
	return nil
}

// GetActiveGroups allows us to identify which and how many groups
// a given staker belongs to. This is our local cache. SyncActiveGroups
// is responsible for getting new groups and memeberships.
func (gm *GroupManager) GetActiveGroups() []*Group {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	var groups []*Group
	for _, group := range gm.Groups {
		groups = append(groups, group)
	}
	return groups
}

// RemoveGroup is a convenience function that untethers the
// connection bettween the client and a group by removing it from the
// known groups list, closing inbound messages channel, and unsubscribing from
// gossip sub.
func (gm *GroupManager) RemoveGroup(name string) error {
	gm.mu.Lock()
	group, ok := gm.Groups[name]
	if !ok {
		return fmt.Errorf("group with name %s does not exist", name)
	}

	// TODO: kill inbound messages list with topic as group name
	group.sub.Cancel()
	delete(gm.Groups, group.name)

	gm.mu.Unlock()

	return nil
}

// GroupDissolution is called when we get the DISSOVE_GROUP message in a our event loop.
// This function is responsible for unsubscribing us from our floodsub subscription, cancelling
// outbound connections to peers, deleting references and other teardown tasks.
func (gm *GroupManager) GroupDissolution(ctx context.Context, name string) error {
	return gm.RemoveGroup(name)
}

// TODO: Will this fail as we need (?) to rendezvous with providers by ensuring peers
// are linked before hand?
func (gm *GroupManager) BroadcastGroupMessage(ctx context.Context, pk ci.PrivKey, topic string, msg []byte) error {
	// TODO: get and increment a seq?
	// TODO: get some sort of protobuf structure from our datastore for the above?
	// TODO: how do I protocol?
	// TODO: sign message:
	signed, err := signBroadcastMessage(pk, msg)
	if err != nil {
		return err
	}
	// TODO: misnamed - maybe call hand over proofs?
	bmsg := signAndHash(signed, msg)

	gm.mu.Lock()
	if _, ok := gm.Groups[topic]; ok {
		// no group of topic exsists; create the group
		gm.mu.Unlock()
		err := gm.JoinGroup(ctx, topic)
		if err != nil {
			return err
		}
		// TODO: publish an actual message
		return gm.pubsub.Publish(topic, bmsg)
		// callers should wait some time for our messages to propogate
	}
	gm.mu.Unlock()
	return gm.pubsub.Publish(topic, signed)
}

func signAndHash(sign []byte, data []byte) []byte {
	concat := fmt.Sprintf("%s||%s", sign, data)
	return []byte(fmt.Sprintf("%x", concat))
}

// JoinGroup
func (gm *GroupManager) JoinGroup(ctx context.Context, name string) error {
	// TODO: add all members to the group via either AddPeers or
	// TODO: constructing a new connection via swarm.NewStreamToPeer

	gm.mu.Lock()
	defer gm.mu.Unlock() // FIXME: say no to fat locks
	if _, ok := gm.Groups[name]; !ok {
		g := &Group{name: name, incomingMessages: make(chan *Message, 250)}

		sub, err := gm.pubsub.Subscribe(name)
		if err != nil {
			return err
		}
		g.sub = sub
		gm.Groups[name] = g

		// TODO: if we're dumping messages on to a channel, read them in another goroutine
		go g.handleGroupMessages(ctx, gm.dht)
	}
	return nil
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

	raw := msg.GetData()
	// TODO: step 3, verify that the message is coming from a valid group member
	//  - now check that the peer's public key is part of the group
	data := isSignedByGroupMember(pub, raw)
	fmt.Println("WE HAVE DATA: %s", data)

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
	return nil
}

func signBroadcastMessage(pk ci.PrivKey, msg []byte) ([]byte, error) {
	return pk.Sign(msg)
}

func isSignedByGroupMember(pub ci.PubKey, msg []byte) string {
	dst := make([]byte, hex.DecodedLen(len(msg)))
	n, err := hex.Decode(dst, msg)
	if err != nil {
		log.Fatal(err)
	}
	pieces := strings.Split(fmt.Sprintf("%s", dst[:n]), "||")
	ok, err := pub.Verify([]byte(pieces[1]), []byte(pieces[0]))
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		fmt.Errorf("Failed to validate signature\n")
	}
	return string(pieces[1])
}

func (g *Group) assertGroupMembership(pub ci.PubKey) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, m := range g.members {
		if pub.Equals(m.pk) {
			return true
		}
	}
	return false
}
