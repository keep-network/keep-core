package libp2p

import (
	"sync"
	"time"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
	host "github.com/libp2p/go-libp2p-host"
	peer "github.com/libp2p/go-libp2p-peer"
)

type messageCache struct {
	nonceService        map[peer.ID]*nonceService
	nonceServiceLock    sync.Mutex
	nonceServiceTimeout time.Time

	messageBuffer     map[peer.ID][]pb.NetworkMessage
	messageBufferLock sync.Mutex

	p2phost  host.Host
	identity *identity
}

func newMessageCache(messageBufferSize int, p2phost host.Host, identity *identity) *messageCache {
	return &messageCache{
		nonceService:  make(map[peer.ID]*nonceService),
		messageBuffer: make(map[peer.ID][]pb.NetworkMessage),
	}
}
