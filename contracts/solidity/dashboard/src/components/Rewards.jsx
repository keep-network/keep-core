import React, { useContext, useEffect, useState } from 'react'
import { RewardsGroups } from './RewardsGroups'
import { Web3Context } from './WithWeb3Context'
import Button from './Button'
import Loadable from './Loadable'
import NoData from './NoData'
import * as Icons from './Icons'
import rewardsService from '../services/rewards.service'

export const Rewards = () => {
  const { keepRandomBeaconOperatorContract, stakingContract, yourAddress, utils } = useContext(Web3Context)
  const [isFetching, setIsFetching] = useState(true)
  const [showAll, setShowAll] = useState(false)
  const [data, setData] = useState([])

  useEffect(() => {
    let shouldSetState = true
    fetchGroups(keepRandomBeaconOperatorContract, stakingContract, yourAddress, utils).then((groups) => {
      console.lolg('groups', groups)
      if (shouldSetState) {
        setIsFetching(false)
        setData(groups)
      }
    }).catch((error) => {
      shouldSetState && setIsFetching(false)
    })
    return () => {
      shouldSetState = false
    }
  }, [])

  return (
    <Loadable isFetching={isFetching}>
      { data.length === 0 ?
        <NoData
          title='No rewards yet!'
          iconComponent={<Icons.Badge width={100} height={100} />}
          content='You can withdraw any future earned rewards from your delegated stake on this page.'
        /> :
        <>
          <RewardsGroups groups={showAll ? data : data.slice(0, 3)} />
          { data.length > 3 &&
            <Button
              className="btn btn-default btn-sm see-all-btn"
              onClick={() => setShowAll(!showAll)}
            >
              {showAll ? 'SEE LESS' : `SEE ALL ${data.length - 2}`}
            </Button>
          }
        </>
      }
    </Loadable>
  )
}
const fetchGroups = async (keepRandomBeaconOperatorContract, stakingContract, yourAddress, utils) => {
  try {
    const expiredGroupsCount = await keepRandomBeaconOperatorContract.methods.getFirstActiveGroupIndex().call()
    const groups = []
    const groupMemberIndices = {}
    // TODO iterate trough expired groups
    for (let groupIndex = 0; groupIndex < 10; groupIndex++) {
      const groupPubKey = await keepRandomBeaconOperatorContract.methods.getGroupPublicKey(groupIndex).call()
      const isStale = await keepRandomBeaconOperatorContract.methods.isStaleGroup(groupPubKey).call()
      if (isStale) {
        continue
      }

      const groupMembers = new Set(await keepRandomBeaconOperatorContract.methods.getGroupMembers(groupPubKey).call())
      groupMemberIndices[groupPubKey] = {}
      for (const memberAddress of groupMembers) {
        const beneficiaryAddressForMember = await stakingContract.methods.magpieOf(memberAddress).call()
        if (utils.toChecksumAddress(yourAddress) !== utils.toChecksumAddress(beneficiaryAddressForMember)) {
          continue
        }
        console.log('groupMemberIndices', groupMemberIndices)
        groupMemberIndices[groupPubKey][memberAddress] = await keepRandomBeaconOperatorContract.methods.getGroupMemberIndices(groupPubKey, memberAddress).call()
      }
      if (Object.keys(groupMemberIndices[groupPubKey]).length === 0) {
        continue
      }
      const memberAddressesInGroup = Object.keys(groupMemberIndices[groupPubKey])
      const multipleReward = memberAddressesInGroup.length === 1 ? groupMemberIndices[groupPubKey][memberAddressesInGroup[0]].length : memberAddressesInGroup
        .reduce((prev, current) => (groupMemberIndices[groupPubKey][prev].length + groupMemberIndices[groupPubKey][current].length))
      const reward = utils.toBN((await keepRandomBeaconOperatorContract.methods.getGroupMemberRewards(groupPubKey).call())).mul(utils.toBN(multipleReward))
      groups.push({ groupIndex, groupPubKey, membersIndeces: groupMemberIndices[groupPubKey], reward })
    }
    return Promise.resolve(groups)
  } catch (error) {
    return Promise.reject(error)
  }
}
