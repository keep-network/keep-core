package tbtc

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/keep-network/keep-core/pkg/internal/pb"
	"golang.org/x/exp/slices"
	"math/rand"
	"sort"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"golang.org/x/sync/semaphore"
)

const (
	// coordinationFrequencyBlocks is the number of blocks between two
	// consecutive coordination windows.
	coordinationFrequencyBlocks = 900
	// coordinationActivePhaseDurationBlocks is the number of blocks in the
	// active phase of the coordination window. The active phase is the
	// phase during which the communication between the coordination leader and
	// their followers is allowed.
	coordinationActivePhaseDurationBlocks = 80
	// coordinationPassivePhaseDurationBlocks is the number of blocks in the
	// passive phase of the coordination window. The passive phase is the
	// phase during which communication is not allowed. Participants are
	// expected to validate the result of the coordination and prepare for
	// execution of the proposed wallet action.
	coordinationPassivePhaseDurationBlocks = 20
	// coordinationDurationBlocks is the number of blocks in a single
	// coordination window.
	coordinationDurationBlocks = coordinationActivePhaseDurationBlocks +
		coordinationPassivePhaseDurationBlocks
	// coordinationSafeBlockShift is the number of blocks by which the
	// coordination block is shifted to obtain a safe block whose 32-byte
	// hash can be used as an ingredient for the coordination seed, computed
	// for the given coordination window.
	coordinationSafeBlockShift = 32
	// coordinationHeartbeatProbability is the probability of proposing a
	// heartbeat action during the coordination procedure, assuming no other
	// higher-priority action is proposed.
	coordinationHeartbeatProbability = float64(0.125)
)

// errCoordinationExecutorBusy is an error returned when the coordination
// executor cannot execute the requested coordination due to an ongoing one.
var errCoordinationExecutorBusy = fmt.Errorf("coordination executor is busy")

// coordinationWindow represents a single coordination window. The coordination
// block is the first block of the window.
type coordinationWindow struct {
	// coordinationBlock is the first block of the coordination window.
	coordinationBlock uint64
}

// newCoordinationWindow creates a new coordination window for the given
// coordination block.
func newCoordinationWindow(coordinationBlock uint64) *coordinationWindow {
	return &coordinationWindow{
		coordinationBlock: coordinationBlock,
	}
}

// ActivePhaseEndBlock returns the block number at which the active phase
// of the coordination window ends.
func (cw *coordinationWindow) activePhaseEndBlock() uint64 {
	return cw.coordinationBlock + coordinationActivePhaseDurationBlocks
}

// EndBlock returns the block number at which the coordination window ends.
func (cw *coordinationWindow) endBlock() uint64 {
	return cw.coordinationBlock + coordinationDurationBlocks
}

// isAfter returns true if this coordination window is after the other
// window.
func (cw *coordinationWindow) isAfter(other *coordinationWindow) bool {
	if other == nil {
		return true
	}

	return cw.coordinationBlock > other.coordinationBlock
}

// index returns the index of the coordination window. The index is computed
// by dividing the coordination block number by the coordination frequency.
// A valid index is a positive integer.
//
// For example:
// - window starting at block 900 has index 1
// - window starting at block 1800 has index 2
// - window starting at block 2700 has index 3
//
// If the coordination block number is not a multiple of the coordination
// frequency, the index is 0.
func (cw *coordinationWindow) index() uint64 {
	if cw.coordinationBlock%coordinationFrequencyBlocks == 0 {
		return cw.coordinationBlock / coordinationFrequencyBlocks
	}

	return 0
}

