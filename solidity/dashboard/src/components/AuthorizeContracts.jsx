import React from "react"
import AddressShortcut from "./AddressShortcut"
import { SubmitButton } from "./Button"
import { ETHERSCAN_DEFAULT_URL } from "../constants/constants"
import { DataTable, Column } from "./DataTable"

const AuthorizeContracts = ({
  contracts,
  onAuthorizeBtn,
  onAuthorizeSuccessCallback,
}) => {
  return (
    <section className="tile">
      <h3 className="text-grey-60">Authorize Contracts</h3>
      <DataTable data={contracts || []} itemFieldId={"contractAddress"}>
        <Column
          header="contract address"
          field="contractAddress"
          renderContent={({ contractAddress }) => (
            <AddressShortcut address={contractAddress} />
          )}
        />
        <Column header="added to the registry" field="blockNumber" />
        <Column
          header="contract details"
          field="details"
          renderContent={({ contractAddress }) => (
            <a
              href={ETHERSCAN_DEFAULT_URL + contractAddress}
              rel="noopener noreferrer"
              target="_blank"
            >
              View in Block Explorer
            </a>
          )}
        />
        <Column
          header=""
          field=""
          renderContent={(contract) => (
            <SubmitButton
              className="btn btn-primary btn-lg flex-1"
              onSubmitAction={(onTransactionHashCallback) =>
                onAuthorizeBtn(contract, onTransactionHashCallback)
              }
              successCallback={onAuthorizeSuccessCallback}
            >
              authorize
            </SubmitButton>
          )}
        />
      </DataTable>
    </section>
  )
}

export default AuthorizeContracts
