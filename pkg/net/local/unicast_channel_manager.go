package local

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/operator"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
)

var unicastChannelManagersMutex = &sync.RWMutex{}
var unicastChannelManagers = make(map[string]*unicastChannelManager)

type unicastChannelManager struct {
	transportID       net.TransportIdentifier
	operatorPublicKey *operator.PublicKey

	channelsMutex *sync.RWMutex
	channels      map[net.TransportIdentifier]*unicastChannel

	onChannelOpenedHandlersMutex *sync.RWMutex
	onChannelOpenedHandlers      []*onChannelOpenedHandler
}

type onChannelOpenedHandler struct {
	ctx      context.Context
	handleFn func(remote net.UnicastChannel)
}

func newUnicastChannelManager(
	operatorPublicKey *operator.PublicKey,
) *unicastChannelManager {
	unicastChannelManagersMutex.Lock()
	defer unicastChannelManagersMutex.Unlock()

	transportID, err := createLocalIdentifier(operatorPublicKey)
	if err != nil {
		panic(err)
	}

	existingChannelManager, ok := unicastChannelManagers[transportID.String()]
	if ok {
		return existingChannelManager
	}

	channelManager := &unicastChannelManager{
		transportID:                  transportID,
		operatorPublicKey:            operatorPublicKey,
		channelsMutex:                &sync.RWMutex{},
		channels:                     make(map[net.TransportIdentifier]*unicastChannel),
		onChannelOpenedHandlersMutex: &sync.RWMutex{},
		onChannelOpenedHandlers:      make([]*onChannelOpenedHandler, 0),
	}

	logger.Debugf("registering as [%v]", transportID)
	unicastChannelManagers[transportID.String()] = channelManager

	return channelManager
}

// UnicastChannelWith creates a unicast channel with the given remote peer.
// If peer is not known or connection could not be open, function returns error.
func (up *unicastChannelManager) UnicastChannelWith(
	peer net.TransportIdentifier,
) (net.UnicastChannel, error) {
	channel := up.createUnicastChannelWith(peer, false)

	unicastChannelManagersMutex.RLock()
	remote, ok := unicastChannelManagers[peer.String()]
	unicastChannelManagersMutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("remote peer not known [%v]", peer)
	}

	remote.createUnicastChannelWith(up.transportID, true)

	return channel, nil
}

// OnUnicastChannelOpened registers UnicastChannelHandler that will be called
// for each incoming unicast channel opened by remote peers against this one.
// The handlers is active for the entire lifetime of the provided context.
// When the context is done, handler is never called again.
func (up *unicastChannelManager) OnUnicastChannelOpened(
	ctx context.Context,
	handler func(remote net.UnicastChannel),
) {
	up.onChannelOpenedHandlersMutex.Lock()
	defer up.onChannelOpenedHandlersMutex.Unlock()

	up.onChannelOpenedHandlers = append(
		up.onChannelOpenedHandlers,
		&onChannelOpenedHandler{ctx, handler},
	)
}

func (up *unicastChannelManager) createUnicastChannelWith(
	peer net.TransportIdentifier,
	notify bool,
) net.UnicastChannel {
	channel, ok := up.getUnicastChannel(peer)
	if ok {
		return channel
	}

	channel = newUnicastChannel(up.transportID, up.operatorPublicKey, peer)
	up.addUnicastChannel(channel)

	if notify {
		up.notifyNewChannel(channel)
	}

	return channel
}

func (up *unicastChannelManager) getUnicastChannel(
	receiver net.TransportIdentifier,
) (*unicastChannel, bool) {
	up.channelsMutex.RLock()
	defer up.channelsMutex.RUnlock()

	channel, ok := up.channels[receiver]
	return channel, ok
}

func (up *unicastChannelManager) addUnicastChannel(channel *unicastChannel) {
	up.channelsMutex.Lock()
	defer up.channelsMutex.Unlock()

	up.channels[channel.receiverTransportID] = channel
}

func (up *unicastChannelManager) notifyNewChannel(channel net.UnicastChannel) {
	// first cleanup
	up.onChannelOpenedHandlersMutex.Lock()
	defer up.onChannelOpenedHandlersMutex.Unlock()

	i := 0
	for _, handler := range up.onChannelOpenedHandlers {
		if handler.ctx.Err() == nil {
			// still active, should remain in the slice
			up.onChannelOpenedHandlers[i] = handler
			i++

			// firing handler asynchronously to
			// do not block the loop
			go handler.handleFn(channel)
		}
	}

	// cleaning up those no longer active
	up.onChannelOpenedHandlers = up.onChannelOpenedHandlers[:i]
}

func deliverMessage(
	sender net.TransportIdentifier,
	receiver net.TransportIdentifier,
	messagePayload []byte,
	messageType string,
) error {
	unicastChannelManagersMutex.RLock()
	receiverChannelManager, ok := unicastChannelManagers[receiver.String()]
	unicastChannelManagersMutex.RUnlock()

	if !ok {
		return fmt.Errorf("peer [%v] not known", receiver)
	}

	channel, ok := receiverChannelManager.getUnicastChannel(sender)
	if !ok {
		return fmt.Errorf("peer [%v] could not find channel for [%v]", receiver, sender)
	}

	return channel.receiveMessage(messagePayload, messageType)
}
