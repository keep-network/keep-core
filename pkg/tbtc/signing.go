package tbtc

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tbtc/gen/pb"
	"github.com/keep-network/keep-core/pkg/tecdsa/retry"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"math/big"
	"math/rand"
	"sort"
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
	logger log.StandardLogger

	message *big.Int

	signingGroupMemberIndex group.MemberIndex
	signingGroupOperators   chain.Addresses

	chainConfig *ChainConfig

	announcementDelayBlocks  uint64
	announcementActiveBlocks uint64

	attemptCounter     uint
	attemptStartBlock  uint64
	attemptSeed        int64
	attemptDelayBlocks uint64

	broadcastChannel net.BroadcastChannel

	membershipValidator *group.MembershipValidator
}

func newSigningRetryLoop(
	logger log.StandardLogger,
	message *big.Int,
	initialStartBlock uint64,
	signingGroupMemberIndex group.MemberIndex,
	signingGroupOperators chain.Addresses,
	chainConfig *ChainConfig,
	broadcastChannel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
) *signingRetryLoop {
	// Compute the 8-byte seed needed for the random retry algorithm. We take
	// the first 8 bytes of the hash of the signed message. This allows us to
	// not care in this piece of the code about the length of the message and
	// how this message is proposed.
	messageSha256 := sha256.Sum256(message.Bytes())
	attemptSeed := int64(binary.BigEndian.Uint64(messageSha256[:8]))

	broadcastChannel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &signingAnnouncementMessage{}
	})

	return &signingRetryLoop{
		logger:                   logger,
		message:                  message,
		signingGroupMemberIndex:  signingGroupMemberIndex,
		signingGroupOperators:    signingGroupOperators,
		chainConfig:              chainConfig,
		announcementDelayBlocks:  1,
		announcementActiveBlocks: 5,
		attemptCounter:           0,
		attemptStartBlock:        initialStartBlock,
		attemptSeed:              attemptSeed,
		attemptDelayBlocks:       5,
		broadcastChannel:         broadcastChannel,
		membershipValidator:      membershipValidator,
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

// waitForBlockFn represents a function blocking the execution until the given
// block height.
type waitForBlockFn func(context.Context, uint64) error

// start begins the signing retry loop using the given signing attempt function.
// The retry loop terminates when the signing result is produced or the ctx
// parameter is done, whatever comes first.
func (srl *signingRetryLoop) start(
	ctx context.Context,
	waitForBlockFn waitForBlockFn,
	signingAttemptFn signingAttemptFn,
) (*signing.Result, error) {
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
				srl.announcementDelayBlocks +
				srl.announcementActiveBlocks +
				signing.ProtocolBlocks() +
				srl.attemptDelayBlocks
		}

		announcementStartBlock := srl.attemptStartBlock + srl.announcementDelayBlocks
		err := waitForBlockFn(ctx, announcementStartBlock)
		if err != nil {
			return nil, fmt.Errorf(
				"failed waiting for announcement start block [%v] "+
					"for attempt [%v]: [%v]",
				announcementStartBlock,
				srl.attemptCounter,
				err,
			)
		}

		// Check the loop stop signal.
		if ctx.Err() != nil {
			return nil, nil
		}

		// Set up the announcement phase stop signal.
		announceCtx, cancelAnnounceCtx := context.WithCancel(ctx)
		announcementEndBlock := announcementStartBlock + srl.announcementActiveBlocks
		go func() {
			defer cancelAnnounceCtx()

			if err := waitForBlockFn(ctx, announcementEndBlock); err != nil {
				srl.logger.Errorf(
					"[member:%v] failed waiting for announcement end "+
						"block [%v] for attempt [%v]: [%v]",
					srl.signingGroupMemberIndex,
					announcementEndBlock,
					srl.attemptCounter,
					err,
				)
			}
		}()

		srl.logger.Infof(
			"[member:%v] starting announcement phase for attempt [%v]",
			srl.signingGroupMemberIndex,
			srl.attemptCounter,
		)

		announcements, err := srl.announceAttempt(announceCtx)
		if err != nil {
			srl.logger.Warnf(
				"[member:%v] announcement for attempt [%v] "+
					"failed: [%v]; starting next attempt",
				srl.signingGroupMemberIndex,
				srl.attemptCounter,
				err,
			)
			continue
		}

		srl.logger.Infof(
			"[member:%v] completed announcement phase for attempt [%v] "+
				"and [%v] other members announced readiness",
			srl.signingGroupMemberIndex,
			srl.attemptCounter,
			len(announcements),
		)

		qualifiedOperatorsSet, err := srl.qualifiedOperatorsSet(announcements)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot get qualified operators for attempt [%v]: [%w]",
				srl.attemptCounter,
				err,
			)
		}

		// Exclude all members controlled by the operators that were not
		// qualified for the current attempt.
		excludedMembersIndexes := srl.excludedMembersIndexes(qualifiedOperatorsSet)

		attemptSkipped := slices.Contains(
			excludedMembersIndexes,
			srl.signingGroupMemberIndex,
		)

		var result *signing.Result
		var attemptErr error

		if !attemptSkipped {
			result, attemptErr = signingAttemptFn(&signingAttemptParams{
				number:                 srl.attemptCounter,
				startBlock:             announcementEndBlock,
				excludedMembersIndexes: excludedMembersIndexes,
			})
		} else {
			srl.logger.Infof(
				"[member:%v] attempt [%v] skipped",
				srl.signingGroupMemberIndex,
				srl.attemptCounter,
			)
		}

		if attemptSkipped || attemptErr != nil {
			continue
		}

		return result, nil
	}
}

