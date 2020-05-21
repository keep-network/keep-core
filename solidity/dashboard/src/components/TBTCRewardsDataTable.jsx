import React, { useState, useCallback } from "react"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import TokenAmount from "./TokenAmount"
import * as Icons from "./Icons"
import { ViewInBlockExplorer } from "./ViewInBlockExplorer"

const TBTCRewardsDataTable = ({ rewards, fetchOperatorByDepositId }) => {
  return (
    <DataTable data={rewards} itemFieldId="depositTokenId">
      <Column
        header="amount"
        field="amount"
        renderContent={({ amount }) => (
          <TokenAmount
            currencyIcon={Icons.TBTC}
            currencyIconProps={{ width: 15, height: 15 }}
            amount={amount}
            amountClassName="text-big text-grey-70"
            withMetricSuffix={false}
          />
        )}
      />
      <Column
        header="transaction hash"
        field="transactionHash"
        renderContent={({ transactionHash }) => (
          <ViewInBlockExplorer type="tx" id={transactionHash} />
        )}
      />
      <Column
        header="deposit token id"
        field="depositTokenId"
        renderContent={({ depositTokenId }) => (
          <AddressShortcut address={depositTokenId} />
        )}
      />
      <Column
        header="operator"
        field="operator"
        renderContent={({ depositTokenId, operatorAddress }) => (
          <OperatorCell
            depositTokenId={depositTokenId}
            operatorAddress={operatorAddress}
            fetchOperatorByDepositId={fetchOperatorByDepositId}
          />
        )}
      />
    </DataTable>
  )
}

const OperatorCell = React.memo(
  ({ depositTokenId, fetchOperatorByDepositId, operatorAddress }) => {
    const [isFetching, setIsFetching] = useState(false)

    const fetchOperatorByDeposit = useCallback(async () => {
      setIsFetching(true)
      await fetchOperatorByDepositId(depositTokenId)
      setIsFetching(false)
    }, [fetchOperatorByDepositId, depositTokenId])

    if (operatorAddress) {
      return <AddressShortcut address={operatorAddress} />
    }

    return (
      <div className="flex row center">
        <span className="text-fade-out">0x000</span>
        <span
          className="flex row center text-secondary cursor-pointer"
          onClick={fetchOperatorByDeposit}
        >
          <Icons.Load style={{ marginRight: "0.5rem" }} />{" "}
          {isFetching ? "Loading" : "Load address"}
        </span>
      </div>
    )
  }
)

export default React.memo(TBTCRewardsDataTable)
