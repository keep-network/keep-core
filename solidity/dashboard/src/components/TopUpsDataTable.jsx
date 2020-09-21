import React, { useCallback } from "react"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { commitTopUp } from "../actions/web3"
import { displayAmount } from "../utils/token.utils"
import moment from "moment"
import { connect } from "react-redux"

export const TopUpsDataTable = ({
  topUps,
  initializationPeriod,
  commitTopUp,
}) => {
  const onCommitTopUpBtn = useCallback(
    async (operator, awaitingPromise) => {
      commitTopUp(operator, awaitingPromise)
    },
    [commitTopUp]
  )

  return (
    <DataTable
      data={topUps}
      itemFieldId="operatorAddress"
      title="Top-ups"
      noDataMessage="No available top-ups."
    >
      <Column
        header="top-up amount"
        field="availableTopUpAmount"
        renderContent={({ availableTopUpAmount }) =>
          `${displayAmount(availableTopUpAmount)} KEEP`
        }
      />
      <Column
        header="operator"
        field="operatorAddress"
        renderContent={({ operatorAddress }) => (
          <AddressShortcut address={operatorAddress} />
        )}
      />
      <Column
        header=""
        field="operatorAddress"
        renderContent={({ operatorAddress, createdAt, isInUndelegation }) => (
          <SubmitButton
            onSubmitAction={(awaitingPromise) =>
              onCommitTopUpBtn(operatorAddress, awaitingPromise)
            }
            className="btn btn-secondary btn-sm"
            disabled={
              isInUndelegation ||
              !moment
                .unix(createdAt)
                .add(initializationPeriod, "seconds")
                .isBefore(moment.now())
            }
          >
            commit top-up
          </SubmitButton>
        )}
      />
    </DataTable>
  )
}

const mapDispatchToProps = {
  commitTopUp,
}

export default connect(null, mapDispatchToProps)(TopUpsDataTable)