// watchCoordinationWindows watches for new coordination windows and runs
// the given callback when a new window is detected. The callback is run
// in a separate goroutine. It is guaranteed that the callback is not run
// twice for the same window. The context passed as the first parameter
// is used to cancel the watch.
func watchCoordinationWindows(
	ctx context.Context,
	watchBlocksFn func(ctx context.Context) <-chan uint64,
	onWindowFn func(window *coordinationWindow),
) {
	blocksChan := watchBlocksFn(ctx)
	var lastWindow *coordinationWindow

	for {
		select {
		case block := <-blocksChan:
			if window := newCoordinationWindow(block); window.index() > 0 {
				// Make sure the current window is not the same as the last one.
				// There is no guarantee that the block channel will not emit
				// the same block again.
				if window.isAfter(lastWindow) {
					lastWindow = window
					// Run the callback in a separate goroutine to avoid blocking
					// this loop and potentially missing the next block.
					go onWindowFn(window)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

// CoordinationFaultType represents a type of the coordination fault.
type CoordinationFaultType uint8

const (
	// FaultUnknown is a fault type used when the fault type is unknown.
	FaultUnknown CoordinationFaultType = iota
	// FaultLeaderIdleness is a fault type used when the leader was idle, i.e.
	// missed their turn to propose a wallet action.
	FaultLeaderIdleness
	// FaultLeaderMistake is a fault type used when the leader's proposal
	// turned out to be invalid.
	FaultLeaderMistake
	// FaultLeaderImpersonation is a fault type used when the leader was
	// impersonated by another operator who raised their own proposal.
	FaultLeaderImpersonation
)

func (cft CoordinationFaultType) String() string {
	switch cft {
	case FaultUnknown:
		return "Unknown"
	case FaultLeaderIdleness:
		return "LeaderIdleness"
	case FaultLeaderMistake:
		return "FaultLeaderMistake"
	case FaultLeaderImpersonation:
		return "LeaderImpersonation"
	default:
		panic("unknown coordination fault type")
	}
}

// coordinationFault represents a single coordination fault.
type coordinationFault struct {
	culprit   chain.Address // address of the operator responsible for the fault
	faultType CoordinationFaultType
}

func (cf *coordinationFault) String() string {
	return fmt.Sprintf(
		"operator [%s], fault [%s]",
		cf.culprit,
		cf.faultType,
	)
}

// coordinationProposalGenerator is a function that generates a coordination
// proposal based on the given checklist of possible wallet actions.
// The checklist is a list of actions that should be checked for the given
// coordination window. The generator is expected to return a proposal
// for the first action from the checklist that is valid for the given
// wallet's state. If none of the actions are valid, the generator
// should return a noopProposal.
type coordinationProposalGenerator func(
	actionsChecklist []WalletActionType,
) (coordinationProposal, error)

// coordinationProposal represents a single action proposal for the given wallet.
type coordinationProposal interface {
	pb.Marshaler
	pb.Unmarshaler

	// actionType returns the specific type of the walletAction being subject
	// of this proposal.
	actionType() WalletActionType
	// validityBlocks returns the number of blocks for which the proposal is
	// valid. This value SHOULD NOT be marshaled/unmarshaled.
	validityBlocks() uint64
}

// noopProposal is a proposal that does not propose any action.
type noopProposal struct{}

func (np *noopProposal) actionType() WalletActionType {
	return ActionNoop
}

func (np *noopProposal) validityBlocks() uint64 {
	// Panic to make sure that the proposal is not processed by the node.
	panic("noop proposal does not have validity blocks")
}

// coordinationResult represents the result of the coordination procedure
// executed for the given wallet in the given coordination window.
type coordinationResult struct {
	wallet   wallet
	window   *coordinationWindow
	leader   chain.Address
	proposal coordinationProposal
	faults   []*coordinationFault
}

func (cr *coordinationResult) String() string {
	return fmt.Sprintf(
		"wallet [%s], window [%v], leader [%s], proposal [%s], faults [%s]",
		&cr.wallet,
		cr.window.coordinationBlock,
		cr.leader,
		cr.proposal.actionType(),
		cr.faults,
	)
}

// coordinationMessage represents a coordination message sent by the leader
// to their followers during the active phase of the coordination window.
type coordinationMessage struct {
	senderID            group.MemberIndex
	coordinationBlock   uint64
	walletPublicKeyHash [20]byte
	proposal            coordinationProposal
}

func (cm *coordinationMessage) Type() string {
	return "tbtc/coordination_message"
}

// coordinationExecutor is responsible for executing the coordination
// procedure for the given wallet.
type coordinationExecutor struct {
	lock *semaphore.Weighted

	chain Chain

	coordinatedWallet wallet
	membersIndexes    []group.MemberIndex
	operatorAddress   chain.Address

	proposalGenerator coordinationProposalGenerator

	broadcastChannel    net.BroadcastChannel
	membershipValidator *group.MembershipValidator
	protocolLatch       *generator.ProtocolLatch

	waitForBlockFn waitForBlockFn
}

// newCoordinationExecutor creates a new coordination executor for the
// given wallet.
func newCoordinationExecutor(
	chain Chain,
	coordinatedWallet wallet,
	membersIndexes []group.MemberIndex,
	operatorAddress chain.Address,
	proposalGenerator coordinationProposalGenerator,
	broadcastChannel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
	protocolLatch *generator.ProtocolLatch,
	waitForBlockFn waitForBlockFn,
) *coordinationExecutor {
	return &coordinationExecutor{
		lock:                semaphore.NewWeighted(1),
		chain:               chain,
		coordinatedWallet:   coordinatedWallet,
		membersIndexes:      membersIndexes,
		operatorAddress:     operatorAddress,
		proposalGenerator:   proposalGenerator,
		broadcastChannel:    broadcastChannel,
		membershipValidator: membershipValidator,
		protocolLatch:       protocolLatch,
		waitForBlockFn:      waitForBlockFn,
	}
}

// walletPublicKeyHash returns the 20-byte public key hash of the
// coordinated wallet.
func (ce *coordinationExecutor) walletPublicKeyHash() [20]byte {
	return bitcoin.PublicKeyHash(ce.coordinatedWallet.publicKey)
}

// coordinate executes the coordination procedure for the given coordination
// window.
//
// TODO: Add logging and cover with unit tests.
func (ce *coordinationExecutor) coordinate(
	window *coordinationWindow,
) (*coordinationResult, error) {
	if lockAcquired := ce.lock.TryAcquire(1); !lockAcquired {
		return nil, errCoordinationExecutorBusy
	}
	defer ce.lock.Release(1)

	ce.protocolLatch.Lock()
	defer ce.protocolLatch.Unlock()

	// Just in case, check if the window is valid.
	if window.index() == 0 {
		return nil, fmt.Errorf(
			"invalid coordination block [%v]",
			window.coordinationBlock,
		)
	}

	seed, err := ce.coordinationSeed(window.coordinationBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to compute coordination seed: [%v]", err)
	}

	leader := ce.coordinationLeader(seed)

	actionsChecklist := ce.actionsChecklist(window.index(), seed)

	// Set up a context that is cancelled when the active phase of the
	// coordination window ends.
	ctx, cancelCtx := withCancelOnBlock(
		context.Background(),
		window.activePhaseEndBlock(),
		ce.waitForBlockFn,
	)
	defer cancelCtx()

	var proposal coordinationProposal
	if leader == ce.operatorAddress {
		proposal, err = ce.leaderRoutine(
			ctx,
			window.coordinationBlock,
			actionsChecklist,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to execute leader's routine: [%v]",
				err,
			)
		}
	} else {
		proposal, err = ce.followerRoutine()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to execute follower's routine: [%v]",
				err,
			)
		}
	}

	// Just in case, if the proposal is nil, set it to noop.
	if proposal == nil {
		proposal = &noopProposal{}
	}

	result := &coordinationResult{
		wallet:   ce.coordinatedWallet,
		window:   window,
		leader:   leader,
		proposal: proposal,
		faults:   nil, // TODO: Fill coordination faults.
	}

	return result, nil
}

// coordinationSeed computes the coordination seed for the given coordination
// window.
func (ce *coordinationExecutor) coordinationSeed(
	coordinationBlock uint64,
) ([32]byte, error) {
	walletPublicKeyHash := ce.walletPublicKeyHash()

	safeBlockNumber := coordinationBlock - coordinationSafeBlockShift
	safeBlockHash, err := ce.chain.GetBlockHashByNumber(safeBlockNumber)
	if err != nil {
		return [32]byte{}, fmt.Errorf(
			"failed to get safe block hash: [%v]",
			err,
		)
	}

	return sha256.Sum256(
		append(
			walletPublicKeyHash[:],
			safeBlockHash[:]...,
		),
	), nil
}

// coordinationLeader returns the address of the coordination leader for the
// given coordination seed.
func (ce *coordinationExecutor) coordinationLeader(seed [32]byte) chain.Address {
	// First, take all operators backing the wallet.
	allOperators := chain.Addresses(ce.coordinatedWallet.signingGroupOperators)

	// Determine a list of unique operators.
	uniqueOperators := make([]chain.Address, 0)
	for operator := range allOperators.Set() {
		uniqueOperators = append(uniqueOperators, operator)
	}

	// Sort the list of unique operators in ascending order.
	sort.Slice(
		uniqueOperators,
		func(i, j int) bool {
			return uniqueOperators[i] < uniqueOperators[j]
		},
	)

	// #nosec G404 (insecure random number source (rand))
	// Shuffling operators does not require secure randomness.
	// Use first 8 bytes of the seed to initialize the RNG.
	rng := rand.New(rand.NewSource(int64(binary.BigEndian.Uint64(seed[:8]))))

	// Shuffle the list of unique operators.
	rng.Shuffle(
		len(uniqueOperators),
		func(i, j int) {
			uniqueOperators[i], uniqueOperators[j] =
				uniqueOperators[j], uniqueOperators[i]
		},
	)

	// The first operator in the shuffled list is the leader.
	return uniqueOperators[0]
}

// actionsChecklist returns a list of wallet actions that should be checked
// for the given coordination window. Returns nil for incorrect coordination
// windows whose index is 0.
func (ce *coordinationExecutor) actionsChecklist(
	windowIndex uint64,
	seed [32]byte,
) []WalletActionType {
	// Return nil checklist for incorrect coordination windows.
	if windowIndex == 0 {
		return nil
	}

	var actions []WalletActionType

	// Redemption action is a priority action and should be checked on every
	// coordination window.
	actions = append(actions, ActionRedemption)

	// Other actions should be checked with a lower frequency. The default
	// frequency is every 16 coordination windows.
	frequencyWindows := uint64(16)

	// TODO: Increase frequency for the active wallet.
	if windowIndex%frequencyWindows == 0 {
		actions = append(actions, ActionDepositSweep)
	}

	if windowIndex%frequencyWindows == 0 {
		actions = append(actions, ActionMovedFundsSweep)
	}

	// TODO: Increase frequency for old wallets.
	if windowIndex%frequencyWindows == 0 {
		actions = append(actions, ActionMovingFunds)
	}

	// #nosec G404 (insecure random number source (rand))
	// Drawing a decision about heartbeat does not require secure randomness.
	// Use first 8 bytes of the seed to initialize the RNG.
	rng := rand.New(rand.NewSource(int64(binary.BigEndian.Uint64(seed[:8]))))
	if rng.Float64() < coordinationHeartbeatProbability {
		actions = append(actions, ActionHeartbeat)
	}

	return actions
}

// leaderRoutine executes the leader's routine for the given coordination
// window. The routine generates a proposal and broadcasts it to the followers.
// It returns the generated proposal or an error if the routine failed.
func (ce *coordinationExecutor) leaderRoutine(
	ctx context.Context,
	coordinationBlock uint64,
	actionsChecklist []WalletActionType,
) (coordinationProposal, error) {
	proposal, err := ce.proposalGenerator(actionsChecklist)
	if err != nil {
		return nil, fmt.Errorf("failed to generate proposal: [%v]", err)
	}

	// Sort members indexes in ascending order, just in case. Choose the first
	// member as the sender of the coordination message.
	membersIndexes := append([]group.MemberIndex{}, ce.membersIndexes...)
	slices.Sort(membersIndexes)
	senderID := membersIndexes[0]

	message := &coordinationMessage{
		senderID:            senderID,
		coordinationBlock:   coordinationBlock,
		walletPublicKeyHash: ce.walletPublicKeyHash(),
		proposal:            proposal,
	}

	err = ce.broadcastChannel.Send(
		ctx,
		message,
		net.BackoffRetransmissionStrategy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send coordination message: [%v]", err)
	}

	return proposal, nil
}

// followerRoutine executes the follower's routine for the given coordination
// window. The routine listens for the coordination message from the leader and
// validates it. If the leader's proposal is valid, it returns the received
// proposal. Returns an error if the routine failed.
func (ce *coordinationExecutor) followerRoutine() (coordinationProposal, error) {
	// TODO: Implement the follower routine.
	return nil, nil
}
