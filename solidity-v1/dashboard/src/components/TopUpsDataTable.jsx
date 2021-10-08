import React, { useCallback } from "react"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { commitTopUp } from "../actions/web3"
import { KEEP } from "../utils/token.utils"
import { connect } from "react-redux"

export const TopUpsDataTable = ({ topUps, commitTopUp }) => {
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
          `${KEEP.displayAmountWithSymbol(availableTopUpAmount)}`
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
        renderContent={({
          operatorAddress,
          readyToBeCommitted,
          isInUndelegation,
        }) => (
          <SubmitButton
            onSubmitAction={(awaitingPromise) =>
              onCommitTopUpBtn(operatorAddress, awaitingPromise)
            }
            className="btn btn-secondary btn-sm"
            disabled={isInUndelegation || !readyToBeCommitted}
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
