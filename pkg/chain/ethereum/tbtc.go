package ethereum

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"time"

	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/pkg/bitcoin"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/chain"
	ecdsaabi "github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen/abi"
	ecdsacontract "github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen/contract"
	tbtcabi "github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/abi"
	tbtccontract "github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/contract"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/subscription"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
)

// Definitions of contract names.
const (
	// TODO: The WalletRegistry address is taken from the Bridge contract.
	//       Remove the possibility of passing it through the config.
	WalletRegistryContractName    = "WalletRegistry"
	BridgeContractName            = "Bridge"
	MaintainerProxyContractName   = "MaintainerProxy"
	WalletCoordinatorContractName = "WalletCoordinator"
)

// TbtcChain represents a TBTC-specific chain handle.
type TbtcChain struct {
	*baseChain

	bridge            *tbtccontract.Bridge
	maintainerProxy   *tbtccontract.MaintainerProxy
	walletRegistry    *ecdsacontract.WalletRegistry
	sortitionPool     *ecdsacontract.EcdsaSortitionPool
	walletCoordinator *tbtccontract.WalletCoordinator
}

// NewTbtcChain construct a new instance of the TBTC-specific Ethereum
// chain handle.
func newTbtcChain(
	config ethereum.Config,
	baseChain *baseChain,
) (*TbtcChain, error) {
	bridgeAddress, err := config.ContractAddress(BridgeContractName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve %s contract address: [%v]",
			BridgeContractName,
			err,
		)
	}

	bridge, err :=
		tbtccontract.NewBridge(
			bridgeAddress,
			baseChain.chainID,
			baseChain.key,
			baseChain.client,
			baseChain.nonceManager,
			baseChain.miningWaiter,
			baseChain.blockCounter,
			baseChain.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach to Bridge contract: [%v]",
			err,
		)
	}

	maintainerProxyAddress, err := config.ContractAddress(MaintainerProxyContractName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve %s contract address: [%v]",
			MaintainerProxyContractName,
			err,
		)
	}

	maintainerProxy, err :=
		tbtccontract.NewMaintainerProxy(
			maintainerProxyAddress,
			baseChain.chainID,
			baseChain.key,
			baseChain.client,
			baseChain.nonceManager,
			baseChain.miningWaiter,
			baseChain.blockCounter,
			baseChain.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach to MaintainerProxy contract: [%v]",
			err,
		)
	}

	references, err := bridge.ContractReferences()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get contract references from Bridge: [%v]",
			err,
		)
	}

	walletRegistryAddress := references.EcdsaWalletRegistry

	walletRegistry, err :=
		ecdsacontract.NewWalletRegistry(
			walletRegistryAddress,
			baseChain.chainID,
			baseChain.key,
			baseChain.client,
			baseChain.nonceManager,
			baseChain.miningWaiter,
			baseChain.blockCounter,
			baseChain.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach to WalletRegistry contract: [%v]",
			err,
		)
	}

	sortitionPoolAddress, err := walletRegistry.SortitionPool()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get sortition pool address: [%v]",
			err,
		)
	}

	sortitionPool, err :=
		ecdsacontract.NewEcdsaSortitionPool(
			sortitionPoolAddress,
			baseChain.chainID,
			baseChain.key,
			baseChain.client,
			baseChain.nonceManager,
			baseChain.miningWaiter,
			baseChain.blockCounter,
			baseChain.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach to EcdsaSortitionPool contract: [%v]",
			err,
		)
	}

	walletCoordinatorAddress, err := config.ContractAddress(
		WalletCoordinatorContractName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve %s contract address: [%v]",
			WalletCoordinatorContractName,
			err,
		)
	}

	walletCoordinator, err :=
		tbtccontract.NewWalletCoordinator(
			walletCoordinatorAddress,
			baseChain.chainID,
			baseChain.key,
			baseChain.client,
			baseChain.nonceManager,
			baseChain.miningWaiter,
			baseChain.blockCounter,
			baseChain.transactionMutex,
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach to WalletCoordinator contract: [%v]",
			err,
		)
	}

	return &TbtcChain{
		baseChain:         baseChain,
		bridge:            bridge,
		maintainerProxy:   maintainerProxy,
		walletRegistry:    walletRegistry,
		sortitionPool:     sortitionPool,
		walletCoordinator: walletCoordinator,
	}, nil
}

// Staking returns address of the TokenStaking contract the WalletRegistry is
// connected to.
func (tc *TbtcChain) Staking() (chain.Address, error) {
	stakingContractAddress, err := tc.walletRegistry.Staking()
	if err != nil {
		return "", fmt.Errorf(
			"failed to get the token staking address: [%w]",
			err,
		)
	}

	return chain.Address(stakingContractAddress.String()), nil
}

// IsRecognized checks whether the given operator is recognized by the TbtcChain
// as eligible to join the network. If the operator has a stake delegation or
// had a stake delegation in the past, it will be recognized.
func (tc *TbtcChain) IsRecognized(operatorPublicKey *operator.PublicKey) (bool, error) {
	operatorAddress, err := operatorPublicKeyToChainAddress(operatorPublicKey)
	if err != nil {
		return false, fmt.Errorf(
			"cannot convert from operator key to chain address: [%v]",
			err,
		)
	}

	stakingProvider, err := tc.walletRegistry.OperatorToStakingProvider(
		operatorAddress,
	)
	if err != nil {
		return false, fmt.Errorf(
			"failed to map operator [%v] to a staking provider: [%v]",
			operatorAddress,
			err,
		)
	}

	if (stakingProvider == common.Address{}) {
		return false, nil
	}

	// Check if the staking provider has an owner. This check ensures that there
	// is/was a stake delegation for the given staking provider.
	_, _, _, hasStakeDelegation, err := tc.baseChain.RolesOf(
		chain.Address(stakingProvider.Hex()),
	)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check stake delegation for staking provider [%v]: [%v]",
			stakingProvider,
			err,
		)
	}

	if !hasStakeDelegation {
		return false, nil
	}

	return true, nil
}

