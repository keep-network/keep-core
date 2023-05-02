import { fork, take, call, put, select } from "redux-saga/effects"
import moment from "moment"
import { createSubcribeToContractEventChannel } from "./web3"
import {
  getContractsContext,
  getWeb3Context,
  getLPRewardsWrapper,
  subscribeToEventAndEmitData,
} from "./utils"
import { createManagedGrantContractInstance } from "../contracts"
import { add, sub } from "../utils/arithmetics.utils"
import { isSameEthAddress } from "../utils/general.utils"
import { getEventsFromTransaction, ZERO_ADDRESS } from "../utils/ethereum.utils"
import { LIQUIDITY_REWARD_PAIRS, MODAL_TYPES } from "../constants/constants"
/** @typedef { import("../services/liquidity-rewards").LiquidityRewards} LiquidityRewards */
import { showMessage } from "../actions/messages"
import { messageType } from "../components/Message"
import {
  OPERATOR_DELEGATION_UNDELEGATED,
  FETCH_OPERATOR_DELEGATIONS_SUCCESS,
  REMOVE_STAKE_FROM_THRESHOLD_AUTH_DATA,
  ADD_STAKE_TO_THRESHOLD_AUTH_DATA,
} from "../actions"
import {
  assetPoolDepositedEventEmitted,
  COVERAGE_POOL_FETCH_COV_POOL_DATA_SUCCESS,
  coveragePoolWithdrawalCompletedEventEmitted,
  coveragePoolWithdrawalInitiatedEventEmitted,
  riskManagerAuctionClosedEventEmitted,
  riskManagerAuctionCreatedEventEmitted,
} from "../actions/coverage-pool"
import { keepBalanceActions } from "../actions"
import { Keep } from "../contracts"
import { EVENTS } from "../constants/events"
import { showModal } from "../actions/modal"
import { thresholdStakeKeepEventEmitted } from "../actions/keep-to-t-staking"

export function* subscribeToKeepTokenTransferEvent() {
  yield take(keepBalanceActions.KEEP_TOKEN_BALANCE_REQUEST_SUCCESS)
  yield fork(observeKeepTokenTransferFrom)
  yield fork(observeKeepTokenTransferTo)
}

function* observeKeepTokenTransferFrom() {
  const { token: keepTokenContractInstance } = yield getContractsContext()
  const {
    eth: { defaultAccount },
  } = yield getWeb3Context()

  const options = {
    filter: {
      from: defaultAccount,
    },
  }

  yield fork(
    subscribeToEventAndEmitData,
    keepTokenContractInstance,
    "Transfer",
    keepBalanceActions.keepTokenTransferFromEventEmitted,
    "KeepToken.Transfer",
    options
  )
}

function* observeKeepTokenTransferTo() {
  const { token: keepTokenContractInstance } = yield getContractsContext()
  const {
    eth: { defaultAccount },
  } = yield getWeb3Context()

  const options = {
    filter: {
      to: defaultAccount,
    },
  }

  yield fork(
    subscribeToEventAndEmitData,
    keepTokenContractInstance,
    "Transfer",
    keepBalanceActions.keepTokenTransferToEventEmitted,
    "KeepToken.Transfer",
    options
  )
}

export function* subscribeToStakedEvent() {
  yield take("staking/fetch_delegations_success")
  yield fork(observeStakedEvents)
}

