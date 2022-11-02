package tbtc

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/tbtc/gen/pb"
	"google.golang.org/protobuf/proto"
)

// TODO: This file should be gone once the contract integration is implemented.

// StopPill is a temporary workaround for a missing chain integration. When
// a group member is not selected for the current attempt of key generation or
// signing and there is no other member from the same client selected for the
// protocol execution, the member does not know what was the result of the
// protocol execution and if it completed or not. In other words, the member
// will stay hung on the block waiter, waiting for their turn. The StopPill is
// sent via broadcast channel on a successful protocol execution and tells all
// members waiting for their turn in the retry loop to stop because the result
// was produced.
type StopPill struct {
	attemptNumber uint64
	dkgSeed       string // empty if the stop pill is sent for signing
	messageToSign string // empty if the stop pill is sent for DKG
}

func (sp *StopPill) Type() string {
	return "tecdsa/stop_pill"
}

func (sp *StopPill) Marshal() ([]byte, error) {
	return proto.Marshal(&pb.StopPill{
		AttemptNumber: sp.attemptNumber,
		DkgSeed:       sp.dkgSeed,
		MessageToSign: sp.messageToSign,
	})
}

func (sp *StopPill) Unmarshal(bytes []byte) error {
	pbStopPill := pb.StopPill{}
	if err := proto.Unmarshal(bytes, &pbStopPill); err != nil {
		return fmt.Errorf("failed to unmarshal StopPill: [%v]", err)
	}

	sp.attemptNumber = pbStopPill.AttemptNumber
	sp.dkgSeed = pbStopPill.DkgSeed
	sp.messageToSign = pbStopPill.MessageToSign

	return nil
}

func registerStopPillUnmarshaller(channel net.BroadcastChannel) {
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &StopPill{}
	})
}

func sendDkgStopPill(
	ctx context.Context,
	broadcastChannel net.BroadcastChannel,
	dkgSeed string,
	attemptNumber uint,
) error {
	stopPill := &StopPill{
		attemptNumber: uint64(attemptNumber),
		dkgSeed:       dkgSeed,
	}
	return broadcastChannel.Send(ctx, stopPill)
}

func sendSigningStopPill(
	ctx context.Context,
	broadcastChannel net.BroadcastChannel,
	messageToSign string,
	attemptNumber uint,
) error {
	stopPill := &StopPill{
		attemptNumber: uint64(attemptNumber),
		messageToSign: messageToSign,
	}
	return broadcastChannel.Send(ctx, stopPill)
}

func cancelDkgContextOnStopSignal(
	ctx context.Context,
	cancelFn func(),
	broadcastChannel net.BroadcastChannel,
	dkgSeed string,
) {
	broadcastChannel.Recv(ctx, func(msg net.Message) {
		switch stopPill := msg.Payload().(type) {
		case *StopPill:
			if stopPill.dkgSeed == dkgSeed {
				cancelFn()
			}
		}
	})
}

func cancelSigningContextOnStopSignal(
	ctx context.Context,
	cancelFn func(),
	broadcastChannel net.BroadcastChannel,
	messageToSign string,
) {
	broadcastChannel.Recv(ctx, func(msg net.Message) {
		switch stopPill := msg.Payload().(type) {
		case *StopPill:
			if stopPill.messageToSign == messageToSign {
				cancelFn()
			}
		}
	})
}
