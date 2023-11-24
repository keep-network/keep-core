package tbtc

import (
	"context"
	"fmt"
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
			if block % coordinationFrequencyBlocks == 0 {
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

// coordinationResult represents the result of the coordination procedure
// executed for the given wallet in the given coordination window.
type coordinationResult struct {
	actionType WalletActionType
	proposal   interface{}
}

func (cr *coordinationResult) String() string {
	return "" // TODO: Implementation.
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
		lock: semaphore.NewWeighted(1),
		signers: signers,
		broadcastChannel: broadcastChannel,
		membershipValidator: membershipValidator,
		protocolLatch: protocolLatch,
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

	// TODO: Implement coordination logic. Determine how to handle window's
	//       context.

	return nil, nil
}