function* observeStakedEvents() {
  const {
    grantContract,
    tokenStakingEscrow,
    stakingContract,
    stakingPortBackerContract,
  } = yield getContractsContext()
  const web3 = yield getWeb3Context()
  const yourAddress = web3.eth.defaultAccount

  // Other events may also be emitted with the `StakeDelegated` event.
  const eventsToCheck = [
    [stakingContract, "OperatorStaked"],
    [grantContract, "TokenGrantStaked"],
    [tokenStakingEscrow, "DepositRedelegated"],
    [stakingPortBackerContract, "StakeCopied"],
  ]

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    stakingContract,
    "StakeDelegated"
  )

  // Observe and dispatch an action that updates staking reducer
  while (true) {
    try {
      const {
        transactionHash,
        returnValues: { owner, operator },
      } = yield take(contractEventCahnnel)

      const { initializationPeriod } = yield select((state) => state.staking)

      const emittedEvents = yield call(
        getEventsFromTransaction,
        eventsToCheck,
        transactionHash
      )

      let isAddressedToCurrentAccount = isSameEthAddress(owner, yourAddress)
      // The `OperatorStaked` is always emitted with the `StakeDelegated` event.
      const { authorizer, beneficiary, value } = emittedEvents.OperatorStaked
      const delegation = {
        createdAt: moment().unix(),
        operatorAddress: operator,
        authorizerAddress: authorizer,
        beneficiary,
        amount: value,
        isInInitializationPeriod: true,
        initializationOverAt: moment
          .unix(moment().unix())
          .add(initializationPeriod, "seconds"),
      }

      if (emittedEvents.StakeCopied) {
        const { owner } = emittedEvents.StakeCopied
        delegation.isCopiedStake = true
        isAddressedToCurrentAccount = isSameEthAddress(owner, yourAddress)

        // Check if the copied delegation is from grant.
        if (isAddressedToCurrentAccount) {
          try {
            const { grantId } = yield call(
              grantContract.methods.getGrantStakeDetails(operator).call
            )

            delegation.isFromGrant = true
            delegation.grantId = grantId
          } catch (error) {
            delegation.isFromGrant = false
          }
        }
      }

      if (
        (emittedEvents.TokenGrantStaked || emittedEvents.DepositRedelegated) &&
        !isAddressedToCurrentAccount
      ) {
        // If the `TokenGrantStaked` or `DepositRedelegated` event exists, it means that a delegation is from grant.
        const { grantId } =
          emittedEvents.TokenGrantStaked || emittedEvents.DepositRedelegated
        delegation.grantId = grantId
        delegation.isFromGrant = true
        const { grantee } = yield call(
          grantContract.methods.getGrant(grantId).call
        )

        isAddressedToCurrentAccount = isSameEthAddress(grantee, yourAddress)

        if (!isAddressedToCurrentAccount) {
          // check if current address is a grantee in the managed grant
          try {
            const managedGrantContractInstance =
              createManagedGrantContractInstance(web3, grantee)
            const granteeAddressInManagedGrant = yield call(
              managedGrantContractInstance.methods.grantee().call
            )
            delegation.managedGrantContractInstance =
              managedGrantContractInstance
            delegation.isManagedGrant = true

            // compere a current address with a grantee address from the ManagedGrant contract
            isAddressedToCurrentAccount = isSameEthAddress(
              yourAddress,
              granteeAddressInManagedGrant
            )
          } catch (error) {
            isAddressedToCurrentAccount = false
          }
        }
      }

      if (!isAddressedToCurrentAccount) {
        return
      }

      if (!delegation.isCopiedStake) {
        if (!delegation.isFromGrant) {
          yield put({
            type: "staking/update_owned_delegated_tokens_balance",
            payload: { operation: add, value },
          })
        } else {
          yield put({
            type: "token-grant/grant_staked",
            payload: {
              grantId: delegation.grantId,
              value,
            },
          })
        }
      }

      yield put({ type: "staking/add_delegation", payload: delegation })
      if (isSameEthAddress(yourAddress, authorizer)) {
        yield put({
          type: ADD_STAKE_TO_THRESHOLD_AUTH_DATA,
          payload: {
            ...delegation,
            owner: yourAddress,
            operatorContractAddress: Keep.thresholdStakingContract.address,
          },
        })
      }
    } catch (error) {
      console.error(`Failed subscribing to StakeDelegated event`, error)
      contractEventCahnnel.close()
    }
  }
}

export function* subscribeToUndelegatedEvent() {
  yield take("staking/fetch_delegations_success")
  yield fork(observeUndelegatedEvent)
}

