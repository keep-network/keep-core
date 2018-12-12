package dkg2

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/net"
)

// Init initializes a given broadcast channel to be able to perform distributed
// key generation interactions.
func Init(channel net.BroadcastChannel) {
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.JoinMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.EphemeralPublicKeyMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.MemberCommitmentsMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.PeerSharesMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.SecretSharesAccusationsMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.MemberPublicKeySharePointsMessage{}
	})
	channel.RegisterUnmarshaler(func() net.TaggedUnmarshaler {
		return &gjkr.PointsAccusationsMessage{}
	})
}
