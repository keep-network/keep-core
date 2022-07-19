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
        title="Threshold Staking"
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
          renderContent={({ stakeAmount, isFromGrant }) => {
            return (
              <>
                <div>{KEEP.displayAmountWithSymbol(stakeAmount)}</div>
                <div className={"text-grey-50"} style={{ fontSize: "14px" }}>
                  {isFromGrant ? "Grant Tokens" : "Wallet Tokens"}
                </div>
              </>
            )
          }}
        />
        <Column
          header="status"
          field="status"
          renderContent={({ isStakedToT, isAuthorized, operatorAddress }) => (
            <div className={"flex column center"}>
              {isStakedToT ? (
                <StatusBadge
                  className="self-start mb-1"
                  status={BADGE_STATUS.COMPLETE}
                  text="confirmed"
                />
              ) : (
                <>
                  <StatusBadge
                    className="self-start mb-1"
                    status={BADGE_STATUS.ERROR}
                    text="missing stake confirmation"
                    withTooltip
                    tooltipId={`missing-stake-confirmation-for-operator-${operatorAddress}`}
                    tooltipProps={{
                      place: "top",
                      type: "dark",
                      effect: "solid",
                      className: "react-tooltip-base",
                    }}
                  />
                </>
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
          renderContent={({ operatorAddress, isPRESetUp }) => (
            <AuthorizationHistoryActions
              operatorAddress={operatorAddress}
              isPRESetUp={isPRESetUp}
            />
          )}
        />
      </DataTable>
    </Tile>
  )
}

const AuthorizationHistoryActions = ({ operatorAddress, isPRESetUp }) => {
  const tooltipText = isPRESetUp ? (
    <span>
      Go to the{" "}
      <a
        className={"no-arrow"}
        href={LINK.thresholdDapp}
        rel="noopener noreferrer"
        target="_blank"
      >
        Threshold dashboard
      </a>{" "}
      to manage and claim your rewards. Rewards will be distributed at the end
      of every month.
    </span>
  ) : (
    <span>
      To be eligible to earn monthly rewards you will need to{" "}
      <a
        className={"no-arrow"}
        href={LINK.setUpPRE}
        rel="noopener noreferrer"
        target="_blank"
      >
        set up and run a PRE node
      </a>
      .
    </span>
  )

  const link = isPRESetUp ? LINK.thresholdDapp : LINK.setUpPRE

  return (
    <>
      <ReactTooltip
        id={`set up pre-for-operator-${operatorAddress}`}
        delayHide={300}
        place="top"
        type="dark"
        effect={"solid"}
        className={
          "react-tooltip-base react-tooltip-base--arrow-right react-tooltip-base--stay-on-hover"
        }
        offset={{ left: "100%!important" }}
      >
        <span>{tooltipText}</span>
      </ReactTooltip>
      <a
        href={link}
        rel="noopener noreferrer"
        target="_blank"
        className={`btn btn-secondary btn-semi-sm`}
        style={{
          marginLeft: "auto",
          fontFamily: `"Work-Sans", sans-serif`,
        }}
      >
        {isPRESetUp ? (
          <Icons.QuestionFill
            data-tip
            data-for={`set up pre-for-operator-${operatorAddress}`}
            className={"tooltip--button-corner"}
          />
        ) : (
          <Icons.AlertFill
            data-tip
            data-for={`set up pre-for-operator-${operatorAddress}`}
            className={"tooltip--button-corner"}
          />
        )}
        {isPRESetUp ? (
          <span className={"flex row center"}>
            <Icons.TTokenSymbol
              width={12}
              height={12}
              style={{ marginRight: "0.5rem" }}
            />{" "}
            rewards
          </span>
        ) : (
          <span>
            set up pre <Icons.ArrowTopRight />
          </span>
        )}
      </a>
    </>
  )
}

export default ThresholdAuthorizationHistory
