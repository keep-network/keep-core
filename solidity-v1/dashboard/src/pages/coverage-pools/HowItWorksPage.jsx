import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import { fetchTvlRequest, fetchAPYRequest } from "../../actions/coverage-pool"
import { MetricsSection } from "../../components/coverage-pools"
import * as Icons from "../../components/Icons"
import NavLink from "../../components/NavLink"
import List from "../../components/List"
import Timeline from "../../components/Timeline"
import Chip from "../../components/Chip"
import { colors } from "../../constants/colors"
import {
  Accordion,
  AccordionItem,
  AccordionItemPanel,
} from "react-accessible-accordion"
import AccordionDefaultItemHeading from "../../components/AccordionDefaultItemHeading"

const triggers = [
  {
    label:
      "When ETH-BTC price drops and there is not enough ETH as collateral for a deposit and in the liquidation state and no buyer is buying the ETH bonded by the stakers, the Coverage Pool will sell KEEP to buy BTC and cover the peg.",
  },
  {
    label:
      "When no valid redemption signature is provided in the required time frame and the deposit enters the liquidation state and no buyer is buying the ETH bonded by the stakers, the Coverage Pool will sell KEEP to buy BTC and cover the peg.",
  },
]
const about = [
  {
    label:
      "Coverage tokens (or covTOKENs) are ERC20 which reflect your share in the coverage pool. For each KEEP token deposited, you will get an amount of covKEEP that represents your share in the pool.",
  },
  {
    label:
      "To withdraw deposited tokens, you must have your covKEEP in the connected wallet. ",
  },
]

