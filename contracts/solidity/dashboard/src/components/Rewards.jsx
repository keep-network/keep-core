import React, { useContext, useEffect, useState } from 'react'
import { RewardsGroups } from './RewardsGroups'
import { Web3Context } from './WithWeb3Context'
import Button from './Button'

export const Rewards = () => {
  const { keepRandomBeaconOperatorContract, stakingContract, yourAddress, utils } = useContext(Web3Context)
  const [isFetching, setIsFetching] = useState(true)
  const [data, setData] = useState([])

  useEffect(() => {
    Promise.all([stakingContract.methods.operatorsOfMagpie(yourAddress).call(), keepRandomBeaconOperatorContract.methods.numberOfGroups().call()])
      .then(async ([operators, numberOfGroups]) => {
        const groups = []
        for (let groupIndex=0; groupIndex < numberOfGroups; groupIndex++) {
          const groupPubKey = await keepRandomBeaconOperatorContract.methods.getGroupPublicKey(groupIndex).call()
          const indices = await keepRandomBeaconOperatorContract.methods.getGroupMemberIndices(groupPubKey, operators[1]).call()
          const reward = utils.toBN(await keepRandomBeaconOperatorContract.methods.getGroupMemberRewards(groupPubKey).call()).mul(utils.toBN(indices.length))
          const isStale = await keepRandomBeaconOperatorContract.methods.isStaleGroup(groupPubKey).call()
          groups.push({ groupIndex, groupPubKey, indices, reward, isStale })
        }
        setIsFetching(false)
        setData(groups)
      }).catch((error) => {
        setIsFetching(false)
      })
  }, [])
  return (
    <>
      <RewardsGroups groups={data} />
      <Button
        className="btn btn-primary brn-sm"
      >
        SEE ALL (10)
      </Button>
    </>
  )
}
