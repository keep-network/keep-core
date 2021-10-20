import React from "react"
import NavLink from "../../components/NavLink"
import * as Icons from "../../components/Icons"
import List from "../../components/List"
import { LINK } from "../../constants/constants"

const HowItWorksPage = () => {
  return (
    <section className="threshold-how-it-works-page">
      <section className="tile threshold-how-it-works-page__explanation">
        <Icons.CoveragePool className="threshold-how-it-works-page__explanation__icon" />
        <header>
          <h2 className="text-grey-70">
            How to move KEEP to Threshold Network
          </h2>
          <h3 className="text-grey-50">
            Learn more below about upgrading your KEEP to T and get started with
            the Threshold Network.
          </h3>
        </header>
        <NavLink
          to="/threshold/upgrade"
          className="btn btn-primary btn-md explanation__upgrade-btn"
        >
          upgrade KEEP
        </NavLink>
      </section>

      <section className="tile threshold-how-it-works-page__resources">
        <h3 className="mb-1">Upgrade Granted KEEP Tokens</h3>
        <List>
          <List.Title className="h5 text-violet-80">Released Tokens</List.Title>
          <List.Content className="bullets bullets--violet-80 text-grey-60">
            <List.Item className="mb-1">
              Available released tokens need to be withdrawn to your wallet.
              There are no partial withdrawals.
            </List.Item>
            <List.Item className="mb-1">
              Liquid KEEP in your wallet can be upgraded to T.
            </List.Item>
            <List.Item className="mb-1">
              The withdraw will be done in one transaction.
            </List.Item>
          </List.Content>
        </List>

        <List className="mt-3">
          <List.Title className="h5 text-violet-80">
            Withdrawn Tokens
          </List.Title>
          <List.Content className="bullets bullets--violet-80 text-grey-60">
            <List.Item className="mb-1">
              Withdrawn tokens are liquid KEEP in your wallet.
            </List.Item>
            <List.Item className="mb-1">
              Migrate your liquid KEEP to T in the portal.
            </List.Item>
            <List.Item className="mb-1">
              You can upgrade whatever amount of tokens you want in one single
              transaction.
            </List.Item>
          </List.Content>
        </List>

        <List className="mt-3">
          <List.Title className="h5 text-violet-80">
            Staked Granted Tokens
          </List.Title>
          <List.Content className="bullets bullets--violet-80 text-grey-60">
            <List.Item className="mb-1">
              If you are staking a part of your granted tokens you can either
              opt to undelegate your tokens or continue staking on Keep Network.
            </List.Item>
            <List.Item className="mb-1">
              If you undelegate, there is a 60 day undelegation period.
            </List.Item>
            <List.Item className="mb-1">
              After 60 days, you can withdraw your tokens and upgrade your KEEP
              to T if they are available to release from your token grant.
            </List.Item>
          </List.Content>
        </List>
      </section>

      <section className="tile threshold-how-it-works-page__upgrade-liquid-keep">
        <h3 className="mb-1">Upgrade Liquid KEEP Tokens</h3>
        <p className="text-grey-60">
          You can upgrade any amount of your liquid KEEP tokens in one
          transaction on the Threshold dapp.
        </p>
        <div className="threshold-how-it-works-page__upgrade-liquid-keep__button-container">
          <Icons.TTokenSymbol />
          <h4 className="button-container__title">Threshold</h4>
          <a
            href={LINK.tbtcDapp}
            rel="noopener noreferrer"
            target="_blank"
            className="btn btn-secondary btn-md"
            style={{ marginLeft: "auto" }}
          >
            go to dapp ↗
          </a>
        </div>
      </section>

      <section className="tile threshold-how-it-works-page__upgrade-staked-keep">
        <h3 className="mb-1">Upgrade Staked KEEP Tokens</h3>
        <List>
          <List.Title className="h5 text-violet-80">Staked Tokens</List.Title>
          <List.Content className="bullets bullets--violet-80 text-grey-60">
            <List.Item className="mb-1">
              You can either undelegate your tokens or continue staking on Keep
              Network.
            </List.Item>
            <List.Item className="mb-1">
              If you undelegate, there is a 60 day undelegation period until you
              can withdraw your tokens.
            </List.Item>
            <List.Item className="mb-1">
              After 60 days, you can withdraw your tokens and upgrade your KEEP
              to T.
            </List.Item>
          </List.Content>
        </List>
      </section>
    </section>
  )
}

HowItWorksPage.route = {
  title: "How it Works",
  path: "/threshold/how-it-works",
  exact: true,
}

export default HowItWorksPage