function* observeUndelegatedEvent() {
  const { stakingContract } = yield getContractsContext()

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    stakingContract,
    "Undelegated"
  )

  // Observe and dispatch an action that updates keep token balance.
  while (true) {
    try {
      const {
        transactionHash,
        returnValues: { operator, undelegatedAt },
      } = yield take(contractEventCahnnel)

      // Find the existing delegation by operatorAddress in the app store.
      const delegations = yield select((state) => state.staking.delegations)
      const delegation = delegations.find(({ operatorAddress }) =>
        isSameEthAddress(operatorAddress, operator)
      )

      if (!delegation) {
        return
      }

      // If the delegation exists, we create a undelegation based on the existing delegation.
      const undelegationPeriod = yield select(
        (state) => state.staking.undelegationPeriod
      )

      yield put(
        showModal({
          modalType: MODAL_TYPES.UndelegationInitiated,
          modalProps: {
            txHash: transactionHash,
            undelegatedAt,
            undelegationPeriod,
          },
        })
      )

      const undelegation = {
        ...delegation,
        undelegatedAt: moment.unix(undelegatedAt),
        undelegationCompleteAt: moment
          .unix(undelegatedAt)
          .add(undelegationPeriod, "seconds"),
        canRecoverStake: false,
      }

      if (!undelegation.isFromGrant) {
        yield put({
          type: "staking/update_owned_delegated_tokens_balance",
          payload: { operation: sub, value: undelegation.amount },
        })
        yield put({
          type: "staking/update_owned_undelegations_tokens_balance",
          payload: { operation: add, value: undelegation.amount },
        })
      }

      yield put({ type: "staking/remove_delegation", payload: operator })
      yield put({ type: "staking/add_undelegation", payload: undelegation })
      yield put({
        type: REMOVE_STAKE_FROM_THRESHOLD_AUTH_DATA,
        payload: operator,
      })
    } catch (error) {
      console.error(`Failed subscribing to Undelegated event`, error)
      contractEventCahnnel.close()
    }
  }
}

export function* subscribeToRecoveredStakeEvent() {
  yield take("staking/fetch_delegations_success")
  yield fork(observeRecoveredStakeEvent)
}

function* observeRecoveredStakeEvent() {
  const { stakingContract } = yield getContractsContext()

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    stakingContract,
    "RecoveredStake"
  )

  // Observe and dispatch an action that updates keep token balance.
  while (true) {
    try {
      const {
        transactionHash,
        returnValues: { operator },
      } = yield take(contractEventCahnnel)

      const undelegations = yield select((state) => state.staking.undelegations)
      const recoveredUndelegation = undelegations.find((undelegation) =>
        isSameEthAddress(undelegation.operatorAddress, operator)
      )

      if (!recoveredUndelegation) {
        return
      }

      yield put(
        showModal({
          modalType: MODAL_TYPES.StakingTokensClaimed,
          modalProps: { txHash: transactionHash },
        })
      )

      if (!recoveredUndelegation.isFromGrant) {
        yield put({ type: "staking/remove_undelegation", payload: operator })

        yield put({
          type: "staking/update_owned_undelegations_tokens_balance",
          payload: { operation: sub, value: recoveredUndelegation.amount },
        })
      }
    } catch (error) {
      console.error(`Failed subscribing to RecoveredStake event`, error)
      contractEventCahnnel.close()
    }
  }
}

export function* subscribeToTokenGrantWithdrawnEvent() {
  yield take("token-grant/fetch_grants_success")
  yield fork(observeTokenGrantWithdrawnEvent)
}

function* observeTokenGrantWithdrawnEvent() {
  const { grantContract } = yield getContractsContext()

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    grantContract,
    "TokenGrantWithdrawn"
  )

  // Observe and dispatch an action that updates grants reducer.
  while (true) {
    try {
      const {
        transactionHash,
        returnValues: { grantId, amount },
      } = yield take(contractEventCahnnel)

      const availableToStake = yield call(
        grantContract.methods.availableToStake(grantId).call
      )
      yield put({
        type: "token-grant/grant_withdrawn",
        payload: { grantId, amount, availableToStake },
      })

      yield put(
        showModal({
          modalType: MODAL_TYPES.GrantTokensWithdrawn,
          modalProps: { txHash: transactionHash },
        })
      )
    } catch (error) {
      console.error(`Failed subscribing to TokenGrantWithdrawn event`, error)
      contractEventCahnnel.close()
    }
  }
}

export function* subscribeToDepositWithdrawEvent() {
  yield take("token-grant/fetch_grants_success")
  yield fork(observeDepositWithdrawnEvent)
}

