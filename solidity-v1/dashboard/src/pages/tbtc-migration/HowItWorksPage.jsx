import React from "react"
import List from "../../components/List"
import * as Icons from "../../components/Icons"
import NavLink from "../../components/NavLink"
import Divider from "../../components/Divider"
import { TBTC_TOKEN_VERSION, LINK } from "../../constants/constants"

const tbtcPools = [
  {
    id: TBTC_TOKEN_VERSION.v2,
    icon: Icons.SaddleWhite,
    type: "Saddle",
    link: LINK.pools.saddle.tbtcV2,
  },
  {
    id: TBTC_TOKEN_VERSION.v1,
    icon: Icons.UniswapLogo,
    type: "Uniswap",
    link: LINK.pools.uniswap.tbtcETH,
  },
]

const HowItWorksPage = () => {
  return (
    <section className="tbtc-migration">
      <section className="tile tbtc-migration__explanation">
        <header>
          <h2 className="text-grey-70 mb-1">What is token migration?</h2>
          <h3 className="text-grey-50">
            Token migration is a method which upgrades your tBTC v1 tokens to
            tBTC v2 tokens in one&nbsp;
            <strong className="text-secondary">single transaction.</strong>
          </h3>
        </header>
        <NavLink
          to="/tbtc-migration/portal"
          className="btn btn-primary btn-md explanation__upgrade-btn"
        >
          upgrade tokens
        </NavLink>
      </section>

      <section className="tile tbtc-migration__resources">
        <h3 className="mb-1">Why upgrade your tBTC?</h3>
        <List>
          <List.Title className="h5 text-violet-80">Upgrade perks</List.Title>
          <List.Content className="bullets bullets--violet-80 text-grey-60">
            <List.Item className="mb-1">
              <strong>
                The upgrade from v1 to v2, or v2 to v1 is reversible.
              </strong>
              &nbsp;This means you can always go back from v2 to v1, v1 to v2 if
              you need to. The reversibility will be working until the tBTC v2
              bridge is launched.
            </List.Item>
            <List.Item className="mb-1">
              <strong>
                The upgrade process is handled in one single transaction.
              </strong>
              &nbsp;By using ApproveAndCall transaction, we will be able to save
              gas for you.
            </List.Item>
            <List.Item className="mb-1">
              <strong>
                It costs you zero to upgrade or downgrade your tokens.
              </strong>
              &nbsp;There will be zero Keep Network fees subsidiesed by the Keep
              Governance, but you will need to pay the Ethereum Network gas
              costs.
            </List.Item>
            <List.Item className="mb-1">
              <strong>You like yield farming?</strong>&nbsp;Existing tBTC pools
              which will soon transition to utilizing tBTCv2 only.
            </List.Item>
          </List.Content>
        </List>

        <List className="mt-3">
          <List.Title className="h5 text-violet-80">downgrade perks</List.Title>
          <List.Content className="bullets bullets--violet-80 text-grey-60">
            <List.Item className="mb-1">
              <strong>Why would I downgrade?</strong>&nbsp;You might need to
              downgrade from v2 to v1 if you need to interact with the v1
              bridge.
            </List.Item>
          </List.Content>
        </List>

        <section className="mt-3">
          <h5 className="text-violet-80">documentation</h5>
          <p className="text-grey-60 mb-1">
            If you want to know how the migration works under the hood, please
            check our documentation.
          </p>
        </section>
        <a
          href={LINK.tbtcMigration.docs}
          rel="noopener noreferrer"
          target="_blank"
          className="text-smaller"
        >
          Read the documentation
        </a>
      </section>

      <section className="tbtc-migration__mint-tbtc-v1">
        <div className="mint-tbtc-v1__icon-wrapper">
          <Icons.TBTC />
        </div>
        <h3>Mint tBTC v1</h3>
        <a
          href={LINK.tbtcDapp}
          rel="noopener noreferrer"
          target="_blank"
          className="btn btn-primary btn-md"
        >
          go to dapp ↗
        </a>
      </section>

      <List
        className="tile tbtc-migration__tbtc-pools"
        items={tbtcPools}
        renderItem={renderPoolItem}
      >
        <List.Title className="h2--alt text-grey-70 mb-2">
          tBTC Pools
        </List.Title>
        <List.Content className="tbtc-pools" />
      </List>
    </section>
  )
}

const PoolItem = ({ id, icon: IconComponent, type, link }) => {
  return (
    <List.Item className={`tbtc-pools__item tbtc-pools__item--${id}`}>
      <Divider className="divider divider--tile-fluid" />
      <div className="item__content">
        <IconComponent className="item__icon" />
        <h3 className="item__title">
          tBTC {id} Pool<p className="mb-0">on {type}</p>
        </h3>
        <a
          href={link}
          rel="noopener noreferrer"
          target="_blank"
          className="btn btn-secondary btn-lg"
        >
          go to pool ↗
        </a>
      </div>
    </List.Item>
  )
}

const renderPoolItem = (item) => <PoolItem key={item.id} {...item} />

HowItWorksPage.route = {
  title: "How it Works",
  path: "/tbtc-migration/how-it-works",
  exact: true,
}

export default HowItWorksPage
