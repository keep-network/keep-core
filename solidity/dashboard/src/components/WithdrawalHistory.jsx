import React, { useState, useContext, useEffect } from "react"
import { Web3Context } from "./WithWeb3Context"
import { SeeAllButton } from "./SeeAllButton"
import { LoadingOverlay } from "./Loadable"
import { useFetchData } from "../hooks/useFetchData"
import rewardsService from "../services/rewards.service"
import { DataTable, Column } from "./DataTable"
import { COMPLETE_STATUS } from "../constants/constants"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import AddressShortcut from "./AddressShortcut"
import moment from "moment"
import web3Utils from "web3-utils"
import {
  formatDate,
  isSameEthAddress,
  isEmptyObj,
} from "../utils/general.utils"
import Tile from "./Tile"
import { usePrevious } from "../hooks/usePrevious"

const previewDataCount = 3
const initialData = []

export const WithdrawalHistory = ({ latestWithdrawalEvent }) => {
  const [state, updateData] = useFetchData(
    rewardsService.fetchWithdrawalHistory,
    initialData
  )
  const { isFetching, data } = state
  const [showAll, setShowAll] = useState(false)
  const { yourAddress, eth, keepRandomBeaconOperatorContract } = useContext(
    Web3Context
  )
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
      blockNumber,
      returnValues: { groupIndex, amount, beneficiary },
    } = latestWithdrawalEvent
    if (!isSameEthAddress(yourAddress, beneficiary)) {
      return
    }
    Promise.all([
      eth.getBlock(blockNumber),
      keepRandomBeaconOperatorContract.methods
        .getGroupPublicKey(groupIndex)
        .call(),
    ]).then(([block, groupPublicKey]) => {
      const withdrawal = {
        blockNumber,
        groupPublicKey,
        date: formatDate(moment.unix(block.timestamp)),
        amount: web3Utils.fromWei(amount, "ether"),
      }
      updateData([withdrawal, ...data])
    })
  })

  return (
    <LoadingOverlay isFetching={isFetching}>
      <Tile title="Rewards History">
        <DataTable
          data={showAll ? data : data.slice(0, previewDataCount)}
          itemFieldId="transactionHash"
        >
          <Column
            header="amount"
            field="amount"
            renderContent={({ amount }) => `${amount.toString()} ETH`}
          />
          <Column
            header="status"
            field="date"
            renderContent={({ date }) => (
              <StatusBadge
                status={BADGE_STATUS[COMPLETE_STATUS]}
                text={date}
                onlyIcon
              />
            )}
          />
          <Column
            header="group key"
            field="groupPublicKey"
            renderContent={({ groupPublicKey }) => (
              <AddressShortcut
                address={groupPublicKey}
                classNames="text-smaller"
              />
            )}
          />
        </DataTable>
        <SeeAllButton
          dataLength={data.length}
          previewDataCount={previewDataCount}
          onClickCallback={() => setShowAll(!showAll)}
          showAll={showAll}
        />
      </Tile>
    </LoadingOverlay>
  )
}
