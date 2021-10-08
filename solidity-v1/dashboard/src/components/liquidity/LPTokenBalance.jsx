import React, { useMemo } from "react"
import * as Icons from "../Icons"
import CountUp from "react-countup"
import BigNumber from "bignumber.js"
import { formatValue } from "../../utils/general.utils"
import Tooltip from "../Tooltip"
import { LPToken } from "../../utils/token.utils"

const LPTokenBalance = ({ lpTokens, lpTokenBalance }) => {
  const formattedLPTokenBalance = useMemo(() => {
    const token0BN = new BigNumber(lpTokenBalance.token0)
    const token1BN = new BigNumber(lpTokenBalance.token1)

    const token0 = formatValue(token0BN, 0)
    const token1 = formatValue(token1BN, 0)

    return {
      token0: LPToken.toTokenUnit(token0),
      token1: LPToken.toTokenUnit(token1),
    }
  }, [lpTokenBalance.token0, lpTokenBalance.token1])

  if (lpTokens && lpTokens.length === 0) return null
  return (
    <div className={"lp-balance"}>
      <h4 className={"text-grey-70 mb-1"}>
        Your LP Token Balance &nbsp;
        <Tooltip
          simple
          delay={0}
          triggerComponent={Icons.MoreInfo}
          className={"lp-balance__tooltip"}
        >
          Estimated value of the LP tokens deposited on the LPRewards contract
        </Tooltip>
      </h4>
      {lpTokens.map((lpToken, i) => {
        const IconComponent = Icons[lpToken.iconName]
        return (
          <div key={`lpToken-${i}`}>
            <div className={"lp-balance__value-container text-grey-70"}>
              <h3 className={"lp-balance__value-label"}>
                <IconComponent />
                <span>{lpToken.tokenName}</span>
              </h3>
              <h3>
                <CountUp
                  end={Object.values(formattedLPTokenBalance)[i].toNumber()}
                  separator={","}
                  preserveValue
                />
              </h3>
            </div>
            {i !== lpTokens.length - 1 && (
              <div
                className={"lp-balance__plus-separator-container text-grey-70"}
              >
                <span className={"lp-balance__plus-separator"}>+</span>
              </div>
            )}
          </div>
        )
      })}
    </div>
  )
}

export default LPTokenBalance
