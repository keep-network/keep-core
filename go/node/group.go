package node

import (
	"context"
	"fmt"
	"log"
	"sync"

	floodsub "github.com/libp2p/go-floodsub"
	ci "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
)

// Only to be used in the core process loop
type GroupManager struct {
	Groups   map[string]*Group
	mu       sync.Mutex       // guards Group
	floodsub *floodsub.PubSub // pub/sub for group communication
	host     host.Host
}

type Group struct {
	name             string
	sub              *floodsub.Subscription
	incomingMessages chan *Message
}

type Message struct {
}

func NewGroupManager() *GroupManager {
	return &GroupManager{}
}

func (gm *GroupManager) GetGroups() []*Group {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	var groups []*Group
	for _, group := range gm.Groups {
		groups = append(groups, group)
	}
	return groups
}

func NewGroup(name string, s *floodsub.Subscription) *Group {
	g := &Group{
		name: name,
		sub:  s,
	}
	return g
}

func (gm *GroupManager) AddGroup(g *Group) {
	gm.mu.Lock()
	if _, ok := gm.Groups[g.name]; !ok {
		gm.Groups[g.name] = g
	}
	gm.mu.Unlock()
	// TODO: kick off the group listening process
}

func (gm *GroupManager) RemoveGroup(g *Group) error {
	// TODO: ensure we cancel listening to the group
	gm.mu.Lock()
	if _, ok := gm.Groups[g.name]; !ok {
		return fmt.Errorf("group with name %s does not exist", g.name)
	}
	delete(gm.Groups, g.name)
	gm.mu.Unlock()

	return nil
}

// GroupDissolution is called when we get the DISSOVE_GROUP message in a our event loop.
// This function is responsible for unsubscribing us from our floodsub subscription, cancelling
// outbound connections to peers, deleting references and other teardown tasks.
func (gm *GroupManager) GroupDissolution(ctx context.Context) error {
	return nil
}

// TODO: this will fail as we need to rendezvous with providers by ensuring peers
// are linked before hand.
func (gm *GroupManager) BroadcastGroupMessage(ctx context.Context, pk ci.PrivKey) error {
	// we need the peer id
	// id, err := peer.IDFromPrivateKey(pk)

	// TODO: Create secret material, sign secret material, get and increment a seq?
	// TODO: get some sort of protobuf structure from our datastore for the above?
	topic := "relay/group/"

	gm.mu.Lock()
	if _, ok := gm.Groups[topic]; !ok {
		// create the group
		gm.mu.Unlock()
		gm.SubscribeToGroupMessages(ctx, topic)
		return gm.floodsub.Publish(topic, []byte(""))
		// wait some time for our messages to propogate
	}
	gm.mu.Unlock()
	return gm.floodsub.Publish(topic, []byte(""))
}

func (gm *GroupManager) SubscribeToGroupMessages(ctx context.Context, name string) error {
	// TODO: retrieve a public key for verifying messages
	gm.mu.Lock()
	defer gm.mu.Unlock()
	if _, ok := gm.Groups[name]; !ok {
		sub, err := gm.floodsub.Subscribe(name)
		if err != nil {
			return err
		}
		// create the group
		g := &Group{
			name: name,
			sub:  sub,
		}
		gm.Groups[name] = g
		go g.handleGroupMessages(ctx, nil)
	}
	// TODO: if we're dumping messages on to a channel, read them in another goroutine
	return nil
}

func (g *Group) handleGroupMessages(ctx context.Context, pub ci.PubKey) {
	defer g.sub.Cancel()
	// TODO: obey ctx.Done()
	for {
		msg, err := g.sub.Next(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		if err := g.handleMessage(msg, pub); err != nil {
			log.Println(err)
		}
	}
}

func (g *Group) handleMessage(msg *floodsub.Message, pub ci.PubKey) error {
	// handle reading the pubsub message via protobuf
	// verify that the message is coming from a valid group member
	// per the bradfield class, if these messages have a ttl, we need to check that?
	// where am I storing these messages? Are messages ordered? Might I have recv this before?
	// slap these messages onto a channel
	return nil
}