// OperatorToStakingProvider returns the staking provider address for the
// operator. If the staking provider has not been registered for the
// operator, the returned address is empty and the boolean flag is set to
// false. If the staking provider has been registered, the address is not
// empty and the boolean flag indicates true.
func (tc *TbtcChain) OperatorToStakingProvider() (chain.Address, bool, error) {
	stakingProvider, err := tc.walletRegistry.OperatorToStakingProvider(tc.key.Address)
	if err != nil {
		return "", false, fmt.Errorf(
			"failed to map operator [%v] to a staking provider: [%v]",
			tc.key.Address,
			err,
		)
	}

	if (stakingProvider == common.Address{}) {
		return "", false, nil
	}

	return chain.Address(stakingProvider.Hex()), true, nil
}

// EligibleStake returns the current value of the staking provider's
// eligible stake. Eligible stake is defined as the currently authorized
// stake minus the pending authorization decrease. Eligible stake
// is what is used for operator's weight in the sortition pool.
// If the authorized stake minus the pending authorization decrease
// is below the minimum authorization, eligible stake is 0.
func (tc *TbtcChain) EligibleStake(stakingProvider chain.Address) (*big.Int, error) {
	eligibleStake, err := tc.walletRegistry.EligibleStake(
		common.HexToAddress(stakingProvider.String()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get eligible stake for staking provider %s: [%w]",
			stakingProvider,
			err,
		)
	}

	return eligibleStake, nil
}

// IsPoolLocked returns true if the sortition pool is locked and no state
// changes are allowed.
func (tc *TbtcChain) IsPoolLocked() (bool, error) {
	return tc.sortitionPool.IsLocked()
}

// IsOperatorInPool returns true if the operator is registered in
// the sortition pool.
func (tc *TbtcChain) IsOperatorInPool() (bool, error) {
	return tc.walletRegistry.IsOperatorInPool(tc.key.Address)
}

// IsOperatorUpToDate checks if the operator's authorized stake is in sync
// with operator's weight in the sortition pool.
// If the operator's authorized stake is not in sync with sortition pool
// weight, function returns false.
// If the operator is not in the sortition pool and their authorized stake
// is non-zero, function returns false.
func (tc *TbtcChain) IsOperatorUpToDate() (bool, error) {
	return tc.walletRegistry.IsOperatorUpToDate(tc.key.Address)
}

// JoinSortitionPool executes a transaction to have the operator join the
// sortition pool.
func (tc *TbtcChain) JoinSortitionPool() error {
	_, err := tc.walletRegistry.JoinSortitionPool()
	return err
}

// UpdateOperatorStatus executes a transaction to update the operator's
// state in the sortition pool.
func (tc *TbtcChain) UpdateOperatorStatus() error {
	_, err := tc.walletRegistry.UpdateOperatorStatus(tc.key.Address)
	return err
}

// IsEligibleForRewards checks whether the operator is eligible for rewards
// or not.
func (tc *TbtcChain) IsEligibleForRewards() (bool, error) {
	return tc.sortitionPool.IsEligibleForRewards(tc.key.Address)
}

// Checks whether the operator is able to restore their eligibility for
// rewards right away.
func (tc *TbtcChain) CanRestoreRewardEligibility() (bool, error) {
	return tc.sortitionPool.CanRestoreRewardEligibility(tc.key.Address)
}

// Restores reward eligibility for the operator.
func (tc *TbtcChain) RestoreRewardEligibility() error {
	_, err := tc.sortitionPool.RestoreRewardEligibility(tc.key.Address)
	return err
}

// Returns true if the chaosnet phase is active, false otherwise.
func (tc *TbtcChain) IsChaosnetActive() (bool, error) {
	return tc.sortitionPool.IsChaosnetActive()
}

// Returns true if operator is a beta operator, false otherwise.
// Chaosnet status does not matter.
func (tc *TbtcChain) IsBetaOperator() (bool, error) {
	return tc.sortitionPool.IsBetaOperator(tc.key.Address)
}

// GetOperatorID returns the ID number of the given operator address. An ID
// number of 0 means the operator has not been allocated an ID number yet.
func (tc *TbtcChain) GetOperatorID(
	operatorAddress chain.Address,
) (chain.OperatorID, error) {
	return tc.sortitionPool.GetOperatorID(
		common.HexToAddress(operatorAddress.String()),
	)
}

// SelectGroup returns the group members selected for the current group
// selection. The function returns an error if the chain's state does not allow
// for group selection at the moment.
func (tc *TbtcChain) SelectGroup() (*tbtc.GroupSelectionResult, error) {
	operatorsIDs, err := tc.walletRegistry.SelectGroup()
	if err != nil {
		return nil, fmt.Errorf(
			"cannot select group in the sortition pool: [%v]",
			err,
		)
	}

	operatorsAddresses, err := tc.sortitionPool.GetIDOperators(operatorsIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot convert operators' IDs to addresses: [%v]",
			err,
		)
	}

	// Should not happen as this is guaranteed by the contract but, just in case.
	if len(operatorsIDs) != len(operatorsAddresses) {
		return nil, fmt.Errorf("operators IDs and addresses mismatch")
	}

	ids := make([]chain.OperatorID, len(operatorsIDs))
	addresses := make([]chain.Address, len(operatorsIDs))
	for i := range ids {
		ids[i] = operatorsIDs[i]
		addresses[i] = chain.Address(operatorsAddresses[i].String())
	}

	return &tbtc.GroupSelectionResult{
		OperatorsIDs:       ids,
		OperatorsAddresses: addresses,
	}, nil
}

func (tc *TbtcChain) OnDKGStarted(
	handler func(event *tbtc.DKGStartedEvent),
) subscription.EventSubscription {
	onEvent := func(
		seed *big.Int,
		blockNumber uint64,
	) {
		handler(&tbtc.DKGStartedEvent{
			Seed:        seed,
			BlockNumber: blockNumber,
		})
	}

	return tc.walletRegistry.DkgStartedEvent(nil, nil).OnEvent(onEvent)
}

