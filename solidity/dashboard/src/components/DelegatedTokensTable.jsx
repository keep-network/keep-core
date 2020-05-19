import React from "react"
import { formatDate } from "../utils/general.utils"
import { displayAmount } from "../utils/token.utils"
import AddressShortcut from "./AddressShortcut"
import UndelegateStakeButton from "./UndelegateStakeButton"
import StatusBadge, { BADGE_STATUS } from "./StatusBadge"
import { PENDING_STATUS, COMPLETE_STATUS } from "../constants/constants"
import { DataTable, Column } from "./DataTable"
import Tile from "./Tile"

const DelegatedTokensTable = ({
  delegatedTokens,
  cancelStakeSuccessCallback,
}) => {
  return (
    <Tile title="Delegations">
      <DataTable data={delegatedTokens} itemFieldId="operatorAddress">
        <Column
          header="amount"
          field="amount"
          renderContent={({ amount }) => `${displayAmount(amount)} KEEP`}
        />
        <Column
          header="status"
          field="delegationStatus"
          renderContent={(delegation) => {
            const delegationStatus = delegation.isInInitializationPeriod
              ? PENDING_STATUS
              : COMPLETE_STATUS
            const statusBadgeText =
              delegationStatus === PENDING_STATUS
                ? `${delegationStatus.toLowerCase()}, ${delegation.initializationOverAt.fromNow(
                    true
                  )}`
                : formatDate(delegation.initializationOverAt)

            return (
              <StatusBadge
                status={BADGE_STATUS[delegationStatus]}
                className="self-start"
                text={statusBadgeText}
                onlyIcon={delegationStatus === COMPLETE_STATUS}
              />
            )
          }}
        />
        <Column
          header="beneficiary"
          field="beneficiary"
          renderContent={({ beneficiary }) => (
            <AddressShortcut address={beneficiary} />
          )}
        />
        <Column
          header="operator"
          field="operatorAddress"
          renderContent={({ operatorAddress }) => (
            <AddressShortcut address={operatorAddress} />
          )}
        />
        <Column
          header="authorizer"
          field="authorizerAddress"
          renderContent={({ authorizerAddress }) => (
            <AddressShortcut address={authorizerAddress} />
          )}
        />
        <Column
          header=""
          field=""
          renderContent={(delegation) => (
            <UndelegateStakeButton
              isInInitializationPeriod={delegation.isInInitializationPeriod}
              isFromGrant={delegation.isFromGrant}
              btnClassName="btn btn-sm btn-secondary"
              operator={delegation.operatorAddress}
              isManagedGrant={delegation.isManagedGrant}
              managedGrantContractInstance={
                delegation.managedGrantContractInstance
              }
              successCallback={
                delegation.isInInitializationPeriod
                  ? cancelStakeSuccessCallback
                  : () => {}
              }
            />
          )}
        />
      </DataTable>
    </Tile>
  )
}

export default DelegatedTokensTable
