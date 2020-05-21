import React, { useCallback } from "react"
import Tile from "./Tile"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import Button, { SubmitButton } from "./Button"
import { displayAmount } from "../utils/token.utils"
import { useModal } from "../hooks/useModal"
import AddEthModal from "./AddETHModal"
import { tbtcAuthorizationService } from "../services/tbtc-authorization.service"
import { useWeb3Context } from "./WithWeb3Context"
import { useShowMessage, messageType } from "./Message"
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
    const { ModalComponent, openModal, closeModal } = useModal()

    const web3Context = useWeb3Context()
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
      [operatorAddress, availableETH, showMessage, web3Context]
    )

    return (
      <>
        <ModalComponent title="Add ETH">
          <AddEthModal
            operatorAddress={operatorAddress}
            closeModal={closeModal}
          />
        </ModalComponent>
        <div
          className="flex row center space-between"
          style={{ marginLeft: "auto" }}
        >
          <Button onClick={openModal} className="btn btn-secondary btn-sm">
            add eth
          </Button>
          <SubmitButton
            onSubmitAction={withdrawAllUnbondedAmount}
            className="btn btn-secondary btn-sm"
            disabled={!(isWithdrawable && gt(availableETHInWei, 0))}
          >
            withdraw
          </SubmitButton>
        </div>
      </>
    )
  }
)
