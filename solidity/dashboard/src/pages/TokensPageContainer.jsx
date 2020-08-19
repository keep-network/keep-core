import React, { useEffect } from "react"
import { useWeb3Context } from "../components/WithWeb3Context"
import TokensPage from "./TokensPage"
import TokenGrantsPage from "./TokenGrantsPage"
import TokensPageContextProvider, {
  useTokensPageContext,
} from "../contexts/TokensPageContext"
import { Route, Switch, Redirect, useLocation } from "react-router-dom"
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
import {
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
} from "../constants/constants"
import { isSameEthAddress } from "../utils/general.utils"
import { sub, add } from "../utils/arithmetics.utils"
import moment from "moment"
import {
  createManagedGrantContractInstance,
  ContractsLoaded,
  Web3Loaded,
} from "../contracts"
import TokenOverviewPage from "./TokenOverviewPage"
import { getEventsFromTransaction } from "../utils/ethereum.utils"

const TokensPageContainer = () => {
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

  return (
    <Switch>
      <Route exact path="/tokens/overview" component={TokenOverviewPage} />
      <Route exact path="/tokens/delegate" component={TokensPage} />
      <Route exact path="/tokens/grants" component={TokenGrantsPage} />
      <Redirect to="/tokens/overview" />
    </Switch>
  )
}

const TokensPageContainerWithContext = () => (
  <TokensPageContextProvider>
    <TokensPageContainer />
  </TokensPageContextProvider>
)

export default React.memo(TokensPageContainerWithContext)

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
    const { grantContract, tokenStakingEscrow } = await ContractsLoaded
    const {
      transactionHash,
      returnValues: { owner, operator, authorizer, beneficiary, value },
    } = event

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

    // Other events may also be emitted with the `Staked` event.
    const eventsToCheck = [
      [grantContract, "TokenGrantStaked"],
      [tokenStakingEscrow, "DepositRedelegated"],
    ]

    const emittedEvents = await getEventsFromTransaction(
      eventsToCheck,
      transactionHash
    )
    let isAddressedToCurrentAccount = isSameEthAddress(owner, yourAddress)

    if (
      (emittedEvents.TokenGrantStaked || emittedEvents.DepositRedelegated) &&
      !isAddressedToCurrentAccount
    ) {
      // If the `TokenGrantStaked` or `DepositRedelegated` event exists, it means that a delegation is from grant.
      const { grantId } =
        emittedEvents.TokenGrantStaked || emittedEvents.DepositReedelegated
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

    if (!delegation.isFromGrant) {
      refreshKeepTokenBalance()
      dispatch({
        type: UPDATE_OWNED_DELEGATED_TOKENS_BALANCE,
        payload: { operation: add, value },
      })
    } else {
      grantStaked(delegation.grantId, value)
    }

    dispatch({ type: ADD_DELEGATION, payload: delegation })
  }
  useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "Staked",
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
    const { tokenStakingEscrow } = await ContractsLoaded

    // Other events may also be emitted with the `TopUpInitiated` event.
    const eventsToCheck = [[tokenStakingEscrow, "DepositRedelegated"]]
    const emmittedEvents = await getEventsFromTransaction(
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

      if (emmittedEvents.DepositRedelegated) {
        const { grantId, amount } = emmittedEvents.DepositRedelegated
        grantStaked(grantId, amount)
      }
    }
  }

  const subscribeToTopUpCompleted = async (event) => {
    const { tokenStakingEscrow } = await ContractsLoaded

    // Other events may also be emitted with the `TopUpCompleted` event.
    const eventsToCheck = [[tokenStakingEscrow, "DepositRedelegated"]]
    const emmittedEvents = await getEventsFromTransaction(
      eventsToCheck,
      event.transactionHash
    )

    dispatch({ type: TOP_UP_COMPLETED, payload: event.returnValues })

    if (emmittedEvents.DepositRedelegated) {
      const { grantId, amount } = emmittedEvents.DepositRedelegated
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
