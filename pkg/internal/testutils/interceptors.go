package testutils

import (
	"github.com/keep-network/keep-core/pkg/net"
)

// NetworkMessageInterceptor defines a rules of intercepting network messages.
// Messages can be returned unmodified, they may be modified on the fly and they
// can be dropped by returning nil.
type NetworkMessageInterceptor = func(msg net.TaggedMarshaler) net.TaggedMarshaler

// InterceptingNetwork is the local test network implementation capable of
// intercepting network messages and modifying/dropping them based on rules
// passed to the network.
type InterceptingNetwork interface {
	ChannelFor(name string) (net.BroadcastChannel, error)
}

// NewInterceptingNetwork creates a new instance of InterceptingNetwork
// interface implementation with message filtering rules passed as a parameter.
func NewInterceptingNetwork(
	provider net.Provider,
	interceptor NetworkMessageInterceptor,
) InterceptingNetwork {
	return &interceptingNetwork{
		provider:    provider,
		interceptor: interceptor,
	}
}

type interceptingNetwork struct {
	provider    net.Provider
	interceptor NetworkMessageInterceptor
}

func (in *interceptingNetwork) ChannelFor(name string) (net.BroadcastChannel, error) {
	delegate, err := in.provider.ChannelFor(name)
	if err != nil {
		return nil, err
	}

	return &interceptingChannel{
		delegate,
		in.interceptor,
	}, nil
}

type interceptingChannel struct {
	delegate    net.BroadcastChannel
	interceptor NetworkMessageInterceptor
}

func (ic *interceptingChannel) Name() string {
	return ic.delegate.Name()
}

func (ic *interceptingChannel) Send(m net.TaggedMarshaler) error {
	altered := ic.interceptor(m)
	if altered == nil {
		// drop the message
		return nil
	}

	return ic.delegate.Send(ic.interceptor(m))
}

func (ic *interceptingChannel) Recv(h net.HandleMessageFunc) error {
	return ic.delegate.Recv(h)
}

func (ic *interceptingChannel) UnregisterRecv(handlerType string) error {
	return ic.delegate.UnregisterRecv(handlerType)
}

func (ic *interceptingChannel) RegisterUnmarshaler(
	unmarshaler func() net.TaggedUnmarshaler,
) error {
	return ic.delegate.RegisterUnmarshaler(unmarshaler)
}
