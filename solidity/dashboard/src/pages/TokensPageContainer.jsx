import React, { useContext } from "react"
import { Web3Context } from "../components/WithWeb3Context"
import TokensPage from "./TokensPage"
import TokenGrantsPage from "./TokenGrantsPage"
import TokensPageContextProvider, {
  useTokensPageContext,
} from "../contexts/TokensPageContext"
import { Route, Switch, Redirect } from "react-router-dom"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent.js"
import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"
import {
  ADD_DELEGATION,
  UPDATE_OWNED_DELEGATED_TOKENS_BALANCE,
  REMOVE_DELEGATION,
  ADD_UNDELEGATION,
  UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE,
  REMOVE_UNDELEGATION,
} from "../reducers/tokens-page.reducer.js"
import {
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
} from "../constants/constants"
import { isSameEthAddress } from "../utils/general.utils"
import { sub, add } from "../utils/arithmetics.utils"
import moment from "moment"

const TokensPageContainer = () => {
  useSubscribeToStakedEvent()
  useSubscribeToUndelegatedEvent()
  useSubscribeToRecoveredStakeEvent()
  useSubscribeToTokenGrantEvents()

  return (
    <TokensPageContextProvider>
      <Switch>
        <Route exact path="/tokens/delegate" component={TokensPage} />
        <Route exact path="/tokens/grants" component={TokenGrantsPage} />
        <Redirect to="/tokens/delegate" />
      </Switch>
    </TokensPageContextProvider>
  )
}

export default React.memo(TokensPageContainer)

const useSubscribeToStakedEvent = async () => {
  const web3Context = useContext(Web3Context)
  const { grantContract, stakingContract, eth } = web3Context

  const {
    initializationPeriod,
    dispatch,
    refreshKeepTokenBalance,
  } = useTokensPageContext()

  const subscribeToEventCallback = async (event) => {
    const {
      blockNumber,
      returnValues: { from, value },
    } = event
    const grantStakeDetails = await getGrantDetails(from, grantContract)
    const isFromGrant = grantStakeDetails !== null

    if (!isAddressedToCurrentAccount(from, web3Context, grantStakeDetails)) {
      return
    }
    const createdAt = (await eth.getBlock(blockNumber)).timestamp

    const delegation = {
      createdAt,
      operatorAddress: from,
      authorizerAddress: await stakingContract.methods
        .authorizerOf(from)
        .call(),
      beneficiary: await stakingContract.methods.beneficiaryOf(from).call(),
      amount: value,
      isInInitializationPeriod: true,
      initializationOverAt: moment
        .unix(createdAt)
        .add(initializationPeriod, "seconds"),
      isFromGrant,
      grantId: isFromGrant ? grantStakeDetails.grantId : null,
    }

    if (!isFromGrant) {
      refreshKeepTokenBalance()
      dispatch({
        type: UPDATE_OWNED_DELEGATED_TOKENS_BALANCE,
        payload: { operation: add, value },
      })
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
  const web3Context = useContext(Web3Context)
  const { grantContract, stakingContract } = web3Context

  const { undelegationPeriod, dispatch } = useTokensPageContext()

  const subscribeToEventCallback = async (event) => {
    const {
      returnValues: { operator, undelegatedAt },
    } = event
    const grantStakeDetails = await getGrantDetails(operator, grantContract)
    const isFromGrant = grantStakeDetails !== null

    if (
      !isAddressedToCurrentAccount(operator, web3Context, grantStakeDetails)
    ) {
      return
    }

    const { amount } = await stakingContract.methods
      .getDelegationInfo(operator)
      .call()

    const undelegation = {
      operatorAddress: operator,
      authorizerAddress: await stakingContract.methods
        .authorizerOf(operator)
        .call(),
      beneficiary: await stakingContract.methods.beneficiaryOf(operator).call(),
      amount,
      undelegatedAt: moment.unix(undelegatedAt),
      undelegationCompleteAt: moment
        .unix(undelegatedAt)
        .add(undelegationPeriod, "seconds"),
      canRecoverStake: false,
      isFromGrant,
      grantId: isFromGrant ? grantStakeDetails.grantId : null,
    }
    dispatch({ type: REMOVE_DELEGATION, payload: operator })

    if (!isFromGrant) {
      dispatch({
        type: UPDATE_OWNED_DELEGATED_TOKENS_BALANCE,
        payload: { operation: sub, value: amount },
      })
      dispatch({
        type: UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE,
        payload: { operation: add, value: amount },
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
    refreshGrants,
  } = useTokensPageContext()

  const subscribeToEventCallback = async (event) => {
    const {
      returnValues: { operator },
    } = event

    const { indexInArray, obj: recoveredUndelegation } = findIndexAndObject(
      "operatorAddress",
      operator,
      undelegations,
      compareEthAddresses
    )

    if (indexInArray === null) {
      return
    }

    dispatch({ type: REMOVE_UNDELEGATION, payload: operator })

    if (!recoveredUndelegation.isFromGrant) {
      refreshKeepTokenBalance()
      dispatch({
        type: UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE,
        payload: { operation: sub, value: recoveredUndelegation.amount },
      })
    } else {
      refreshGrants()
    }
  }

  useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "RecoveredStake",
    subscribeToEventCallback
  )
}

const getGrantDetails = async (operator, grantContract) => {
  let grantStakeDetails = null
  try {
    grantStakeDetails = await grantContract.methods
      .getGrantStakeDetails(operator)
      .call()
  } catch (error) {
    return grantStakeDetails
  }
  return grantStakeDetails
}

const isAddressedToCurrentAccount = async (
  operator,
  web3Context,
  grantStakeDetails
) => {
  const { yourAddress, grantContract, stakingContract } = web3Context
  const isFromGrant = grantStakeDetails !== null
  if (isFromGrant) {
    const { grantId } = grantStakeDetails
    const { grantee } = await grantContract.methods.getGrant(grantId).call()
    return !isSameEthAddress(grantee, yourAddress)
  } else {
    const owner = await stakingContract.methods.ownerOf(operator).call()
    return !isSameEthAddress(owner, yourAddress)
  }
}

const useSubscribeToTokenGrantEvents = () => {
  const {
    refreshGrantTokenBalance,
    refreshKeepTokenBalance,
    grantStaked,
    grantWithdrawn,
  } = useTokensPageContext()

  const subscribeToStakedEventCallback = (stakedEvent) => {
    const {
      returnValues: { grantId, amount },
    } = stakedEvent
    grantStaked(grantId, amount)
  }

  const subscribeToWithdrawanEventCallback = (withdrawanEvent) => {
    const {
      returnValues: { grantId, amount },
    } = withdrawanEvent
    grantWithdrawn(grantId, amount)
    refreshGrantTokenBalance()
    refreshKeepTokenBalance()
  }

  useSubscribeToContractEvent(
    TOKEN_GRANT_CONTRACT_NAME,
    "TokenGrantStaked",
    subscribeToStakedEventCallback
  )
  useSubscribeToContractEvent(
    TOKEN_GRANT_CONTRACT_NAME,
    "TokenGrantWithdrawn",
    subscribeToWithdrawanEventCallback
  )
}