func (tc *TbtcChain) PastDKGStartedEvents(
	filter *tbtc.DKGStartedEventFilter,
) ([]*tbtc.DKGStartedEvent, error) {
	var startBlock uint64
	var endBlock *uint64
	var seed []*big.Int

	if filter != nil {
		startBlock = filter.StartBlock
		endBlock = filter.EndBlock
		seed = filter.Seed
	}

	events, err := tc.walletRegistry.PastDkgStartedEvents(
		startBlock,
		endBlock,
		seed,
	)
	if err != nil {
		return nil, err
	}

	dkgStartedEvents := make([]*tbtc.DKGStartedEvent, len(events))
	for i, event := range events {
		dkgStartedEvents[i] = &tbtc.DKGStartedEvent{
			Seed:        event.Seed,
			BlockNumber: event.Raw.BlockNumber,
		}
	}

	sort.SliceStable(dkgStartedEvents, func(i, j int) bool {
		return dkgStartedEvents[i].BlockNumber < dkgStartedEvents[j].BlockNumber
	})

	return dkgStartedEvents, err
}

func (tc *TbtcChain) OnDKGResultSubmitted(
	handler func(event *tbtc.DKGResultSubmittedEvent),
) subscription.EventSubscription {
	onEvent := func(
		resultHash [32]byte,
		seed *big.Int,
		result ecdsaabi.EcdsaDkgResult,
		blockNumber uint64,
	) {
		tbtcResult, err := convertDkgResultFromAbiType(result)
		if err != nil {
			logger.Errorf(
				"unexpected DKG result in DKGResultSubmitted event: [%v]",
				err,
			)
			return
		}

		handler(&tbtc.DKGResultSubmittedEvent{
			Seed:        seed,
			ResultHash:  resultHash,
			Result:      tbtcResult,
			BlockNumber: blockNumber,
		})
	}

	return tc.walletRegistry.
		DkgResultSubmittedEvent(nil, nil, nil).
		OnEvent(onEvent)
}

// convertDkgResultFromAbiType converts the WalletRegistry-specific DKG
// result to the format applicable for the TBTC application.
func convertDkgResultFromAbiType(
	result ecdsaabi.EcdsaDkgResult,
) (*tbtc.DKGChainResult, error) {
	if err := validateMemberIndex(result.SubmitterMemberIndex); err != nil {
		return nil, fmt.Errorf(
			"unexpected submitter member index: [%v]",
			err,
		)
	}

	signingMembersIndexes := make(
		[]group.MemberIndex,
		len(result.SigningMembersIndices),
	)
	for i, memberIndex := range result.SigningMembersIndices {
		if err := validateMemberIndex(memberIndex); err != nil {
			return nil, fmt.Errorf(
				"unexpected signing member index: [%v]",
				err,
			)
		}

		signingMembersIndexes[i] = group.MemberIndex(memberIndex.Uint64())
	}

	return &tbtc.DKGChainResult{
		SubmitterMemberIndex:     group.MemberIndex(result.SubmitterMemberIndex.Uint64()),
		GroupPublicKey:           result.GroupPubKey,
		MisbehavedMembersIndexes: result.MisbehavedMembersIndices,
		Signatures:               result.Signatures,
		SigningMembersIndexes:    signingMembersIndexes,
		Members:                  result.Members,
		MembersHash:              result.MembersHash,
	}, nil
}

// convertDkgResultToAbiType converts the TBTC-specific DKG result to
// the format applicable for the WalletRegistry ABI.
func convertDkgResultToAbiType(
	result *tbtc.DKGChainResult,
) ecdsaabi.EcdsaDkgResult {
	signingMembersIndices := make([]*big.Int, len(result.SigningMembersIndexes))
	for i, memberIndex := range result.SigningMembersIndexes {
		signingMembersIndices[i] = big.NewInt(int64(memberIndex))
	}

	return ecdsaabi.EcdsaDkgResult{
		SubmitterMemberIndex:     big.NewInt(int64(result.SubmitterMemberIndex)),
		GroupPubKey:              result.GroupPublicKey,
		MisbehavedMembersIndices: result.MisbehavedMembersIndexes,
		Signatures:               result.Signatures,
		SigningMembersIndices:    signingMembersIndices,
		Members:                  result.Members,
		MembersHash:              result.MembersHash,
	}
}

func validateMemberIndex(chainMemberIndex *big.Int) error {
	maxMemberIndex := big.NewInt(group.MaxMemberIndex)
	if chainMemberIndex.Cmp(maxMemberIndex) > 0 {
		return fmt.Errorf("invalid member index value: [%v]", chainMemberIndex)
	}

	return nil
}

func (tc *TbtcChain) OnDKGResultChallenged(
	handler func(event *tbtc.DKGResultChallengedEvent),
) subscription.EventSubscription {
	onEvent := func(
		resultHash [32]byte,
		challenger common.Address,
		reason string,
		blockNumber uint64,
	) {
		handler(&tbtc.DKGResultChallengedEvent{
			ResultHash:  resultHash,
			Challenger:  chain.Address(challenger.Hex()),
			Reason:      reason,
			BlockNumber: blockNumber,
		})
	}

	return tc.walletRegistry.
		DkgResultChallengedEvent(nil, nil, nil).
		OnEvent(onEvent)
}

func (tc *TbtcChain) OnDKGResultApproved(
	handler func(event *tbtc.DKGResultApprovedEvent),
) subscription.EventSubscription {
	onEvent := func(
		resultHash [32]byte,
		approver common.Address,
		blockNumber uint64,
	) {
		handler(&tbtc.DKGResultApprovedEvent{
			ResultHash:  resultHash,
			Approver:    chain.Address(approver.Hex()),
			BlockNumber: blockNumber,
		})
	}

	return tc.walletRegistry.
		DkgResultApprovedEvent(nil, nil, nil).
		OnEvent(onEvent)
}

