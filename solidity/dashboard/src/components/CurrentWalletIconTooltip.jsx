import React from "react"
import { WALLETS } from "../constants/constants"
import * as Icons from "./Icons"
import { useWeb3Context } from "./WithWeb3Context"
import Tooltip from "./Tooltip"

const CurrentWalletIconTooltip = () => {
  const { connector } = useWeb3Context()

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
    <div className={"web3-status__wallet-connected-icon-wrapper"}>
      {renderWalletTypeIcon()}
    </div>
  )
}

export default CurrentWalletIconTooltip
