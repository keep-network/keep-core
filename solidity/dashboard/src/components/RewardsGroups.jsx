import React, { useState, useContext, useMemo, useEffect } from "react"
import { SeeAllButton } from "./SeeAllButton"
import { LoadingOverlay } from "./Loadable"
import { useFetchData } from "../hooks/useFetchData"
import rewardsService from "../services/rewards.service"
import Dropdown from "./Dropdown"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { useShowMessage, messageType, useCloseMessage } from "./Message"
import { Web3Context } from "./WithWeb3Context"
import { findIndexAndObject } from "../utils/array.utils"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { COMPLETE_STATUS, PENDING_STATUS } from "../constants/constants"
import SelectedRewardDropdown from "./SelectedRewardDropdown"
import { isEmptyObj, isSameEthAddress } from "../utils/general.utils"
import { sub, lte, gt } from "../utils/arithmetics.utils"
import web3Utils from "web3-utils"
import { usePrevious } from "../hooks/usePrevious"

const previewDataCount = 3
const initialData = [[], "0"]

export const RewardsGroups = React.memo(({ latestWithdrawalEvent }) => {
  const { yourAddress } = useContext(Web3Context)
  const [state, updateData] = useFetchData(
    rewardsService.fetchAvailableRewards,
    initialData
  )
  const {
    isFetching,
    data: [groups, totalRewardsBalance],
  } = state
  const [showAll, setShowAll] = useState(false)
  const [selectedReward, setSelectedReward] = useState({})
  const [withdrawAction] = useWithdrawAction()
  const previousWithdrawalEvent = usePrevious(latestWithdrawalEvent)

  useEffect(() => {
    if (isEmptyObj(latestWithdrawalEvent)) {
      return
    } else if (
      previousWithdrawalEvent.transactionHash ===
      latestWithdrawalEvent.transactionHash
    ) {
      return
    }

    const {
      returnValues: { groupIndex, amount, operator, beneficiary },
    } = latestWithdrawalEvent
    if (!isSameEthAddress(yourAddress, beneficiary)) {
      return
    }
    const { indexInArray, obj } = findIndexAndObject(
      "groupIndex",
      groupIndex,
      groups
    )
    if (indexInArray === null) {
      return
    }

    const updatedGroups = [...groups]
    const updatedGroupReward = sub(web3Utils.toWei(obj.reward, "ether"), amount)
    let updateTotalRewardsBalance = sub(
      web3Utils.toWei(totalRewardsBalance, "ether"),
      amount
    )
    updateTotalRewardsBalance = gt(updateTotalRewardsBalance, 0)
      ? updateTotalRewardsBalance
      : "0"

    if (lte(updatedGroupReward, 0)) {
      updatedGroups.splice(indexInArray, 1)
      setSelectedReward({})
    } else {
      const updatedMembersIndeces = { ...obj.membersIndeces }
      delete updatedMembersIndeces[operator]
      const updatedGroup = {
        ...obj,
        membersIndeces: updatedMembersIndeces,
        reward: web3Utils.fromWei(updatedGroupReward, "ether"),
      }
      updatedGroups[indexInArray] = updatedGroup
      setSelectedReward(updatedGroup)
    }

    updateData([
      updatedGroups,
      web3Utils.fromWei(updateTotalRewardsBalance, "ether"),
    ])
  }, [
    groups,
    latestWithdrawalEvent,
    previousWithdrawalEvent,
    yourAddress,
    totalRewardsBalance,
    updateData,
  ])

  const updateWithdrawStatus = (status) => {
    const { groupIndex } = selectedReward
    const { indexInArray } = findIndexAndObject(
      "groupIndex",
      groupIndex,
      groups
    )
    if (indexInArray === null) {
      return
    }
    const updatedGroups = [...groups]
    updatedGroups[indexInArray].status = status

    updateData([updatedGroups, totalRewardsBalance])
  }

  const dropdownOptions = useMemo(() => {
    return groups.filter((group) => group.isStale)
  }, [groups])

  return (
    <>
      <LoadingOverlay isFetching={isFetching}>
        <section className="tile total-rewards-section">
          <div className="total-rewards-balance">
            <h3 className="text-grey-70 pb-2">Total Balance</h3>
            <h2 className="balance">{`${totalRewardsBalance} ETH`}</h2>
          </div>
          <section className="withdraw-dropdown-section">
            <h4 className="text-grey-70 text-normal">Withdraw</h4>
            <div className="withdraw-dropdown">
              <div className="dropdown">
                <Dropdown
                  options={dropdownOptions}
                  onSelect={(reward) => setSelectedReward(reward)}
                  valuePropertyName="groupPublicKey"
                  labelPropertyName="groupPublicKey"
                  selectedItem={selectedReward}
                  labelPrefix="Group:"
                  noItemSelectedText="Select Group"
                  label="Choose Amount"
                  selectedItemComponent={
                    <SelectedRewardDropdown groupReward={selectedReward} />
                  }
                  renderOptionComponent={(groupReward) => (
                    <SelectedRewardDropdown groupReward={groupReward} />
                  )}
                />
              </div>
              <SubmitButton
                className="btn btn-primary btn-lg flex-1"
                onSubmitAction={() =>
                  withdrawAction(selectedReward, updateWithdrawStatus)
                }
                disabled={isEmptyObj(selectedReward)}
              >
                withdraw
              </SubmitButton>
            </div>
          </section>
        </section>
      </LoadingOverlay>
      <LoadingOverlay
        isFetching={isFetching}
        classNames="group-items self-start"
      >
        <section className="group-items tile">
          <h3 className="text-grey-70 mb-2">Totals</h3>
          <DataTable
            data={showAll ? groups : groups.slice(0, previewDataCount)}
            itemFieldId="groupPublicKey"
          >
            <Column
              header="amount"
              field="reward"
              renderContent={({ reward }) => `${reward.toString()} ETH`}
            />
            <Column
              header="status"
              field="isStale"
              renderContent={({ isStale, status }) => {
                if (status && status === PENDING_STATUS) {
                  return (
                    <StatusBadge
                      text="pending"
                      status={BADGE_STATUS[PENDING_STATUS]}
                    />
                  )
                } else if (isStale) {
                  return (
                    <StatusBadge
                      text="available"
                      status={BADGE_STATUS[COMPLETE_STATUS]}
                    />
                  )
                } else {
                  return (
                    <div className="text-big text-grey-70">
                      Active
                      <div className="text-smaller">
                        Signing group still working.
                      </div>
                    </div>
                  )
                }
              }}
            />
            <Column
              header="group key"
              field="groupPublicKey"
              renderContent={({ groupPublicKey }) => (
                <AddressShortcut address={groupPublicKey} />
              )}
            />
          </DataTable>
          <SeeAllButton
            dataLength={groups.length}
            previewDataCount={previewDataCount}
            onClickCallback={() => setShowAll(!showAll)}
            showAll={showAll}
          />
        </section>
      </LoadingOverlay>
    </>
  )
})

