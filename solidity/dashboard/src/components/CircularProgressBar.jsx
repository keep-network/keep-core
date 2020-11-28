import React, { useMemo } from "react"
import { colors } from "../constants/colors"
import * as Icons from "./Icons"
import BigNumber from "bignumber.js"
import { renderProgressBarLegendItem } from "./ProgressBar"

const countCircumference = (radius) => {
  return new BigNumber(2 * Math.PI * radius)
}

const countProgressValue = (value, total, radius) => {
  const valueInBn = new BigNumber(value || 0)
  const totalInBn = new BigNumber(total || 1)
  const progress = valueInBn.div(totalInBn)

  return countCircumference(radius)
    .multipliedBy(new BigNumber(1).minus(progress))
    .toFixed(2)
    .toString()
}

const barWidth = 10

const CircularProgressBar = ({
  radius,
  value,
  backgroundStroke,
  color,
  withBackgroundStroke,
  total,
}) => {
  return (
    <svg className="circular-progress-bar" width={120} height={120}>
      {withBackgroundStroke && (
        <circle
          fill="none"
          className="background"
          cx={60}
          cy={60}
          r={radius - barWidth / 2}
          strokeWidth={barWidth}
          stroke={backgroundStroke}
        />
      )}
      <circle
        fill="none"
        strokeDashoffset={countProgressValue(value, total, radius - 5)}
        strokeDasharray={countCircumference(radius - 5)}
        className="value"
        cx={60}
        cy={60}
        r={radius - barWidth / 2}
        strokeWidth={barWidth}
        stroke={color}
        strokeLinecap="round"
      />
    </svg>
  )
}

CircularProgressBar.defaultProps = {
  radius: 60,
  value: 0,
  backgroundStroke: colors.grey,
  color: colors.primary,
  withBackgroundStroke: true,
}

export const CircularProgressBars = React.memo(
  ({ withLegend, total, items }) => {
    const bars = useMemo(() => {
      return items.map((item, index) => (
        <CircularProgressBar key={index} {...item} total={total} />
      ))
    }, [total, items])

    return (
      <>
        <svg
          className="wrapper-circular-progress-bar"
          width={120}
          height={120}
          viewBox="0 0 120 120"
        >
          {bars}
          <g className="keep-circle">
            <Icons.KeepCircle />
          </g>
        </svg>
        <div className="mb-1">
          {withLegend && items.map(renderProgressBarLegendItem)}
        </div>
      </>
    )
  }
)

export const CircularProgressBarEth = React.memo(
  ({ withLegend, total, items }) => {
    const bars = useMemo(() => {
      return items.map((item, index) => (
        <CircularProgressBar key={index} {...item} total={total} />
      ))
    }, [total, items])

    return (
      <>
        <svg
          className="wrapper-circular-progress-bar"
          width={120}
          height={120}
          viewBox="0 0 120 120"
        >
          {bars}
          <g className="keep-circle">
            <Icons.ETH width={58}
                    height={58} />
          </g>
        </svg>
        <div className="mb-1">
          {withLegend && items.map(renderProgressBarLegendItem)}
        </div>
      </>
    )
  }
)

export default React.memo(CircularProgressBar)
