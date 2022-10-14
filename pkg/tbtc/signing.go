package tbtc

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"sort"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tbtc/gen/pb"
	"github.com/keep-network/keep-core/pkg/tecdsa/retry"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
	"google.golang.org/protobuf/proto"
	"golang.org/x/exp/slices"
)

// signingAnnouncementMessage represents a message that is used to announce
// member's participation in the given signing attempt for the given message.
type signingAnnouncementMessage struct {
	senderID      group.MemberIndex
	message       *big.Int
	attemptNumber uint64
}

func (sam *signingAnnouncementMessage) Marshal() ([]byte, error) {
	return proto.Marshal(&pb.SigningAnnouncementMessage{
		SenderID:      uint32(sam.senderID),
		Message:       sam.message.Bytes(),
		AttemptNumber: sam.attemptNumber,
	})
}

func (sam *signingAnnouncementMessage) Unmarshal(bytes []byte) error {
	pbMessage := pb.SigningAnnouncementMessage{}
	if err := proto.Unmarshal(bytes, &pbMessage); err != nil {
		return fmt.Errorf(
			"failed to unmarshal SigningAnnouncementMessage: [%v]",
			err,
		)
	}

	if senderID := pbMessage.SenderID; senderID > group.MaxMemberIndex {
		return fmt.Errorf("invalid member index value: [%v]", senderID)
	} else {
		sam.senderID = group.MemberIndex(senderID)
	}

	sam.message = new(big.Int).SetBytes(pbMessage.Message)
	sam.attemptNumber = pbMessage.AttemptNumber

	return nil
}

func (sam *signingAnnouncementMessage) Type() string {
	return "tbtc/signing_announcement_message"
}

// signingRetryLoop is a struct that encapsulates the signing retry logic.
type signingRetryLoop struct {
	signingGroupMemberIndex group.MemberIndex
	signingGroupOperators   chain.Addresses

	chainConfig *ChainConfig

	attemptCounter    uint
	attemptStartBlock uint64
	attemptSeed       int64

	delayBlocks uint64
}

func newSigningRetryLoop(
	message *big.Int,
	initialStartBlock uint64,
	signingGroupMemberIndex group.MemberIndex,
	signingGroupOperators chain.Addresses,
	chainConfig *ChainConfig,
) *signingRetryLoop {
	// Compute the 8-byte seed needed for the random retry algorithm. We take
	// the first 8 bytes of the hash of the signed message. This allows us to
	// not care in this piece of the code about the length of the message and
	// how this message is proposed.
	messageSha256 := sha256.Sum256(message.Bytes())
	attemptSeed := int64(binary.BigEndian.Uint64(messageSha256[:8]))

	return &signingRetryLoop{
		signingGroupMemberIndex: signingGroupMemberIndex,
		signingGroupOperators:   signingGroupOperators,
		chainConfig:             chainConfig,
		attemptCounter:          0,
		attemptStartBlock:       initialStartBlock,
		attemptSeed:             attemptSeed,
		delayBlocks:             5,
	}
}

// signingAttemptParams represents parameters of a signing attempt.
type signingAttemptParams struct {
	number                 uint
	startBlock             uint64
	excludedMembersIndexes []group.MemberIndex
}

// signingAttemptFn represents a function performing a signing attempt.
type signingAttemptFn func(*signingAttemptParams) (*signing.Result, error)

// waitWithSigningAttemptFn represents a function blocking the attempt execution
// until the given block height.
type waitWithSigningAttemptFn func(context.Context, uint64) error

