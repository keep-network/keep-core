package node

import (
	"fmt"
	"sync"

	floodsub "github.com/libp2p/go-floodsub"
	host "github.com/libp2p/go-libp2p-host"
)

// Only to be used in the core process loop
type GroupManager struct {
	Groups   map[string]*Group
	mu       sync.Mutex // guards Grouop
	floodsub *floodsub.PubSub
	host     host.Host
}

type Group struct {
	name string
	sub  *floodsub.Subscription
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

// TODO: this will fail as we need to rendezvous with providers by ensuring peers
// are linked before hand.
// func (g *GroupManager) BroadcastGroupMessage(ctx context.Context, pk ci.PrivKey) error {
// 	// we need the peer id
// 	id, err := peer.IDFromPrivateKey(pk)

// 	// TODO: Create secret material, sign secret material, get and increment a seq?
// 	// TODO: get some sort of protobuf structure from our datastore for the above?
// 	// topic := "relay/group/" + g.name

// 	g.mu.Lock()
// 	// if _, ok := g.Groups[topic]
// 	// get a group
// 	// if we have the group
// 	g.mu.Unlock()
// }

// func (g *) broadcastToListeners(ctx context.Context, host host.Host, name string) {
// topic := "keep:" + name
// 	// use multihash to has the topic
// }
