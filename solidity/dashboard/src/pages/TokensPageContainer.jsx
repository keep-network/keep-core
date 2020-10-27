import React, { useEffect } from "react"
import { useWeb3Context } from "../components/WithWeb3Context"
import TokensPage from "./TokensPage"
import TokenGrantsPage from "./TokenGrantsPage"
import TokensPageContextProvider, {
  useTokensPageContext,
} from "../contexts/TokensPageContext"
import { useLocation } from "react-router-dom"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent.js"
import {
  ADD_DELEGATION,
  UPDATE_OWNED_DELEGATED_TOKENS_BALANCE,
  REMOVE_DELEGATION,
  ADD_UNDELEGATION,
  UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE,
  REMOVE_UNDELEGATION,
  SET_TOKENS_CONTEXT,
  TOP_UP_INITIATED,
  TOP_UP_COMPLETED,
} from "../reducers/tokens-page.reducer.js"
import { FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST } from "../actions"
import { connect } from "react-redux"
import {
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
} from "../constants/constants"
import { isSameEthAddress } from "../utils/general.utils"
import { sub, add } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"
import moment from "moment"
import {
  createManagedGrantContractInstance,
  ContractsLoaded,
  Web3Loaded,
} from "../contracts"
import { getEventsFromTransaction } from "../utils/ethereum.utils"
import Banner, { BANNER_TYPE } from "../components/Banner"
import Button from "../components/Button"
import { useModal } from "../hooks/useModal"
import CopyStakePage from "./CopyStakePage"
import PageWrapper from "../components/PageWrapper"

const TokensPageContainer = ({
  title,
  routes,
  oldDelegations,
  fetchOldDelegations,
}) => {
  useSubscribeToStakedEvent()
  useSubscribeToUndelegatedEvent()
  useSubscribeToRecoveredStakeEvent()
  useSubscribeToTokenGrantEvents()
  useSubscribeToTopUpsEvents()

  const { hash } = useLocation()
  const { dispatch } = useTokensPageContext()

  useEffect(() => {
    const tokenContext = hash.substring(1)
    if (tokenContext === "owned" || tokenContext === "granted") {
      dispatch({ type: SET_TOKENS_CONTEXT, payload: tokenContext })
    }
  }, [hash, dispatch])

  useEffect(() => {
    fetchOldDelegations()
  }, [fetchOldDelegations])

  const { openModal } = useModal()

  return (
    <PageWrapper title={title} routes={routes}>
      {!isEmptyArray(oldDelegations) && (
        <Banner
          type={BANNER_TYPE.NOTIFICATION}
          withIcon
          title="New upgrade available for your stake delegations!"
          titleClassName="h4"
          subtitle="Upgrade now to keep earning rewards on your stake."
        >
          <Button
            className="btn btn-tertiary btn-sm ml-a"
            onClick={() => openModal(<CopyStakePage />, { isFullScreen: true })}
          >
            upgrade my stake
          </Button>
        </Banner>
      )}
    </PageWrapper>
  )
}

const useSubscribeToStakedEvent = () => {
  const {
    initializationPeriod,
    dispatch,
    refreshKeepTokenBalance,
    grantStaked,
  } = useTokensPageContext()

  const subscribeToEventCallback = async (event) => {
    const web3 = await Web3Loaded
    const yourAddress = web3.eth.defaultAccount
    const {
      grantContract,
      tokenStakingEscrow,
      stakingContract,
      stakingPortBackerContract,
    } = await ContractsLoaded
    const {
      transactionHash,
      returnValues: { owner, operator },
    } = event

    // Other events may also be emitted with the `StakeDelegated` event.
    const eventsToCheck = [
      [stakingContract, "OperatorStaked"],
      [grantContract, "TokenGrantStaked"],
      [tokenStakingEscrow, "DepositRedelegated"],
      [stakingPortBackerContract, "StakeCopied"],
    ]

    const emittedEvents = await getEventsFromTransaction(
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
          const { grantId } = await grantContract.methods
            .getGrantStakeDetails(operator)
            .call()

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
      const { grantee } = await grantContract.methods.getGrant(grantId).call()

      isAddressedToCurrentAccount = isSameEthAddress(grantee, yourAddress)

      if (!isAddressedToCurrentAccount) {
        // check if current address is a grantee in the managed grant
        try {
          const managedGrantContractInstance = createManagedGrantContractInstance(
            web3,
            grantee
          )
          const granteeAddressInManagedGrant = await managedGrantContractInstance.methods
            .grantee()
            .call()
          delegation.managedGrantContractInstance = managedGrantContractInstance
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
        refreshKeepTokenBalance()
        dispatch({
          type: UPDATE_OWNED_DELEGATED_TOKENS_BALANCE,
          payload: { operation: add, value },
        })
      } else {
        grantStaked(delegation.grantId, value)
      }
    }

    dispatch({ type: ADD_DELEGATION, payload: delegation })
  }
  useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "StakeDelegated",
    subscribeToEventCallback
  )
}

