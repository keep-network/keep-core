package tbtc

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
	"math/big"
)

// signingSyncMessage is a message used to synchronize members of a signing
// group upon successful signature calculation.
type signingSyncMessage struct {
	senderID      group.MemberIndex
	message       *big.Int
	attemptNumber uint
	signature     *tecdsa.Signature
	endBlock      uint64
}

func (ssm *signingSyncMessage) Type() string {
	return "tbtc/signing_sync_message"
}

// signingSyncer is a component that is responsible for synchronization of
// the signing result across all signing group members.
type signingSyncer struct {
	groupSize           int
	broadcastChannel    net.BroadcastChannel
	membershipValidator *group.MembershipValidator
}

// newSigningSyncer creates a new instance of the signingSyncer.
func newSigningSyncer(
	groupSize int,
	broadcastChannel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
) *signingSyncer {
	broadcastChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &signingSyncMessage{}
	})

	return &signingSyncer{
		groupSize:           groupSize,
		broadcastChannel:    broadcastChannel,
		membershipValidator: membershipValidator,
	}
}

// syncAttemptParticipant runs the attempt participant sync routine. This
// function broadcasts the signing result along with information necessary to
// attribute the result to the given signing attempt and listens for similar
// messages from other members participating in the given attempt. This function
// blocks until it receives all the required sync messages from other members
// or until the passed context is done. In the first case, it returns the
// block at which the slowest signer completed the signature computation process.
// Otherwise, it returns an error.
func (ss *signingSyncer) syncAttemptParticipant(
	ctx context.Context,
	memberIndex group.MemberIndex,
	message *big.Int,
	attemptNumber uint,
	attemptMembersIndexes []group.MemberIndex,
	result *signing.Result,
	endBlock uint64,
) (uint64, error) {
	messagesChan := make(chan net.Message, ss.groupSize)
	ss.broadcastChannel.Recv(ctx, func(message net.Message) {
		messagesChan <- message
	})

	err := ss.broadcastChannel.Send(ctx, &signingSyncMessage{
		senderID:      memberIndex,
		message:       message,
		attemptNumber: attemptNumber,
		signature:     result.Signature,
		endBlock:      endBlock,
	})
	if err != nil {
		return 0, fmt.Errorf("cannot send sync message: [%v]", err)
	}

	awaitingSenders := make(map[group.MemberIndex]bool)
	for _, attemptMemberIndex := range attemptMembersIndexes {
		if memberIndex == attemptMemberIndex {
			continue
		}

		awaitingSenders[attemptMemberIndex] = true
	}

	latestEndBlock := endBlock

	for {
		select {
		case netMessage := <-messagesChan:
			syncMessage, ok := netMessage.Payload().(*signingSyncMessage)
			if !ok {
				continue
			}

			if !ss.isValidSyncMessage(
				syncMessage,
				awaitingSenders,
				netMessage.SenderPublicKey(),
				message,
				attemptNumber,
				result.Signature,
			) {
				continue
			}

			if syncMessage.endBlock > latestEndBlock {
				latestEndBlock = syncMessage.endBlock
			}

			delete(awaitingSenders, syncMessage.senderID)

			if len(awaitingSenders) == 0 {
				return latestEndBlock, nil
			}
		case <-ctx.Done():
			return 0, fmt.Errorf("cannot receive sync messages on time")
		}
	}
}

// syncAttemptObserver runs the attempt observer sync routine. This function
// is kind of a "passive" version of the syncAttemptParticipant function.
// It only listens for signing sync messages from members participating in the
// given signing attempt. This function blocks until it receives all the
// required sync messages from signing members or until the passed context is
// done. In the first case, it returns the signature computed by the signing
// members and the block at which the slowest signer completed the signature
// computation process. Otherwise, it returns an error.
func (ss *signingSyncer) syncAttemptObserver(
	ctx context.Context,
	message *big.Int,
	attemptNumber uint,
	attemptMembersIndexes []group.MemberIndex,
) (*signing.Result, uint64, error) {
	messagesChan := make(chan net.Message, ss.groupSize)
	ss.broadcastChannel.Recv(ctx, func(message net.Message) {
		messagesChan <- message
	})

	awaitingSenders := make(map[group.MemberIndex]bool)
	for _, attemptMemberIndex := range attemptMembersIndexes {
		awaitingSenders[attemptMemberIndex] = true
	}

	var signature *tecdsa.Signature
	latestEndBlock := uint64(0)

	for {
		select {
		case netMessage := <-messagesChan:
			syncMessage, ok := netMessage.Payload().(*signingSyncMessage)
			if !ok {
				continue
			}

			if !ss.isValidSyncMessage(
				syncMessage,
				awaitingSenders,
				netMessage.SenderPublicKey(),
				message,
				attemptNumber,
				signature,
			) {
				continue
			}

			if syncMessage.endBlock > latestEndBlock {
				latestEndBlock = syncMessage.endBlock
			}

			if signature == nil {
				signature = syncMessage.signature
			}

			delete(awaitingSenders, syncMessage.senderID)

			if len(awaitingSenders) == 0 {
				return &signing.Result{Signature: signature}, latestEndBlock, nil
			}
		case <-ctx.Done():
			return nil, 0, fmt.Errorf("cannot receive sync messages on time")
		}
	}
}

// isValidSyncMessage validates the given signingSyncMessage in the context
// of the given signing attempt.
func (ss *signingSyncer) isValidSyncMessage(
	syncMessage *signingSyncMessage,
	awaitingSenders map[group.MemberIndex]bool,
	senderPublicKey []byte,
	message *big.Int,
	attemptNumber uint,
	signature *tecdsa.Signature,
) bool {
	if !awaitingSenders[syncMessage.senderID] {
		return false
	}

	if !ss.membershipValidator.IsValidMembership(
		syncMessage.senderID,
		senderPublicKey,
	) {
		return false
	}

	if syncMessage.message.Cmp(message) != 0 {
		return false
	}

	if syncMessage.attemptNumber != attemptNumber {
		return false
	}

	if signature != nil {
		if !syncMessage.signature.Equals(signature) {
			return false
		}
	}

	return true
}