function* observeDepositWithdrawnEvent() {
  const { tokenStakingEscrow, grantContract } = yield getContractsContext()
  const {
    eth: { defaultAccount },
  } = yield getWeb3Context()

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    tokenStakingEscrow,
    "DepositWithdrawn"
  )

  // Observe and dispatch an action that updates grants reducer.
  while (true) {
    try {
      const {
        returnValues: { grantee, operator, amount },
      } = yield take(contractEventCahnnel)

      // A `grantee` param in the `DepositWithdrawn` event always points to the "right" grantee address.
      // No needed additional check if it's about a managed grant.
      if (!isSameEthAddress(grantee, defaultAccount)) {
        return
      }

      const grantId = yield call(
        tokenStakingEscrow.methods.depositGrantId(operator).call
      )
      const availableToStake = yield call(
        grantContract.methods.availableToStake(grantId).call
      )
      yield put({
        type: "token-grant/grant_withdrawn",
        payload: { grantId, amount, operator, availableToStake },
      })
    } catch (error) {
      console.error(`Failed subscribing to DepositWithdrawn event`, error)
      contractEventCahnnel.close()
    }
  }
}

export function* subscribeToDepositedEvent() {
  yield take("token-grant/fetch_grants_success")
  yield fork(observeDepositedEvent)
}

function* observeDepositedEvent() {
  const { tokenStakingEscrow, grantContract } = yield getContractsContext()

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    tokenStakingEscrow,
    "Deposited"
  )

  // Observe and dispatch an action that updates grants reducer.
  while (true) {
    try {
      const {
        transactionHash,
        returnValues: { operator, grantId, amount },
      } = yield take(contractEventCahnnel)

      const grants = yield select((state) => state.tokenGrants.grants)
      if (grants.find((grant) => grant.id === grantId)) {
        yield put({ type: "staking/remove_delegation", payload: operator })
        yield put({ type: "staking/remove_undelegation", payload: operator })

        const availableToWitdrawGrant = yield call(
          grantContract.methods.withdrawable(grantId).call
        )

        yield put(
          showModal({
            modalType: MODAL_TYPES.StakingTokensClaimed,
            modalProps: { txHash: transactionHash },
          })
        )

        yield put({
          type: "token-grant/grant_deposited",
          payload: {
            grantId,
            availableToWitdrawGrant,
            amount,
            operator,
          },
        })
      }
    } catch (error) {
      console.error(`Failed subscribing to Deposited event`, error)
      contractEventCahnnel.close()
    }
  }
}

export function* subscribeToTopUpInitiatedEvent() {
  yield take("staking/fetch_delegations_success")
  yield fork(observeTopUpInitiatedEvent)
}

function* observeTopUpInitiatedEvent() {
  const { stakingContract, tokenStakingEscrow, grantContract } =
    yield getContractsContext()

  // Other events may also be emitted with the `TopUpInitiated` event.
  const eventsToCheck = [
    [grantContract, "TokenGrantStaked"],
    [tokenStakingEscrow, "DepositRedelegated"],
  ]

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    stakingContract,
    "TopUpInitiated"
  )
  while (true) {
    try {
      const {
        transactionHash,
        returnValues: { operator, topUp },
      } = yield take(contractEventCahnnel)

      const emittedEvents = yield call(
        getEventsFromTransaction,
        eventsToCheck,
        transactionHash
      )

      // Find existing delegation in the app context
      const delegations = yield select((state) => state.staking.delegations)
      const delegation = delegations.find(({ operatorAddress }) =>
        isSameEthAddress(operatorAddress, operator)
      )

      if (delegation) {
        yield put(
          showModal({
            modalType: MODAL_TYPES.TopUpInitiatedConfirmation,
            modalProps: {
              addedAmount: topUp,
              currentAmount: delegation.amount,
              authorizerAddress: delegation.authorizerAddress,
              beneficiary: delegation.beneficiary,
              operatorAddress: delegation.operatorAddress,
            },
          })
        )
        yield put({
          type: "staking/top_up_initiated",
          payload: { operator, topUp },
        })

        if (
          emittedEvents.DepositRedelegated ||
          emittedEvents.TokenGrantStaked
        ) {
          const { grantId, amount } =
            emittedEvents.DepositRedelegated || emittedEvents.TokenGrantStaked
          yield put({
            type: "token-grant/grant_staked",
            payload: { grantId, value: amount },
          })
        }
      }
    } catch (error) {
      console.error(`Failed subscribing to TopUpInitiated event`, error)
      contractEventCahnnel.close()
    }
  }
}

