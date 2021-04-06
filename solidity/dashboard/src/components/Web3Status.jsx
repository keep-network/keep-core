import React from "react"
import { useSelector } from "react-redux"
import { useWeb3Context } from "./WithWeb3Context"
import { NetworkStatus } from "./NetworkStatus"
import * as Icons from "./Icons"
import { shortenAddress } from "../utils/general.utils"
import WalletOptions from "./WalletOptions"
import CopyToClipboard from "./CopyToClipboard"
import { displayAmount } from "../utils/token.utils"
import { WALLETS } from "../constants/constants"
import Tooltip from "./Tooltip"

export const Web3Status = () => {
  const { yourAddress, isConnected, connector } = useWeb3Context()

  const renderWalletTypeIcon = () => {
    let tooltipText = ""
    let iconComponent = <></>

    switch (connector?.name) {
      case WALLETS.METAMASK.name:
        tooltipText = "You are connected to MetaMask"
        iconComponent = <Icons.MetaMask className="flex" />
        break
      case WALLETS.LEDGER:
        tooltipText = "You are connected to Ledger"
        iconComponent = (
          <Icons.Ledger className="ledger-logo ledger-logo--black flex" />
        )
        break
      case WALLETS.TREZOR:
        tooltipText = "You are connected to Trezor"
        iconComponent = (
          <Icons.Trezor className="trezor-logo trezor-logo--black flex" />
        )
        break
      case WALLETS.WALLET_CONNECT:
        tooltipText = "You are connected to WalletConnect"
        iconComponent = (
          <Icons.WalletConnect className="wallet-connect-logo wallet-connect-logo--black flex" />
        )
        break
      case WALLETS.EXPLORER_MODE.name:
        tooltipText = "You are viewing the site in an Explorer Mode"
        iconComponent = (
          <Icons.Wallet className="wallet-connect-logo wallet-connect-logo--black flex" />
        )
        break
      default:
        tooltipText = ""
        iconComponent = <></>
    }

    return (
      <Tooltip
        simple
        delay={0}
        triggerComponent={() => {
          return iconComponent
        }}
        className={"web3-status__wallet-connected-tooltip"}
        tooltipContentWrapperClassName={
          "tooltip__content-wrapper--lower-position"
        }
      >
        {tooltipText}
      </Tooltip>
    )
  }

  return (
    <div className="web3-status">
      <div className="web3-status__content-wrapper">
        <div className="web3-status__network-status">
          <NetworkStatus />
        </div>
        {connector && (
          <div className={"web3-status__wallet-connected-icon-wrapper"}>
            {renderWalletTypeIcon()}
          </div>
        )}
        <div
          className={`web3-status__wallet ${
            connector?.name === WALLETS.EXPLORER_MODE.name
              ? "web3-status__wallet--explorer-mode"
              : ""
          }`}
        >
          <Icons.Wallet
            className={`wallet__icon${isConnected ? "--active" : ""}`}
          />
          <div className={`wallet__address${isConnected ? "--active" : ""}`}>
            {isConnected ? shortenAddress(yourAddress) : "connect wallet"}
          </div>
          <div className="wallet__menu-container">
            <div className="wallet__menu">
              {isConnected ? <WalletMenu /> : <WalletOptions />}
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