const useSubscribeToUndelegatedEvent = () => {
  const { undelegationPeriod, dispatch, delegations } = useTokensPageContext()

  const subscribeToEventCallback = async (event) => {
    const {
      returnValues: { operator, undelegatedAt },
    } = event

    // Find the existing delegation by operatorAddress in the app context.
    const delegation = delegations.find(({ operatorAddress }) =>
      isSameEthAddress(operatorAddress, operator)
    )

    if (!delegation) {
      return
    }
    // If the delegation exists, we create a undelegation based on the existing delegation.
    const undelegation = {
      ...delegation,
      undelegatedAt: moment.unix(undelegatedAt),
      undelegationCompleteAt: moment
        .unix(undelegatedAt)
        .add(undelegationPeriod, "seconds"),
      canRecoverStake: false,
    }

    if (!undelegation.isFromGrant) {
      dispatch({
        type: UPDATE_OWNED_DELEGATED_TOKENS_BALANCE,
        payload: { operation: sub, value: undelegation.amount },
      })
      dispatch({
        type: UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE,
        payload: { operation: add, value: undelegation.amount },
      })
    }

    dispatch({ type: REMOVE_DELEGATION, payload: operator })
    dispatch({ type: ADD_UNDELEGATION, payload: undelegation })
  }

  useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "Undelegated",
    subscribeToEventCallback
  )
}

const useSubscribeToRecoveredStakeEvent = async () => {
  const {
    refreshKeepTokenBalance,
    dispatch,
    undelegations,
  } = useTokensPageContext()

  const subscribeToEventCallback = async (event) => {
    const {
      returnValues: { operator },
    } = event

    const recoveredUndelegation = undelegations.find((undelegation) =>
      isSameEthAddress(undelegation.operatorAddress, operator)
    )

    if (!recoveredUndelegation) {
      return
    }

    if (!recoveredUndelegation.isFromGrant) {
      dispatch({ type: REMOVE_UNDELEGATION, payload: operator })
      refreshKeepTokenBalance()
      dispatch({
        type: UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE,
        payload: { operation: sub, value: recoveredUndelegation.amount },
      })
    }
  }

  useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "RecoveredStake",
    subscribeToEventCallback
  )
}

const useSubscribeToTokenGrantEvents = () => {
  const {
    refreshGrantTokenBalance,
    refreshKeepTokenBalance,
    grantWithdrawn,
    grants,
    dispatch,
    grantDeposited,
  } = useTokensPageContext()
  const { yourAddress, tokenStakingEscrow } = useWeb3Context()

  const subscribeToWithdrawanEventCallback = (withdrawanEvent) => {
    const {
      returnValues: { grantId, amount },
    } = withdrawanEvent
    grantWithdrawn(grantId, amount)
    refreshGrantTokenBalance()
    refreshKeepTokenBalance()
  }

  const subscribeToDepositWithdrawn = async (depositWithdrawnEvent) => {
    const {
      returnValues: { grantee, operator, amount },
    } = depositWithdrawnEvent
    // A `grantee` param in the `DepositWithdrawn` event always points to the "right" grantee address.
    // No needed additional check if it's about a managed grant.
    if (!isSameEthAddress(grantee, yourAddress)) {
      return
    }

    const grantId = await tokenStakingEscrow.methods
      .depositGrantId(operator)
      .call()
    grantWithdrawn(grantId, amount, operator)
    refreshGrantTokenBalance()
    refreshKeepTokenBalance()
  }

  const subscribeToDepositedEvent = async (depositedEvent) => {
    const {
      returnValues: { operator, grantId, amount },
    } = depositedEvent

    if (grants.find((grant) => grant.id === grantId)) {
      dispatch({ type: REMOVE_DELEGATION, payload: operator })
      dispatch({ type: REMOVE_UNDELEGATION, payload: operator })
      grantDeposited(grantId, operator, amount)
    }
  }

  useSubscribeToContractEvent(
    TOKEN_GRANT_CONTRACT_NAME,
    "TokenGrantWithdrawn",
    subscribeToWithdrawanEventCallback
  )

  useSubscribeToContractEvent(
    TOKEN_STAKING_ESCROW_CONTRACT_NAME,
    "DepositWithdrawn",
    subscribeToDepositWithdrawn
  )

  useSubscribeToContractEvent(
    TOKEN_STAKING_ESCROW_CONTRACT_NAME,
    "Deposited",
    subscribeToDepositedEvent
  )
}

