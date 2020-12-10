import React from "react"
import { DataTable, Column } from "./DataTable"

const TBTCRewardsDataTable = ({ data = [] }) => {
  return (
    <DataTable data={data} itemFieldId="id" noDataMessage="No rewards history.">
      <Column header="amount" field="amount" />
      <Column header="status" field="amount" />
      <Column header="rewards period" field="period" />
      <Column header="operator" field="operator" />
    </DataTable>
  )
}

export default TBTCRewardsDataTable
