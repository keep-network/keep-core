import React from "react"
import { WALLETS } from "../constants/constants"
import * as Icons from "./Icons"
import { useWeb3Context } from "./WithWeb3Context"
import Tooltip from "./Tooltip"

export const CurrentWalletIconTooltipView = ({ connector }) => {
  const getTooltipText = (walletLabel) => {
    return walletLabel === WALLETS.EXPLORER_MODE.label
      ? `You are viewing the site in an ${walletLabel}`
      : `You are connected to ${walletLabel}`
  }

  const renderWalletTypeIcon = () => {
    let tooltipText = ""
    let iconComponent = <></>

    switch (connector?.name) {
      case WALLETS.METAMASK.name:
        tooltipText = getTooltipText(WALLETS.METAMASK.label)
        iconComponent = <Icons.MetaMask className="flex" />
        break
      case WALLETS.LEDGER.name:
        tooltipText = getTooltipText(WALLETS.LEDGER.label)
        iconComponent = (
          <Icons.Ledger className="ledger-logo ledger-logo--black flex" />
        )
        break
      case WALLETS.TREZOR.name:
        tooltipText = getTooltipText(WALLETS.TREZOR.label)
        iconComponent = (
          <Icons.Trezor className="trezor-logo trezor-logo--black flex" />
        )
        break
      case WALLETS.WALLET_CONNECT.name:
        tooltipText = getTooltipText(WALLETS.WALLET_CONNECT.label)
        iconComponent = (
          <Icons.WalletConnect className="wallet-connect-logo wallet-connect-logo--black wallet-connect-logo--small flex" />
        )
        break
      case WALLETS.EXPLORER_MODE.name:
        tooltipText = getTooltipText(WALLETS.EXPLORER_MODE.label)
        iconComponent = (
          <Icons.Explore className="wallet-connect-logo wallet-connect-logo--black wallet-connect-logo--small flex" />
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

const CurrentWalletIconTooltip = () => {
  const { connector } = useWeb3Context()

  return <CurrentWalletIconTooltipView connector={connector} />
}

export default CurrentWalletIconTooltip