const HowItWorksPage = () => {
  const dispatch = useDispatch()

  useEffect(() => {
    dispatch(fetchTvlRequest())
    dispatch(fetchAPYRequest())
  }, [dispatch])

  const {
    totalValueLocked,
    totalValueLockedInUSD,
    isTotalValueLockedFetching,
    apy,
    isApyFetching,
    totalAllocatedRewards,
    totalCoverageClaimed,
  } = useSelector((state) => state.coveragePool)

  return (
    <>
      <MetricsSection
        tvl={totalValueLocked}
        tvlInUSD={totalValueLockedInUSD}
        rewardRate={apy}
        isRewardRateFetching={isApyFetching}
        totalAllocatedRewards={totalAllocatedRewards}
        isTotalAllocatedRewardsFetching={isTotalValueLockedFetching}
        lifetimeCovered={totalCoverageClaimed}
        isLifetimeCoveredFetching={isTotalValueLockedFetching}
      />
      <section className="cov-how-it-works__info-section">
        <Icons.CoveragePool className="info-section__icon" />
        <header className="info-section__header">
          <h2 className="text-grey-70">What’s a coverage pool?</h2>
          <h3 className="text-grey-50">
            A coverage pool functions as a form of insurance. It helps secure
            the network and is an opportunity to earn rewards.
          </h3>
        </header>
        <div className="info-section__cta">
          <NavLink
            to="/coverage-pools/deposit"
            className="btn btn-primary btn-md mb-1 w-100"
          >
            get started
          </NavLink>
          <a
            href="https://github.com/keep-network/coverage-pools/blob/main/docs/design.adoc"
            className="btn btn-tertiary btn-md w-100"
            rel="noopener noreferrer"
            target="_blank"
          >
            read the docs
          </a>
        </div>
      </section>

      <section className="cov-how-it-works__diagram-section">
        <Icons.CovPoolsHowItWorksDiagram />
      </section>

      <section className="coverage-pool-resources-grid">
        <section className="tile bg-violet-10">
          <h3 className="mb-2">Overview</h3>
          <Timeline>
            <Timeline.Element>
              <Timeline.Breakpoint>
                <Timeline.BreakpointDot>1</Timeline.BreakpointDot>
                <Timeline.BreakpointLine active />
              </Timeline.Breakpoint>
              <Timeline.Content>
                <Timeline.ElementDefaultCard>
                  <h4 className="text-violet-80">Deposit your tokens</h4>
                  <span className="text-grey-60">
                    There is no minimum KEEP amount for your deposit and no
                    minimum time lock.
                  </span>
                </Timeline.ElementDefaultCard>
              </Timeline.Content>
            </Timeline.Element>

            <Timeline.Element>
              <Timeline.Breakpoint>
                <Timeline.BreakpointDot>2</Timeline.BreakpointDot>
                <Timeline.BreakpointLine active />
              </Timeline.Breakpoint>
              <Timeline.Content>
                <Timeline.ElementDefaultCard>
                  <h4 className="text-violet-80">Withdraw deposit</h4>
                  <span className="text-grey-60">
                    <strong>Withdrawing is a two step action.</strong>
                    &nbsp;First, you initiate your withdrawal. After that there
                    is a 21 day cooldown period. During cooldown, your tokens
                    are still accumulating rewards but are also subject to risk
                    to cover for a hit. After 21 days, you can claim your token.
                  </span>
                </Timeline.ElementDefaultCard>
              </Timeline.Content>
            </Timeline.Element>
            <Timeline.Element>
              <Timeline.Breakpoint>
                <Timeline.BreakpointDot
                  lineBreaker
                  lineBreakerColor="violet-80"
                  style={{ backgroundColor: colors.brandViolet10 }}
                />
                <Timeline.BreakpointLine active />
              </Timeline.Breakpoint>
              <Timeline.Content>
                <Chip text="21 day cooldown" color="strong" size="big" />
              </Timeline.Content>
            </Timeline.Element>

            <Timeline.Element>
              <Timeline.Breakpoint>
                <Timeline.BreakpointDot>3</Timeline.BreakpointDot>
                <Timeline.BreakpointLine active />
              </Timeline.Breakpoint>
              <Timeline.Content>
                <Timeline.ElementDefaultCard>
                  <h4 className="text-violet-80">Claim tokens</h4>
                  <span className="text-grey-60">
                    <strong>
                      You have a 2 day claim window to claim your tokens and
                      rewards.
                    </strong>
                    &nbsp;Your deposit and rewards will be sent in one
                    transaction. If you do not claim your tokens within 2 days,
                    your tokens will return to the pool and you will have to
                    re-withdraw them.
                  </span>
                </Timeline.ElementDefaultCard>
              </Timeline.Content>
            </Timeline.Element>
          </Timeline>
        </section>

        <section className="tile cov-pools-accordion-section">
          <Accordion allowZeroExpanded>
            <AccordionItem>
              <AccordionDefaultItemHeading>
                Coverage pool
              </AccordionDefaultItemHeading>
              <AccordionItemPanel>
                <h5 className="text-violet-80 mb-1">what is it</h5>
                <div className="text-grey-60">
                  <p>
                    A pool of capital that serves as external aid to maintain
                    the 1:1 peg between tBTC and BTC deposits. Those who deposit
                    into the pool are effectively underwriting a rare event
                    where there&apos;s not enough money in the system to
                    purchase BTC.
                  </p>
                </div>
                <h5 className="text-violet-80 mb-1">how it works</h5>
                <div className="text-grey-60">
                  <p>
                    Coverage pools serve as a ‘buyer of last resort’. A buyer of
                    last resort is the buyer that will purchase enough tBTC when
                    no-one else will, to make a depositor whole in the event
                    liquidation if the stakers collateral is insufficient.
                  </p>
                  <p>
                    When coverage is demanded some part of the collateral pool
                    must be sold to obtain enough of the covered asset to
                    fulfill the claim. Liquidating the coverage pool fairly
                    means selling a basket of assets, in a fixed ratio, with
                    good price discovery. For this reason, collateral is
                    liquidated using a Dutch auction.
                  </p>
                </div>
              </AccordionItemPanel>
            </AccordionItem>

            <AccordionItem>
              <AccordionDefaultItemHeading>
                Being an underwriter
              </AccordionDefaultItemHeading>
              <AccordionItemPanel>
                <h5 className="text-violet-80 mb-1">
                  becoming and underwriter
                </h5>
                <div className="text-grey-60">
                  <p>
                    When you make a deposit into the pool, you become an
                    underwriter by securing the network. Because you provide
                    capital and put your funds to risk you earn rewards. The job
                    of an underwriter is quite passive, you don’t need to
                    monitor the network or run a node.
                  </p>
                </div>
                <h5 className="text-violet-80 mb-1">what to expect</h5>
                <div className="text-grey-60">
                  <p>
                    As an underwriter you are entering into a position relative
                    to tBTC, the asset that the coverage pool is backing, and
                    you are exposed to insurance events (liquidation).
                  </p>
                  <p>
                    Due to the sound math of the Coverage Pools a liquidation
                    event is extremely rare and considered a black swan event.
                  </p>
                </div>
              </AccordionItemPanel>
            </AccordionItem>

            <AccordionItem>
              <AccordionDefaultItemHeading>
                Insurance events
              </AccordionDefaultItemHeading>
              <AccordionItemPanel>
                <List items={triggers}>
                  <List.Title className="h5 text-violet-80">
                    Triggers
                  </List.Title>
                  <List.Content className="bullets bullets--violet-80 text-grey-60" />
                </List>
              </AccordionItemPanel>
            </AccordionItem>

            <AccordionItem>
              <AccordionDefaultItemHeading>
                Earning rewards
              </AccordionDefaultItemHeading>
              <AccordionItemPanel>
                <h5 className="text-violet-80 mb-1">how you earn</h5>
                <div className="text-grey-60">
                  <p>
                    There are weekly rewards emissions. The rewards emitted are
                    deposited in the Coverage Pool. Rewards are KEEP tokens.
                  </p>
                  <p>
                    They will be calculated based on a variable APY. You can
                    withdraw you rewards alongside with your deposit in a single
                    transaction. You can withdraw partial amounts of the deposit
                    and rewards.
                  </p>
                  <p>
                    As long as you keep your tokens in the pool your{" "}
                    <strong>rewards will be autocompounded</strong> and earn
                    rewards as well.
                  </p>
                </div>
              </AccordionItemPanel>
            </AccordionItem>

            <AccordionItem>
              <AccordionDefaultItemHeading>covKEEP</AccordionDefaultItemHeading>
              <AccordionItemPanel>
                <List items={about}>
                  <List.Title className="h5 text-violet-80">About</List.Title>
                  <List.Content className="bullets bullets--violet-80 text-grey-60" />
                </List>
              </AccordionItemPanel>
            </AccordionItem>
          </Accordion>
        </section>
      </section>
    </>
  )
}

HowItWorksPage.route = {
  title: "How it Works",
  path: "/coverage-pools/how-it-works",
  exact: true,
}

export default HowItWorksPage
