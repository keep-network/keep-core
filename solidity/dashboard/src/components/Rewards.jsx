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
import { PENDING_STATUS } from "../constants/constants"
import SelectedRewardDropdown from "./SelectedRewardDropdown"
import { isEmptyObj, isSameEthAddress } from "../utils/general.utils"
import { sub, lte, gt } from "../utils/arithmetics.utils"
import web3Utils from "web3-utils"
import Tile from "./Tile"
import { usePrevious } from "../hooks/usePrevious"
import TokenAmount from "./TokenAmount"
import * as Icons from "./Icons"
import RewardsStatus from "./RewardsStatus"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent"
import { OPERATOR_CONTRACT_NAME } from "../constants/constants"

const previewDataCount = 10
const initialData = [[], "0"]

export const Rewards = React.memo(() => {
  const { yourAddress, keepRandomBeaconOperatorContract } = useContext(
    Web3Context
  )
  const [state, updateData] = useFetchData(
    rewardsService.fetchAvailableRewards,
    initialData
  )
  const [withdrawalHistoryState, updateWithdrawalHistoryData] = useFetchData(
    rewardsService.fetchWithdrawalHistory,
    []
  )
  const {
    isFetching,
    data: [groups, totalRewardsBalance],
  } = state
  const [showAll, setShowAll] = useState(false)
  const [selectedReward, setSelectedReward] = useState({})
  const [withdrawAction] = useWithdrawAction()
  const { data: withdrawals } = withdrawalHistoryState
  const { latestEvent } = useSubscribeToContractEvent(
    OPERATOR_CONTRACT_NAME,
    "GroupMemberRewardsWithdrawn"
  )
  const previousWithdrawalEvent = usePrevious(latestEvent)

  useEffect(() => {
    if (isEmptyObj(latestEvent)) {
      return
    } else if (
      previousWithdrawalEvent.transactionHash === latestEvent.transactionHash
    ) {
      return
    }
    const {
      transactionHash,
      blockNumber,
      returnValues: { groupIndex, amount, beneficiary },
    } = latestEvent
    if (!isSameEthAddress(yourAddress, beneficiary)) {
      return
    }
    keepRandomBeaconOperatorContract.methods
      .getGroupPublicKey(groupIndex)
      .call()
      .then((groupPublicKey) => {
        const withdrawal = {
          blockNumber,
          groupPublicKey,
          transactionHash,
          reward: web3Utils.fromWei(amount, "ether"),
          status: "WITHDRAWN",
        }
        updateWithdrawalHistoryData([withdrawal, ...withdrawals])
      })
  })

  useEffect(() => {
    if (isEmptyObj(latestEvent)) {
      return
    } else if (
      previousWithdrawalEvent.transactionHash === latestEvent.transactionHash
    ) {
      return
    }

    const {
      returnValues: { groupIndex, amount, operator, beneficiary },
    } = latestEvent
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
    latestEvent,
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

  const rewardsData = useMemo(() => {
    console.log("updategin rewards data")
    const data = [...groups, ...withdrawals]
    return showAll ? data : data.slice(0, previewDataCount)
  }, [groups, withdrawals, showAll])

  return (
    <>
      <LoadingOverlay isFetching={isFetching}>
        <section className="flex row wrap">
          <Tile
            title="Balance"
            titleClassName="text-grey-70 h2"
            id="rewards-total-balance"
          >
            <h1 className="balance">
              {totalRewardsBalance}
              <span className="h3 mr-1">&nbsp;ETH</span>
            </h1>
          </Tile>
          <Tile title="Available to Withdraw" id="withdraw-dropdown-section">
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
          </Tile>
        </section>
      </LoadingOverlay>
      <LoadingOverlay
        isFetching={isFetching}
        classNames="group-items self-start"
      >
        <Tile title="Rewards Status" className="group-items tile">
          <DataTable data={rewardsData} itemFieldId="groupPublicKey">
            <Column
              header="amount"
              field="reward"
              renderContent={({ reward }) => (
                <TokenAmount
                  currencyIcon={Icons.ETH}
                  currencyIconProps={{ width: 20, height: 20 }}
                  withMetricSuffix={false}
                  amount={reward}
                  amountClassName="text-big text-grey-70"
                  currencySymbol="ETH"
                  displayAmountFunction={(amount) => amount}
                />
              )}
            />
            <Column
              header="status"
              field="isStale"
              renderContent={(rewards) => <RewardsStatus {...rewards} />}
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
        </Tile>
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
