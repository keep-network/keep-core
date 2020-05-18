import React from "react"
import * as Icons from "./Icons"
import Button from "./Button"
import ChooseWalletAddress from "./ChooseWalletAddress"

const LedgerModal = ({ onSelectProvider, accounts, onSelectAccount }) => {
  return (
    <div className="flex column center">
      <div className="flex full-center mb-3">
        <Icons.Ledger />
        <h3 className="ml-1">Ledger</h3>
      </div>
      <Icons.LedgerDevice className="mb3" />
      <span className="text-center">Plug in Ledger device and unlock.</span>
      <div
        className="flex mt-1"
        style={{ alignSelf: "normal", justifyContent: "space-around" }}
      >
        <Button
          onClick={() => onSelectProvider("LEDGER_LIVE")}
          className="btn btn-primary btn-md"
        >
          ledger live
        </Button>
        <Button
          onClick={() => onSelectProvider("LEDGER_LEGACY")}
          className="btn btn-primary btn-md"
        >
          ledger legacy
        </Button>
      </div>
      {accounts && (
        <ChooseWalletAddress
          onSelectAccount={onSelectAccount}
          addresses={accounts}
        />
      )}
    </div>
  )
}

export default LedgerModal