// start begins the signing retry loop using the given signing attempt function.
// The retry loop terminates when the signing result is produced or the ctx
// parameter is done, whatever comes first.
func (srl *signingRetryLoop) start(
	ctx context.Context,
	waitWithSigningAttemptFn waitWithSigningAttemptFn,
	signingAttemptFn signingAttemptFn,
) (*signing.Result, error) {
	// We want to take the random subset right away for the first attempt.
	qualifiedOperatorsSet, err := srl.qualifiedOperatorsSet()
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get qualified operators for attempt [%v]: [%w]",
			srl.attemptCounter+1,
			err,
		)
	}

	for {
		srl.attemptCounter++

		// In order to start attempts >1 in the right place, we need to
		// determine how many blocks were taken by previous attempts. We assume
		// the worst case that each attempt failed at the end of the signing
		// protocol.
		//
		// That said, we need to increment the previous attempt start
		// block by the number of blocks equal to the protocol duration and
		// by some additional delay blocks. We need a small fixed delay in
		// order to mitigate all corner cases where the actual attempt duration
		// was slightly longer than the expected duration determined by the
		// signing.ProtocolBlocks function.
		//
		// For example, the attempt may fail at the end of the protocol but the
		// error is returned after some time and more blocks than expected are
		// mined in the meantime.
		if srl.attemptCounter > 1 {
			srl.attemptStartBlock = srl.attemptStartBlock +
				signing.ProtocolBlocks() +
				srl.delayBlocks
		}

		// Exclude all members controlled by the operators that were not
		// qualified for the current attempt.
		includedMembersIndexes := make([]group.MemberIndex, 0)
		excludedMembersIndexes := make([]group.MemberIndex, 0)
		for i, operator := range srl.signingGroupOperators {
			memberIndex := group.MemberIndex(i + 1)

			if qualifiedOperatorsSet[operator] {
				includedMembersIndexes = append(
					includedMembersIndexes,
					memberIndex,
				)
			} else {
				excludedMembersIndexes = append(
					excludedMembersIndexes,
					memberIndex,
				)
			}
		}

		// Make sure we always use just the smallest required count of
		// signing members for performance reasons
		if len(includedMembersIndexes) > srl.chainConfig.HonestThreshold {
			// #nosec G404 (insecure random number source (rand))
			// Shuffling does not require secure randomness.
			rng := rand.New(rand.NewSource(
				srl.attemptSeed + int64(srl.attemptCounter),
			))
			// Sort in ascending order just in case.
			sort.Slice(includedMembersIndexes, func(i, j int) bool {
				return includedMembersIndexes[i] < includedMembersIndexes[j]
			})
			// Shuffle the included members slice to randomize the
			// selection of additionally excluded members.
			rng.Shuffle(len(includedMembersIndexes), func(i, j int) {
				includedMembersIndexes[i], includedMembersIndexes[j] =
					includedMembersIndexes[j], includedMembersIndexes[i]
			})
			// Get the surplus of included members and add them to
			// the excluded members list.
			excludedMembersIndexes = append(
				excludedMembersIndexes,
				includedMembersIndexes[srl.chainConfig.HonestThreshold:]...,
			)
			// Sort the resulting excluded members list in ascending order.
			sort.Slice(excludedMembersIndexes, func(i, j int) bool {
				return excludedMembersIndexes[i] < excludedMembersIndexes[j]
			})
		}

		attemptSkipped := slices.Contains(
			excludedMembersIndexes,
			srl.signingGroupMemberIndex,
		)

		// Wait for the right moment to execute the signingAttemptFn, as
		// calculated in srl.attemptStartBlock.
		err := waitWithSigningAttemptFn(ctx, srl.attemptStartBlock)
		if err != nil {
			return nil, fmt.Errorf(
				"failed waiting for block [%v] for attempt [%v]: [%v]",
				srl.attemptStartBlock,
				srl.attemptCounter,
				err,
			)
		}

		// Check the loop stop signal.
		if ctx.Err() != nil {
			return nil, nil
		}

		var result *signing.Result
		var attemptErr error

		if !attemptSkipped {
			result, attemptErr = signingAttemptFn(&signingAttemptParams{
				number:                 srl.attemptCounter,
				startBlock:             srl.attemptStartBlock,
				excludedMembersIndexes: excludedMembersIndexes,
			})
		}

		if attemptSkipped || attemptErr != nil {
			var err error
			qualifiedOperatorsSet, err = srl.qualifiedOperatorsSet()
			if err != nil {
				return nil, fmt.Errorf(
					"cannot get qualified operators for attempt [%v]: [%w]",
					srl.attemptCounter+1,
					err,
				)
			}

			continue
		}

		return result, nil
	}
}

// qualifiedOperatorsSet returns a set of operators qualified to participate
// in the given signing attempt.
func (srl *signingRetryLoop) qualifiedOperatorsSet() (
	map[chain.Address]bool,
	error,
) {
	qualifiedOperators, err := retry.EvaluateRetryParticipantsForSigning(
		srl.signingGroupOperators,
		srl.attemptSeed,
		srl.attemptCounter,
		uint(srl.chainConfig.HonestThreshold),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"random operator selection failed: [%w]",
			err,
		)
	}

	return chain.Addresses(qualifiedOperators).Set(), nil
}
