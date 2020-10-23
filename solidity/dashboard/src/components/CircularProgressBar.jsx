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

const CircularProgressBar = ({
  radius,
  value,
  backgroundStroke,
  color,
  withBackgroundStroke,
  total,
  barWidth = 10,
}) => {
  const normalizedRadius = useMemo(() => radius - barWidth / 2, [
    radius,
    barWidth,
  ])

  const circumference = useMemo(() => {
    return countCircumference(normalizedRadius)
  }, [normalizedRadius])

  const progress = useMemo(() => {
    return countProgressValue(value, total, normalizedRadius)
  }, [value, total, normalizedRadius])

  return (
    <svg
      className="circular-progress-bar"
      width={radius * 2}
      height={radius * 2}
    >
      {withBackgroundStroke && (
        <circle
          fill="none"
          className="background"
          cx={radius}
          cy={radius}
          r={normalizedRadius}
          strokeWidth={barWidth}
          stroke={backgroundStroke}
        />
      )}
      <circle
        fill="none"
        strokeDashoffset={progress}
        strokeDasharray={`${circumference} ${circumference}`}
        className="value"
        cx={radius}
        cy={radius}
        r={normalizedRadius}
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

export default React.memo(CircularProgressBar)