// AssembleDKGResult assembles the DKG chain result according to the rules
// expected by the given chain.
func (tc *TbtcChain) AssembleDKGResult(
	submitterMemberIndex group.MemberIndex,
	groupPublicKey *ecdsa.PublicKey,
	operatingMembersIndexes []group.MemberIndex,
	misbehavedMembersIndexes []group.MemberIndex,
	signatures map[group.MemberIndex][]byte,
	groupSelectionResult *tbtc.GroupSelectionResult,
) (*tbtc.DKGChainResult, error) {
	serializedGroupPublicKey, err := convertPubKeyToChainFormat(groupPublicKey)
	if err != nil {
		return nil, fmt.Errorf(
			"could not convert group public key to chain format: [%v]",
			err,
		)
	}

	// Sort misbehavedMembersIndexes slice in ascending order as expected
	// by the on-chain contract.
	sort.Slice(misbehavedMembersIndexes[:], func(i, j int) bool {
		return misbehavedMembersIndexes[i] < misbehavedMembersIndexes[j]
	})

	signingMemberIndices, signatureBytes, err := convertSignaturesToChainFormat(
		signatures,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not convert signatures to chain format: [%v]",
			err,
		)
	}

	// Sort operatingOperatorsIDs slice in ascending order as the slice
	// holding the operators IDs used to compute the members hash is
	// expected to be sorted in the same way.
	sort.Slice(operatingMembersIndexes[:], func(i, j int) bool {
		return operatingMembersIndexes[i] < operatingMembersIndexes[j]
	})

	operatingOperatorsIDs := make([]chain.OperatorID, len(operatingMembersIndexes))
	for i, operatingMemberIndex := range operatingMembersIndexes {
		operatingOperatorsIDs[i] =
			groupSelectionResult.OperatorsIDs[operatingMemberIndex-1]
	}

	membersHash, err := computeOperatorsIDsHash(operatingOperatorsIDs)
	if err != nil {
		return nil, fmt.Errorf("could not compute members hash: [%v]", err)
	}

	return &tbtc.DKGChainResult{
		SubmitterMemberIndex:     submitterMemberIndex,
		GroupPublicKey:           serializedGroupPublicKey[:],
		MisbehavedMembersIndexes: misbehavedMembersIndexes,
		Signatures:               signatureBytes,
		SigningMembersIndexes:    signingMemberIndices,
		Members:                  groupSelectionResult.OperatorsIDs,
		MembersHash:              membersHash,
	}, nil
}

func (tc *TbtcChain) SubmitDKGResult(
	dkgResult *tbtc.DKGChainResult,
) error {
	_, err := tc.walletRegistry.SubmitDkgResult(
		convertDkgResultToAbiType(dkgResult),
	)

	return err
}

// computeOperatorsIDsHash computes the keccak256 hash for the given list
// of operators IDs.
func computeOperatorsIDsHash(operatorsIDs chain.OperatorIDs) ([32]byte, error) {
	uint32SliceType, err := abi.NewType("uint32[]", "uint32[]", nil)
	if err != nil {
		return [32]byte{}, err
	}

	bytes, err := abi.Arguments{{Type: uint32SliceType}}.Pack(operatorsIDs)
	if err != nil {
		return [32]byte{}, err
	}

	return crypto.Keccak256Hash(bytes), nil
}

// convertSignaturesToChainFormat converts signatures map to two slices. The
// first slice contains indices of members from the map, sorted in ascending order
// as required by the contract. The second slice is a slice of concatenated
// signatures. Signatures and member indices are returned in the matching order.
// It requires each signature to be exactly 65-byte long.
func convertSignaturesToChainFormat(
	signatures map[group.MemberIndex][]byte,
) ([]group.MemberIndex, []byte, error) {
	membersIndexes := make([]group.MemberIndex, 0)
	for memberIndex := range signatures {
		membersIndexes = append(membersIndexes, memberIndex)
	}

	sort.Slice(membersIndexes, func(i, j int) bool {
		return membersIndexes[i] < membersIndexes[j]
	})

	signatureSize := 65

	var signaturesSlice []byte

	for _, memberIndex := range membersIndexes {
		signature := signatures[memberIndex]

		if len(signature) != signatureSize {
			return nil, nil, fmt.Errorf(
				"invalid signature size for member [%v] got [%d] bytes but [%d] bytes required",
				memberIndex,
				len(signature),
				signatureSize,
			)
		}

		signaturesSlice = append(signaturesSlice, signature...)
	}

	return membersIndexes, signaturesSlice, nil
}

// convertPubKeyToChainFormat takes X and Y coordinates of a signer's public key
// and concatenates it to a 64-byte long array. If any of coordinates is shorter
// than 32-byte it is preceded with zeros.
func convertPubKeyToChainFormat(publicKey *ecdsa.PublicKey) ([64]byte, error) {
	var serialized [64]byte

	x, err := byteutils.LeftPadTo32Bytes(publicKey.X.Bytes())
	if err != nil {
		return serialized, err
	}

	y, err := byteutils.LeftPadTo32Bytes(publicKey.Y.Bytes())
	if err != nil {
		return serialized, err
	}

	serializedBytes := append(x, y...)

	copy(serialized[:], serializedBytes)

	return serialized, nil
}

func (tc *TbtcChain) GetDKGState() (tbtc.DKGState, error) {
	walletCreationState, err := tc.walletRegistry.GetWalletCreationState()
	if err != nil {
		return 0, err
	}

	var state tbtc.DKGState

	switch walletCreationState {
	case 0:
		state = tbtc.Idle
	case 1:
		state = tbtc.AwaitingSeed
	case 2:
		state = tbtc.AwaitingResult
	case 3:
		state = tbtc.Challenge
	default:
		err = fmt.Errorf(
			"unexpected wallet creation state: [%v]",
			walletCreationState,
		)
	}

	return state, err
}

