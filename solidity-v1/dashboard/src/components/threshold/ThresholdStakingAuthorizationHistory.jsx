import React from "react"
import AddressShortcut from "./../AddressShortcut"
import StatusBadge, { BADGE_STATUS } from "./../StatusBadge"
import { DataTable, Column } from "../DataTable"
import Tile from "./../Tile"
import { KEEP } from "../../utils/token.utils"
import OnlyIf from "../OnlyIf"
import { LINK } from "../../constants/constants"
import * as Icons from "../Icons"
import ReactTooltip from "react-tooltip"

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
          renderContent={({ operatorAddress }) => (
            <AuthorizationHistoryActions operatorAddress={operatorAddress} />
          )}
        />
      </DataTable>
    </Tile>
  )
}

const AuthorizationHistoryActions = ({ operatorAddress }) => {
  return (
    <a
      href={LINK.setUpPRE}
      rel="noopener noreferrer"
      target="_blank"
      className={`btn btn-secondary btn-semi-sm`}
      style={{
        marginLeft: "auto",
        fontFamily: `"Work-Sans", sans-serif`,
      }}
    >
      <Icons.QuestionFill
        data-tip
        data-for={`set up pre-for-operator-${operatorAddress}`}
        className={"tooltip--button-corner"}
      />
      <ReactTooltip
        id={`set up pre-for-operator-${operatorAddress}`}
        place="top"
        type="dark"
        effect={"solid"}
        className={"react-tooltip-base react-tooltip-base--arrow-right"}
        offset={{ left: "100%!important" }}
      >
        <span>
          To be eligible to earn monthly rewards you will need to set up and run
          a PRE node.
        </span>
      </ReactTooltip>
      set up pre <Icons.ArrowTopRight />
    </a>
  )
}

export default ThresholdAuthorizationHistory