export function* subsribeToTopUpCompletedEvent() {
  yield take("staking/fetch_delegations_success")
  yield fork(observeTopUpCompletedEvent)
}

function* observeTopUpCompletedEvent() {
  const { stakingContract, tokenStakingEscrow, grantContract } =
    yield getContractsContext()
  const eventsToCheck = [
    [grantContract, "TokenGrantStaked"],
    [tokenStakingEscrow, "DepositRedelegated"],
  ]

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    stakingContract,
    "TopUpCompleted"
  )

  while (true) {
    try {
      const {
        transactionHash,
        returnValues: { operator, newAmount },
      } = yield take(contractEventCahnnel)

      const emittedEvents = yield call(
        getEventsFromTransaction,
        eventsToCheck,
        transactionHash
      )

      yield put({
        type: "staking/top_up_completed",
        payload: { operator, newAmount },
      })
      if (emittedEvents.DepositRedelegated || emittedEvents.TokenGrantStaked) {
        const { grantId, amount } =
          emittedEvents.DepositRedelegated || emittedEvents.TokenGrantStaked
        yield put({
          type: "token-grant/grant_satked",
          payload: { grantId, value: amount },
        })
      }
    } catch (error) {
      console.error(`Failed subscribing to TopUpCompleted event`, error)
      contractEventCahnnel.close()
    }
  }
}

function* observeECDSARewardsClaimedEvent(data) {
  const { ECDSARewardsDistributorContract } = yield getContractsContext()

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    ECDSARewardsDistributorContract,
    "RewardsClaimed"
  )

  while (true) {
    try {
      const {
        returnValues: { merkleRoot, index, operator, amount },
      } = yield take(contractEventCahnnel)

      yield put({
        type: "rewards/ecdsa_withdrawn",
        payload: { merkleRoot, index, operator, amount },
      })
    } catch (error) {
      console.error(`Failed subscribing to RewardsClaimed event`, error)
      contractEventCahnnel.close()
    }
  }
}

export function* subsribeToECDSARewardsClaimedEvent() {
  yield take("rewards/ecdsa_fetch_rewards_data_success")
  yield fork(observeECDSARewardsClaimedEvent)
}

function* observeLiquidityTokenStakedEvent(liquidityRewardPair) {
  /** @type LiquidityRewards */
  const LiquidityRewards = yield getLPRewardsWrapper(liquidityRewardPair)

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    LiquidityRewards.LPRewardsContract,
    LiquidityRewards.stakedEventName
  )

  while (true) {
    try {
      const eventData = yield take(contractEventCahnnel)

      yield* lpTokensStakedOrWithdrawn(
        eventData.returnValues,
        LiquidityRewards,
        liquidityRewardPair.name,
        `liquidity_rewards/${liquidityRewardPair.name}_staked`
      )
    } catch (error) {
      console.error(`Failed subscribing to Staked event`, error)
      contractEventCahnnel.close()
    }
  }
}

function* observeLiquidityTokenWithdrawnEvent(liquidityRewardPair) {
  /** @type LiquidityRewards */
  const LiquidityRewards = yield getLPRewardsWrapper(liquidityRewardPair)

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    LiquidityRewards.LPRewardsContract,
    LiquidityRewards.depositWithdrawnEventName
  )

  while (true) {
    try {
      const eventData = yield take(contractEventCahnnel)
      yield* lpTokensStakedOrWithdrawn(
        eventData.returnValues,
        LiquidityRewards,
        liquidityRewardPair.name,
        `liquidity_rewards/${liquidityRewardPair.name}_withdrawn`
      )
    } catch (error) {
      console.error(`Failed subscribing to Withdrawn event`, error)
      contractEventCahnnel.close()
    }
  }
}

