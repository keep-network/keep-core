package relayclient

import (
	"context"
	"fmt"
	"sync"

	floodsub "github.com/libp2p/go-floodsub"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

type GroupManager struct {
	Groups map[string]*Group
	mu     sync.Mutex // guards Groups

	pubsub *floodsub.PubSub // pub/sub for group communication
	dht    *dht.IpfsDHT
	host   host.Host

	id  *Identity
	ctx context.Context
}

// NewGroupManager gives us the client's GroupManager with a new floodsub
func NewGroupManager(ctx context.Context, id *Identity, h host.Host, d *dht.IpfsDHT) (*GroupManager, error) {
	gs, err := floodsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}
	gm := &GroupManager{Groups: make(map[string]*Group), pubsub: gs, host: h, dht: d, id: id, ctx: ctx}
	return gm, nil
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

// TODO: just for demo, remove this
func (gm *GroupManager) BroadcastGroupMessage(ctx context.Context, topic string, msg *Message) error {
	signed, err := signBroadcastMessage(gm.id.privKey, msg)
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
		return gm.pubsub.Publish(topic, []byte(signed.Data))
		// callers should wait some time for our messages to propogate
	}
	gm.mu.Unlock()
	return gm.pubsub.Publish(topic, []byte(signed.Data))
}

func (gm *GroupManager) GetGroup(ctx context.Context, name string) (*Group, error) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	group, ok := gm.Groups[name]
	if !ok {
		// no group of topic exsists; create the group
		err := gm.JoinGroup(ctx, name)
		if err != nil {
			return nil, err
		}
		// callers should wait some time for our messages to propogate
		return gm.Groups[name], nil
	}
	return group, nil
}

// JoinGroup is not threadsafe - ensure it's called with a lock!
func (gm *GroupManager) JoinGroup(ctx context.Context, name string) error {
	// TODO: add all members to the group via either AddPeers or
	// TODO: constructing a new connection via swarm.NewStreamToPeer
	if _, ok := gm.Groups[name]; !ok {
		g := &Group{name: name, incomingMessages: make(chan *Message, 250)}

		sub, err := gm.pubsub.Subscribe(name)
		if err != nil {
			return err
		}
		g.sub = sub
		gm.Groups[name] = g

		go g.handleGroupMessages(ctx, gm.dht)
		go g.flushMessages(ctx, gm.pubsub)
	}
	return nil
}