// CalculateDKGResultSignatureHash calculates a 32-byte hash that is used
// to produce a signature supporting the given groupPublicKey computed
// as result of the given DKG process. The misbehavedMembersIndexes parameter
// should contain indexes of members that were considered as misbehaved
// during the DKG process. The startBlock argument is the block at which
// the given DKG process started.
func (tc *TbtcChain) CalculateDKGResultSignatureHash(
	groupPublicKey *ecdsa.PublicKey,
	misbehavedMembersIndexes []group.MemberIndex,
	startBlock uint64,
) (dkg.ResultSignatureHash, error) {
	groupPublicKeyBytes := elliptic.Marshal(
		groupPublicKey.Curve,
		groupPublicKey.X,
		groupPublicKey.Y,
	)
	// Crop the 04 prefix as the calculateDKGResultSignatureHash function
	// expects an unprefixed 64-byte public key,
	unprefixedGroupPublicKeyBytes := groupPublicKeyBytes[1:]

	// Sort misbehavedMembersIndexes slice in ascending order as expected
	// by the calculateDKGResultSignatureHash function.
	sort.Slice(misbehavedMembersIndexes[:], func(i, j int) bool {
		return misbehavedMembersIndexes[i] < misbehavedMembersIndexes[j]
	})

	return calculateDKGResultSignatureHash(
		tc.chainID,
		unprefixedGroupPublicKeyBytes,
		misbehavedMembersIndexes,
		big.NewInt(int64(startBlock)),
	)
}

// calculateDKGResultSignatureHash computes the keccak256 hash for the given DKG
// result parameters. It expects that the groupPublicKey is a 64-byte uncompressed
// public key without the 04 prefix and misbehavedMembersIndexes slice is
// sorted in ascending order. Those expectations are forced by the contract.
func calculateDKGResultSignatureHash(
	chainID *big.Int,
	groupPublicKey []byte,
	misbehavedMembersIndexes []group.MemberIndex,
	startBlock *big.Int,
) (dkg.ResultSignatureHash, error) {
	publicKeySize := 64

	if len(groupPublicKey) != publicKeySize {
		return dkg.ResultSignatureHash{}, fmt.Errorf(
			"wrong group public key length",
		)
	}

	uint256Type, err := abi.NewType("uint256", "uint256", nil)
	if err != nil {
		return dkg.ResultSignatureHash{}, err
	}
	bytesType, err := abi.NewType("bytes", "bytes", nil)
	if err != nil {
		return dkg.ResultSignatureHash{}, err
	}
	uint8SliceType, err := abi.NewType("uint8[]", "uint8[]", nil)
	if err != nil {
		return dkg.ResultSignatureHash{}, err
	}

	bytes, err := abi.Arguments{
		{Type: uint256Type},
		{Type: bytesType},
		{Type: uint8SliceType},
		{Type: uint256Type},
	}.Pack(
		chainID,
		groupPublicKey,
		misbehavedMembersIndexes,
		startBlock,
	)
	if err != nil {
		return dkg.ResultSignatureHash{}, err
	}

	return dkg.ResultSignatureHash(crypto.Keccak256Hash(bytes)), nil
}

func (tc *TbtcChain) IsDKGResultValid(
	dkgResult *tbtc.DKGChainResult,
) (bool, error) {
	outcome, err := tc.walletRegistry.IsDkgResultValid(
		convertDkgResultToAbiType(dkgResult),
	)
	if err != nil {
		return false, fmt.Errorf("cannot check result validity: [%v]", err)
	}

	return parseDkgResultValidationOutcome(&outcome)
}

// parseDkgResultValidationOutcome parses the DKG validation outcome and returns
// a boolean indicating whether the result is valid or not. The outcome parameter
// must be a pointer to a struct containing a boolean flag as the first field.
//
// TODO: Find a better way to get the validity flag. This would require changes
// in the contracts binding generator.
func parseDkgResultValidationOutcome(
	outcome interface{},
) (bool, error) {
	value := reflect.ValueOf(outcome)
	switch value.Kind() {
	case reflect.Pointer:
	default:
		return false, fmt.Errorf("result validation outcome is not a pointer")
	}

	field := value.Elem().Field(0)
	switch field.Kind() {
	case reflect.Bool:
		return field.Bool(), nil
	default:
		return false, fmt.Errorf("cannot parse result validation outcome")
	}
}

func (tc *TbtcChain) ChallengeDKGResult(dkgResult *tbtc.DKGChainResult) error {
	_, err := tc.walletRegistry.ChallengeDkgResult(
		convertDkgResultToAbiType(dkgResult),
	)

	return err
}

func (tc *TbtcChain) ApproveDKGResult(dkgResult *tbtc.DKGChainResult) error {
	result := convertDkgResultToAbiType(dkgResult)

	gasEstimate, err := tc.walletRegistry.ApproveDkgResultGasEstimate(result)
	if err != nil {
		return err
	}

	// The original estimate for this contract call turned out to be too low.
	// Here we add a 20% margin to overcome the gas problems.
	gasEstimateWithMargin := float64(gasEstimate) * float64(1.2)

	_, err = tc.walletRegistry.ApproveDkgResult(
		result,
		ethutil.TransactionOptions{
			GasLimit: uint64(gasEstimateWithMargin),
		},
	)

	return err
}

func (tc *TbtcChain) DKGParameters() (*tbtc.DKGParameters, error) {
	parameters, err := tc.walletRegistry.DkgParameters()
	if err != nil {
		return nil, err
	}

	return &tbtc.DKGParameters{
		SubmissionTimeoutBlocks:       parameters.ResultSubmissionTimeout.Uint64(),
		ChallengePeriodBlocks:         parameters.ResultChallengePeriodLength.Uint64(),
		ApprovePrecedencePeriodBlocks: parameters.SubmitterPrecedencePeriodLength.Uint64(),
	}, nil
}

