import React, { useState, useContext, useMemo, useEffect } from "react"
import Button from "./Button"
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
import { OPERATOR_CONTRACT_NAME, REWARD_STATUS } from "../constants/constants"

const previewDataCount = 10
const initialData = [[], "0"]
const rewardsStatusFilterOptions = [
  { status: REWARD_STATUS.AVAILABLE },
  { status: REWARD_STATUS.ACTIVE },
  { status: REWARD_STATUS.WITHDRAWN },
]

export const Rewards = React.memo(() => {
  const { yourAddress, keepRandomBeaconOperatorContract } = useContext(
    Web3Context
  )
  // fetch rewards
  const [state, updateData] = useFetchData(
    rewardsService.fetchAvailableRewards,
    initialData
  )
  const {
    isFetching,
    data: [rewards, totalRewardsBalance],
  } = state

  // fetch withdrawals
  const [withdrawalHistoryState, updateWithdrawalHistoryData] = useFetchData(
    rewardsService.fetchWithdrawalHistory,
    []
  )
  const { data: withdrawals } = withdrawalHistoryState

  // see more/less button state
  const [showAll, setShowAll] = useState(false)

  // selected reward to withdraw
  const [selectedReward, setSelectedReward] = useState({})

  // filter dropdown
  const [rewardFilter, setRewardFilter] = useState({})

  const [withdrawAction] = useWithdrawAction()

  // subscribe to `GroupMemberRewardsWithdrawn` event
  const { latestEvent } = useSubscribeToContractEvent(
    OPERATOR_CONTRACT_NAME,
    "GroupMemberRewardsWithdrawn"
  )
  const previousWithdrawalEvent = usePrevious(latestEvent)

  useEffect(() => {
    const isSameEvent =
      previousWithdrawalEvent.transactionHash === latestEvent.transactionHash
    if (isEmptyObj(latestEvent) || isSameEvent) {
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

    updateRewards(latestEvent)
    keepRandomBeaconOperatorContract.methods
      .getGroupPublicKey(groupIndex)
      .call()
      .then((groupPublicKey) => {
        const withdrawal = {
          blockNumber,
          groupPublicKey,
          transactionHash,
          reward: web3Utils.fromWei(amount, "ether"),
          status: REWARD_STATUS.WITHDRAWN,
        }
        updateWithdrawalHistoryData([withdrawal, ...withdrawals])
      })
  })

  const updateRewards = (latestEvent) => {
    const {
      returnValues: { groupIndex, amount, operator },
    } = latestEvent
    const { indexInArray, obj } = findIndexAndObject(
      "groupIndex",
      groupIndex,
      rewards
    )
    if (indexInArray === null) {
      return
    }

    const updatedRewards = [...rewards]
    const updatedReward = sub(web3Utils.toWei(obj.reward, "ether"), amount)
    let updateTotalRewardsBalance = sub(
      web3Utils.toWei(totalRewardsBalance, "ether"),
      amount
    )
    updateTotalRewardsBalance = gt(updateTotalRewardsBalance, 0)
      ? updateTotalRewardsBalance
      : "0"

    if (lte(updatedReward, 0)) {
      updatedRewards.splice(indexInArray, 1)
      setSelectedReward({})
    } else {
      const updatedMembersIndeces = { ...obj.membersIndeces }
      delete updatedMembersIndeces[operator]
      const updatedGroup = {
        ...obj,
        membersIndeces: updatedMembersIndeces,
        reward: web3Utils.fromWei(updatedReward, "ether"),
      }
      updatedRewards[indexInArray] = updatedGroup
      setSelectedReward(updatedGroup)
    }

    updateData([
      updatedRewards,
      web3Utils.fromWei(updateTotalRewardsBalance, "ether"),
    ])
  }

  const updateWithdrawStatus = (status) => {
    const { groupIndex } = selectedReward
    const { indexInArray } = findIndexAndObject(
      "groupIndex",
      groupIndex,
      rewards
    )
    if (indexInArray === null) {
      return
    }
    const updatedGroups = [...rewards]
    updatedGroups[indexInArray].status = status

    updateData([updatedGroups, totalRewardsBalance])
  }

  const availableRewardsOptions = useMemo(() => {
    return rewards.filter((group) => group.isStale)
  }, [rewards])

  const rewardsData = useMemo(() => {
    const allRewards = [...rewards, ...withdrawals]
    let rewardsToReturn
    switch (rewardFilter.status) {
      case REWARD_STATUS.AVAILABLE:
        rewardsToReturn = allRewards.filter(({ isStale }) => isStale)
        break
      case REWARD_STATUS.ACTIVE:
        rewardsToReturn = allRewards.filter(
          ({ isStale, status }) => !isStale && !status
        )
        break
      case REWARD_STATUS.WITHDRAWN:
        rewardsToReturn = allRewards.filter(
          ({ status }) => status === REWARD_STATUS.WITHDRAWN
        )
        break
      default:
        rewardsToReturn = allRewards
    }

    return showAll
      ? rewardsToReturn
      : rewardsToReturn.slice(0, previewDataCount)
  }, [rewards, withdrawals, showAll, rewardFilter.status])

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
                  options={availableRewardsOptions}
                  onSelect={setSelectedReward}
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
        <section className="group-items tile">
          <div className="flex row space-between">
            <h4 className="text-grey-70">Rewards Status</h4>
            <Dropdown
              withLabel={false}
              options={rewardsStatusFilterOptions}
              onSelect={setRewardFilter}
              valuePropertyName="status"
              labelPropertyName="status"
              selectedItem={rewardFilter}
              noItemSelectedText="All rewards"
              selectedItemComponent={rewardFilter.status}
              renderOptionComponent={({ status }) => status}
              isFilterDropdow
              allItemsFilterText="All rewards"
            />
          </div>
          <DataTable data={rewardsData} itemFieldId="groupPublicKey">
            <Column
              header="amount"
              field="reward"
              renderContent={({ reward, status }) => (
                <TokenAmount
                  currencyIcon={Icons.ETH}
                  currencyIconProps={{
                    width: 20,
                    height: 20,
                    className: `eth-icon${
                      status === REWARD_STATUS.WITHDRAWN ? " grey-40" : ""
                    }`,
                  }}
                  displayWithMetricSuffix={false}
                  amount={reward}
                  amountClassName={`text-big text-grey-${
                    status === REWARD_STATUS.WITHDRAWN ? "40" : "70"
                  }`}
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
              renderContent={({ groupPublicKey, status }) => (
                <AddressShortcut
                  address={groupPublicKey}
                  classNames={
                    status === REWARD_STATUS.WITHDRAWN ? "text-grey-40" : ""
                  }
                />
              )}
            />
          </DataTable>
          <div className="flex full-center">
            {rewards.length + withdrawals.length > previewDataCount && (
              <Button
                className="btn btn-secondary"
                onClick={() => setShowAll(!showAll)}
              >
                {showAll ? "see less" : "see more"}
              </Button>
            )}
          </div>
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