const useSubscribeToTopUpsEvents = () => {
  const {
    dispatch,
    delegations,
    refreshKeepTokenBalance,
    grantStaked,
  } = useTokensPageContext()

  const subscribeToTopUpInitiated = async (event) => {
    const {
      transactionHash,
      returnValues: { operator },
    } = event
    const { tokenStakingEscrow, grantContract } = await ContractsLoaded

    // Other events may also be emitted with the `TopUpInitiated` event.
    const eventsToCheck = [
      [grantContract, "TokenGrantStaked"],
      [tokenStakingEscrow, "DepositRedelegated"],
    ]
    const emittedEvents = await getEventsFromTransaction(
      eventsToCheck,
      transactionHash
    )

    // Find existing delegation in the app context
    const delegation = delegations.find(({ operatorAddress }) =>
      isSameEthAddress(operatorAddress, operator)
    )

    if (delegation) {
      dispatch({ type: TOP_UP_INITIATED, payload: event.returnValues })
      if (!delegation.isFromGrant) {
        refreshKeepTokenBalance()
      }

      if (emittedEvents.DepositRedelegated || emittedEvents.TokenGrantStaked) {
        const { grantId, amount } =
          emittedEvents.DepositRedelegated || emittedEvents.TokenGrantStaked
        grantStaked(grantId, amount)
      }
    }
  }

  const subscribeToTopUpCompleted = async (event) => {
    const { tokenStakingEscrow, grantContract } = await ContractsLoaded

    // Other events may also be emitted with the `TopUpCompleted` event.
    const eventsToCheck = [
      [grantContract, "TokenGrantStaked"],
      [tokenStakingEscrow, "DepositRedelegated"],
    ]
    const emittedEvents = await getEventsFromTransaction(
      eventsToCheck,
      event.transactionHash
    )

    dispatch({ type: TOP_UP_COMPLETED, payload: event.returnValues })
    if (emittedEvents.DepositRedelegated || emittedEvents.TokenGrantStaked) {
      const { grantId, amount } =
        emittedEvents.DepositRedelegated || emittedEvents.TokenGrantStaked
      grantStaked(grantId, amount)
    }
  }

  useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "TopUpInitiated",
    subscribeToTopUpInitiated
  )
  useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "TopUpCompleted",
    subscribeToTopUpCompleted
  )
}

const mapStateToProps = ({ copyStake }) => {
  const { oldDelegations } = copyStake

  return { oldDelegations }
}

const mapDispatchToProps = (dispatch) => {
  return {
    fetchOldDelegations: () =>
      dispatch({ type: FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST }),
  }
}

const TokensPageContainerWithRedux = connect(
  mapStateToProps,
  mapDispatchToProps
)(TokensPageContainer)

const TokensPageContainerWithContext = React.memo((props) => (
  <TokensPageContextProvider>
    <TokensPageContainerWithRedux {...props} />
  </TokensPageContextProvider>
))

TokensPageContainerWithContext.route = {
  title: "Tokens",
  path: "/tokens",
  pages: [TokensPage, TokenGrantsPage],
}

export default TokensPageContainerWithContext
