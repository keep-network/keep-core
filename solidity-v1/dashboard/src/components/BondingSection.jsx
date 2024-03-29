import React from "react"
import Tile from "./Tile"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import Button from "./Button"
import { KEEP, ETH } from "../utils/token.utils"
import { useModal } from "../hooks/useModal"
import AvailableEthAmount from "./AvailableEthAmount"
import { gt } from "../utils/arithmetics.utils"
import TokenAmount from "./TokenAmount"
import { MODAL_TYPES } from "../constants/constants"

export const BondingSection = ({ data }) => {
  return (
    <Tile>
      <DataTable
        data={data}
        itemFieldId="operatorAddress"
        title="Add ETH for Bonding"
        subtitle={
          <div>
            Add an amount of ETH to the available balance to be eligible for
            signing group selection.&nbsp;
            <span className="text-bold text-validation">
              NOTE: Withdrawn ETH will go to the beneficiary address.
            </span>
          </div>
        }
        noDataMessage="No bonding data."
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
          header="bonded eth"
          field="bondedETH"
          renderContent={renderBondedETHContent}
        />
        <Column
          header="available eth"
          field="availableETH"
          renderContent={renderAvailableEthContent}
        />
        <Column
          header=""
          headerStyle={{ textAlign: "right", colSpan: "2" }}
          field="availableETH"
          renderContent={(item) => {
            return <ActionCell {...item} />
          }}
        />
      </DataTable>
    </Tile>
  )
}
const renderBondedETHContent = ({ bondedETHInWei }) => (
  <TokenAmount
    token={ETH}
    amount={bondedETHInWei}
    amountClassName=""
    symbolClassName=""
  />
)
const renderAvailableEthContent = (data) => <AvailableEthAmount {...data} />

export default React.memo(BondingSection)

const ActionCell = React.memo(
  ({
    availableETH,
    availableETHInWei,
    operatorAddress,
    managedGrantAddress,
    isWithdrawableForOperator,
  }) => {
    const { openModal } = useModal()

    const onBtnClick = (event) => {
      const action = event.currentTarget.id

      const modalType =
        action === "add"
          ? MODAL_TYPES.BondingAddETH
          : MODAL_TYPES.BondingWithdrawETH

      openModal(modalType, {
        operatorAddress,
        availableETH,
        availableETHInWei,
        managedGrantAddress,
      })
    }

    return (
      <>
        <div className="flex-gap">
          <Button
            id="add"
            onClick={onBtnClick}
            className="btn btn-secondary btn-sm"
          >
            add eth
          </Button>
          <Button
            id="withdraw"
            onClick={onBtnClick}
            className="btn btn-secondary btn-sm"
            disabled={!(isWithdrawableForOperator && gt(availableETHInWei, 0))}
          >
            withdraw
          </Button>
        </div>
      </>
    )
  }
)
