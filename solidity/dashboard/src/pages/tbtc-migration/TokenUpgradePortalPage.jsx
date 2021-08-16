import React from "react"
import List from "../../components/List"
import TokenAmount from "../../components/TokenAmount"
import { MigrationPortalForm } from "../../components/tbtc-migration"
import { TBTC } from "../../utils/token.utils"

const TokenUpgradePortalPage = () => {
  const tbtcV1Balance = "0"
  const tbtcV2Balance = "0"
  return (
    <section className="tbtc-migration-portal">
      <List className="tbtc-migration-portal__tbtc-balances">
        <List.Title className="h3 text-grey-70">Balance</List.Title>
        <List.Content className="tbtc-balances">
          <List.Item className="tbtc-balance tbtc-balance--v1">
            <TokenAmount
              token={TBTC}
              amount={tbtcV1Balance}
              symbol="tBTC v1"
              amountClassName="h2 text-white"
              symbolClassName="h3 text-white"
              withIcon
            />
          </List.Item>
          <List.Item className="tbtc-balance tbtc-balance--v2">
            <TokenAmount
              token={TBTC}
              amount={tbtcV2Balance}
              symbol="tBTC v2"
              amountClassName="h2 text-black"
              symbolClassName="h3 text-black"
              withIcon
            />
          </List.Item>
        </List.Content>
      </List>
      <section className="tbtc-migration-portal__form-wrapper">
        <h3 className="text-grey-70 mb-1">Migration Portal</h3>
        {/* TODO: Pass props */}
        <MigrationPortalForm />
      </section>
    </section>
  )
}

TokenUpgradePortalPage.route = {
  title: "Token Upgrade Portal",
  path: "/tbtc-migration/portal",
  exact: true,
}

export default TokenUpgradePortalPage
