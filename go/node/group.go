package node

import (
	"context"
	"fmt"
	"log"
	"sync"

	floodsub "github.com/libp2p/go-floodsub"
	ci "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
)

// Only to be used in the core process loop
type GroupManager struct {
	Groups   map[string]*Group
	mu       sync.Mutex       // guards Group
	floodsub *floodsub.PubSub // pub/sub for group communication
	ps       pstore.Peerstore
	host     host.Host
}

type Group struct {
	name             string
	sub              *floodsub.Subscription
	members          []*Member  // value from on-chain
	mu               sync.Mutex // guards members
	incomingMessages chan *Message
}

type Member struct {
	id peer.ID
	pk ci.PubKey
}

type Message struct {
	from   peer.ID
	seqno  int
	data   string
	topics string
	msg    string
}

func NewGroupManager(fs *floodsub.PubSub, h host.Host) *GroupManager {
	return &GroupManager{Groups: make(map[string]*Group), floodsub: fs, host: h}
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
func (gm *GroupManager) BroadcastGroupMessage(ctx context.Context, pk ci.PrivKey, topic string, msg string) error {
	// TODO: get and increment a seq?
	// TODO: get some sort of protobuf structure from our datastore for the above?
	// TODO: how do I protocol?
	// TODO: sign message:
	signed, err := signBroadcastMessage(pk, []byte(msg))
	if err != nil {
		return err
	}

	gm.mu.Lock()
	if _, ok := gm.Groups[topic]; ok {
		// no group of topic exsists; create the group
		gm.mu.Unlock()
		err := gm.JoinGroup(ctx, topic)
		if err != nil {
			return err
		}
		// TODO: publish an actual message
		return gm.floodsub.Publish(topic, signed)
		// callers should wait some time for our messages to propogate
	}
	gm.mu.Unlock()
	return gm.floodsub.Publish(topic, signed)
}

// JoinGroup
func (gm *GroupManager) JoinGroup(ctx context.Context, name string) error {
	// TODO: add all members to the group via either AddPeers or
	// TODO: constructing a new connection via swarm.NewStreamToPeer

	gm.mu.Lock()
	defer gm.mu.Unlock() // FIXME: say no to fat locks
	if _, ok := gm.Groups[name]; !ok {
		g := &Group{name: name, incomingMessages: make(chan *Message, 250)}

		sub, err := gm.floodsub.Subscribe(name)
		if err != nil {
			return err
		}
		g.sub = sub
		gm.Groups[name] = g

		// TODO: if we're dumping messages on to a channel, read them in another goroutine
		go g.handleGroupMessages(ctx, gm.ps)
	}
	return nil
}

func (g *Group) handleGroupMessages(ctx context.Context, ps pstore.Peerstore) {
	defer g.sub.Cancel()
	// TODO: obey ctx.Done()
	for {
		msg, err := g.sub.Next(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		if err := g.handleMessage(msg, ps); err != nil {
			log.Println(err)
		}
	}
}

func (g *Group) handleMessage(msg *floodsub.Message, ps pstore.Peerstore) error {
	// TODO:
	// Step one, given the message, see who the from is
	// sender := msg.GetFrom()
	// fmt.Println("SENDER: ", sender)

	// TODO:
	// step two, look up that person in your peerstore
	// pinfo := ps.PeerInfo(sender)
	// fmt.Println("PINFO: ", pinfo)

	// TODO: step 3, are they a group member

	// TODO: step 4, don't know them? add the peer
	// if pinfo.Addrs == nil {
	// 	pub, err := sender.ExtractEd25519PublicKey()
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// n.Sub.AddPeer(peerid, floodsub.GossipSubID)
	// TODO: do we need to construct a new connection via swarm.NewStreamToPeer
	// TODO: How can I measure Peer grafting?

	// TODO: step 5, verify that the message is coming from a valid group member
	//  - now check that the peer's public key is part of the group
	// isSignedByGroupMember(pub, msg.data)

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

func isSignedByGroupMember(ci.PubKey, string) {}

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
