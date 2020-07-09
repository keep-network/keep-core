import React, { useState, useMemo } from "react"
import Button from "./Button"
import { LoadingOverlay } from "./Loadable"
import { useFetchData } from "../hooks/useFetchData"
import rewardsService from "../services/rewards.service"
import DataTableSkeleton from "./skeletons/DataTableSkeleton"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { useShowMessage, messageType } from "./Message"
import { useWeb3Context } from "./WithWeb3Context"
import { findIndexAndObject } from "../utils/array.utils"
import { PENDING_STATUS } from "../constants/constants"
import { isSameEthAddress } from "../utils/general.utils"
import { sub } from "../utils/arithmetics.utils"
import web3Utils from "web3-utils"
import Tile from "./Tile"
import TokenAmount from "./TokenAmount"
import * as Icons from "./Icons"
import RewardsStatus from "./RewardsStatus"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent"
import {
  OPERATOR_CONTRACT_NAME,
  REWARD_STATUS,
  SIGNING_GROUP_STATUS,
} from "../constants/constants"
import { SpeechBubbleTooltip } from "./SpeechBubbleTooltip"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import Skeleton from "./skeletons/Skeleton"

const previewDataCount = 10
const initialData = [[], "0"]
const rewardsStatusFilterOptions = [
  { status: REWARD_STATUS.AVAILABLE },
  { status: REWARD_STATUS.ACCUMULATING },
  { status: REWARD_STATUS.WITHDRAWN },
]

