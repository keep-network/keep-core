package entry

import (
	"fmt"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/beacon/relay/state"
	"github.com/keep-network/keep-core/pkg/bls"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

type signingState = state.State

type signingStateBase struct {
	channel      net.BroadcastChannel
	relayChain   relayChain.Interface
	blockCounter chain.BlockCounter

	signer *dkg.ThresholdSigner

	previousEntry []byte

	honestThreshold int
}

type signatureShareState struct {
	signingStateBase

	selfSignatureShare     *bn256.G1
	signatureShareMessages []*SignatureShareMessage

	signingStartBlockHeight uint64
}

func (sss *signatureShareState) DelayBlocks() uint64 {
	return state.MessagingStateDelayBlocks
}

func (sss *signatureShareState) ActiveBlocks() uint64 {
	return state.MessagingStateActiveBlocks
}

func (sss *signatureShareState) Initiate() error {
	share, err := sss.signer.CalculateSignatureShare(sss.previousEntry)
	if err != nil {
		return fmt.Errorf("could not evaluate signature share: [%v]", err)
	}

	sss.selfSignatureShare = share

	message := &SignatureShareMessage{
		sss.MemberIndex(),
		sss.selfSignatureShare.Marshal(),
	}
	if err := sss.channel.Send(message); err != nil {
		return err
	}
	return nil
}

func (sss *signatureShareState) Receive(msg net.Message) error {
	switch signatureShareMessage := msg.Payload().(type) {
	case *SignatureShareMessage:
		if !group.IsMessageFromSelf(
			sss.MemberIndex(),
			signatureShareMessage,
		) {
			sss.signatureShareMessages = append(
				sss.signatureShareMessages,
				signatureShareMessage,
			)
		}
	}

	return nil
}

func (sss *signatureShareState) Next() signingState {
	return &signatureCompleteState{
		signingStateBase:      sss.signingStateBase,
		selfSignatureShare:    sss.selfSignatureShare,
		previousPhaseMessages: sss.signatureShareMessages,
		signatureCompletionStartBlockHeight: sss.signingStartBlockHeight +
			sss.DelayBlocks() +
			sss.ActiveBlocks(),
	}
}

func (sss *signatureShareState) MemberIndex() group.MemberIndex {
	return sss.signer.MemberID()
}

type signatureCompleteState struct {
	signingStateBase

	selfSignatureShare    *bn256.G1
	previousPhaseMessages []*SignatureShareMessage
	fullSignature         []byte

	signatureCompletionStartBlockHeight uint64
}

func (scs *signatureCompleteState) DelayBlocks() uint64 {
	return state.SilentStateDelayBlocks
}

func (scs *signatureCompleteState) ActiveBlocks() uint64 {
	return state.SilentStateActiveBlocks
}

func (scs *signatureCompleteState) Initiate() error {
	seenShares := make(map[group.MemberIndex]*bn256.G1)
	seenShares[scs.MemberIndex()] = scs.selfSignatureShare
	logger.Debugf(
		"[member:%v] auto-accepting self signature share: [%v]",
		scs.MemberIndex(),
		scs.MemberIndex(),
	)

	for _, message := range scs.previousPhaseMessages {
		share := new(bn256.G1)
		_, err := share.Unmarshal(message.shareBytes)
		if err != nil {
			logger.Errorf(
				"[member:%v] failed to unmarshal signature share from member [%v]: [%v]",
				scs.MemberIndex(),
				message.senderID,
				err,
			)
		} else {
			logger.Debugf(
				"[member:%v] accepting signature share from member [%v]",
				scs.MemberIndex(),
				message.senderID,
			)
			seenShares[message.senderID] = share
		}
	}

	seenSharesSlice := make([]*bls.SignatureShare, 0)
	for memberID, share := range seenShares {
		signatureShare := &bls.SignatureShare{I: int(memberID), V: share}
		seenSharesSlice = append(seenSharesSlice, signatureShare)
	}

	logger.Infof(
		"[member:%v] restoring signature from [%v] shares",
		scs.MemberIndex(),
		len(seenSharesSlice),
	)
	signature, err := scs.signer.CompleteSignature(
		seenSharesSlice,
		scs.honestThreshold,
	)
	if err != nil {
		return err
	}

	scs.fullSignature = signature.Marshal()

	return nil
}

func (scs *signatureCompleteState) Receive(msg net.Message) error {
	return nil
}

func (scs *signatureCompleteState) Next() signingState {
	return &entrySubmissionState{
		signingStateBase: scs.signingStateBase,
		signature:        scs.fullSignature,
		entrySubmissionStartBlockHeight: scs.signatureCompletionStartBlockHeight +
			scs.DelayBlocks() +
			scs.ActiveBlocks(),
	}
}

func (scs *signatureCompleteState) MemberIndex() group.MemberIndex {
	return scs.signer.MemberID()
}

type entrySubmissionState struct {
	signingStateBase

	signature []byte

	entrySubmissionStartBlockHeight uint64
}

func (ess *entrySubmissionState) DelayBlocks() uint64 {
	return state.SilentStateDelayBlocks
}

func (ess *entrySubmissionState) ActiveBlocks() uint64 {
	// We do not exchange any messages in this phase. We publish entry to the
	// chain but it is an action blocking all group members for the same time
	// - members exit when the first valid entry is accepted by the chain.
	// How long it takes depends on the block step and group size.
	return state.SilentStateActiveBlocks
}

func (ess *entrySubmissionState) Initiate() error {
	submitter := &relayEntrySubmitter{
		chain:        ess.relayChain,
		blockCounter: ess.blockCounter,
		index:        ess.MemberIndex(),
	}

	return submitter.submitRelayEntry(
		ess.signature,
		ess.signer.GroupPublicKeyBytes(),
		ess.entrySubmissionStartBlockHeight,
	)
}

func (ess *entrySubmissionState) Receive(msg net.Message) error {
	return nil
}

func (ess *entrySubmissionState) Next() signingState {
	return nil
}

func (ess *entrySubmissionState) MemberIndex() group.MemberIndex {
	return ess.signer.MemberID()
}
