import React, { useState, useEffect } from "react"
import * as Icons from "./Icons"
import Button from "./Button"
import ChooseWalletAddress from "./ChooseWalletAddress"
import { isEmptyArray } from "../utils/array.utils"
import { useShowMessage, messageType } from "./Message"

const LedgerModal = ({ connector, connectAppWithWallet, closeModal }) => {
  const [accounts, setAccounts] = useState([])
  const [ledgerVersion, setLedgerVersion] = useState("")
  const showMessage = useShowMessage()

  useEffect(() => {
    if (ledgerVersion === "LEDGER_LIVE" || ledgerVersion === "LEDGER_LEGACY") {
      connector[ledgerVersion]
        .getAccounts()
        .then(setAccounts)
        .catch((error) => {
          showMessage({ type: messageType.ERROR, title: error.message })
        })
    }
  }, [ledgerVersion, connector, showMessage])

  const onSelectAccount = async (account) => {
    connector[ledgerVersion].defaultAccount = account
    await connectAppWithWallet(connector[ledgerVersion], "LEDGER")
    closeModal()
  }

  return (
    <div className="flex column center">
      <div className="flex full-center mb-3">
        <Icons.Ledger />
        <h3 className="ml-1">Ledger</h3>
      </div>
      <Icons.LedgerDevice className="mb3" />
      <span className="text-center">Plug in Ledger device and unlock.</span>
      <div
        className="flex column mt-1"
        style={{
          alignSelf: "normal",
          justifyContent: "space-around",
        }}
      >
        <Button
          onClick={() => setLedgerVersion("LEDGER_LIVE")}
          className="btn btn-primary btn-md mb-1"
        >
          ledger live
        </Button>
        <Button
          onClick={() => setLedgerVersion("LEDGER_LEGACY")}
          className="btn btn-primary btn-md"
        >
          ledger legacy
        </Button>
      </div>
      {!isEmptyArray(accounts) && (
        <ChooseWalletAddress
          onSelectAccount={onSelectAccount}
          addresses={accounts}
        />
      )}
    </div>
  )
}

export default LedgerModal
