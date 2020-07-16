import React, { useContext, useEffect } from "react"
import { Web3Context, useWeb3Context } from "../components/WithWeb3Context"
import TokensPage from "./TokensPage"
import TokenGrantsPage from "./TokenGrantsPage"
import TokensPageContextProvider, {
  useTokensPageContext,
} from "../contexts/TokensPageContext"
import { Route, Switch, Redirect, useLocation } from "react-router-dom"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent.js"
import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"
import {
  ADD_DELEGATION,
  UPDATE_OWNED_DELEGATED_TOKENS_BALANCE,
  REMOVE_DELEGATION,
  ADD_UNDELEGATION,
  UPDATE_OWNED_UNDELEGATIONS_TOKEN_BALANCE,
  REMOVE_UNDELEGATION,
  SET_TOKENS_CONTEXT,
} from "../reducers/tokens-page.reducer.js"
import {
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
} from "../constants/constants"
import { isSameEthAddress } from "../utils/general.utils"
import { sub, add } from "../utils/arithmetics.utils"
import moment from "moment"
import { createManagedGrantContractInstance, isCodeValid } from "../contracts"
import TokenOverviewPage from "./TokenOverviewPage"

const TokensPageContainer = () => {
  useSubscribeToStakedEvent()
  useSubscribeToUndelegatedEvent()
  useSubscribeToRecoveredStakeEvent()
  useSubscribeToTokenGrantEvents()
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

const useSubscribeToStakedEvent = async () => {
  const web3Context = useContext(Web3Context)
  const { grantContract, eth, web3, yourAddress } = web3Context

  const {
    initializationPeriod,
    dispatch,
    refreshKeepTokenBalance,
  } = useTokensPageContext()

  const subscribeToEventCallback = async (event) => {
    const {
      blockNumber,
      returnValues: { owner, operator, authorizer, beneficiary, value },
    } = event

    const grantStakeDetails = await getGrantDetails(operator, grantContract)
    const isFromGrant = grantStakeDetails !== null
    const { grantee } = isFromGrant
      ? await grantContract.methods.getGrant(grantStakeDetails.grantId).call()
      : {}
    let isManagedGrant
    let managedGrantContractInstance
    if (isFromGrant && (await isGranteeInManagedGrant(web3Context, grantee))) {
      isManagedGrant = true
      managedGrantContractInstance = createManagedGrantContractInstance(
        web3,
        grantee
      )
    } else if (!isAddressedToCurrentAccount(owner, yourAddress, grantee)) {
      return
    }

    const createdAt = (await eth.getBlock(blockNumber)).timestamp

    const delegation = {
      createdAt,
      operatorAddress: operator,
      authorizerAddress: authorizer,
      beneficiary: beneficiary,
      amount: value,
      isInInitializationPeriod: true,
      initializationOverAt: moment
        .unix(createdAt)
        .add(initializationPeriod, "seconds"),
      grantId: isFromGrant ? grantStakeDetails.grantId : null,
      isFromGrant,
      isManagedGrant,
      managedGrantContractInstance,
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
  const { grantContract, stakingContract, web3, yourAddress } = web3Context

  const { undelegationPeriod, dispatch, delegations } = useTokensPageContext()

  const subscribeToEventCallback = async (event) => {
    const {
      returnValues: { operator, undelegatedAt },
    } = event
    const delegation = delegations.find(({ operatorAddress }) =>
      isSameEthAddress(operatorAddress, operator)
    )
    let undelegation = {}
    if (delegation) {
      undelegation = {
        ...delegation,
        undelegatedAt: moment.unix(undelegatedAt),
        undelegationCompleteAt: moment
          .unix(undelegatedAt)
          .add(undelegationPeriod, "seconds"),
        canRecoverStake: false,
      }
    } else {
      const grantStakeDetails = await getGrantDetails(operator, grantContract)
      const owner = await stakingContract.methods.ownerOf(operator).call()
      const isFromGrant = grantStakeDetails !== null
      const { grantee } = isFromGrant
        ? await grantContract.methods.getGrant(grantStakeDetails.grantId).call()
        : {}
      let isManagedGrant
      let managedGrantContractInstance
      if (
        isFromGrant &&
        (await isGranteeInManagedGrant(web3Context, grantee))
      ) {
        isManagedGrant = true
        managedGrantContractInstance = createManagedGrantContractInstance(
          web3,
          grantee
        )
      } else if (!isAddressedToCurrentAccount(owner, yourAddress, grantee)) {
        return
      }
      const { amount } = await stakingContract.methods
        .getDelegationInfo(operator)
        .call()
      undelegation = {
        operatorAddress: operator,
        authorizerAddress: await stakingContract.methods
          .authorizerOf(operator)
          .call(),
        beneficiary: await stakingContract.methods
          .beneficiaryOf(operator)
          .call(),
        amount,
        undelegatedAt: moment.unix(undelegatedAt),
        undelegationCompleteAt: moment
          .unix(undelegatedAt)
          .add(undelegationPeriod, "seconds"),
        canRecoverStake: false,
        isFromGrant,
        grantId: isFromGrant ? grantStakeDetails.grantId : null,
        isManagedGrant,
        managedGrantContractInstance,
      }
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
  } catch {
    return grantStakeDetails
  }
  return grantStakeDetails
}

const isAddressedToCurrentAccount = async (owner, yourAddress, grantee) => {
  const isFromGrant = !!grantee

  if (isFromGrant) {
    return isSameEthAddress(grantee, yourAddress)
  }

  return isSameEthAddress(owner, yourAddress)
}

const useSubscribeToTokenGrantEvents = () => {
  const {
    refreshGrantTokenBalance,
    refreshKeepTokenBalance,
    grantStaked,
    grantWithdrawn,
    grants,
    dispatch,
  } = useTokensPageContext()
  const { yourAddress, tokenStakingEscrow } = useWeb3Context()

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

  const subscribeToDepositWithdrawn = async (depositWithdrawnEvent) => {
    const {
      returnValues: { grantee, operator, amount },
    } = depositWithdrawnEvent
    if (!isSameEthAddress(grantee, yourAddress)) {
      return
    }

    const grantId = await tokenStakingEscrow.methods.depositGrantId(operator)
    grantWithdrawn(grantId, amount)
    refreshGrantTokenBalance()
    refreshKeepTokenBalance()
  }

  const subscribeToDepositedEvent = async (depositedEvent) => {
    const {
      returnValues: { operator, grantId },
    } = depositedEvent

    if (grants.find((grant) => grant.id === grantId)) {
      dispatch({ type: REMOVE_UNDELEGATION, payload: operator })
    }
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

const isGranteeInManagedGrant = async (web3Context, grantee) => {
  const { web3, yourAddress } = web3Context
  const managedGrantContractInstance = createManagedGrantContractInstance(
    web3,
    grantee
  )

  // check if grantee is a contract
  const code = await web3.eth.getCode(grantee)

  if (!isCodeValid(code)) {
    return false
  }

  const granteeAddressInManagedGrant = await managedGrantContractInstance.methods
    .grantee()
    .call()

  return isSameEthAddress(yourAddress, granteeAddressInManagedGrant)
}
