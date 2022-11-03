package interception

import (
	"context"

	"github.com/keep-network/keep-core/pkg/net"
)

// Rules defines the rules of intercepting network messages. Messages can be
// returned unmodified, they may be modified on the fly and they can be dropped
// by returning nil.
type Rules = func(msg net.TaggedMarshaler) net.TaggedMarshaler

// Network is the local test network implementation capable of
// intercepting network messages and modifying/dropping them based on rules
// passed to the network.
type Network interface {
	BroadcastChannelFor(name string) (net.BroadcastChannel, error)
}

// NewNetwork creates a new instance of Network interface implementation with
// message filtering rules passed as a parameter.
func NewNetwork(
	provider net.Provider,
	rules Rules,
) Network {
	return &network{
		provider: provider,
		rules:    rules,
	}
}

type network struct {
	provider net.Provider
	rules    Rules
}

func (n *network) BroadcastChannelFor(name string) (net.BroadcastChannel, error) {
	delegate, err := n.provider.BroadcastChannelFor(name)
	if err != nil {
		return nil, err
	}

	return &channel{
		delegate,
		n.rules,
	}, nil
}

type channel struct {
	delegate net.BroadcastChannel
	rules    Rules
}

func (c *channel) Name() string {
	return c.delegate.Name()
}

func (c *channel) Send(
	ctx context.Context,
	m net.TaggedMarshaler,
	strategy ...net.RetransmissionStrategy,
) error {
	altered := c.rules(m)
	if altered == nil {
		// drop the message
		return nil
	}

	return c.delegate.Send(ctx, c.rules(m))
}

func (c *channel) Recv(ctx context.Context, handler func(m net.Message)) {
	c.delegate.Recv(ctx, handler)
}

func (c *channel) SetUnmarshaler(unmarshaler func() net.TaggedUnmarshaler) {
	c.delegate.SetUnmarshaler(unmarshaler)
}

func (c *channel) SetFilter(filter net.BroadcastChannelFilter) error {
	return nil // no-op
}
