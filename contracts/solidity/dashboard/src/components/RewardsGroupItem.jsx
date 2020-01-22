import React, { useContext } from 'react'
import Button from './Button'
import AddressShortcut from './AddressShortcut'
import rewardsService from '../services/rewards.service'
import { Web3Context } from './WithWeb3Context'

export const RewardsGroupItem = ({ groupIndex, groupPublicKey, membersIndeces, reward }) => {
  const web3Context = useContext(Web3Context)

  const withdraw = async () => {
    try {
      await rewardsService.withdrawRewardFromGroup(groupIndex, membersIndeces, web3Context)
    } catch (error) {
      console.log('errro', error)
    }
  }
  return (
    <div className='group-item'>
      <div className='group-key'>
        <AddressShortcut address={groupPublicKey} />
        <span>GROUP PUBLIC KEY</span>
      </div>
      <div className='group-reward'>
        <span className='reward-value'>
          {reward.toString()}
        </span>
        <span className='reward-currency'>ETH</span>
      </div>
      <Button
        className='btn btn-primary'
        onClick={withdraw}
      >
        WITHDRAW
      </Button>
    </div>
  )
}
