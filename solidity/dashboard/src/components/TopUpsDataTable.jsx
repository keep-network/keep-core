import React, { useCallback } from "react"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { useShowMessage, messageType } from "./Message"
import { commitTopUp } from "../services/top-ups.service"
import { displayAmount } from "../utils/token.utils"

export const TopUpsDataTable = ({ topUps }) => {
  const showMessage = useShowMessage()

  const onCommitTopUpBtn = useCallback(
    async (operator, transactionHashCallback) => {
      try {
        await commitTopUp(operator, transactionHashCallback)
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content: "Top up committed successfully",
        })
      } catch (error) {
        showMessage({
          type: messageType.ERROR,
          title: "Commit action has failed ",
          content: error.message,
        })
      }
    },
    [showMessage]
  )

  return (
    <DataTable
      data={topUps}
      itemFieldId="operatorAddress"
      title="Available top-ups"
      noDataMessage="No available top-ups."
    >
      <Column
        header="amount"
        field="amount"
        renderContent={({ amount }) => `${displayAmount(amount)} KEEP`}
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
        renderContent={({ operatorAddress }) => (
          <SubmitButton
            onSubmitAction={(transactionHashCallback) =>
              onCommitTopUpBtn(operatorAddress, transactionHashCallback)
            }
            className="btn btn-primary btn-sm"
          >
            commit top-up
          </SubmitButton>
        )}
      />
    </DataTable>
  )
}

export default TopUpsDataTable
