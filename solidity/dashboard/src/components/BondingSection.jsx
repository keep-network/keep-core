import React from "react"
import Tile from "./Tile"
import { DataTable, Column } from "./DataTable"
import AddressShortcut from "./AddressShortcut"
import Button from "./Button"
import { displayAmount } from "../utils/token.utils"
import { useModal } from "../hooks/useModal"
import AddBondingTokenModal from "./AddBondingTokenModal"
import WithdrawBondingTokenModal from "./WithdrawBondingTokenModal"
import AvailableBondingTokenAmount from "./AvailableBondingTokenAmount"
import { gt } from "../utils/arithmetics.utils"

export const BondingSection = ({ data }) => {
  return (
    <Tile>
      <DataTable
        data={data}
        itemFieldId="operatorAddress"
        title="Add ERC20 for Bonding"
        subtitle={
          <div>
            Add an amount of ERC20 to the available balance to be eligible for
            signing group selection.&nbsp;
            <span className="text-bold text-validation">
              NOTE: Withdrawn ERC20 will go to the beneficiary address.
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
            `${displayAmount(stakeAmount)} KEEP`
          }
        />
        <Column header="bonded ERC20" field="bondedTokens" />
        <Column
          header="available ERC20"
          field="availableTokens"
          renderContent={renderAvailableBondingTokensContent}
        />
        <Column
          header=""
          headerStyle={{ textAlign: "right", colSpan: "2" }}
          field="availableTokens"
          renderContent={(item) => {
            return <ActionCell {...item} />
          }}
        />
      </DataTable>
    </Tile>
  )
}

const renderAvailableBondingTokensContent = (data) => <AvailableBondingTokenAmount {...data} />

export default React.memo(BondingSection)

const ActionCell = React.memo(
  ({
    availableTokens,
    availableTokensInWei,
    operatorAddress,
    managedGrantAddress,
    isWithdrawableForOperator,
  }) => {
    const { openModal, closeModal } = useModal()

    const onBtnClick = (event) => {
      const action = event.currentTarget.id
      const title = action === "add" ? "Add ERC20" : "Withdraw ERC20"

      const component =
        action === "add" ? (
          <AddBondingTokenModal
            operatorAddress={operatorAddress}
            closeModal={closeModal}
          />
        ) : (
          <WithdrawBondingTokenModal
            operatorAddress={operatorAddress}
            availableTokens={availableTokens}
            availableTokensInWei={availableTokensInWei}
            closeModal={closeModal}
            managedGrantAddress={managedGrantAddress}
          />
        )
      openModal(component, { title })
    }

    return (
      <>
        <div className="flex-gap">
          <Button
            id="add"
            onClick={onBtnClick}
            className="btn btn-secondary btn-sm"
          >
            add ERC20
          </Button>
          <Button
            id="withdraw"
            onClick={onBtnClick}
            className="btn btn-secondary btn-sm"
            disabled={!(isWithdrawableForOperator && gt(availableTokensInWei, 0))}
          >
            withdraw
          </Button>
        </div>
      </>
    )
  }
)
