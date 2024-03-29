import Banner from "../Banner"
import * as Icons from "../Icons"
import React from "react"
import Chip from "../Chip"
import { LINK } from "../../constants/constants"
import OnlyIf from "../OnlyIf"
import { KEEP, ThresholdToken } from "../../utils/token.utils"
import BigNumber from "bignumber.js"
import { Keep } from "../../contracts"
import List from "../List"

const AuthorizeStakesBanner = ({ stakesToAuthOrMoveToT = [] }) => {
  const numberOfStakesToAuthorize = stakesToAuthOrMoveToT.filter((stake) => {
    return !stake.contract.isAuthorized
  })

  const listItems = [
    {
      icon: Icons.Rewards,
      iconProps: {
        className: "authorize-stakes-banner__base-info-icon",
      },
      label: "Earn rewards on your KEEP stake with Threshold.",
      className: "mb-1",
    },
    {
      icon: Icons.Money,
      iconProps: { className: "authorize-stakes-banner__base-info-icon" },
      label: `Exchange rate is 1 KEEP =
                    ${ThresholdToken.displayAmountWithSymbol(
                      Keep.keepToTStaking.toThresholdTokenAmount(
                        KEEP.fromTokenUnit(1)
                      ),
                      3,
                      (amount) =>
                        new BigNumber(amount).toFormat(3, BigNumber.ROUND_DOWN)
                    )}
                    .`,
    },
  ]

  return (
    <Banner className="banner authorize-stakes-banner">
      <div className="banner__content-wrapper">
        <Banner.Icon icon={Icons.EarnThresholdTokens} />
        <div className="authorize-stakes-banner__content">
          <Banner.Title className="text-white banner__title--font-weight-600">
            {numberOfStakesToAuthorize.length > 0 ? (
              <h4 className="mb-1">
                Authorize your{" "}
                <OnlyIf condition={stakesToAuthOrMoveToT.length > 0}>
                  {stakesToAuthOrMoveToT.length}
                </OnlyIf>{" "}
                stake
                <OnlyIf condition={stakesToAuthOrMoveToT.length !== 1}>
                  {"s"}
                </OnlyIf>{" "}
                below to get started staking on Threshold.
              </h4>
            ) : (
              <h4 className="mb-1">
                Manage your Threshold stake in the Threshold Staking table
                below.
              </h4>
            )}
          </Banner.Title>
          <Banner.Description>
            <div className={"flex row space-between"}>
              <div className={"authorize-stakes-banner__base-info mr-1"}>
                <List items={listItems}>
                  <List.Content />
                </List>
              </div>
              <div
                className={"flex row authorize-stakes-banner__pre-node-info"}
              >
                <Chip
                  size={"small"}
                  text={"NOTE"}
                  className={"authorize-stakes-banner__pre-node-info-chip"}
                />
                <div>
                  You will need to run a PRE node to qualify for rewards.{" "}
                  <a
                    href={LINK.setUpPRE}
                    rel="noopener noreferrer"
                    target="_blank"
                    className={"authorize-stakes-banner__pre-node-info-link"}
                  >
                    Set up PRE
                  </a>
                </div>
              </div>
            </div>
          </Banner.Description>
        </div>
      </div>
    </Banner>
  )
}

export default AuthorizeStakesBanner
