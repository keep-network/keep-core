import React, { useState } from "react"
import Tile from "./Tile"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import Button from "./Button"
import { displayAmount } from "../utils/token.utils"
import { useModal } from "../hooks/useModal"
import AddEthModal from "./AddETHModal"
import WithdrawEthModal from "./WithdrawETHModal"
import { colors } from "../constants/colors"
import { gt } from "../utils/arithmetics.utils"

export const BondingSection = ({ data }) => {
  return (
    <Tile
      title="Add ETH for Bonding"
      subtitle={
        <>
          Add an amount of ETH to the available balance to be eligible for
          signing group selection.&nbsp;
          <span className="text-bold text-validation">
            NOTE: Withdrawn ETH will go to the beneficiary address.
          </span>
        </>
      }
    >
      <DataTable data={data} itemFieldId="operatorAddress">
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
            `${displayAmount(stakeAmount)} KEEP`
          }
        />
        <Column header="bonded eth" field="bondedETH" />
        <Column
          header="available eth"
          field="availableETH"
          renderContent={({ availableETH }) => {
            return <AvailableEthCell availableETH={availableETH} />
          }}
        />
        <Column
          header=""
          headerStyle={{ textAlign: "right" }}
          field="availableETH"
          renderContent={({
            availableETH,
            operatorAddress,
            isWithdrawable,
            availableETHInWei,
          }) => {
            return (
              <ActionCell
                availableETH={availableETH}
                operatorAddress={operatorAddress}
                isWithdrawable={isWithdrawable}
                availableETHInWei={availableETHInWei}
              />
            )
          }}
        />
      </DataTable>
    </Tile>
  )
}

export default React.memo(BondingSection)

const AvailableEthCell = React.memo(({ availableETH }) => {
  return (
    <>
      <span
        className="text-big text-grey-70"
        style={{
          textAlign: "right",
          padding: "0.25rem 1rem",
          paddingLeft: "2rem",
          borderRadius: "100px",
          border: `1px solid ${colors.grey20}`,
          backgroundColor: `${colors.grey10}`,
        }}
      >
        {availableETH}
      </span>
      &nbsp;ETH
    </>
  )
})

const ActionCell = React.memo(
  ({ availableETH, availableETHInWei, operatorAddress, isWithdrawable }) => {
    const { openModal, closeModal, ModalComponent } = useModal()
    const [action, setAction] = useState("withdraw")
    const title = action === "add" ? "Add ETH" : "Withdraw ETH"

    return (
      <>
        <ModalComponent title={title}>
          {action === "add" ? (
            <AddEthModal
              operatorAddress={operatorAddress}
              closeModal={closeModal}
            />
          ) : (
            <WithdrawEthModal
              operatorAddress={operatorAddress}
              availableETH={availableETH}
              closeModal={closeModal}
            />
          )}
        </ModalComponent>
        <div
          className="flex row center space-between"
          style={{ marginLeft: "auto" }}
        >
          <Button
            onClick={() => {
              setAction("add")
              openModal()
            }}
            className="btn btn-secondary btn-sm"
          >
            add eth
          </Button>
          <Button
            onClick={() => {
              setAction("withdraw")
              openModal()
            }}
            className="btn btn-secondary btn-sm"
            disabled={!(isWithdrawable && gt(availableETHInWei, 0))}
          >
            withdraw
          </Button>
        </div>
      </>
    )
  }
)
