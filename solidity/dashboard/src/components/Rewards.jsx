import React, { useEffect, useState, useMemo } from "react"
import { useDispatch } from "react-redux"
import Button from "./Button"
import { LoadingOverlay } from "./Loadable"
import { useFetchData } from "../hooks/useFetchData"
import rewardsService from "../services/rewards.service"
import DataTableSkeleton from "./skeletons/DataTableSkeleton"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { useWeb3Context, useWeb3Address } from "./WithWeb3Context"
import { findIndexAndObject } from "../utils/array.utils"
import { PENDING_STATUS } from "../constants/constants"
import { isSameEthAddress } from "../utils/general.utils"
import { sub, lt, gt } from "../utils/arithmetics.utils"
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
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import Skeleton from "./skeletons/Skeleton"
import { withdrawGroupMemberRewards } from "../actions/web3"
import {
  displayEthAmount,
  MIN_ETH_AMOUNT_TO_DISPLAY_IN_WEI,
} from "../utils/ethereum.utils"
import ResourceTooltip from "./ResourceTooltip"

const previewDataCount = 10
const initialRewardsData = [[], "0"]
const rewardsStatusFilterOptions = [
  { status: REWARD_STATUS.AVAILABLE },
  { status: REWARD_STATUS.ACCUMULATING },
  { status: REWARD_STATUS.WITHDRAWN },
]

export const Rewards = () => {
  const dispatch = useDispatch()
  const { keepRandomBeaconOperatorContract } = useWeb3Context()
  const address = useWeb3Address()

  // fetch rewards
  const [state, updateData, , setFetchAvailableRewardsArgs] = useFetchData(
    rewardsService.fetchAvailableRewards,
    initialRewardsData,
    address
  )
  const {
    isFetching,
    data: [rewards, totalRewardsBalance],
  } = state

  useEffect(() => {
    if (address) {
      setFetchAvailableRewardsArgs([address])
    }
  }, [setFetchAvailableRewardsArgs, address])

  // fetch withdrawals
  const [
    withdrawalHistoryState,
    updateWithdrawalHistoryData,
    ,
    setFetchWithdrawalHistoryArgs,
  ] = useFetchData(rewardsService.fetchWithdrawalHistory, [])
  const { data: withdrawals } = withdrawalHistoryState

  useEffect(() => {
    if (address) {
      setFetchWithdrawalHistoryArgs([address])
    }
  }, [setFetchWithdrawalHistoryArgs, address])

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

    if (!isSameEthAddress(address, beneficiary)) {
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
          reward: amount,
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

    const updateTotalRewardsBalance = sub(totalRewardsBalance, amount)
    const updatedRewards = [...rewards]
    updatedRewards.splice(indexInArray, 1)

    updateData([updatedRewards, updateTotalRewardsBalance])
  }

  const withdrawReward = async (
    operatorAddress,
    groupIndex,
    awaitingPromise
  ) => {
    updateRewardStatus(PENDING_STATUS, groupIndex, operatorAddress)
    dispatch(
      withdrawGroupMemberRewards(operatorAddress, groupIndex, awaitingPromise)
    )
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

  const amountTooltipText = (amount) => {
    return `${displayEthAmount(amount, "gwei", null)} gwei`
  }

  return (
    <>
      <Tile title="Total Balance" titleClassName="text-grey-70 h2">
        <header className="flex row center">
          {isFetching ? (
            <Skeleton className="h1 mb-1" styles={{ width: "25%" }} />
          ) : (
            <div className="flex row mb-1 mt-1">
              <TokenAmount
                currencyIcon={Icons.ETH}
                currencyIconProps={{
                  width: 64,
                  height: 64,
                  className: "eth-icon primary",
                }}
                currencySymbol="ETH"
                displayWithMetricSuffix={false}
                amount={totalRewardsBalance}
                amountClassName="h1 text-primary"
                displayAmountFunction={displayEthAmount}
                withTooltip={
                  lt(totalRewardsBalance, MIN_ETH_AMOUNT_TO_DISPLAY_IN_WEI) &&
                  gt(totalRewardsBalance, 0)
                }
                tooltipText={amountTooltipText(totalRewardsBalance)}
              />
              <div className="ml-1 self-center">
                <ResourceTooltip
                  title="Beacon earnings"
                  content="The total balance reflects the total available and accumulating rewards. Available rewards are ready to be withdrawn. Accumulating rewards become available after a signing group expires."
                  withRedirectButton={false}
                />
              </div>
            </div>
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
                    width: 32,
                    height: 32,
                    className: "eth-icon grey-60",
                  }}
                  displayWithMetricSuffix={false}
                  amount={reward}
                  amountClassName={`text-big text-grey-${
                    status === REWARD_STATUS.WITHDRAWN ? "40" : "70"
                  }`}
                  currencySymbol="ETH"
                  displayAmountFunction={displayEthAmount}
                  withTooltip={lt(reward, MIN_ETH_AMOUNT_TO_DISPLAY_IN_WEI)}
                  tooltipText={amountTooltipText(reward)}
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
                    onSubmitAction={async (awaitingPromise) =>
                      await withdrawReward(
                        operatorAddress,
                        groupIndex,
                        awaitingPromise
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
}

const isSameRewardRecord = (reward, groupIndex, operator) =>
  reward.groupIndex === groupIndex &&
  isSameEthAddress(operator, reward.operatorAddress)
