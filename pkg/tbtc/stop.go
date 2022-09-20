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
}

func (sp *StopPill) Type() string {
	return "tecdsa/stop_pill"
}

func (sp *StopPill) Marshal() ([]byte, error) {
	return proto.Marshal(&pb.StopPill{
		AttemptNumber: sp.attemptNumber,
	})
}

func (sp *StopPill) Unmarshal(bytes []byte) error {
	pbStopPill := pb.StopPill{}
	if err := proto.Unmarshal(bytes, &pbStopPill); err != nil {
		return fmt.Errorf("failed to unmarshal StopPill: [%v]", err)
	}

	sp.attemptNumber = pbStopPill.AttemptNumber

	return nil
}

func registerStopPillUnmarshaller(channel net.BroadcastChannel) {
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &StopPill{}
	})
}

func sendStopPill(
	ctx context.Context,
	broadcastChannel net.BroadcastChannel,
	attemptNumber uint,
) {
	stopPill := &StopPill{uint64(attemptNumber)}
	broadcastChannel.Send(ctx, stopPill)
}

func cancelContextOnStopSignal(
	ctx context.Context,
	cancelFn func(),
	broadcastChannel net.BroadcastChannel,
) {
	broadcastChannel.Recv(ctx, func(msg net.Message) {
		switch msg.Payload().(type) {
		case *StopPill:
			cancelFn()
		}
	})
}
