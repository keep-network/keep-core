package ethereum

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"

	gethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen/abi"
	ecdsacontract "github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen/contract"
	tbtccontract "github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/contract"
	"github.com/keep-network/keep-core/pkg/internal/byteutils"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/subscription"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
)

// TODO: implement DKG result challenge functions

// Definitions of contract names.
const (
	// Deprecated: The wallet registry address is taken from the Bridge.
	// TODO: Remove that field from the config template.
	WalletRegistryContractName = "WalletRegistry"
	BridgeContractName         = "Bridge"
)

// TbtcChain represents a TBTC-specific chain handle.
type TbtcChain struct {
	*baseChain

	bridge         *tbtccontract.Bridge
	walletRegistry *ecdsacontract.WalletRegistry
	sortitionPool  *ecdsacontract.EcdsaSortitionPool
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

	return &TbtcChain{
		baseChain:      baseChain,
		walletRegistry: walletRegistry,
		sortitionPool:  sortitionPool,
	}, nil
}

// GetConfig returns the expected configuration of the TBTC module.
func (tc *TbtcChain) GetConfig() *tbtc.ChainConfig {
	groupSize := 100
	groupQuorum := 90
	honestThreshold := 51

	return &tbtc.ChainConfig{
		GroupSize:       groupSize,
		GroupQuorum:     groupQuorum,
		HonestThreshold: honestThreshold,
	}
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

	ids := make([]chain.OperatorID, len(operatorsIDs))
	addresses := make([]chain.Address, len(operatorsIDs))
	for i := range ids {
		ids[i] = chain.OperatorID(operatorsIDs[i])
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

func (tc *TbtcChain) OnDKGResultSubmitted(
	handler func(event *tbtc.DKGResultSubmittedEvent),
) subscription.EventSubscription {
	onEvent := func(
		resultHash [32]byte,
		seed *big.Int,
		result abi.EcdsaDkgResult,
		blockNumber uint64,
	) {
		if err := validateMemberIndex(result.SubmitterMemberIndex); err != nil {
			logger.Errorf(
				"unexpected submitter member index in DKGResultSubmitted event: [%v]",
				err,
			)
			return
		}

		handler(&tbtc.DKGResultSubmittedEvent{
			MemberIndex:         uint32(result.SubmitterMemberIndex.Uint64()),
			GroupPublicKeyBytes: result.GroupPubKey,
			Misbehaved:          result.MisbehavedMembersIndices,
			BlockNumber:         blockNumber,
		})
	}

	return tc.walletRegistry.DkgResultSubmittedEvent(nil, nil, nil).OnEvent(onEvent)
}

func validateMemberIndex(chainMemberIndex *big.Int) error {
	maxMemberIndex := big.NewInt(group.MaxMemberIndex)
	if chainMemberIndex.Cmp(maxMemberIndex) > 0 {
		return fmt.Errorf("invalid member index value: [%v]", chainMemberIndex)
	}

	return nil
}

func (tc *TbtcChain) SubmitDKGResult(
	memberIndex group.MemberIndex,
	dkgResult *dkg.Result,
	signatures map[group.MemberIndex][]byte,
	groupSelectionResult *tbtc.GroupSelectionResult,
) error {
	serializedKey, err := convertPubKeyToChainFormat(
		dkgResult.PrivateKeyShare.PublicKey(),
	)
	if err != nil {
		return fmt.Errorf("could not serialize the public key: [%v]", err)
	}

	signingMemberIndices, signatureBytes, err := convertSignaturesToChainFormat(
		signatures,
	)
	if err != nil {
		return fmt.Errorf("could not convert signatures to chain format: [%v]", err)
	}

	operatingMembersIndexes := dkgResult.Group.OperatingMemberIndexes()
	operatingOperatorsIDs := make([]chain.OperatorID, len(operatingMembersIndexes))
	for i, operatingMemberID := range operatingMembersIndexes {
		operatingOperatorsIDs[i] = groupSelectionResult.OperatorsIDs[operatingMemberID-1]
	}

	membersHash, err := computeOperatorsIDsHash(operatingOperatorsIDs)
	if err != nil {
		return fmt.Errorf("could not compute members hash: [%v]", err)
	}

	_, err = tc.walletRegistry.SubmitDkgResult(abi.EcdsaDkgResult{
		SubmitterMemberIndex:     big.NewInt(int64(memberIndex)),
		GroupPubKey:              serializedKey[:],
		MisbehavedMembersIndices: dkgResult.MisbehavedMembersIndexes(),
		Signatures:               signatureBytes,
		SigningMembersIndices:    signingMemberIndices,
		Members:                  groupSelectionResult.OperatorsIDs,
		MembersHash:              membersHash,
	})

	return err
}

// computeOperatorsIDsHash computes the keccak256 hash for the given list
// of operators IDs.
func computeOperatorsIDsHash(operatorsIDs chain.OperatorIDs) ([32]byte, error) {
	uint32SliceType, err := gethabi.NewType("uint32[]", "uint32[]", nil)
	if err != nil {
		return [32]byte{}, err
	}

	bytes, err := gethabi.Arguments{{Type: uint32SliceType}}.Pack(operatorsIDs)
	if err != nil {
		return [32]byte{}, err
	}

	return crypto.Keccak256Hash(bytes), nil
}

// convertSignaturesToChainFormat converts signatures map to two slices. First
// slice contains indices of members from the map, second slice is a slice of
// concatenated signatures. Signatures and member indices are returned in the
// matching order. It requires each signature to be exactly 65-byte long.
func convertSignaturesToChainFormat(
	signatures map[group.MemberIndex][]byte,
) ([]*big.Int, []byte, error) {
	signatureSize := 65

	var membersIndices []*big.Int
	var signaturesSlice []byte

	for memberIndex, signature := range signatures {
		if len(signatures[memberIndex]) != signatureSize {
			return nil, nil, fmt.Errorf(
				"invalid signature size for member [%v] got [%d] bytes but [%d] bytes required",
				memberIndex,
				len(signatures[memberIndex]),
				signatureSize,
			)
		}
		membersIndices = append(membersIndices, big.NewInt(int64(memberIndex)))
		signaturesSlice = append(signaturesSlice, signature...)
	}

	return membersIndices, signaturesSlice, nil
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

// CalculateDKGResultHash calculates Keccak-256 hash of the DKG result. Operation
// is performed off-chain.
//
// It first encodes the result using solidity ABI and then calculates Keccak-256
// hash over it. This corresponds to the DKG result hash calculation on-chain.
// Hashes calculated off-chain and on-chain must always match.
func (tc *TbtcChain) CalculateDKGResultHash(
	result *dkg.Result,
) (dkg.ResultHash, error) {
	groupPublicKeyBytes, err := result.GroupPublicKeyBytes()
	if err != nil {
		return dkg.ResultHash{}, err
	}

	// Encode DKG result to the format matched with Solidity keccak256(abi.encodePacked(...))
	// TODO: Adjust the message structure to the format needed by the wallet
	//       registry contract:
	//       \x19Ethereum signed message:\n${keccak256(groupPubKey,misbehavedIndices,startBlock)}
	hash := crypto.Keccak256(groupPublicKeyBytes, result.MisbehavedMembersIndexes())
	return dkg.ResultHashFromBytes(hash)
}

// OnHeartbeatRequested runs a heartbeat loop that produces a heartbeat
// request every ~8 hours. A single heartbeat request consists of 5 messages
// that must be signed sequentially.
func (tc *TbtcChain) OnHeartbeatRequested(
	handler func(event *tbtc.HeartbeatRequestedEvent),
) subscription.EventSubscription {
	ctx, cancelCtx := context.WithCancel(context.Background())
	blocksChan := tc.blockCounter.WatchBlocks(ctx)

	go func() {
		for {
			select {
			case block := <-blocksChan:
				// Generate a heartbeat every 2400 block, i.e. ~8 hours.
				if block%2400 == 0 {
					walletPublicKey, ok, err := tc.activeWalletPublicKey()
					if err != nil {
						logger.Errorf(
							"cannot get active wallet for heartbeat request: [%v]",
							err,
						)
						continue
					}

					if !ok {
						logger.Infof("there is no active wallet for heartbeat at the moment")
						continue
					}

					if len(walletPublicKey) > 0 {
						prefixBytes := make([]byte, 8)
						binary.BigEndian.PutUint64(
							prefixBytes,
							0xffffffffffffffff,
						)

						messages := make([]*big.Int, 5)
						for i := range messages {
							suffixBytes := make([]byte, 8)
							binary.BigEndian.PutUint64(
								suffixBytes,
								block+uint64(i),
							)

							preimage := append(prefixBytes, suffixBytes...)
							preimageSha256 := sha256.Sum256(preimage)
							message := sha256.Sum256(preimageSha256[:])

							messages[i] = new(big.Int).SetBytes(message[:])
						}

						go handler(&tbtc.HeartbeatRequestedEvent{
							WalletPublicKey: walletPublicKey,
							Messages:        messages,
							BlockNumber:     block,
						})
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return subscription.NewEventSubscription(func() {
		cancelCtx()
	})
}

func (tc *TbtcChain) activeWalletPublicKey() ([]byte, bool, error) {
	walletPublicKeyHash, err := tc.bridge.ActiveWalletPubKeyHash()
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot get active wallet public key hash: [%v]",
			err,
		)
	}

	if walletPublicKeyHash == [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} {
		return nil, false, nil
	}

	bridgeWalletData, err := tc.bridge.Wallets(walletPublicKeyHash)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot get active wallet data from Bridge: [%v]",
			err,
		)
	}

	registryWalletData, err := tc.walletRegistry.GetWallet(bridgeWalletData.EcdsaWalletID)
	if err != nil {
		return nil, false, fmt.Errorf(
			"cannot get active wallet data from WalletRegistry: [%v]",
			err,
		)
	}

	publicKeyBytes := []byte{0x04} // pre-fill with uncompressed ECDSA public key prefix
	publicKeyBytes = append(publicKeyBytes, registryWalletData.PublicKeyX[:]...)
	publicKeyBytes = append(publicKeyBytes, registryWalletData.PublicKeyY[:]...)

	return publicKeyBytes, true, nil
}