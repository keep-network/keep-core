package entry

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/beacon/event"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/ipfs/go-log/v2"
	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/beacon/dkg"
	"github.com/keep-network/keep-core/pkg/bls"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// RegisterUnmarshallers initializes the given broadcast channel to be able to
// perform relay entry signing protocol interactions by registering all the
// required protocol message unmarshallers.
// The channel has to be initialized before the SignAndSubmit is called.
func RegisterUnmarshallers(channel net.BroadcastChannel) {
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &SignatureShareMessage{}
	})
}

// SignAndSubmit triggers the threshold signature process for the
// previous relay entry and publishes the signature to the chain as
// a new relay entry.
func SignAndSubmit(
	logger log.StandardLogger,
	blockCounter chain.BlockCounter,
	channel net.BroadcastChannel,
	beaconChain beaconchain.Interface,
	previousEntryBytes []byte,
	honestThreshold int,
	signer *dkg.ThresholdSigner,
	startBlockHeight uint64,
) error {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	relayEntrySubmittedChannel := make(chan uint64)
	subscription := beaconChain.OnRelayEntrySubmitted(
		func(event *event.RelayEntrySubmitted) {
			relayEntrySubmittedChannel <- event.BlockNumber
		},
	)
	defer subscription.Unsubscribe()

	chainConfig := beaconChain.GetConfig()

	relayEntryTimeoutChannel, err := blockCounter.BlockHeightWaiter(
		startBlockHeight + chainConfig.RelayEntryTimeout,
	)
	if err != nil {
		return err
	}

	previousEntry := new(bn256.G1)
	_, err = previousEntry.Unmarshal(previousEntryBytes)
	if err != nil {
		return err
	}

	selfShare := signer.CalculateSignatureShare(previousEntry)

	go broadcastShare(ctx, logger, signer.MemberID(), selfShare, channel)

	receiveChannel := make(chan net.Message, 64)
	channel.Recv(ctx, func(netMessage net.Message) {
		receiveChannel <- netMessage
	})

	receivedValidShares := map[group.MemberIndex]*bn256.G1{
		signer.MemberID(): selfShare,
	}

	// Run the message loop until the number of received and valid signature
	// shares is equal to the honest threshold. Message loop will be also
	// terminated if an other member submits the result or the relay entry
	// timeout block is reached.
	for len(receivedValidShares) < honestThreshold {
		select {
		case netMessage := <-receiveChannel:
			message, ok := netMessage.Payload().(*SignatureShareMessage)
			if !ok || signer.MemberID() == message.SenderID() {
				continue
			}

			share, err := extractAndValidateShare(
				message,
				signer.GroupPublicKeyShares(),
				previousEntry,
			)
			if err != nil {
				logger.Warnf(
					"[member:%v] rejecting signature share from "+
						"member [%v]: [%v]",
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

			receivedValidShares[message.senderID] = share
		case blockNumber := <-relayEntrySubmittedChannel:
			logger.Infof(
				"[member:%v] leaving message loop; "+
					"relay entry submitted by other member at block [%v]",
				signer.MemberID(),
				blockNumber,
			)
			return nil
		case blockNumber := <-relayEntryTimeoutChannel:
			return fmt.Errorf(
				"relay entry timed out at block [%v]; received [%v] valid signature shares",
				blockNumber,
				len(receivedValidShares),
			)
		}
	}

	signature, err := completeSignature(logger, signer, receivedValidShares, honestThreshold)
	if err != nil {
		return err
	}

	submitter := &relayEntrySubmitter{
		logger:       logger,
		chain:        beaconChain,
		blockCounter: blockCounter,
		index:        signer.MemberID(),
	}

	// relayEntrySubmittedChannel and relayEntryTimeoutChannel are passed to
	// the submitter. This should be done because no entry submission or
	// timeout signal appeared while executing the message loop. There is
	// still a possibility those signals appear in the future so the submitter
	// must be aware of them and break the execution if they occur.
	return submitter.submitRelayEntry(
		signature.Marshal(),
		signer.GroupPublicKeyBytes(),
		startBlockHeight,
		relayEntrySubmittedChannel,
		relayEntryTimeoutChannel,
	)
}

func broadcastShare(
	ctx context.Context,
	logger log.StandardLogger,
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

func extractAndValidateShare(
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
				"group public key share for sender not found",
		)
	}

	if !bls.VerifyG1(publicKeyShare, previousEntry, share) {
		return nil, fmt.Errorf("invalid signature share")
	}

	return share, nil
}

func completeSignature(
	logger log.StandardLogger,
	signer *dkg.ThresholdSigner,
	shares map[group.MemberIndex]*bn256.G1,
	honestThreshold int,
) (*bn256.G1, error) {
	signatureShares := make([]*bls.SignatureShare, 0)
	for memberID, share := range shares {
		signatureShare := &bls.SignatureShare{I: int(memberID), V: share}
		signatureShares = append(signatureShares, signatureShare)
	}

	logger.Infof(
		"[member:%v] restoring signature from [%v] shares",
		signer.MemberID(),
		len(signatureShares),
	)

	signature, err := signer.CompleteSignature(
		signatureShares,
		honestThreshold,
	)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
