import React from "react"
import { DataTable, Column } from "./DataTable"
import Chip from "./Chip"
import AddressShortcut from "./AddressShortcut"
import * as Icons from "./Icons"
import { displayAmount } from "../utils/token.utils"
import { REWARD_STATUS } from "../constants/constants"

const TBTCRewardsDataTable = ({ data = [] }) => {
  return (
    <DataTable data={data} itemFieldId="id" noDataMessage="No rewards history.">
      <Column header="amount" field="amount" renderContent={AmountCell} />
      <Column header="status" field="status" renderContent={renderStatus} />
      <Column header="rewards period" field="rewardsPeriod" />
      <Column
        header="operator"
        field="operator"
        renderContent={renderOperator}
      />
    </DataTable>
  )
}

const AmountCell = ({ amount }) => {
  return (
    <div className="flex row center">
      <Icons.KeepOutline
        className="keep-outline--grey-40 mr-1"
        width={32}
        height={32}
      />
      <div>
        <div>{displayAmount(amount)}&nbsp;KEEP</div>
        <div className="flex row center">
          <Icons.Rewards width={8} height={8} />
          <span className="text-small text-grey-40">&nbsp;Reward</span>
        </div>
      </div>
    </div>
  )
}

const renderStatus = ({ status }) => {
  switch (status) {
    case REWARD_STATUS.AVAILABLE:
      return <Chip text={status} />
    case REWARD_STATUS.WITHDRAWN:
    default:
      return <Chip text={status} color="disabled" />
  }
}
const renderOperator = ({ operator }) => <AddressShortcut address={operator} />

export default TBTCRewardsDataTable
