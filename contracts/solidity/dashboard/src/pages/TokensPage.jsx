import React, { useContext } from 'react'
import DelegateStakeForm from '../components/DelegateStakeForm'
import TokensOverview from '../components/TokensOverview'
import Undelegations from '../components/Undelegations'
import { useFetchData } from '../hooks/useFetchData'
import { tokensPageService } from '../services/tokens-page.service'
import DelegatedTokensList from '../components/DelegatedTokensList'
import { Web3Context } from '../components/WithWeb3Context'
import { useShowMessage, messageType } from '../components/Message'
import SpeechBubbleInfo from '../components/SpeechBubbleInfo'
import { LoadingOverlay } from '../components/Loadable'
import { useSubscribeToContractEvent } from '../hooks/useSubscribeToContractEvent.js'
import { TOKEN_STAKING_CONTRACT_NAME } from '../constants/constants'
import { isSameEthAddress } from '../utils/general.utils'
import { sub, add } from '../utils/arithmetics.utils'
import { findIndexAndObject, compareEthAddresses } from '../utils/array.utils'

const initialData = {
  ownedKeepBalance: '',
  tokenStakingBalance: '',
  pendingUndelegationBalance: '',
  tokenGrantsBalance: '',
  tokenGrantsStakeBalance: '',
  minimumStake: '',
  delegations: [],
  undelegations: [],
}

const TokensPage = () => {
  const web3Context = useContext(Web3Context)
  const showMessage = useShowMessage()
  const [state, setData, refreshData] = useFetchData(tokensPageService.fetchTokensPageData, initialData)
  useSubscribeToStakedEvent(state.data, setData)
  useSubscribeToUndelegatedEvent(state.data, setData)
  useSubscribeToRecoveredStakeEvent(state.data, setData)

  const {
    undelegationPeriod,
    ownedKeepBalance,
    pendingUndelegationBalance,
    tokenStakingBalance,
    tokenGrantsBalance,
    tokenGrantsStakeBalance,
    minimumStake,
    delegations,
    undelegations,
  } = state.data

  const handleSubmit = async (values, onTransactionHashCallback) => {
    try {
      await tokensPageService.delegateStake(web3Context, values, onTransactionHashCallback)
      showMessage({ type: messageType.SUCCESS, title: 'Success', content: 'Staking delegate transaction has been successfully completed' })
    } catch (error) {
      showMessage({ type: messageType.ERROR, title: 'Staking delegate action has been failed ', content: error.message })
      throw error
    }
  }

  return (
    <LoadingOverlay isFetching={state.isFetching}>
      <h2 className="mb-2">My Tokens</h2>
      <div className="tokens-wrapper flex wrap flex-1 row space-between">
        <section id="delegate-stake-section" className="tile">
          <h2 className="text-grey-70 mb-1">Delegate Stake</h2>
          <div className="text-big text-black">
            Earn ETH rewards by delegating stake to an operator address.
            All ETH rewards will be sent to the address you set as the beneficiary.
          </div>
          <SpeechBubbleInfo>
            A&nbsp;<span className="text-bold">stake</span>&nbsp;is an amount of KEEP
            that’s bonded in order to participate in the threshold relay and, optionally, the Keep network.
          </SpeechBubbleInfo>
          <hr/>
          <DelegateStakeForm
            onSubmit={handleSubmit}
            minStake={minimumStake}
            keepBalance={ownedKeepBalance}
            grantBalance={tokenGrantsBalance}
          />
        </section>
        <TokensOverview
          keepBalance={ownedKeepBalance}
          stakingBalance={tokenStakingBalance}
          pendingUndelegationBalance={pendingUndelegationBalance}
          grantBalance={tokenGrantsBalance}
          tokenGrantsStakeBalance={tokenGrantsStakeBalance}
          undelegationPeriod={undelegationPeriod}
        />
      </div>
      <Undelegations
        undelegations={undelegations}
        successUndelegationCallback={refreshData}
      />
      <DelegatedTokensList
        delegatedTokens={delegations}
        cancelStakeSuccessCallback={refreshData}
      />
    </LoadingOverlay>
  )
}

export default TokensPage