function* lpTokensStakedOrWithdrawn(
  eventValues,
  /** @type LiquidityRewards */
  LiquidityRewards,
  liquidityRewardPairName,
  actionType
) {
  const {
    eth: { defaultAccount },
  } = yield getWeb3Context()

  const { user, amount } = eventValues
  const totalSupply = yield call([
    LiquidityRewards,
    LiquidityRewards.totalSupply,
  ])

  const apy = yield call(
    [LiquidityRewards, LiquidityRewards.calculateAPY],
    totalSupply
  )

  const { lpBalance } = yield select(
    (state) => state.liquidityRewards[liquidityRewardPairName]
  )

  let updatedlpBalance = lpBalance
  let emittedAmountValue = 0
  // Update only if this transaction relates to the current logged account.
  if (isSameEthAddress(defaultAccount, user)) {
    emittedAmountValue = amount
    const arithmeticOperation = actionType.includes("withdrawn") ? sub : add
    updatedlpBalance = arithmeticOperation(lpBalance, amount).toString()
  }

  yield* updateLPTokenBalance(
    LiquidityRewards,
    liquidityRewardPairName,
    defaultAccount,
    updatedlpBalance
  )

  const reward = yield call(
    [LiquidityRewards, LiquidityRewards.rewardBalance],
    defaultAccount,
    updatedlpBalance
  )

  const rewardMultiplier = yield call(
    [LiquidityRewards, LiquidityRewards.calculateRewardMultiplier],
    defaultAccount
  )

  // If the `Withdrawn` or `Staked` event was emitted the total pool of the LPRewards,
  // APY and reward value have changed.
  yield put({
    type: actionType,
    payload: {
      amount: emittedAmountValue,
      lpBalance: updatedlpBalance,
      totalSupply,
      reward,
      apy,
      liquidityRewardPairName,
      rewardMultiplier,
    },
  })
}

function* observeLiquidityRewardPaidEvent(liquidityRewardPair) {
  /** @type LiquidityRewards */
  const LiquidityRewards = yield getLPRewardsWrapper(liquidityRewardPair)
  const {
    eth: { defaultAccount },
  } = yield getWeb3Context()

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    LiquidityRewards.LPRewardsContract,
    LiquidityRewards.rewardClaimedEventName
  )

  while (true) {
    try {
      const { returnValues } = yield take(contractEventCahnnel)
      // LPRewards and TokenGeyser contract have different param names in an
      // emitted event which is triggered when the reward is claimed but param
      // which points to claimed reward amount is at the same index- 1. So we
      // can get claimed amount by index eg. `event.returnValues["1"]`.
      const reward = returnValues["1"]

      if (isSameEthAddress(defaultAccount, returnValues.user)) {
        yield put({
          type: `liquidity_rewards/${liquidityRewardPair.name}_reward_paid`,
          payload: {
            reward,
            liquidityRewardPairName: liquidityRewardPair.name,
          },
        })
      }
    } catch (error) {
      console.error(`Failed subscribing to RewardPaid event`, error)
      contractEventCahnnel.close()
    }
  }
}

