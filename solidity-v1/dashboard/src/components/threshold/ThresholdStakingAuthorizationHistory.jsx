import React from "react"
import AddressShortcut from "./../AddressShortcut"
import StatusBadge, { BADGE_STATUS } from "./../StatusBadge"
import { DataTable, Column } from "../DataTable"
import Tile from "./../Tile"
import { KEEP } from "../../utils/token.utils"
import { SubmitButton } from "../Button"
import OnlyIf from "../OnlyIf"

const ThresholdAuthorizationHistory = ({ contracts }) => {
  return (
    <Tile>
      <DataTable
        data={contracts || []}
        title="Threshold staking"
        itemFieldId="contractAddress"
        noDataMessage="No authorization history."
        centered
      >
        <Column
          header="operator"
          field="operatorAddress"
          renderContent={({ operatorAddress }) => (
            <AddressShortcut address={operatorAddress} />
          )}
        />
        <Column
          header="stake"
          field="stakeAmount"
          renderContent={({ stakeAmount }) =>
            `${KEEP.displayAmountWithSymbol(stakeAmount)}`
          }
        />
        <Column
          header="status"
          field="status"
          renderContent={({ isStakedToT, isAuthorized }) => (
            <div className={"flex column center"}>
              {isStakedToT ? (
                <StatusBadge
                  className="self-start mb-1"
                  status={BADGE_STATUS.COMPLETE}
                  text="confirmed"
                />
              ) : (
                <StatusBadge
                  className="self-start mb-1"
                  status={BADGE_STATUS.ERROR}
                  text="missing stake confirmation"
                />
              )}
              <OnlyIf condition={isAuthorized}>
                <StatusBadge
                  className="self-start"
                  status={BADGE_STATUS.COMPLETE}
                  text="authorized"
                />
              </OnlyIf>
            </div>
          )}
        />
        <Column
          headerStyle={{ width: "20%", textAlign: "right" }}
          header="actions"
          tdStyles={{ textAlign: "right" }}
          field=""
          renderContent={() => <AuthorizationHistoryActions />}
        />
      </DataTable>
    </Tile>
  )
}

const AuthorizationHistoryActions = () => {
  return (
    <SubmitButton
      className="btn btn-secondary btn-semi-sm"
      style={{ marginLeft: "auto" }}
    >
      set up pre
    </SubmitButton>
  )
}

export default ThresholdAuthorizationHistory
