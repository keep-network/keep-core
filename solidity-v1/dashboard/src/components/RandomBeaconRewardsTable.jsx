import React from "react"
import { DataTable, Column } from "./DataTable"
import * as Icons from "./Icons"
import Chip from "./Chip"
import { shortenAddress } from "../utils/general.utils"

const RandomBeaconRewardsTable = ({ data }) => {
  return (
    <DataTable
      data={data}
      itemFieldId="groupPublicKey"
      title="Rewards Status"
      noDataMessage="No rewards"
    >
      <Column header="amount" field="amount" renderContent={renderAmount} />
      <Column header="status" field="status" renderContent={renderStatus} />
      <Column header="rewards period" field="period" />
      <Column
        header="group key"
        field="groupPublicKey"
        renderContent={renderGroupPublicKey}
      />
      {/* <Column header="operator" /> */}
    </DataTable>
  )
}

const renderAmount = ({ amount }) => <AmountCell amount={amount} />

const AmountCell = ({ amount }) => {
  return (
    <div className="flex row center">
      <Icons.KeepOutline
        className="keep-outline--grey-40 mr-1"
        width={32}
        height={32}
      />
      <div>
        <div>{amount} KEEP</div>
        <div className="flex row center">
          <Icons.Rewards width={8} height={8} />
          <span className="text-small text-grey-40">&nbsp;Reward</span>
        </div>
      </div>
    </div>
  )
}

const renderStatus = ({ status }) => (
  <Chip text={status} size="small" color="success" />
)

const renderGroupPublicKey = ({ groupPublicKey }) =>
  shortenAddress(groupPublicKey)

export default RandomBeaconRewardsTable