const useWithdrawAction = () => {
  const web3Context = useContext(Web3Context)
  const showMessage = useShowMessage()
  const closeMessage = useCloseMessage()

  const withdraw = async (group, updateWithdrawStatus) => {
    const { groupIndex, membersIndeces } = group
    try {
      updateWithdrawStatus(PENDING_STATUS)
      const message = showMessage({
        type: messageType.PENDING_ACTION,
        sticky: true,
        title: "Withdrawal in progress",
      })
      const result = await rewardsService.withdrawRewardFromGroup(
        groupIndex,
        membersIndeces,
        web3Context
      )
      closeMessage(message)
      const unacceptedTransactions = result.filter((reward) => reward.isError)
      const errorTransactionCount = unacceptedTransactions.length

      if (errorTransactionCount === 0) {
        showMessage({
          type: messageType.SUCCESS,
          title: "Reward withdrawal completed",
        })
      } else if (errorTransactionCount === result.length) {
        throw new Error("Reward withdrawal failed")
      } else {
        updateWithdrawStatus(null)
        showMessage({
          type: messageType.INFO,
          title: `${errorTransactionCount} of ${result.length} transactions have been not approved`,
        })
      }
    } catch (error) {
      updateWithdrawStatus(null)
      showMessage({ type: messageType.ERROR, title: error.message })
      throw error
    }
  }

  return [withdraw]
}
