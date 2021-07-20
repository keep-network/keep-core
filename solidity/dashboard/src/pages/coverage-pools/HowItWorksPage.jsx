import React from "react"
import * as Icons from "../../components/Icons"
import NavLink from "../../components/NavLink"
import List from "../../components/List"
// import Timeline from "../../components/Timeline"

const triggers = [
  {
    label:
      "When ETH-BTC price drops and undercollateralized deposit enters the liquidation state and no buyer will take the signer bonds auction.",
  },
  {
    label:
      "When no valid redemption signature is provided in the required time frame and the deposit enters the liquidation state and no buyer will take the signer bonds auction.",
  },
]

const expectations = [
  {
    label:
      "You will deposit funds in the coverage pool which will secure the network and be used as coverage. Coverage is akin to filing a claim in traditional insurance and processing your own claim. The triggers of coverage are listed above.",
  },
  {
    label:
      "By providing coverage funds you will get to earn rewards which are allocated weekly and calculated based on your share in the pool.",
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
  return (
    <>
      <section className="tile flex row full-center">
        <Icons.CoveragePool style={{ width: 200, height: 150 }} />
        <header className="ml-1 mr-1">
          <h2 className="text-grey-70">What’s a coverage pool?</h2>
          <h3 className="text-grey-50">
            A coverage pool functions as a form of insurance. It helps secure
            the network and is an opportunity to earn rewards.
          </h3>
        </header>
        <div>
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
      <section className="tile">
        <h3 className="mb-1">Rewards pool</h3>
        <h5 className="text-violet-80">how it works</h5>
        <div className="text-grey-60">
          <p>
            The <strong>rewards are in KEEP which are allocated weekly.</strong>
          </p>
          <p>
            They will be calculated based on a variable APY. You can withdraw
            you rewards alongside with your deposit in a single transaction. You
            can withdraw partial amounts of the deposit and rewards.
          </p>
          <p>
            As long as you keep your tokens in the pool your&nbsp;
            <strong>rewards will be autocompounded</strong> and earn rewards as
            well.
          </p>
        </div>
      </section>
      <section className="tile">
        <h3 className="mb-1">Covering a hit</h3>
        <List items={triggers}>
          <List.Title className="h5 text-violet-80">Triggers</List.Title>
          <List.Content className="bullets bullets--violet-80 text-grey-60" />
        </List>

        <List items={expectations} className="mt-1">
          <List.Title className="h5 text-violet-80">
            What to expect as an underwriter
          </List.Title>
          <List.Content className="bullets bullets--violet-80 text-grey-60" />
        </List>
      </section>

      <section className="tile">
        <h3 className="mb-1">covKEEP</h3>
        <List items={about}>
          <List.Title className="h5 text-violet-80">About</List.Title>
          <List.Content className="bullets bullets--violet-80 text-grey-60" />
        </List>
      </section>

      {/* <section className="tile">
        <Timeline>
          <Timeline.Element>
            <Timeline.ElementBreakpoint active>1</Timeline.ElementBreakpoint>
            <h4 className="text-violet-80">Deposit your tokens</h4>
            <span className="text-grey-60">
              There is no minimum KEEP amount for your deposit and no minimum
              time lock.
            </span>
          </Timeline.Element>

          <Timeline.Element>
            <Timeline.ElementBreakpoint active>2</Timeline.ElementBreakpoint>
            <h4 className="text-violet-80">Withdraw deposit </h4>
            <span className="text-grey-60">
              <strong>Withdrawing is a two step action.</strong>&nbsp;From the
              initiation of your withdrawal to claiming your tokens there is a
              21 day cooldown period in which your tokens are still accumulating
              rewards but are also subject to risk to cover for a hit.
            </span>
          </Timeline.Element>

          <Timeline.Element>
            <Timeline.ElementBreakpoint active>3</Timeline.ElementBreakpoint>
            <h4 className="text-violet-80">Claim tokens</h4>
            <span className="text-grey-60">
              <strong>
                You have a 2 day grace period to claim your tokens and rewards.
              </strong>
              Your deposit and rewards will be sent in one transaction. If you
              do not claim your tokens within 2 days, your tokens will return to
              the pool and you will have to re-withdraw them.
            </span>
          </Timeline.Element>
        </Timeline>
      </section> */}
    </>
  )
}

HowItWorksPage.route = {
  title: "How it Works",
  path: "/coverage-pools/how-it-works",
  exact: true,
}

export default HowItWorksPage
