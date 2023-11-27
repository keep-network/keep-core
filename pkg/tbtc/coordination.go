package tbtc

import (
	"context"
	"fmt"
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
			if block%coordinationFrequencyBlocks == 0 {
				// Make sure the current window is not the same as the last one.
				// There is no guarantee that the block channel will not emit
				// the same block again.
				if window := newCoordinationWindow(block); window.isAfter(lastWindow) {
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
	// culprit is the address of the operator that is responsible for the fault.
	culprit chain.Address
	// faultType is the type of the fault.
	faultType CoordinationFaultType
}

func (cf *coordinationFault) String() string {
	return fmt.Sprintf(
		"operator [%s], fault [%s]",
		cf.culprit,
		cf.faultType,
	)
}

// coordinationProposal represents a single action proposal for the given wallet.
type coordinationProposal interface {
	// actionType returns the specific type of the walletAction being subject
	// of this proposal.
	actionType() WalletActionType
	// validityBlocks returns the number of blocks for which the proposal is
	// valid.
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

// coordinationExecutor is responsible for executing the coordination
// procedure for the given wallet.
type coordinationExecutor struct {
	lock *semaphore.Weighted

	signers             []*signer // TODO: Do we need whole signers?
	broadcastChannel    net.BroadcastChannel
	membershipValidator *group.MembershipValidator
	protocolLatch       *generator.ProtocolLatch
}

// newCoordinationExecutor creates a new coordination executor for the
// given wallet.
func newCoordinationExecutor(
	signers []*signer,
	broadcastChannel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
	protocolLatch *generator.ProtocolLatch,
) *coordinationExecutor {
	return &coordinationExecutor{
		lock:                semaphore.NewWeighted(1),
		signers:             signers,
		broadcastChannel:    broadcastChannel,
		membershipValidator: membershipValidator,
		protocolLatch:       protocolLatch,
	}
}

// wallet returns the wallet this executor is responsible for.
func (ce *coordinationExecutor) wallet() wallet {
	// All signers belong to one wallet. Take that wallet from the
	// first signer.
	return ce.signers[0].wallet
}

// coordinate executes the coordination procedure for the given coordination
// window.
func (ce *coordinationExecutor) coordinate(
	window *coordinationWindow,
) (*coordinationResult, error) {
	if lockAcquired := ce.lock.TryAcquire(1); !lockAcquired {
		return nil, errCoordinationExecutorBusy
	}
	defer ce.lock.Release(1)

	// TODO: Implement coordination logic. Remember about:
	//       - Setting up the right context
	//       - Using the protocol latch
	//       - Using the membership validator
	//       Example result:
	result := &coordinationResult{
		wallet:   ce.wallet(),
		window:   window,
		leader:   ce.wallet().signingGroupOperators[0],
		proposal: &noopProposal{},
		faults:   nil,
	}

	return result, nil
}
