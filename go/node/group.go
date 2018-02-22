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
	incomingMessages chan *Message
}

type Message struct {
	from   peer.ID
	seqno  int
	data   string
	topics string
}

func NewGroupManager() *GroupManager {
	return &GroupManager{}
}

// GetActiveGroups allows us to identify which and how many groups
// a given staker belongs to.
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
func (gm *GroupManager) BroadcastGroupMessage(ctx context.Context, pk ci.PrivKey) error {
	// TODO: Create secret material, sign secret material, get and increment a seq?
	// TODO: get some sort of protobuf structure from our datastore for the above?
	topic := "relay/group/"

	gm.mu.Lock()
	if _, ok := gm.Groups[topic]; !ok {
		// no group of topic exsists; create the group
		gm.mu.Unlock()
		err := gm.JoinGroup(ctx, topic)
		if err != nil {
			return err
		}
		// TODO: publish an actual message
		return gm.floodsub.Publish(topic, []byte(""))
		// callers should wait some time for our messages to propogate
	}
	gm.mu.Unlock()
	return gm.floodsub.Publish(topic, []byte(""))
}

// JoinGroup
func (gm *GroupManager) JoinGroup(ctx context.Context, name string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
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

	// look up that person in your peerstore
	// pinfo := ps.PeerInfo(sender)

	// don't know them? add the peer
	// if pinfo.Addrs == nil {
	// 	pub, err := sender.ExtractEd25519PublicKey()
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// TODO: swarm.NewStreamToPeer
	// TODO: How can I measure Peer grafting?

	// verify that the message is coming from a valid group member
	//  - now check that the peer's public key is part of the group

	// per the bradfield class, if these messages have a ttl, we need to check that?
	// where am I storing these messages? Are messages ordered? Might I have recv this before?

	// slap these messages onto our event loop
	return nil
}

func (g *Group) messageListener(ctx context.Context) {
	for {

		select {
		case <-ctx.Done():
			log.Println("group untethered, shuttingdown")
			return

		}
	}

}
