import React from "react"
import * as Icons from "./Icons"
import ChooseWalletAddress from "./ChooseWalletAddress"

const TrezorModal = ({ accounts, onSelectAccount }) => {
  return (
    <div className="flex column center">
      <div className="flex full-center mb-3">
        <Icons.Ledger />
        <h3 className="ml-1">Trezor</h3>
      </div>
      <Icons.TrezorDevice className="mb3" />
      <span className="text-center">
        Plug in your Trezor device and unlock. If the setup screen doesn’t load
        right away, go to Trezor setup:
      </span>
      <a
        href="https://trezor.io/start/</div>"
        className="btn bt-lg btn-primary mt-3 mb-2"
        target="_blank"
        rel="noopener noreferrer"
      >
        go to trezor setup
      </a>
      {accounts && (
        <ChooseWalletAddress
          onSelectAccount={onSelectAccount}
          addresses={accounts}
        />
      )}
    </div>
  )
}

export default TrezorModal