function* observeWrappedTokenMintAndBurnTx(liquidityRewardPair) {
  /** @type LiquidityRewards */
  const LiquidityRewards = yield getLPRewardsWrapper(liquidityRewardPair)

  const {
    eth: { defaultAccount },
  } = yield getWeb3Context()

  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    LiquidityRewards.wrappedToken,
    "Transfer"
  )

  while (true) {
    try {
      const {
        transactionHash,
        returnValues: { from, to },
      } = yield take(contractEventCahnnel)

      // If the `from` is a zero address it means it was a `mint` transaction.
      // If the `to` is a zero address it means it was a `burn` transaction. For
      // these cases we need to update APY value because the total pool value
      // of the wrapped token has been increased / decreased.
      if (from === ZERO_ADDRESS || to === ZERO_ADDRESS) {
        yield* updateAPY(LiquidityRewards, liquidityRewardPair.name)
        yield* updateLPTokenBalance(
          LiquidityRewards,
          liquidityRewardPair.name,
          defaultAccount
        )
      }

      // If the 'to' address is equal to the address of the connected wallet
      // then it means that LP tokens are transferred to this address so we are
      // displaying the proper notification
      if (
        isSameEthAddress(to, defaultAccount) &&
        liquidityRewardPair.label !== LIQUIDITY_REWARD_PAIRS.KEEP_ONLY.label
      ) {
        yield put(
          showMessage({
            messageType: messageType.NEW_LP_TOKENS_IN_WALLET,
            messageProps: {
              liquidityRewardPairName: liquidityRewardPair.label,
              sticky: true,
            },
          })
        )

        const eventsToCheck = [
          [
            LiquidityRewards.LPRewardsContract,
            LiquidityRewards.depositWithdrawnEventName,
          ],
        ]

        const emittedEvents = yield call(
          getEventsFromTransaction,
          eventsToCheck,
          transactionHash
        )

        // Fetch wrappedTokenBalance value only when depositWithdrawnEventName event was not previously emitted (so
        // in other words when user did not withdraw all on the dApp, but still got some lp tokens transferred
        // to his wallet). If the user clicks withdraw all the wrappedTokenBalance will be updated on the redux
        // side (`liquidity_rewards/${liquidityRewardPairName}_withdrawn` case) so we don't need to update it here.
        if (!emittedEvents[LiquidityRewards.depositWithdrawnEventName]) {
          yield* updateWrappedTokenBalance(
            LiquidityRewards,
            liquidityRewardPair.name,
            defaultAccount
          )
        }
      }
    } catch (error) {
      console.error(`Failed subscribing to Transfer event`, error)
      contractEventCahnnel.close()
    }
  }
}

function* updateAPY(
  /** @type LiquidityRewards */
  LiquidityRewards,
  liquidityRewardPairName
) {
  const totalSupply = yield call([
    LiquidityRewards,
    LiquidityRewards.totalSupply,
  ])

  const apy = yield call(
    [LiquidityRewards, LiquidityRewards.calculateAPY],
    totalSupply
  )
  yield put({
    type: `liquidity_rewards/${liquidityRewardPairName}_apy_updated`,
    payload: {
      apy,
      liquidityRewardPairName,
    },
  })
}

function* updateWrappedTokenBalance(
  LiquidityRewards,
  liquidityRewardPairName,
  address
) {
  // Fetching balance of liquidity token for a given uniswap pair.
  const wrappedTokenBalance = yield call(
    [LiquidityRewards, LiquidityRewards.wrappedTokenBalance],
    address
  )

  yield put({
    type: `liquidity_rewards/${liquidityRewardPairName}_wrapped_token_balance_updated`,
    payload: {
      wrappedTokenBalance,
      liquidityRewardPairName,
    },
  })
}

function* updateLPTokenBalance(
  LiquidityRewards,
  liquidityRewardPairName,
  address,
  lpBalance = null
) {
  if (!lpBalance) {
    // Fetching balance of liquidity token for a given uniswap pair deposited in
    // the `LPRewards` contract.
    lpBalance = yield call(
      [LiquidityRewards, LiquidityRewards.stakedBalance],
      address
    )
  }

  const lpTokenBalance = yield call(
    [LiquidityRewards, LiquidityRewards.calculateLPTokenBalance],
    lpBalance
  )

  yield put({
    type: `liquidity_rewards/${liquidityRewardPairName}_lp_balance_updated`,
    payload: { lpTokenBalance, liquidityRewardPairName },
  })
}

export function* subscribeToLiquidityRewardsEvents() {
  for (const [pairName, value] of Object.entries(LIQUIDITY_REWARD_PAIRS)) {
    if (value.contractName) {
      yield fork(
        function* (liquidityRewardPair) {
          yield fork(observeWrappedTokenMintAndBurnTx, liquidityRewardPair)
          yield take(
            `liquidity_rewards/${liquidityRewardPair.name}_fetch_data_success`
          )
          yield fork(observeLiquidityTokenStakedEvent, liquidityRewardPair)
          yield fork(observeLiquidityTokenWithdrawnEvent, liquidityRewardPair)
          yield fork(observeLiquidityRewardPaidEvent, liquidityRewardPair)
        },
        {
          name: pairName,
          ...value,
        }
      )
    }
  }
}

