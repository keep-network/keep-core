package chain

import (
	"github.com/keep-network/keep-core/pkg/beacon/relay/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/gen/async"
)

// Interface represents the interface that the relay expects to interact
// with the anchoring blockchain on.
type Interface interface {
	// GetConfig returns the expected configuration of the threshold relay.
	GetConfig() (config.Chain, error)
	// SubmitGroupPublicKey submits a 96-byte BLS public key to the blockchain,
	// associated with a string groupID. On-chain errors are
	// are reported through the promise.
	SubmitGroupPublicKey(groupID string, key [96]byte) *async.GroupRegistrationPromise
	// SubmitRelayEntry submits an entry in the threshold relay and returns a
	// promise to track the submission result. The promise is fulfilled with
	// the entry as seen on-chain, or failed if there is an error submitting
	// the entry.
	SubmitRelayEntry(entry *event.Entry) *async.RelayEntryPromise
	// OnRelayEntryGenerated is a callback that is invoked when an on-chain
	// notification of a new, valid relay entry is seen.
	OnRelayEntryGenerated(handle func(entry *event.Entry))
	// OnRelayEntryRequested is a callback that is invoked when an on-chain
	// notification of a new, valid relay request is seen.
	OnRelayEntryRequested(func(request *event.Request))
	// AddStaker is a temporary function for Milestone 1 that
	// adds a staker to the group contract.
	AddStaker(groupMemberID string) *async.StakerRegistrationPromise
	// GetStakerList is a temporary function for Milestone 1 that
	// gets back the list of stakers.
	GetStakerList() ([]string, error)
	// OnGroupRegistered is a callback that is invoked when an on-chain
	// notification of a new, valid group being registered is seen.
	OnGroupRegistered(handle func(key *event.GroupRegistration))

	/* Preliminary - the interface for creating a relay request.

		// Will call .../contracts/solidity/contracts/KeepRandomBeaconImplV1.sol;
	    // function requestRelayEntry(uint256 _blockReward, uint256 _seed) public payable returns (uint256 requestID) {
		// through the proxy.
		// Note: using 'seed' as a big.Int will loose/drop any leading 0's in the seed.
		// Note: Having a return value `returns (uint256 requestID) {` will not return this value through the
		//       transaction.  The requestID can be found in the event:
		//       RelayEntryRequested(requestID, msg.value, _blockReward, _seed, block.number);
		RequestRelayEntry(blockReward, seed *big.Int) *async.RelayEntryRequestedPromise

		// Need to call .../contracts/solidity/contracts/KeepRandomBeaconImplV1.sol;
	    // function initialize(address _stakingProxy, uint256 _minPayment, uint256 _minStake, uint256 _withdrawalDelay)
		// This call only needs to happen once for each contract load.
		// through the proxy once - to setup the contract.   This is not happening in the JS setup as not all this
		// data is available and contract build time.
		InitializeKeepRandomBeacon(stakingProxy string, minPayment, minStake, withdrawalDelay *big.Int) *async.RandomBeaconInitialized

	*/
}
