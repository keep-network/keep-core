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
	ChannelFor(name string) (net.BroadcastChannel, error)
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

func (n *network) ChannelFor(name string) (net.BroadcastChannel, error) {
	delegate, err := n.provider.ChannelFor(name)
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

func (c *channel) Send(m net.TaggedMarshaler, ctx ...context.Context) error {
	altered := c.rules(m)
	if altered == nil {
		// drop the message
		return nil
	}

	return c.delegate.Send(c.rules(m), ctx...)
}

func (c *channel) Recv(h net.HandleMessageFunc) error {
	return c.delegate.Recv(h)
}

func (c *channel) UnregisterRecv(handlerType string) error {
	return c.delegate.UnregisterRecv(handlerType)
}

func (c *channel) RegisterUnmarshaler(
	unmarshaler func() net.TaggedUnmarshaler,
) error {
	return c.delegate.RegisterUnmarshaler(unmarshaler)
}

func (c *channel) AddFilter(filter net.BroadcastChannelFilter) error {
	return nil // no-op
}
