import React, { useCallback } from "react"
import Tile from "./Tile"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import Button from "./Button"
import { displayAmount } from "../utils/token.utils"
import { useModal } from "../hooks/useModal"
import AddEthModal from "./AddETHModal"
import { tbtcAuthorizationService } from "../services/tbtc-authorization.service"
import { useWeb3Context } from "./WithWeb3Context"
import { useShowMessage, messageType } from "./Message"

export const BondingSection = ({ data }) => {
  return (
    <Tile
      title="Add ETH for Bonding"
      subtitle="Add an amount of ETH to the available balance to be eligible for signing group selection.
      NOTE: Withdrawn ETH will go to the beneficiary address." // TODO: Make NOTE just like in figma.
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

  const web3Context = useWeb3Context()
  const { yourAddress, web3 } = web3Context
  const showMessage = useShowMessage()

  const withdrawAllUnbondedAmount = useCallback(
    async (onTransactionHashCallback) => {
      try {
        await tbtcAuthorizationService.withdrawAllEthForOperator(
          web3Context,
          { operatorAddress, availableETH },
          onTransactionHashCallback
        )
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content: "Withdrawal of ETH transaction successfully completed",
        })
      } catch (error) {
        showMessage({
          type: messageType.ERROR,
          title: "Withrawal of ETH has failed ",
          content: error.message,
        })
        throw error
      }
    },
    [operatorAddress, showMessage, web3Context]
  )

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
      <div className="flex">
        <div className="flex row center" style={{ marginLeft: "auto" }}>
          <Button onClick={withdrawAllUnbondedAmount} className="btn btn-secondary btn-sm">
            withdraw
          </Button>
        </div>
      </div>
    </>
  )
}