package entry

import (
	"context"
	"fmt"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/ipfs/go-log"
	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/dkg"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/bls"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
)

var logger = log.Logger("keep-entry")

// RegisterUnmarshallers initializes the given broadcast channel to be able to
// perform relay entry signing protocol interactions by registering all the
// required protocol message unmarshallers.
// The channel has to be initialized before the SignAndSubmit is called.
func RegisterUnmarshallers(channel net.BroadcastChannel) {
	channel.RegisterUnmarshaler(
		func() net.TaggedUnmarshaler { return &SignatureShareMessage{} })
}

// SignAndSubmit triggers the threshold signature process for the
// previous relay entry and publishes the signature to the chain as
// a new relay entry.
func SignAndSubmit(
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	relayChain relayChain.Interface,
	previousEntryBytes []byte,
	honestThreshold int,
	signer *dkg.ThresholdSigner,
	startBlockHeight uint64,
) error {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	receiveChannel := make(chan net.Message, 64)
	channel.Recv(ctx, func(netMessage net.Message) {
		receiveChannel <- netMessage
	})

	previousEntry := new(bn256.G1)
	_, err := previousEntry.Unmarshal(previousEntryBytes)
	if err != nil {
		return err
	}

	selfShare := signer.CalculateSignatureShare(previousEntry)

	go broadcastShare(ctx, signer.MemberID(), selfShare, channel)

	seenShares := make(map[group.MemberIndex]*bn256.G1)
	seenShares[signer.MemberID()] = selfShare

	logger.Debugf(
		"[member:%v] auto-accepting self signature share",
		signer.MemberID(),
	)

	config, err := relayChain.GetConfig()
	if err != nil {
		return err
	}

	timeoutChannel, err := blockCounter.BlockHeightWaiter(
		startBlockHeight + config.RelayEntryTimeout,
	)
	if err != nil {
		return err
	}

	for len(seenShares) < honestThreshold {
		select {
		case netMessage := <-receiveChannel:
			message, ok := netMessage.Payload().(*SignatureShareMessage)
			if !ok || group.IsMessageFromSelf(signer.MemberID(), message) {
				continue
			}

			share, err := extractShare(
				message,
				signer.GroupPublicKeyShares(),
				previousEntry,
			)
			if err != nil {
				logger.Warningf(
					"[member:%v] rejecting signature share from member [%v]: [%v]",
					signer.MemberID(),
					message.senderID,
					err,
				)
				continue
			}

			logger.Debugf(
				"[member:%v] accepting signature share from member [%v]",
				signer.MemberID(),
				message.senderID,
			)

			seenShares[message.senderID] = share
		case <-timeoutChannel:
			return fmt.Errorf("relay entry timed out")
		}
	}

	seenSharesSlice := make([]*bls.SignatureShare, 0)
	for memberID, share := range seenShares {
		signatureShare := &bls.SignatureShare{I: int(memberID), V: share}
		seenSharesSlice = append(seenSharesSlice, signatureShare)
	}

	logger.Infof(
		"[member:%v] restoring signature from [%v] shares",
		signer.MemberID(),
		len(seenSharesSlice),
	)

	signature, err := signer.CompleteSignature(
		seenSharesSlice,
		honestThreshold,
	)
	if err != nil {
		return err
	}

	submitter := &relayEntrySubmitter{
		chain:        relayChain,
		blockCounter: blockCounter,
		index:        signer.MemberID(),
	}

	return submitter.submitRelayEntry(
		signature.Marshal(),
		signer.GroupPublicKeyBytes(),
		startBlockHeight,
	)
}

func broadcastShare(
	ctx context.Context,
	memberID group.MemberIndex,
	share *bn256.G1,
	channel net.BroadcastChannel,
) {
	message := &SignatureShareMessage{
		memberID,
		share.Marshal(),
	}

	if err := channel.Send(ctx, message); err != nil {
		logger.Errorf(
			"[member:%v] could not send signature share: [%v]",
			memberID,
			err,
		)
	}
}

func extractShare(
	message *SignatureShareMessage,
	groupPublicKeyShares map[group.MemberIndex]*bn256.G2,
	previousEntry *bn256.G1,
) (*bn256.G1, error) {
	share := new(bn256.G1)
	_, err := share.Unmarshal(message.shareBytes)
	if err != nil {
		return nil, fmt.Errorf(
			"could not unmarshal signature share: [%v]",
			err,
		)
	}

	publicKeyShare, ok := groupPublicKeyShares[message.senderID]
	if !ok {
		return nil, fmt.Errorf(
			"could not validate signature share; " +
				"group public key share not found",
		)
	}

	if !bls.VerifyG1(publicKeyShare, previousEntry, share) {
		return nil, fmt.Errorf("invalid signature share")
	}

	return share, nil
}
