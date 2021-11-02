import React, { useMemo } from "react"
import ProgressBar from "../ProgressBar"
import { colors } from "../../constants/colors"
import OnlyIf from "../OnlyIf"
import { add, gt } from "../../utils/arithmetics.utils"

const AllocationProgressBar = ({
  title,
  currentValue = 0,
  totalValue = 100,
  className = "",
  secondaryValue = null,
  withLegend = false,
  currentValueLegendLabel = "",
  secondaryValueLegendLabel = "",
}) => {
  const percentageValue = useMemo(() => {
    if (gt(totalValue, 0)) {
      return Math.round(
        (add(currentValue, secondaryValue | 0) / totalValue) * 100
      )
    } else {
      return 0
    }
  }, [currentValue, secondaryValue, totalValue])

  return (
    <div className={`allocation-progress-bar ${className}`}>
      <div className="allocation-progress-bar__token-allocation">
        <h5 className="text-gray-60">{title}</h5>
        <div className="allocation-progress-bar__progress-bar-container">
          <ProgressBar
            value={currentValue}
            total={totalValue}
            color={colors.secondary}
            bgColor={colors.grey20}
            secondaryValue={secondaryValue}
            secondaryColor={colors.yellowSecondary}
          >
            <ProgressBar.Inline
              className="allocation-progress-bar__progress-bar-wrapper"
              height={20}
            />
          </ProgressBar>
          <span className="text-grey-70 ml-1 allocation-progress-bar__allocation-percentage-value">
            {/** TODO: 2 decimal places, maybe even print it as >99 % and <1%
              // when there is small difference betweent currentValue and total
              // Value */}
            {percentageValue}%
          </span>
        </div>
        <div className={"allocation-progress-bar__legend"}>
          <OnlyIf condition={withLegend}>
            <div
              className="allocation-progress-bar__legend-dot"
              style={{ backgroundColor: colors.yellowSecondary }}
            />
            <span>{secondaryValueLegendLabel}</span>
            <div
              className="allocation-progress-bar__legend-dot ml-1"
              style={{ backgroundColor: colors.secondary }}
            />
            <span>{currentValueLegendLabel}</span>
          </OnlyIf>
        </div>
      </div>
    </div>
  )
}

export default AllocationProgressBar
