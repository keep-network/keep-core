package libp2p

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

type retransmitter struct {
	cycles   uint32
	interval time.Duration
	cache    *timeCache
}

func newRetransmitter(cycles uint32, interval time.Duration) *retransmitter {
	retransmissionDuration := time.Duration(cycles) * interval
	cacheLifetime := 2 * time.Minute

	cache := newTimeCache(retransmissionDuration + cacheLifetime)

	return &retransmitter{
		cycles:   cycles,
		interval: interval,
		cache:    cache,
	}
}

func (r *retransmitter) scheduleRetransmission(
	message *pb.NetworkMessage,
	sender func(*pb.NetworkMessage) error,
) {
	go func() {
		for i := uint32(1); i <= r.cycles; i++ {
			time.Sleep(r.interval)

			copy := pb.NetworkMessage(*message)
			copy.Retransmission = i

			go func() {
				if err := sender(&copy); err != nil {
					logger.Errorf(
						"could not retransmit message: [%v]",
						err,
					)
				}
			}()
		}
		return
	}()
}

func (r *retransmitter) sweepReceived(
	message *pb.NetworkMessage,
	receive func() error,
) error {
	fingerprint, err := calculateFingerprint(message)
	if err != nil {
		return fmt.Errorf("could not calculate message fingerprint: [%v]", err)
	}

	if r.cache.has(fingerprint) {
		return nil
	}

	if r.cache.add(fingerprint) {
		return receive()
	}

	return nil
}

func calculateFingerprint(message *pb.NetworkMessage) (string, error) {
	copy := pb.NetworkMessage(*message)
	copy.Retransmission = 0

	bytes, err := copy.Marshal()
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(bytes)
	return hex.EncodeToString(sum[:]), nil
}
