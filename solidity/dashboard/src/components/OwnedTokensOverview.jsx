import React, { useMemo } from "react"
import Tile from "./Tile"
import { CircularProgressBars } from "./CircularProgressBar"
import { add } from "../utils/arithmetics.utils"
import {
  displayAmountWithMetricSuffix,
  displayAmount,
} from "../utils/token.utils"
import { colors } from "../constants/colors"
import TokenAmount from "./TokenAmount"

const OwnedTokensOverview = ({ keepBalance, stakedBalance }) => {
  const total = useMemo(() => {
    return add(keepBalance, stakedBalance)
  }, [keepBalance, stakedBalance])

  return (
    <Tile id="tokens-overview" title="Owned Tokens">
      <TokenAmount amount={keepBalance} />
      <hr />
      <div className="flex">
        <div className="flex-1 self-center">
          <CircularProgressBars
            total={total}
            items={[
              {
                value: stakedBalance,
                color: colors.grey70,
                backgroundStroke: colors.grey10,
                label: "Staked",
              },
            ]}
            withLegend
          />
        </div>
        <div className="ml-2 mt-1 self-start flex-1">
          <h5 className="text-grey-70">staked</h5>
          <h4 className="text-grey-70">{displayAmount(stakedBalance)}</h4>
          <div className="text-smaller text-grey-40">
            of {displayAmountWithMetricSuffix(total)} Total
          </div>
        </div>
      </div>
    </Tile>
  )
}

export default React.memo(OwnedTokensOverview)
