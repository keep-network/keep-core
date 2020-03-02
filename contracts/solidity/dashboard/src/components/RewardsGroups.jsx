import React, { useState, useEffect, useContext } from 'react'
import { RewardsGroupItem } from './RewardsGroupItem'
import { SeeAllButton } from './SeeAllButton'
import NoData from './NoData'
import * as Icons from './Icons'
import { LoadingOverlay } from './Loadable'
import { useFetchData } from '../hooks/useFetchData'
import rewardsService from '../services/rewards.service'
import { Web3Context } from './WithWeb3Context'

const previewDataCount = 3

export const RewardsGroups = ({ setTotalRewardsBalance }) => {
  const [state, updateData] = useFetchData(rewardsService.fetchAvailableRewards, [[], '0'])
  const { isFetching, data: [groups, totalRewardsBalance] } = state
  const [showAll, setShowAll] = useState(false)
  const updateGroupsAfterWithdrawalAction = useUpdateGroupRewardAfterWithdrawal()

  useEffect(() => {
    setTotalRewardsBalance(totalRewardsBalance)
  }, [totalRewardsBalance])

  const updateGroupsAfterWithdrawal = (groupToUpdate) => {
    updateData(updateGroupsAfterWithdrawalAction(groupToUpdate, groups, totalRewardsBalance))
  }

  const renderGroupItem = (group) => (
    <RewardsGroupItem
      key={group.groupIndex}
      group={group}
      updateGroupsAfterWithdrawal={updateGroupsAfterWithdrawal}
    />
  )

  return (
    <LoadingOverlay isFetching={isFetching} classNames='group-items self-start'>
      <section className="tile">
        <h3 className='text-grey-70 mb-2'>Available to Withdraw</h3>
        <div className="flex row center">
          <div className="flex-1 text-label">
            group key
          </div>
          <div className="flex-2 text-label">
            amount
          </div>
          <div className="flex-1 text-label"/>
        </div>
        <ul className='group-items self-start'>
          { groups.length === 0 ?
            <NoData
              title='No rewards yet!'
              iconComponent={<Icons.Badge width={100} height={100} />}
              content='You can withdraw any future earned rewards from your delegated stake on this page.'
            /> :
          <>
            {showAll ? groups.map(renderGroupItem) : groups.slice(0, previewDataCount).map(renderGroupItem)}
            <SeeAllButton
              dataLength={groups.length}
              previewDataCount={previewDataCount}
              onClickCallback={() => setShowAll(!showAll)}
              showAll={showAll}
            />
          </>
          }
        </ul>
      </section>
    </LoadingOverlay>
  )
}

const useUpdateGroupRewardAfterWithdrawal = () => {
  const { utils } = useContext(Web3Context)

  const updateGroupAfterWithdrawal = (groupToUpdate, groups, totalRewardsBalance) => {
    const updatedGroups = [...groups]
    let indexInArray
    const currentGroup = groups.find((group, index) => {
      if (group.groupIndex === groupToUpdate.groupIndex) {
        indexInArray = index
        return true
      }

      return false
    })

    let updateTotalRewardsBalance = utils.toBN(utils.toWei(totalRewardsBalance, 'ether'))

    if (Object.keys(groupToUpdate.membersIndeces).length === 0) {
      updateTotalRewardsBalance = updateTotalRewardsBalance.sub(utils.toBN(utils.toWei(currentGroup.reward, 'ether')))
      updatedGroups.splice(indexInArray, 1)
    } else {
      const currentGroupRewardInWei = utils.toBN(utils.toWei(currentGroup.reward, 'ether'))
      const groupToUpdateRewardInWei = utils.toBN(utils.toWei(groupToUpdate.reward, 'ether'))
      updateTotalRewardsBalance = updateTotalRewardsBalance.sub(currentGroupRewardInWei.sub(groupToUpdateRewardInWei))
      updatedGroups[indexInArray] = groupToUpdate
    }
    return [updatedGroups, utils.fromWei(updateTotalRewardsBalance, 'ether')]
  }

  return updateGroupAfterWithdrawal
}
