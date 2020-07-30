import React, { useCallback } from "react"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { useShowMessage, messageType } from "./Message"
import { commitTopUp } from "../services/top-ups.service"
import { displayAmount } from "../utils/token.utils"
import moment from "moment"

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

        throw error
      }
    },
    [showMessage]
  )

  return (
    <DataTable
      data={topUps}
      itemFieldId="operatorAddress"
      title="Available Top-ups"
      noDataMessage="No available top-ups."
    >
      <Column
        header="available top-up amount"
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
        renderContent={({ operatorAddress, createdAt }) => (
          <SubmitButton
            onSubmitAction={async (transactionHashCallback) =>
              await onCommitTopUpBtn(operatorAddress, transactionHashCallback)
            }
            className="btn btn-secondary btn-sm"
            disabled={moment.unix(createdAt).isBefore(moment.now())}
          >
            commit top-up
          </SubmitButton>
        )}
      />
    </DataTable>
  )
}

export default TopUpsDataTable