function* updateOperatorData() {
  const { stakingContract } = yield getContractsContext()

  // Create subscription channel.
  const contractEventCahnnel = yield call(
    createSubcribeToContractEventChannel,
    stakingContract,
    "Undelegated"
  )

  // Observe and dispatch an action that updates keep token balance.
  while (true) {
    try {
      const {
        transactionHash,
        returnValues: { operator, undelegatedAt },
      } = yield take(contractEventCahnnel)
      const {
        eth: { defaultAccount },
      } = yield getWeb3Context()

      if (!isSameEthAddress(defaultAccount, operator)) {
        return
      }

      const { undelegationPeriod } = yield select((state) => state.operator)

      yield put(
        showModal({
          modalType: MODAL_TYPES.UndelegationInitiated,
          modalProps: {
            txHash: transactionHash,
            undelegatedAt,
            undelegationPeriod,
          },
        })
      )

      const undelegationCompletedAt = moment
        .unix(undelegatedAt)
        .add(undelegationPeriod, "seconds")

      yield put({
        type: OPERATOR_DELEGATION_UNDELEGATED,
        payload: { undelegationCompletedAt, delegationStatus: "UNDELEGATED" },
      })
    } catch (error) {
      console.error(`Failed subscribing to Undelegated event`, error)
      contractEventCahnnel.close()
    }
  }
}

export function* subscribeToOperatorUndelegateEvent() {
  yield take(FETCH_OPERATOR_DELEGATIONS_SUCCESS)
  yield fork(updateOperatorData)
}

export function* observeAssetPoolDepositedEvent() {
  yield take(COVERAGE_POOL_FETCH_COV_POOL_DATA_SUCCESS)

  const assetPoolContract = Keep.coveragePoolV1.assetPoolContract.instance

  yield fork(
    subscribeToEventAndEmitData,
    assetPoolContract,
    EVENTS.COVERAGE_POOLS.DEPOSITED,
    assetPoolDepositedEventEmitted,
    `AssetPool.${EVENTS.COVERAGE_POOLS.DEPOSITED}`
  )
}

export function* observeWithdrawalInitiatedEvent() {
  const assetPoolContract = Keep.coveragePoolV1.assetPoolContract.instance

  yield fork(
    subscribeToEventAndEmitData,
    assetPoolContract,
    EVENTS.COVERAGE_POOLS.WITHDRAWAL_INITIATED,
    coveragePoolWithdrawalInitiatedEventEmitted,
    `AssetPool.${EVENTS.COVERAGE_POOLS.WITHDRAWAL_INITIATED}`
  )
}

export function* observeWithdrawalCompletedEvent() {
  const assetPoolContract = Keep.coveragePoolV1.assetPoolContract.instance

  yield fork(
    subscribeToEventAndEmitData,
    assetPoolContract,
    EVENTS.COVERAGE_POOLS.WITHDRAWAL_COMPLETED,
    coveragePoolWithdrawalCompletedEventEmitted,
    `AssetPool.${EVENTS.COVERAGE_POOLS.WITHDRAWAL_COMPLETED}`
  )
}

export function* observeAuctionCreatedEvent() {
  const riskManagerV1Contract =
    Keep.coveragePoolV1.riskManagerV1Contract.instance

  yield fork(
    subscribeToEventAndEmitData,
    riskManagerV1Contract,
    "AuctionCreated",
    riskManagerAuctionCreatedEventEmitted,
    "RiskManagerV1.AuctionCreated"
  )
}

export function* observeAuctionClosedEvent() {
  const riskManagerV1Contract =
    Keep.coveragePoolV1.riskManagerV1Contract.instance

  yield fork(
    subscribeToEventAndEmitData,
    riskManagerV1Contract,
    "AuctionClosed",
    riskManagerAuctionClosedEventEmitted,
    "RiskManagerV1.AuctionClosed"
  )
}

export function* observeThresholdStakeKeepEvent() {
  const thresholdStakingContract =
    Keep.keepToTStaking.thresholdStakingContract.instance

  yield fork(
    subscribeToEventAndEmitData,
    thresholdStakingContract,
    "Staked",
    thresholdStakeKeepEventEmitted,
    `ThresholdTokenStaking.StakeKeep`
  )
}
