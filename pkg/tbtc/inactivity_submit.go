package tbtc

import (
	"context"
	"fmt"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/inactivity"
)

// inactivityClaimSigner is responsible for signing the inactivity claim and
// verification of signatures generated by other group members.
type inactivityClaimSigner struct {
	chain Chain
}

func newInactivityClaimSigner(
	chain Chain,
) *inactivityClaimSigner {
	return &inactivityClaimSigner{
		chain: chain,
	}
}

func (ics *inactivityClaimSigner) SignClaim(claim *inactivity.Claim) (
	*inactivity.SignedClaim,
	error,
) {
	if claim == nil {
		return nil, fmt.Errorf("result is nil")
	}

	claimHash, err := ics.chain.CalculateInactivityClaimSignatureHash(claim)
	if err != nil {
		return nil, fmt.Errorf(
			"inactivity claim hash calculation failed [%w]",
			err,
		)
	}

	signing := ics.chain.Signing()

	signature, err := signing.Sign(claimHash[:])
	if err != nil {
		return nil, fmt.Errorf(
			"inactivity claim hash signing failed [%w]",
			err,
		)
	}

	return &inactivity.SignedClaim{
		PublicKey: signing.PublicKey(),
		Signature: signature,
		ClaimHash: claimHash,
	}, nil
}

// VerifySignature verifies if the signature was generated from the provided
// inactivity claim using the provided public key.
func (ics *inactivityClaimSigner) VerifySignature(
	signedClaim *inactivity.SignedClaim,
) (
	bool,
	error,
) {
	return ics.chain.Signing().VerifyWithPublicKey(
		signedClaim.ClaimHash[:],
		signedClaim.Signature,
		signedClaim.PublicKey,
	)
}

type inactivityClaimSubmitter struct {
	inactivityLogger log.StandardLogger

	chain           Chain
	groupParameters *GroupParameters
	groupMembers    []uint32

	waitForBlockFn waitForBlockFn
}

func newInactivityClaimSubmitter(
	inactivityLogger log.StandardLogger,
	chain Chain,
	groupParameters *GroupParameters,
	groupMembers []uint32,
	waitForBlockFn waitForBlockFn,
) *inactivityClaimSubmitter {
	return &inactivityClaimSubmitter{
		inactivityLogger: inactivityLogger,
		chain:            chain,
		groupParameters:  groupParameters,
		groupMembers:     groupMembers,
		waitForBlockFn:   waitForBlockFn,
	}
}

func (ics *inactivityClaimSubmitter) SubmitClaim(
	ctx context.Context,
	memberIndex group.MemberIndex,
	claim *inactivity.Claim,
	signatures map[group.MemberIndex][]byte,
) error {
	if len(signatures) < ics.groupParameters.HonestThreshold {
		return fmt.Errorf(
			"could not submit inactivity claim with [%v] signatures for "+
				"group honest threshold [%v]",
			len(signatures),
			ics.groupParameters.HonestThreshold,
		)
	}

	// The inactivity nonce at the beginning of the execution process.
	inactivityNonce := claim.Nonce

	walletPublicKeyHash := bitcoin.PublicKeyHash(claim.WalletPublicKey)

	walletRegistryData, err := ics.chain.GetWallet(walletPublicKeyHash)
	if err != nil {
		return fmt.Errorf("could not get registry data on wallet: [%v]", err)
	}

	ecdsaWalletID := walletRegistryData.EcdsaWalletID

	currentNonce, err := ics.chain.GetInactivityClaimNonce(
		ecdsaWalletID,
	)
	if err != nil {
		return fmt.Errorf("could not get nonce for wallet: [%v]", err)
	}

	if currentNonce.Cmp(inactivityNonce) > 0 {
		// Someone who was ahead of us in the queue submitted the claim. Giving up.
		ics.inactivityLogger.Infof(
			"[member:%v] inactivity claim already submitted; "+
				"aborting inactivity claim on-chain submission",
			memberIndex,
		)
		return nil
	}

	inactivityClaim, err := ics.chain.AssembleInactivityClaim(
		ecdsaWalletID,
		claim.GetInactiveMembersIndexes(),
		signatures,
		claim.HeartbeatFailed,
	)
	if err != nil {
		return fmt.Errorf("could not assemble inactivity chain claim [%w]", err)
	}

	blockCounter, err := ics.chain.BlockCounter()
	if err != nil {
		return err
	}

	// We can't determine a common block at which the publication starts.
	// However, all we want here is to ensure the members does not submit
	// in the same time. This can be achieved by simply using the index-based
	// delay starting from the current block.
	currentBlock, err := blockCounter.CurrentBlock()
	if err != nil {
		return fmt.Errorf("cannot get current block: [%v]", err)
	}
	delayBlocks := uint64(memberIndex-1) * inactivityClaimSubmissionDelayStepBlocks
	submissionBlock := currentBlock + delayBlocks

	ics.inactivityLogger.Infof(
		"[member:%v] waiting for block [%v] to submit inactivity claim",
		memberIndex,
		submissionBlock,
	)

	err = ics.waitForBlockFn(ctx, submissionBlock)
	if err != nil {
		return fmt.Errorf(
			"error while waiting for inactivity claim submission block: [%v]",
			err,
		)
	}

	if ctx.Err() != nil {
		// The context was cancelled by the upstream. Regardless of the cause,
		// that means the inactivity execution is no longer awaiting the result,
		//  and we can safely return.
		ics.inactivityLogger.Infof(
			"[member:%v] inactivity execution is no longer awaiting the "+
				"result; aborting inactivity claim on-chain submission",
			memberIndex,
		)
		return nil
	}

	ics.inactivityLogger.Infof(
		"[member:%v] submitting inactivity claim with [%v] supporting "+
			"member signatures",
		memberIndex,
		len(signatures),
	)

	return ics.chain.SubmitInactivityClaim(
		inactivityClaim,
		inactivityNonce,
		ics.groupMembers,
	)
}
