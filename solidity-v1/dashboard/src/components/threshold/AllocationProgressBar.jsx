import React, { useMemo } from "react"
import ProgressBar, { ProgressBarLegendContext } from "../ProgressBar"
import { colors } from "../../constants/colors"
import OnlyIf from "../OnlyIf"
import { add, gt } from "../../utils/arithmetics.utils"
import BigNumber from "bignumber.js"
import TokenAmount from "../TokenAmount"

const AllocationProgressBar = ({
  title,
  currentValue = 0,
  totalValue = 100,
  className = "",
  secondaryValue = null,
  withLegend = false,
  currentValueLabel = "",
  secondaryValueLabel = "",
  isDataFetching = false,
}) => {
  const displayPercentageValue = useMemo(() => {
    if (isDataFetching) {
      return "--"
    }

    if (gt(totalValue, 0)) {
      const currentValueBN = new BigNumber(currentValue)
      const secondaryValueBN = secondaryValue
        ? new BigNumber(secondaryValue)
        : 0
      const actualProgressBarValueBN = currentValueBN.plus(secondaryValueBN)
      const totalValueBN = new BigNumber(totalValue)
      const percentageValue = Math.round(
        actualProgressBarValueBN.div(totalValueBN).multipliedBy(100).toNumber()
      )

      if (
        percentageValue === 100 &&
        !actualProgressBarValueBN.isEqualTo(totalValueBN)
      ) {
        return ">99%"
      } else if (
        percentageValue === 0 &&
        !actualProgressBarValueBN.isEqualTo("0")
      ) {
        return "<1%"
      } else {
        return `${percentageValue}%`
      }
    } else {
      return "0%"
    }
  }, [currentValue, secondaryValue, totalValue, isDataFetching])

  const renderProgressBarTooltipTooltip = (value) => {
    return (
      <TokenAmount
        amount={value}
        symbolClassName="upgrade-token-tile-row__token-amount-symbol"
        amountClassName="upgrade-token-tile-row__token-amount-amount"
      />
    )
  }

  return (
    <div className={`allocation-progress-bar ${className}`}>
      <div>
        <h5 className="text-grey-60">{title}</h5>
        <div className="allocation-progress-bar__progress-bar-container">
          <ProgressBar
            total={isDataFetching ? 100 : totalValue}
            bgColor={colors.grey20}
          >
            <ProgressBar.Inline
              className="allocation-progress-bar__progress-bar-wrapper"
              height={20}
            >
              <ProgressBar.InlineItem
                label={currentValueLabel}
                value={
                  isDataFetching ? 0 : add(currentValue, secondaryValue || 0)
                }
                color={colors.secondary}
                withTooltip
                renderTooltip={renderProgressBarTooltipTooltip(currentValue)}
              />
              <OnlyIf condition={!!secondaryValue}>
                <ProgressBar.InlineItem
                  label={secondaryValueLabel}
                  value={isDataFetching ? 0 : secondaryValue}
                  color={colors.yellowSecondary}
                  withTooltip
                  renderTooltip={renderProgressBarTooltipTooltip(
                    secondaryValue
                  )}
                />
              </OnlyIf>
            </ProgressBar.Inline>
            <OnlyIf condition={withLegend}>
              <div className="allocation-progress-bar__progress-bar-legend">
                <ProgressBarLegendContext.Provider
                  value={{ renderValuePattern: () => <></> }}
                >
                  <ProgressBar.LegendItem
                    value={secondaryValue?.toString()}
                    label={secondaryValueLabel}
                    color={colors.yellowSecondary}
                    className={
                      "allocation-progress-bar__progress-bar-legend-item"
                    }
                  />
                  <ProgressBar.LegendItem
                    value={currentValue?.toString()}
                    label={currentValueLabel}
                    color={colors.secondary}
                    className={
                      "allocation-progress-bar__progress-bar-legend-item"
                    }
                  />
                </ProgressBarLegendContext.Provider>
              </div>
            </OnlyIf>
          </ProgressBar>
          <span className="text-grey-70 ml-1 allocation-progress-bar__percentage-value">
            {displayPercentageValue}
          </span>
        </div>
      </div>
    </div>
  )
}

export default AllocationProgressBar
