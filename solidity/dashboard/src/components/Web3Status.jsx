import React from "react"
import { useWeb3Context } from "./WithWeb3Context"
import { NetworkStatus } from "./NetworkStatus"
import * as Icons from "./Icons"
import { shortenAddress } from "../utils/general.utils"
import WalletOptions from "./WalletOptions"
import CopyToClipboard from "./CopyToClipboard"

export const Web3Status = () => {
  const { yourAddress, provider } = useWeb3Context()

  const isActive = yourAddress && provider

  return (
    <div className="web3-status">
      <div className="web3-status__content-wrapper">
        <div className="web3-status__network-status">
          <NetworkStatus />
        </div>
        <div className="web3-status__wallet">
          <Icons.Wallet
            className={`wallet__icon${isActive ? "--active" : ""}`}
          />
          <div className={`wallet__address${isActive ? "--active" : ""}`}>
            {isActive ? shortenAddress(yourAddress) : "connect wallet"}
          </div>
          <div className="wallet__menu">
            {isActive ? <WalletMenu /> : <WalletOptions />}
          </div>
        </div>
      </div>
    </div>
  )
}

const WalletMenu = () => {
  const { yourAddress, disconnect } = useWeb3Context()
  return (
    <>
      <CopyToClipboard
        toCopy={yourAddress}
        defaultCopyText="copy address"
        render={(copyProps) => {
          return (
            <div style={{ textAlign: "center", lineHeight: 0 }}>
              <span
                className="wallet__menu__copy-address"
                onClick={copyProps.copyToClipboard}
                onMouseOut={copyProps.reset}
              >
                {copyProps.copyStatus}
              </span>
            </div>
          )
        }}
      />
      <hr className="wallet__menu__divider" />
      {/* TODO: here display keep balance */}
      <div className="wallet__menu__balance">1,000 KEEP</div>
      <div className="wallet__menu__disconnect" onClick={disconnect}>
        Disconnect
      </div>
    </>
  )
}