// qualifiedOperatorsSet returns a set of operators qualified to participate
// in the given signing attempt. The set of qualified operators is taken
// from the set of active operators who announced readiness through
// their controlled signing group members.
func (srl *signingRetryLoop) qualifiedOperatorsSet(
	announcements map[group.MemberIndex]bool,
) (map[chain.Address]bool, error) {
	// The retry algorithm expects that we count retries from 0. Since
	// the first invocation of the algorithm will be for `attemptCounter == 1`
	// we need to subtract one while determining the number of the given retry.
	retryCount := srl.attemptCounter - 1

	var announcedSigningGroupOperators []chain.Address
	for i, operator := range srl.signingGroupOperators {
		memberIndex := group.MemberIndex(i + 1)

		if announcements[memberIndex] {
			announcedSigningGroupOperators = append(
				announcedSigningGroupOperators,
				operator,
			)
		}
	}

	qualifiedOperators, err := retry.EvaluateRetryParticipantsForSigning(
		announcedSigningGroupOperators,
		srl.attemptSeed,
		retryCount,
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

// excludedMembersIndexes returns a list of excluded members' indexes for
// the given qualified operators set.
func (srl *signingRetryLoop) excludedMembersIndexes(
	qualifiedOperatorsSet map[chain.Address]bool,
) []group.MemberIndex {
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

	return excludedMembersIndexes
}

// announceAttempt broadcasts the member's readiness announcement for the
// given signing attempt and listens for announcements from other signing
// group members. This function keeps working until the ctx parameter is done
// and returns successfully only if the total number of ready members
// (including the announcing member) is equal to or grater than the honest
// threshold parameter.
func (srl *signingRetryLoop) announceAttempt(ctx context.Context) (
	map[group.MemberIndex]bool,
	error,
) {
	err := srl.broadcastChannel.Send(ctx, &signingAnnouncementMessage{
		senderID:      srl.signingGroupMemberIndex,
		message:       srl.message,
		attemptNumber: uint64(srl.attemptCounter),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot send announcement message: [%w]", err)
	}

	messagesChan := make(chan net.Message, len(srl.signingGroupOperators))
	srl.broadcastChannel.Recv(ctx, func(message net.Message) {
		messagesChan <- message
	})

	announcements := make(map[group.MemberIndex]bool)

loop:
	for {
		select {
		case netMessage := <-messagesChan:
			announcement, ok := netMessage.Payload().(*signingAnnouncementMessage)
			if !ok {
				continue
			}

			if announcement.senderID == srl.signingGroupMemberIndex {
				continue
			}

			if !srl.membershipValidator.IsValidMembership(
				announcement.senderID,
				netMessage.SenderPublicKey(),
			) {
				continue
			}

			if announcement.message.Cmp(srl.message) != 0 {
				continue
			}

			if announcement.attemptNumber != uint64(srl.attemptCounter) {
				continue
			}

			announcements[announcement.senderID] = true
		case <-ctx.Done():
			break loop
		}
	}

	// The total number of operating members for the given attempt is the count
	// of the received announcements plus the member itself.
	operatingMembers := len(announcements) + 1
	if operatingMembers < srl.chainConfig.HonestThreshold {
		return nil, fmt.Errorf(
			"operating members count is lesser than the honest threshold",
		)
	}

	return announcements, nil
}
