import React from "react"
import PageWrapper from "../../components/PageWrapper"
import CardContainer from "../../components/CardContainer"
import Card from "../../components/Card"
import * as Icons from "../../components/Icons"
import DoubleIcon from "../../components/DoubleIcon"
import { SubmitButton } from "../../components/Button"

const LiquidityPage = ({ title }) => {
  return (
    <PageWrapper title={title}>
      <CardContainer className={"flex wrap"}>
        <Card className={"tile"}>
          <div className={"liquidity__card-title-section"}>
            <DoubleIcon
              MainIcon={Icons.KeepBlackGreen}
              SecondaryIcon={Icons.EthToken}
              className={`liquidity__double-icon-container`}
            />
            <h2 className={"h2--alt text-grey-70"}>KEEP + ETH</h2>
          </div>
          <div className={"liquidity-card-subtitle-section"}>
            <span className="text-grey-40">Uniswap Pool&nbsp;</span>
            <a
              href="https://github.com/keep-network/keep-core/blob/master/docs/glossary.adoc"
              className="arrow-link text-small"
              style={{ marginLeft: "auto", marginRight: "2rem" }}
            >
              View pool
            </a>
          </div>
          <div className={"liquidity__info text-grey-60"}>
            <div className={"liquidity__info-tile bg-mint-10"}>
              <h2 className={"liquidity__info-tile__title text-mint-100"}>
                200%
              </h2>
              <h6>Anual % yield</h6>
            </div>
            <div className={"liquidity__info-tile bg-mint-10"}>
              <h2 className={"liquidity__info-tile__title text-mint-100"}>
                10%
              </h2>
              <h6>% of total pool</h6>
            </div>
          </div>
          <div className={"liquidity__token-balance"}>
            <span className={"liquidity__token-balance_title text-grey-70"}>
              Reward
            </span>
            <div className={"liquidity__token-balance_values text-grey-70"}>
              <h3 className={"liquidity__token-balance_values_label"}><Icons.KeepOutline /><span>KEEP</span></h3>
              <h3>1,000,000</h3>
            </div>
          </div>
          <div className={"liquidity__add-more-tokens"}>
            <SubmitButton className={`btn btn-primary btn-lg w-100`}>
              add more lp tokens
            </SubmitButton>
          </div>
          <div className={"liquidity__withdraw"}>
            <SubmitButton className={"btn btn-primary btn-lg w-100 text-black"}>
              withdraw all
            </SubmitButton>
          </div>
        </Card>
        <Card className={"tile"}>KEEP + TBTC</Card>
        <Card className={"tile"}>TBTC + ETH</Card>
      </CardContainer>
    </PageWrapper>
  )
}

export default LiquidityPage