func (tc *TbtcChain) PastDepositRevealedEvents(
	filter *tbtc.DepositRevealedEventFilter,
) ([]*tbtc.DepositRevealedEvent, error) {
	var startBlock uint64
	var endBlock *uint64
	var depositor []common.Address
	var walletPublicKeyHash [][20]byte

	if filter != nil {
		startBlock = filter.StartBlock
		endBlock = filter.EndBlock

		for _, d := range filter.Depositor {
			depositor = append(depositor, common.HexToAddress(d.String()))
		}

		walletPublicKeyHash = filter.WalletPublicKeyHash
	}

	events, err := tc.bridge.PastDepositRevealedEvents(
		startBlock,
		endBlock,
		depositor,
		walletPublicKeyHash,
	)
	if err != nil {
		return nil, err
	}

	convertedEvents := make([]*tbtc.DepositRevealedEvent, 0)
	for _, event := range events {
		var vault *chain.Address
		if event.Vault != [20]byte{} {
			v := chain.Address(event.Vault.Hex())
			vault = &v
		}

		convertedEvent := &tbtc.DepositRevealedEvent{
			// We can map the event.FundingTxHash field directly to the
			// bitcoin.Hash type. This is because event.FundingTxHash is
			// a [32]byte type representing a hash in the bitcoin.InternalByteOrder,
			// just as bitcoin.Hash assumes.
			FundingTxHash:       event.FundingTxHash,
			FundingOutputIndex:  event.FundingOutputIndex,
			Depositor:           chain.Address(event.Depositor.Hex()),
			Amount:              event.Amount,
			BlindingFactor:      event.BlindingFactor,
			WalletPublicKeyHash: event.WalletPubKeyHash,
			RefundPublicKeyHash: event.RefundPubKeyHash,
			RefundLocktime:      event.RefundLocktime,
			Vault:               vault,
			BlockNumber:         event.Raw.BlockNumber,
		}

		convertedEvents = append(convertedEvents, convertedEvent)
	}

	sort.SliceStable(
		convertedEvents,
		func(i, j int) bool {
			return convertedEvents[i].BlockNumber < convertedEvents[j].BlockNumber
		},
	)

	return convertedEvents, err
}

func (tc *TbtcChain) GetDepositRequest(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) (*tbtc.DepositChainRequest, error) {
	depositKey := buildDepositKey(fundingTxHash, fundingOutputIndex)

	depositRequest, err := tc.bridge.Deposits(depositKey)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get deposit request for key [0x%x]: [%v]",
			depositKey.Text(16),
			err,
		)
	}

	// Deposit not found.
	if depositRequest.RevealedAt == 0 {
		return nil, fmt.Errorf(
			"no deposit request for key [0x%x]",
			depositKey.Text(16),
		)
	}

	var vault *chain.Address
	if depositRequest.Vault != [20]byte{} {
		v := chain.Address(depositRequest.Vault.Hex())
		vault = &v
	}

	return &tbtc.DepositChainRequest{
		Depositor:   chain.Address(depositRequest.Depositor.Hex()),
		Amount:      depositRequest.Amount,
		RevealedAt:  time.Unix(int64(depositRequest.RevealedAt), 0),
		Vault:       vault,
		TreasuryFee: depositRequest.TreasuryFee,
		SweptAt:     time.Unix(int64(depositRequest.SweptAt), 0),
	}, nil
}

func (tc *TbtcChain) PastNewWalletRegisteredEvents(
	filter *tbtc.NewWalletRegisteredEventFilter,
) ([]*tbtc.NewWalletRegisteredEvent, error) {
	var startBlock uint64
	var endBlock *uint64
	var ecdsaWalletID [][32]byte
	var walletPublicKeyHash [][20]byte

	if filter != nil {
		startBlock = filter.StartBlock
		endBlock = filter.EndBlock
		ecdsaWalletID = filter.EcdsaWalletID
		walletPublicKeyHash = filter.WalletPublicKeyHash
	}

	events, err := tc.bridge.PastNewWalletRegisteredEvents(
		startBlock,
		endBlock,
		ecdsaWalletID,
		walletPublicKeyHash,
	)
	if err != nil {
		return nil, err
	}

	convertedEvents := make([]*tbtc.NewWalletRegisteredEvent, 0)
	for _, event := range events {
		convertedEvent := &tbtc.NewWalletRegisteredEvent{
			EcdsaWalletID:       event.EcdsaWalletID,
			WalletPublicKeyHash: event.WalletPubKeyHash,
			BlockNumber:         event.Raw.BlockNumber,
		}

		convertedEvents = append(convertedEvents, convertedEvent)
	}

	sort.SliceStable(
		convertedEvents,
		func(i, j int) bool {
			return convertedEvents[i].BlockNumber < convertedEvents[j].BlockNumber
		},
	)

	return convertedEvents, err
}

func (tc *TbtcChain) GetWallet(
	walletPublicKeyHash [20]byte,
) (*tbtc.WalletChainData, error) {
	wallet, err := tc.bridge.Wallets(walletPublicKeyHash)
	if err != nil {
		return nil, fmt.Errorf(
			"cannot get wallet for public key hash [0x%x]: [%v]",
			walletPublicKeyHash,
			err,
		)
	}

	// Wallet not found.
	if wallet.CreatedAt == 0 {
		return nil, fmt.Errorf(
			"no wallet for public key hash [0x%x]",
			wallet,
		)
	}

	return &tbtc.WalletChainData{
		EcdsaWalletID:                          wallet.EcdsaWalletID,
		MainUtxoHash:                           wallet.MainUtxoHash,
		PendingRedemptionsValue:                wallet.PendingRedemptionsValue,
		CreatedAt:                              time.Unix(int64(wallet.CreatedAt), 0),
		MovingFundsRequestedAt:                 time.Unix(int64(wallet.MovingFundsRequestedAt), 0),
		ClosingStartedAt:                       time.Unix(int64(wallet.ClosingStartedAt), 0),
		PendingMovedFundsSweepRequestsCount:    wallet.PendingMovedFundsSweepRequestsCount,
		State:                                  wallet.State,
		MovingFundsTargetWalletsCommitmentHash: wallet.MovingFundsTargetWalletsCommitmentHash,
	}, nil
}

func (tc *TbtcChain) ComputeMainUtxoHash(
	mainUtxo *bitcoin.UnspentTransactionOutput,
) [32]byte {
	return computeMainUtxoHash(mainUtxo)
}

func computeMainUtxoHash(mainUtxo *bitcoin.UnspentTransactionOutput) [32]byte {
	outputIndexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(outputIndexBytes, mainUtxo.Outpoint.OutputIndex)

	valueBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(valueBytes, uint64(mainUtxo.Value))

	mainUtxoHash := crypto.Keccak256Hash(
		append(
			append(
				mainUtxo.Outpoint.TransactionHash[:],
				outputIndexBytes...,
			), valueBytes...,
		),
	)

	return mainUtxoHash
}

func (tc *TbtcChain) BuildDepositKey(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) *big.Int {
	return buildDepositKey(fundingTxHash, fundingOutputIndex)
}