const useSubscribeToStakedEvent = async (data, setData) => {
  const {
    yourAddress,
    grantContract,
    stakingContract,
  } = useContext(Web3Context)

  const {
    ownedKeepBalance,
    tokenStakingBalance,
    delegations,
    initializationPeriod,
  } = data

  const subscribeToEventCallback = async (event) => {
    const { blockNumber, returnValues: { from, value } } = event
    const owner = await stakingContract.methods.ownerOf(from).call()
    let isFromGrant = false

    if (isSameEthAddress(grantContract.options.address, owner)) {
      isFromGrant = true
      const { grantId } = await grantContract.methods.getGrantStakeDetails(from).call()
      const { grantee } = await grantContract.methods.getGrant(grantId).call()
      if (!isSameEthAddress(grantee, yourAddress)) {
        return
      }
    } else if (!isSameEthAddress(owner, yourAddress)) {
      return
    }

    const beneficiary = await stakingContract.methods.magpieOf(from).call()
    const authorizerAddress = await stakingContract.methods.authorizerOf(from).call()

    const delegation = {
      createdAt: blockNumber,
      operatorAddress: from,
      authorizerAddress,
      beneficiary,
      amount: value,
    }
    let keepBalance = ownedKeepBalance
    let keepStakingBalance = tokenStakingBalance

    const initializationOverAt = add(blockNumber || 0, initializationPeriod)
    delegation.isInInitializationPeriod = true
    delegation.initializationOverAt = initializationOverAt.toString()
    if (!isFromGrant) {
      keepBalance = sub(keepBalance, value)
      keepStakingBalance = add(keepStakingBalance, value)
    }

    setData({
      ...data,
      ownedKeepBalance: keepBalance,
      tokenStakingBalance: keepStakingBalance,
      delegations: [delegation, ...delegations],
    })
  }
  useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    'Staked',
    subscribeToEventCallback
  )
}

const useSubscribeToUndelegatedEvent = (data, setData) => {
  const {
    yourAddress,
    grantContract,
    stakingContract,
  } = useContext(Web3Context)

  const {
    pendingUndelegationBalance,
    undelegations,
    undelegationPeriod,
    delegations,
    tokenStakingBalance,
  } = data

  const subscribeToEventCallback = async (event) => {
    const { blockNumber, returnValues: { operator } } = event
    const owner = await stakingContract.methods.ownerOf(operator).call()
    let isFromGrant = false

    if (isSameEthAddress(grantContract.options.address, owner)) {
      isFromGrant = true
      const { grantId } = await grantContract.methods.getGrantStakeDetails(operator).call()
      const { grantee } = await grantContract.methods.getGrant(grantId).call()
      if (!isSameEthAddress(grantee, yourAddress)) {
        return
      }
    } else if (!isSameEthAddress(owner, yourAddress)) {
      return
    }

    const beneficiary = await stakingContract.methods.magpieOf(operator).call()
    const authorizerAddress = await stakingContract.methods.authorizerOf(operator).call()
    const { amount } = await stakingContract.methods.getDelegationInfo(operator).call()

    const undelegation = {
      createdAt: blockNumber,
      operatorAddress: operator,
      authorizerAddress,
      beneficiary,
      amount,
    }
    const updatedDelegations = [...delegations]
    const { indexInArray } = findIndexAndObject('operatorAddress', operator, updatedDelegations, compareEthAddresses)
    if (indexInArray !== null) {
      updatedDelegations.splice(indexInArray, 1)
    }

    let keepStakingBalance = tokenStakingBalance
    let keepPendingUndelegationBalance = pendingUndelegationBalance

    undelegation.undelegationCompleteAt = add(blockNumber, undelegationPeriod)
    undelegation.canRecoverStake = false

    if (!isFromGrant) {
      keepPendingUndelegationBalance = add(pendingUndelegationBalance, amount)
      keepStakingBalance = sub(keepStakingBalance, amount)
    }

    setData({
      ...data,
      tokenStakingBalance: keepStakingBalance,
      pendingUndelegationBalance: keepPendingUndelegationBalance,
      undelegations: [undelegation, ...undelegations],
      delegations: [...updatedDelegations],
    })
  }
  useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    'Undelegated',
    subscribeToEventCallback,
  )
}

const useSubscribeToRecoveredStakeEvent = async (data, setData) => {
  const {
    ownedKeepBalance,
    undelegations,
    pendingUndelegationBalance,
  } = data

  const subscribeToEventCallback = async (event) => {
    const { returnValues: { operator } } = event
    let keepBalance = ownedKeepBalance
    let keepUndelegationBalance = pendingUndelegationBalance

    const updatedUndelegations = [...undelegations]
    const {
      indexInArray,
      obj: recoveredUndelegation,
    } = findIndexAndObject('operatorAddress', operator, updatedUndelegations, compareEthAddresses)

    if (indexInArray !== null) {
      updatedUndelegations.splice(indexInArray, 1)

      if (!recoveredUndelegation.isFromGrant) {
        keepBalance = add(keepBalance, recoveredUndelegation.amount)
        keepUndelegationBalance = sub(keepUndelegationBalance, recoveredUndelegation.amount)
      }
    }

    setData({
      ...data,
      ownedKeepBalance: keepBalance,
      pendingUndelegationBalance: keepUndelegationBalance,
      undelegations: [...updatedUndelegations],
    })
  }

  useSubscribeToContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    'RecoveredStake',
    subscribeToEventCallback
  )
}
