import React from "react"
import Tile from "./Tile"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import Button from "./Button"
import { displayAmount } from "../utils/token.utils"
import { useModal } from "../hooks/useModal"
import AddEthModal from "./AddETHModal"

export const BondingSection = ({ data }) => {
  return (
    <Tile
      title="Add ETH for Bonding"
      subtitle="Add an amount of ETH to the available balance."
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
          headerStyle={{ textAlign: "right" }}
          field="availableETH"
          renderContent={({ availableETH, operatorAddress }) => {
            return (
              <AddEthCell
                availableETH={availableETH}
                operatorAddress={operatorAddress}
              />
            )
          }}
        />
      </DataTable>
    </Tile>
  )
}

export default React.memo(BondingSection)

const AddEthCell = ({ availableETH, operatorAddress }) => {
  const { ModalComponent, openModal, closeModal } = useModal()

  return (
    <>
      <ModalComponent title="Add ETH">
        <AddEthModal
          operatorAddress={operatorAddress}
          closeModal={closeModal}
        />
      </ModalComponent>
      <div className="flex">
        <div className="flex row center" style={{ marginLeft: "auto" }}>
          <span className="mr-1">{availableETH}</span>
          <Button onClick={openModal} className="btn btn-secondary btn-sm">
            add eth
          </Button>
        </div>
      </div>
    </>
  )
}