func (tc *TbtcChain) GetDepositParameters() (
	dustThreshold uint64,
	treasuryFeeDivisor uint64,
	txMaxFee uint64,
	revealAheadPeriod uint32,
	err error,
) {
	parameters, callErr := tc.bridge.DepositParameters()
	if callErr != nil {
		err = callErr
		return
	}

	dustThreshold = parameters.DepositDustThreshold
	treasuryFeeDivisor = parameters.DepositTreasuryFeeDivisor
	txMaxFee = parameters.DepositTxMaxFee
	revealAheadPeriod = parameters.DepositRevealAheadPeriod

	return
}

func (tc *TbtcChain) TxProofDifficultyFactor() (*big.Int, error) {
	return tc.bridge.TxProofDifficultyFactor()
}

func (tc *TbtcChain) SubmitDepositSweepProofWithReimbursement(
	transaction *bitcoin.Transaction,
	proof *bitcoin.SpvProof,
	mainUTXO bitcoin.UnspentTransactionOutput,
	vault common.Address,
) error {
	bitcoinTxInfo := tbtcabi.BitcoinTxInfo3{
		Version:      transaction.SerializeVersion(),
		InputVector:  transaction.SerializeInputs(),
		OutputVector: transaction.SerializeOutputs(),
		Locktime:     transaction.SerializeLocktime(),
	}
	sweepProof := tbtcabi.BitcoinTxProof2{
		MerkleProof:    proof.MerkleProof,
		TxIndexInBlock: big.NewInt(int64(proof.TxIndexInBlock)),
		BitcoinHeaders: proof.BitcoinHeaders,
	}
	utxo := tbtcabi.BitcoinTxUTXO2{
		TxHash:        mainUTXO.Outpoint.TransactionHash,
		TxOutputIndex: mainUTXO.Outpoint.OutputIndex,
		TxOutputValue: uint64(mainUTXO.Value),
	}

	_, err := tc.maintainerProxy.SubmitDepositSweepProof(
		bitcoinTxInfo,
		sweepProof,
		utxo,
		vault,
	)
	return err
}

func buildDepositKey(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
) *big.Int {
	fundingOutputIndexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(fundingOutputIndexBytes, fundingOutputIndex)

	depositKey := crypto.Keccak256Hash(
		append(fundingTxHash[:], fundingOutputIndexBytes...),
	)

	return depositKey.Big()
}

func (tc *TbtcChain) OnHeartbeatRequestSubmitted(
	handler func(event *tbtc.HeartbeatRequestSubmittedEvent),
) subscription.EventSubscription {
	onEvent := func(
		walletPubKeyHash [20]byte,
		message []byte,
		coordinator common.Address,
		blockNumber uint64,
	) {
		handler(&tbtc.HeartbeatRequestSubmittedEvent{
			WalletPublicKeyHash: walletPubKeyHash,
			Message:             message,
			Coordinator:         chain.Address(coordinator.Hex()),
			BlockNumber:         blockNumber,
		})
	}

	return tc.walletCoordinator.HeartbeatRequestSubmittedEvent(nil, nil).OnEvent(onEvent)
}

func (tc *TbtcChain) OnDepositSweepProposalSubmitted(
	handler func(event *tbtc.DepositSweepProposalSubmittedEvent),
) subscription.EventSubscription {
	onEvent := func(
		proposal tbtcabi.WalletCoordinatorDepositSweepProposal,
		coordinator common.Address,
		blockNumber uint64,
	) {
		handler(&tbtc.DepositSweepProposalSubmittedEvent{
			Proposal:    convertDepositSweepProposalFromAbiType(proposal),
			Coordinator: chain.Address(coordinator.Hex()),
			BlockNumber: blockNumber,
		})
	}

	return tc.walletCoordinator.
		DepositSweepProposalSubmittedEvent(nil, nil).
		OnEvent(onEvent)
}

func (tc *TbtcChain) PastDepositSweepProposalSubmittedEvents(
	filter *tbtc.DepositSweepProposalSubmittedEventFilter,
) ([]*tbtc.DepositSweepProposalSubmittedEvent, error) {
	var startBlock uint64
	var endBlock *uint64
	var coordinator []common.Address
	var walletPublicKeyHash [20]byte

	if filter != nil {
		startBlock = filter.StartBlock
		endBlock = filter.EndBlock

		for _, ps := range filter.Coordinator {
			coordinator = append(
				coordinator,
				common.HexToAddress(ps.String()),
			)
		}

		walletPublicKeyHash = filter.WalletPublicKeyHash
	}

	events, err := tc.walletCoordinator.PastDepositSweepProposalSubmittedEvents(
		startBlock,
		endBlock,
		coordinator,
	)
	if err != nil {
		return nil, err
	}

	convertedEvents := make([]*tbtc.DepositSweepProposalSubmittedEvent, 0)
	for _, event := range events {
		// If the wallet PKH filter is set, omit all events that target
		// different wallets.
		if walletPublicKeyHash != [20]byte{} {
			if event.Proposal.WalletPubKeyHash != walletPublicKeyHash {
				continue
			}
		}

		convertedEvent := &tbtc.DepositSweepProposalSubmittedEvent{
			Proposal:    convertDepositSweepProposalFromAbiType(event.Proposal),
			Coordinator: chain.Address(event.Coordinator.Hex()),
			BlockNumber: event.Raw.BlockNumber,
		}

		convertedEvents = append(convertedEvents, convertedEvent)
	}

	sort.SliceStable(
		convertedEvents,
		func(i, j int) bool {
			return convertedEvents[i].BlockNumber < convertedEvents[j].BlockNumber
		},
	)

	return convertedEvents, err
}