export const Rewards = React.memo(() => {
  const web3Context = useWeb3Context()
  const showMessage = useShowMessage()

  const { yourAddress, keepRandomBeaconOperatorContract } = web3Context
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

  // filter dropdown
  const [rewardFilter, setRewardFilter] = useState({})

  const withdrawnEventCallback = (latestEvent) => {
    const {
      transactionHash,
      blockNumber,
      returnValues: { groupIndex, amount, beneficiary, operator },
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
          operatorAddress: operator,
          status: REWARD_STATUS.WITHDRAWN,
          groupStatus: SIGNING_GROUP_STATUS.COMPLETED,
        }
        updateWithdrawalHistoryData([withdrawal, ...withdrawals])
      })
  }

  // subscribe to `GroupMemberRewardsWithdrawn` event
  useSubscribeToContractEvent(
    OPERATOR_CONTRACT_NAME,
    "GroupMemberRewardsWithdrawn",
    withdrawnEventCallback
  )

  const updateRewards = (latestEvent) => {
    const {
      returnValues: { groupIndex, operator, amount },
    } = latestEvent
    const { indexInArray } = findIndexAndObject(
      "groupIndex",
      groupIndex,
      rewards,
      (object) => isSameRewardRecord(object, groupIndex, operator)
    )

    if (indexInArray === null) {
      return
    }

    const updateTotalRewardsBalance = sub(
      web3Utils.toWei(totalRewardsBalance, "ether"),
      amount
    )
    const updatedRewards = [...rewards]
    updatedRewards.splice(indexInArray, 1)

    updateData([
      updatedRewards,
      web3Utils.fromWei(updateTotalRewardsBalance, "ether"),
    ])
  }

  const withdrawReward = async (
    operatorAddress,
    groupIndex,
    onTransactionHashCallback
  ) => {
    try {
      updateRewardStatus(PENDING_STATUS, groupIndex, operatorAddress)
      await rewardsService.withdrawRewardFromGroup(
        web3Context,
        { operatorAddress, groupIndex },
        onTransactionHashCallback
      )
      showMessage({
        type: messageType.SUCCESS,
        title: "Success",
        content: "Withdrawal successfully completed",
      })
    } catch (error) {
      showMessage({
        type: messageType.ERROR,
        title: "Withdrawal action has failed ",
        content: error.message,
      })
      throw error
    }
  }

  const updateRewardStatus = (status, groupIndex, operator) => {
    const { indexInArray } = findIndexAndObject(
      "groupIndex",
      groupIndex,
      rewards,
      (object) => isSameRewardRecord(object, groupIndex, operator)
    )
    if (indexInArray === null) {
      return
    }
    const updatedGroups = [...rewards]
    updatedGroups[indexInArray].status = status

    updateData([updatedGroups, totalRewardsBalance])
  }

  const rewardsData = useMemo(() => {
    const allRewards = [...rewards, ...withdrawals]
    let rewardsToReturn = []

    if (!rewardFilter.status) {
      rewardsToReturn = allRewards
    } else {
      rewardsToReturn = allRewards.filter(
        ({ status }) => status === rewardFilter.status
      )
    }

    return showAll
      ? rewardsToReturn
      : rewardsToReturn.slice(0, previewDataCount)
  }, [rewards, withdrawals, showAll, rewardFilter.status])

  return (
    <>
      <Tile title="Total Balance" titleClassName="text-grey-70 h2">
        <header className="flex row center">
          {isFetching ? (
            <Skeleton className="h1 mb-1" styles={{ width: "25%" }} />
          ) : (
            <>
              <h1 className="balance">
                {totalRewardsBalance}
                <span className="h3 mr-1">&nbsp;ETH</span>
              </h1>
              <SpeechBubbleTooltip text="The total balance reflects the total Available and Active rewards. Available rewards are ready to be withdrawn. Active rewards become available after a signing group expires." />
            </>
          )}
        </header>
        <div className="flex row wrap">
          <StatusBadge
            className="mr-1"
            text={REWARD_STATUS.AVAILABLE}
            status={BADGE_STATUS.COMPLETE}
          />
          <StatusBadge
            text={REWARD_STATUS.ACCUMULATING}
            status={BADGE_STATUS.ACTIVE}
          />
        </div>
      </Tile>
      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<DataTableSkeleton columns={6} />}
      >
        <section className="group-items tile">
          <DataTable
            data={rewardsData}
            itemFieldId="groupPublicKey"
            title="Rewards Status"
            withFilterDropdown
            filterDropdownProps={{
              withLabel: false,
              options: rewardsStatusFilterOptions,
              onSelect: setRewardFilter,
              valuePropertyName: "status",
              labelPropertyName: "status",
              selectedItem: rewardFilter,
              noItemSelectedText: "All rewards",
              selectedItemComponent: rewardFilter.status,
              renderOptionComponent: ({ status }) => status,
              isFilterDropdow: true,
              allItemsFilterText: "All rewards",
            }}
          >
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
              header="signing group"
              field="groupStatus"
              renderContent={({ status, groupStatus }) => (
                <span
                  className={
                    status === REWARD_STATUS.WITHDRAWN ? "text-grey-40" : ""
                  }
                >
                  {groupStatus}
                </span>
              )}
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
            <Column
              header="operator"
              field="operatorAddress"
              renderContent={({ operatorAddress, status }) => (
                <AddressShortcut
                  address={operatorAddress}
                  classNames={
                    status === REWARD_STATUS.WITHDRAWN ? "text-grey-40" : ""
                  }
                />
              )}
            />
            <Column
              header=""
              field="operatorAddress"
              renderContent={({ status, operatorAddress, groupIndex }) =>
                status !== REWARD_STATUS.WITHDRAWN && (
                  <SubmitButton
                    className="btn btn-secondary btn-sm"
                    pendingMessageTitle="Pending rewards withdrawal"
                    disabled={status !== REWARD_STATUS.AVAILABLE}
                    onSubmitAction={(onTransactionHashCallback) =>
                      withdrawReward(
                        operatorAddress,
                        groupIndex,
                        onTransactionHashCallback
                      )
                    }
                  >
                    withdraw
                  </SubmitButton>
                )
              }
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

const isSameRewardRecord = (reward, groupIndex, operator) =>
  reward.groupIndex === groupIndex &&
  isSameEthAddress(operator, reward.operatorAddress)
