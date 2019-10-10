package libp2p

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/key"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type channelManager struct {
	ctx context.Context

	identity  *identity
	peerStore peerstore.Peerstore

	channelsMutex sync.Mutex
	channels      map[string]*channel

	pubsub *pubsub.PubSub
}

type channelConfig struct {
	authorizedAddresses [][]byte
}

func (cc *channelConfig) SetAuthorizedAddresses(addresses [][]byte) error {
	if !(len(addresses) > 0) {
		return fmt.Errorf("addresses should have a non-zero length")
	}
	cc.authorizedAddresses = addresses
	return nil
}

func (cc *channelConfig) apply(options ...net.ChannelOption) error {
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(cc); err != nil {
			return err
		}
	}
	return nil
}

func newChannelManager(
	ctx context.Context,
	identity *identity,
	p2phost host.Host,
) (*channelManager, error) {
	floodsub, err := pubsub.NewFloodSub(
		ctx,
		p2phost,
		pubsub.WithMessageAuthor(identity.id),
		pubsub.WithMessageSigning(true),
		pubsub.WithStrictSignatureVerification(true),
	)
	if err != nil {
		return nil, err
	}
	return &channelManager{
		channels:  make(map[string]*channel),
		pubsub:    floodsub,
		peerStore: p2phost.Peerstore(),
		identity:  identity,
		ctx:       ctx,
	}, nil
}

func (cm *channelManager) getChannel(
	name string,
	options ...net.ChannelOption,
) (*channel, error) {
	var (
		channel       *channel
		channelConfig channelConfig
		exists        bool
		err           error
	)

	cm.channelsMutex.Lock()
	channel, exists = cm.channels[name]
	cm.channelsMutex.Unlock()

	if !exists {
		err = channelConfig.apply(options...)
		if err != nil {
			return nil, err
		}

		channel, err = cm.newChannel(name, channelConfig)
		if err != nil {
			return nil, err
		}

		// Ensure we update our cache of known channels
		cm.channelsMutex.Lock()
		cm.channels[name] = channel
		cm.channelsMutex.Unlock()
	}

	return channel, nil
}

func (cm *channelManager) newChannel(
	name string,
	config channelConfig,
) (*channel, error) {
	sub, err := cm.pubsub.Subscribe(name)
	if err != nil {
		return nil, err
	}

	if len(config.authorizedAddresses) > 0 {
		err := cm.pubsub.RegisterTopicValidator(
			name,
			newAuthorizedAddressesValidator(config.authorizedAddresses),
			pubsub.WithValidatorInline(true),
		)
		if err != nil {
			logger.Errorf("validator registration failed [%v]", err)
			return nil, err
		}
	}

	channel := &channel{
		name:               name,
		clientIdentity:     cm.identity,
		peerStore:          cm.peerStore,
		pubsub:             cm.pubsub,
		subscription:       sub,
		messageHandlers:    make([]net.HandleMessageFunc, 0),
		unmarshalersByType: make(map[string]func() net.TaggedUnmarshaler),
	}

	go channel.handleMessages(cm.ctx)

	return channel, nil
}

func newAuthorizedAddressesValidator(authorizedAddresses [][]byte) pubsub.Validator {
	authorizations := make(map[string]bool, len(authorizedAddresses))
	for _, address := range authorizedAddresses {
		encodedAddress := strings.ToLower("0x" + hex.EncodeToString(address))
		authorizations[encodedAddress] = true
	}

	return func(_ context.Context, peerID peer.ID, message *pubsub.Message) bool {
		// TODO `peerID` is probably only the message sender not necessarily
		//  the author (to check in libp2p code). `message.GetFrom()` is
		//  the author (according to libp2p docs). Rethink which id should
		//  be used here.
		peerPublicKey, err := peerID.ExtractPublicKey()
		if err != nil {
			logger.Errorf(
				"cannot extract public key for peer [%v] [%v]",
				peerID,
				err,
			)
			return false
		}
		peerNetworkPublicKey := key.Libp2pKeyToNetworkKey(peerPublicKey)
		peerEthAddress := strings.ToLower(
			key.NetworkPubKeyToEthAddress(peerNetworkPublicKey),
		)

		_, isAuthorized := authorizations[peerEthAddress]
		return isAuthorized
	}
}
