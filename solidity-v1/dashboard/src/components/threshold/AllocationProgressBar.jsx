import React from "react"
import ProgressBar from "../ProgressBar"
import { colors } from "../../constants/colors"

const AllocationProgressBar = ({
  title,
  currentValue = 0,
  totalValue = 100,
  className = "",
}) => {
  if (totalValue <= 0) {
    currentValue = 0
    totalValue = 100
  }
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
          >
            <ProgressBar.Inline
              className="allocation-progress-bar__progress-bar-wrapper"
              height={20}
            />
          </ProgressBar>
          <span className="text-grey-70 ml-2 allocation-progress-bar__allocation-percentage-value">
            {Math.round((currentValue / totalValue) * 100)}%
          </span>
        </div>
      </div>
    </div>
  )
}

export default AllocationProgressBar
