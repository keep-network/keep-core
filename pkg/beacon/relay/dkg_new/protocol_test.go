package dkg

import (
	"math/big"
)

func filterMessagesForReceiver(
	messages []*MemberCommitmentsMessage,
	receiverID *big.Int,
) []*MemberCommitmentsMessage {
	var result []*MemberCommitmentsMessage
	for _, msg := range messages {
		if msg.senderID != receiverID {
			result = append(result, msg)
		}
	}
	return result
}
