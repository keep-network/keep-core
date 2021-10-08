import React, { useState, useCallback } from "react"
import * as Icons from "./Icons"
import Button from "./Button"
import SelectedWalletModal from "./SelectedWalletModal"

const LedgerModal = ({ connector, connectAppWithWallet, closeModal }) => {
  const [ledgerVersion, setLedgerVersion] = useState("")

  const fetchAccounts = useCallback(
    async (numberOfAccounts, accountsOffset) => {
      try {
        const accounts = await connector[ledgerVersion].getAccounts(
          numberOfAccounts,
          accountsOffset
        )
        return accounts
      } catch (error) {
        throw error
      }
    },
    [connector, ledgerVersion]
  )

  return (
    <SelectedWalletModal
      icon={<Icons.Ledger className="ledger-logo ledger-logo--black" />}
      walletName="Ledger"
      descriptionIcon={<Icons.LedgerDevice />}
      description="Plug in Ledger device and unlock."
      connector={connector[ledgerVersion]}
      connectAppWithWallet={connectAppWithWallet}
      closeModal={closeModal}
      fetchAvailableAccounts={fetchAccounts}
      numberOfAccounts={5}
      withAccountPagination={true}
    >
      <>
        <div
          className="flex column mt-1"
          style={{
            alignSelf: "normal",
            justifyContent: "space-around",
          }}
        >
          <Button
            onClick={() => setLedgerVersion("LEDGER_LIVE")}
            className="btn btn-primary btn-lg mb-1"
          >
            ledger live
          </Button>
          <Button
            onClick={() => setLedgerVersion("LEDGER_LEGACY")}
            className="btn btn-primary btn-lg"
          >
            ledger legacy
          </Button>
        </div>
      </>
    </SelectedWalletModal>
  )
}

export default LedgerModal