func convertDepositSweepProposalFromAbiType(
	proposal tbtcabi.WalletCoordinatorDepositSweepProposal,
) *tbtc.DepositSweepProposal {
	depositsKeys := make(
		[]struct {
			FundingTxHash      bitcoin.Hash
			FundingOutputIndex uint32
		},
		len(proposal.DepositsKeys),
	)

	for i, depositKey := range proposal.DepositsKeys {
		// We can map the depositKey.FundingTxHash field directly to the
		// bitcoin.Hash type. This is because depositKey.FundingTxHash is
		// a [32]byte type representing a hash in the bitcoin.InternalByteOrder,
		// just as bitcoin.Hash assumes.
		depositsKeys[i] = struct {
			FundingTxHash      bitcoin.Hash
			FundingOutputIndex uint32
		}{
			FundingTxHash:      depositKey.FundingTxHash,
			FundingOutputIndex: depositKey.FundingOutputIndex,
		}
	}

	return &tbtc.DepositSweepProposal{
		WalletPublicKeyHash:  proposal.WalletPubKeyHash,
		DepositsKeys:         depositsKeys,
		SweepTxFee:           proposal.SweepTxFee,
		DepositsRevealBlocks: proposal.DepositsRevealBlocks,
	}
}

func convertDepositSweepProposalToAbiType(
	proposal *tbtc.DepositSweepProposal,
) tbtcabi.WalletCoordinatorDepositSweepProposal {
	depositsKeys := make(
		[]tbtcabi.WalletCoordinatorDepositKey,
		len(proposal.DepositsKeys),
	)

	for i, depositKey := range proposal.DepositsKeys {
		// We can map the depositKey.FundingTxHash field directly to the
		// [32]byte type. This is because depositKey.FundingTxHash is
		// a bitcoin.Hash type representing a hash in the
		// bitcoin.InternalByteOrder, just as the on-chain contract assumes.
		depositsKeys[i] = tbtcabi.WalletCoordinatorDepositKey{
			FundingTxHash:      depositKey.FundingTxHash,
			FundingOutputIndex: depositKey.FundingOutputIndex,
		}
	}

	return tbtcabi.WalletCoordinatorDepositSweepProposal{
		WalletPubKeyHash:     proposal.WalletPublicKeyHash,
		DepositsKeys:         depositsKeys,
		SweepTxFee:           proposal.SweepTxFee,
		DepositsRevealBlocks: proposal.DepositsRevealBlocks,
	}
}

func (tc *TbtcChain) GetWalletLock(
	walletPublicKeyHash [20]byte,
) (time.Time, tbtc.WalletActionType, error) {
	lock, err := tc.walletCoordinator.WalletLock(walletPublicKeyHash)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("cannot get wallet lock from chain: [%v]", err)
	}

	cause, err := parseWalletActionType(lock.Cause)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("cannot parse wallet lock cause: [%v]", err)
	}

	return time.Unix(int64(lock.ExpiresAt), 0), cause, nil
}

func parseWalletActionType(value uint8) (tbtc.WalletActionType, error) {
	switch value {
	case 0:
		return tbtc.Noop, nil
	case 1:
		return tbtc.Heartbeat, nil
	case 2:
		return tbtc.DepositSweep, nil
	case 3:
		return tbtc.Redemption, nil
	case 4:
		return tbtc.MovingFunds, nil
	case 5:
		return tbtc.MovedFundsSweep, nil
	default:
		return 0, fmt.Errorf("unexpected wallet action value: [%v]", value)
	}
}

func (tc *TbtcChain) ValidateDepositSweepProposal(
	proposal *tbtc.DepositSweepProposal,
	depositsExtraInfo []struct {
		*tbtc.Deposit
		FundingTx *bitcoin.Transaction
	},
) error {
	dei := make([]tbtcabi.WalletCoordinatorDepositExtraInfo, len(depositsExtraInfo))
	for i, depositExtraInfo := range depositsExtraInfo {
		fundingTx := tbtcabi.BitcoinTxInfo2{
			Version:      depositExtraInfo.FundingTx.SerializeVersion(),
			InputVector:  depositExtraInfo.FundingTx.SerializeInputs(),
			OutputVector: depositExtraInfo.FundingTx.SerializeOutputs(),
			Locktime:     depositExtraInfo.FundingTx.SerializeLocktime(),
		}

		dei[i] = tbtcabi.WalletCoordinatorDepositExtraInfo{
			FundingTx:        fundingTx,
			BlindingFactor:   depositExtraInfo.Deposit.BlindingFactor,
			WalletPubKeyHash: depositExtraInfo.Deposit.WalletPublicKeyHash,
			RefundPubKeyHash: depositExtraInfo.Deposit.RefundPublicKeyHash,
			RefundLocktime:   depositExtraInfo.Deposit.RefundLocktime,
		}
	}

	valid, err := tc.walletCoordinator.ValidateDepositSweepProposal(
		convertDepositSweepProposalToAbiType(proposal),
		dei,
	)
	if err != nil {
		return fmt.Errorf("validation failed: [%v]", err)
	}

	// Should never happen because `validateDepositSweepProposal` returns true
	// or reverts (returns an error) but do the check just in case.
	if !valid {
		return fmt.Errorf("unexpected validation result")
	}

	return nil
}

func (tc *TbtcChain) SubmitDepositSweepProposalWithReimbursement(
	proposal *tbtc.DepositSweepProposal,
) error {
	gasEstimate, err := tc.walletCoordinator.SubmitDepositSweepProposalWithReimbursementGasEstimate(
		convertDepositSweepProposalToAbiType(proposal),
	)
	if err != nil {
		return err
	}

	// The original estimate for this contract call turned is low and the call
	// fails on reimbursing the submitter. Examples:
	// 0x5711df32d785140ca6b5b12c87f818a6c5d75d10445a12a7d3d75caadb40c0ac
	// 0xf9a8c0b0ecceb673e19eed7af7c9963cdd929468fb3818e9a8c3b8c59dc6ef85
	// Here we add a 20% margin to overcome the gas problems.
	gasEstimateWithMargin := float64(gasEstimate) * float64(1.2)

	_, err = tc.walletCoordinator.SubmitDepositSweepProposalWithReimbursement(
		convertDepositSweepProposalToAbiType(proposal),
		ethutil.TransactionOptions{
			GasLimit: uint64(gasEstimateWithMargin),
		},
	)

	return err
}

func (tc *TbtcChain) GetDepositSweepMaxSize() (uint16, error) {
	return tc.walletCoordinator.DepositSweepMaxSize()
}
