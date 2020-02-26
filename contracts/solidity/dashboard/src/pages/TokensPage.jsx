import React, { useContext } from 'react'
import DelegateStakeForm from '../components/DelegateStakeForm'
import TokensOverview from '../components/TokensOverview'
import Undelegations from '../components/Undelegations'
import { useFetchData } from '../hooks/useFetchData'
import { tokensPageService } from '../services/tokens-page.service'
import DelegatedTokensList from '../components/DelegatedTokensList'
import { Web3Context } from '../components/WithWeb3Context'
import { useShowMessage, messageType } from '../components/Message'
import web3Utils from 'web3-utils'
import SpeechBubbleInfo from '../components/SpeechBubbleInfo'
import { LoadingOverlay } from '../components/Loadable'

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
    const { stakeTokens, operatorAddress, beneficiaryAddress, context } = values
    try {
      await tokensPageService.delegateStake(web3Context, values, onTransactionHashCallback)
      showMessage({ type: messageType.SUCCESS, title: 'Success', content: 'Staking delegate transaction has been successfully completed' })
      if (context === 'grant') {
        return
      }
      const amount = web3Utils.toBN(stakeTokens).mul(web3Utils.toBN(10).pow(web3Utils.toBN(18)))
      const updatedKeepBalance = web3Utils.toBN(ownedKeepBalance).sub(amount)
      const updatedTokenStakingBalance = web3Utils.toBN(tokenStakingBalance).add(amount)
      const updatedDelegations = [{ operatorAddress, beneficiary: beneficiaryAddress, amount: amount.toString() }, ...delegations]
      setData({
        ...state.data,
        ownedKeepBalance: updatedKeepBalance,
        tokenStakingBalance: updatedTokenStakingBalance,
        delegations: updatedDelegations,
      })
    } catch (error) {
      showMessage({ type: messageType.ERROR, title: 'Staking delegate action has been failed ', content: error.message })
    }
  }

  return (
    <LoadingOverlay isFetching={state.isFetching}>
      <h3>My Tokens</h3>
      <div className="tokens-wrapper flex flex-1 flex-row-space-between flex-wrap">
        <section id="delegate-stake-section" className="tile">
          <h2 className="text-darker-grey">Delegate Stake</h2>
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
        successDelegationCallback={refreshData}
      />
    </LoadingOverlay>
  )
}

export default TokensPage
