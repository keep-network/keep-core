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

// signingDoneMessage is a message used to signal a successful signature
// calculation across all signing group members.
type signingDoneMessage struct {
	senderID      group.MemberIndex
	message       *big.Int
	attemptNumber uint64
	signature     *tecdsa.Signature
	endBlock      uint64
}

func (sdm *signingDoneMessage) Type() string {
	return "tbtc/signing_done_message"
}

// signingDoneCheck is a component that is responsible for signaling a
// successful signature calculation across all signing group members.
type signingDoneCheck struct {
	groupSize           int
	broadcastChannel    net.BroadcastChannel
	membershipValidator *group.MembershipValidator
}

func newSigningDoneCheck(
	groupSize int,
	broadcastChannel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
) *signingDoneCheck {
	broadcastChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &signingDoneMessage{}
	})

	return &signingDoneCheck{
		groupSize:           groupSize,
		broadcastChannel:    broadcastChannel,
		membershipValidator: membershipValidator,
	}
}

// exchange runs the signing done check exchanging routine. This function:
// - broadcasts the signing done check along with information necessary to
//   attribute the result to the given signing attempt.
// - listens for incoming signing done checks from other members participating
//   in the given signing attempt, matching the broadcasted done check.
// This function blocks until it receives all the required done checks from
// other members or until the passed context is done. In the first case, it
// returns the block at which the slowest signer completed the signature
// computation process. However, even after the function return, the done check
// is retransmitted for the lifetime of the passed context. If the expected
// done checks are not received on time, the function returns an error.
func (sdc *signingDoneCheck) exchange(
	ctx context.Context,
	memberIndex group.MemberIndex,
	message *big.Int,
	attemptNumber uint64,
	attemptMembersIndexes []group.MemberIndex,
	result *signing.Result,
	endBlock uint64,
) (uint64, error) {
	// Use a separate context for the message receiver as the receiver must
	// be closed upon function return. Leaving a dangling receiver without
	// the message processing loop will cause warnings on the channel level.
	receiveCtx, cancelReceiveCtx := context.WithCancel(ctx)
	defer cancelReceiveCtx()

	messagesChan := make(chan net.Message, sdc.groupSize)
	sdc.broadcastChannel.Recv(receiveCtx, func(message net.Message) {
		messagesChan <- message
	})

	// Use the original context for the send routine as we want to keep
	// retransmissions on until the context is alive.
	err := sdc.broadcastChannel.Send(ctx, &signingDoneMessage{
		senderID:      memberIndex,
		message:       message,
		attemptNumber: attemptNumber,
		signature:     result.Signature,
		endBlock:      endBlock,
	}, net.BackoffRetransmissionStrategy)
	if err != nil {
		return 0, fmt.Errorf("cannot send signing done message: [%v]", err)
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
			doneMessage, ok := netMessage.Payload().(*signingDoneMessage)
			if !ok {
				continue
			}

			if !sdc.isValidDoneMessage(
				doneMessage,
				awaitingSenders,
				netMessage.SenderPublicKey(),
				message,
				attemptNumber,
				result.Signature,
			) {
				continue
			}

			if doneMessage.endBlock > latestEndBlock {
				latestEndBlock = doneMessage.endBlock
			}

			delete(awaitingSenders, doneMessage.senderID)

			if len(awaitingSenders) == 0 {
				return latestEndBlock, nil
			}
		case <-ctx.Done():
			return 0, fmt.Errorf("cannot receive signing done messages on time")
		}
	}
}

// listen runs the signing done check listening routine. This function listens
// for incoming signing done checks from members participating in the given
// signing attempt. This function blocks until it receives all the required
// done checks from members or until the passed context is done. In the first
// case, it returns the signature computed by the signing members and the block
// at which the slowest signer completed the signature computation process.
// If the expected done checks are not received on time, the function returns
// an error.
func (sdc *signingDoneCheck) listen(
	ctx context.Context,
	message *big.Int,
	attemptNumber uint64,
	attemptMembersIndexes []group.MemberIndex,
) (*signing.Result, uint64, error) {
	// Use a separate context for the message receiver as the receiver must
	// be closed upon function return. Leaving a dangling receiver without
	// the message processing loop will cause warnings on the channel level.
	receiveCtx, cancelReceiveCtx := context.WithCancel(ctx)
	defer cancelReceiveCtx()

	messagesChan := make(chan net.Message, sdc.groupSize)
	sdc.broadcastChannel.Recv(receiveCtx, func(message net.Message) {
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
			doneMessage, ok := netMessage.Payload().(*signingDoneMessage)
			if !ok {
				continue
			}

			if !sdc.isValidDoneMessage(
				doneMessage,
				awaitingSenders,
				netMessage.SenderPublicKey(),
				message,
				attemptNumber,
				signature,
			) {
				continue
			}

			if doneMessage.endBlock > latestEndBlock {
				latestEndBlock = doneMessage.endBlock
			}

			if signature == nil {
				signature = doneMessage.signature
			}

			delete(awaitingSenders, doneMessage.senderID)

			if len(awaitingSenders) == 0 {
				return &signing.Result{Signature: signature}, latestEndBlock, nil
			}
		case <-ctx.Done():
			return nil, 0, fmt.Errorf("cannot receive signing done messages on time")
		}
	}
}

// isValidDoneMessage validates the given signingDoneMessage in the context
// of the given signing attempt.
func (sdc *signingDoneCheck) isValidDoneMessage(
	doneMessage *signingDoneMessage,
	awaitingSenders map[group.MemberIndex]bool,
	senderPublicKey []byte,
	message *big.Int,
	attemptNumber uint64,
	signature *tecdsa.Signature,
) bool {
	if !awaitingSenders[doneMessage.senderID] {
		return false
	}

	if !sdc.membershipValidator.IsValidMembership(
		doneMessage.senderID,
		senderPublicKey,
	) {
		return false
	}

	if doneMessage.message.Cmp(message) != 0 {
		return false
	}

	if doneMessage.attemptNumber != attemptNumber {
		return false
	}

	if signature != nil {
		if !doneMessage.signature.Equals(signature) {
			return false
		}
	}

	return true
}
