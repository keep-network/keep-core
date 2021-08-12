import React from "react"
import List from "../../components/List"
import * as Icons from "../../components/Icons"
import NavLink from "../../components/NavLink"
import Divider from "../../components/Divider"

const tbtcPools = [
  {
    id: "saddle",
    icon: Icons.SaddleWhite,
    label: "tBTC v2 Pool on Saddle",
    btnText: "go to pool ↗",
    link: "https://saddle.exchange/#/pools/btc/deposit",
    external: true,
  },
  {
    id: "keep",
    icon: Icons.KeepBlackGreen,
    label: "tBTC v1 Pool on Keep",
    btnText: "go to pool",
    link: "/liquidity",
    external: false,
  },
]

const HowItWorksPage = () => {
  return (
    <>
      <section className="tile explanation">
        <h2 className="text-grey-70 mb-1">What is token migration?</h2>
        <h3 className="text-grey-50">
          Token migration is a method which upgrades and downgrades your tBTC
          tokens in one single transaction.
        </h3>
      </section>

      <section className="tile resources">
        <h3 className="mb-1">Why upgrade your tBTC?</h3>
        <List>
          <List.Title className="h5 text-violet-80">how it works</List.Title>
          <List.Content className="bullets bullets--violet-80 text-grey-60">
            <List.Item>
              <strong>
                It costs you zero to upgrade or downgrade your tokens.
              </strong>
              &nbsp;There will be zero Keep Network fees.
            </List.Item>
            <List.Item>
              The upgrade from v1 to v2 is&nbsp;<strong>reversible</strong>
              .&nbsp;This means you can always go back from v2 to v1 if you need
              to.
            </List.Item>
            <List.Item>
              <strong>Why would I downgrade?</strong>&nbsp;You might need to
              downgrade from v2 to v1 in the event of liquidation or redemption.
            </List.Item>
            <List.Item>
              <strong>You like yield farming?</strong>&nbsp;From existing tBTC
              pools which will soon/now transition to utilizing tBTCv2 only.
            </List.Item>
          </List.Content>
        </List>

        <section className="mt-2">
          <h5 className="text-violet-80">documentation</h5>
          <p className="text-grey-60 mb-1">
            If you want to know how the migration works under the hood, please
            check our documentation.
          </p>
        </section>
        <a
          href="https://example.com"
          rel="noopener noreferrer"
          target="_blank"
          className="text-smaller"
        >
          Read the documentation
        </a>
      </section>

      <section className="mint-tbtc-v1">
        <div className="mint-tbtc-v1__icon-wrapper">
          <Icons.TBTC />
        </div>
        <h3 className="text-white">Mint tBTC v1</h3>
        <a
          href="https://dapp.tbtc.network"
          rel="noopener noreferrer"
          target="_blank"
          className="btn btn-secondary btn-md"
        >
          go to dapp ↗
        </a>
      </section>

      <List
        className="tile tbtc-pools-section"
        items={tbtcPools}
        renderItem={renderPoolItem}
      >
        <List.Title className="h2--alt text-grey-70 mb-2">
          tBTC Pools
        </List.Title>
        <List.Content className="tbtc-pools" />
      </List>
    </>
  )
}

const PoolItem = ({
  id,
  icon: IconComponent,
  label,
  btnText,
  link,
  external,
}) => {
  return (
    <List.Item className={`tbtc-pools__item tbtc-pools__item--${id}`}>
      <Divider className="divider divider--tile-fluid" />
      <div className="item__content">
        <IconComponent className="item__icon" />
        <h3 className="item__title">{label}</h3>
        {external ? (
          <a
            href={link}
            rel="noopener noreferrer"
            target="_blank"
            className="btn btn-primary btn-md"
          >
            {btnText}
          </a>
        ) : (
          <NavLink to={link} className="btn btn-primary btn-md">
            {" "}
            {btnText}
          </NavLink>
        )}
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
