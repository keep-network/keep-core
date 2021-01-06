import React from "react"
import { useSelector } from "react-redux"
import { useWeb3Context } from "./WithWeb3Context"
import { NetworkStatus } from "./NetworkStatus"
import * as Icons from "./Icons"
import { shortenAddress } from "../utils/general.utils"
import WalletOptions from "./WalletOptions"
import CopyToClipboard from "./CopyToClipboard"
import { displayAmount } from "../utils/token.utils"

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
          <div className="wallet__menu-container">
            <div className="wallet__menu">
              {isActive ? <WalletMenu /> : <WalletOptions />}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

const WalletMenu = () => {
  const { yourAddress } = useWeb3Context()
  const keepTokenBalance = useSelector((state) => state.keepTokenBalance)
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
      <div className="wallet__menu__balance">
        {keepTokenBalance.isFetching
          ? "loading KEEP balance..."
          : `${displayAmount(keepTokenBalance.value)} KEEP`}
      </div>
      {/* TODO add support for dissconnect wallet */}
      {/* <div className="wallet__menu__disconnect" onClick={disconnect}>
        Disconnect
      </div> */}
    </>
  )
}
