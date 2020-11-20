import React, { useEffect, useRef } from "react"
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ReferenceArea,
} from "recharts"
import { colors } from "../constants/colors"
import {
  displayAmount,
  getNumberWithMetricSuffix,
  fromTokenUnit,
} from "../utils/token.utils"
import { formatDate } from "../utils/general.utils"
import BeaconRewardsHelper from "../utils/beaconRewardsHelper"

const data = BeaconRewardsHelper.keepAllocationsInInterval.map(
  (amount, index) => {
    return {
      interval: index,
      amount,
    }
  }
)

// We want to display only first and last tick.
const xAxisTicks = [0, data.length - 1]

const styles = {
  dot: {
    r: 6,
    stroke: "white",
    fill: colors.primary,
    strokeWidth: 2,
  },
  chartWrapper: { marginTop: "1rem", position: "relative" },
}

const StakeDropChart = () => {
  const tooltipRef = useRef(null)

  useEffect(() => {
    const tooltipElement = tooltipRef.current

    return () => {
      if (tooltipElement) {
        tooltipElement.style.opacity = "0"
      }
    }
  })

  const mouseEnterOnDot = (e) => {
    const { interval, amount } = e.payload

    if (!tooltipRef.current) {
      return
    }

    const x = e.cx
    const tooltipContentEl = tooltipRef.current.firstChild
    // 8- it's a height of the triangle at the bottom of a tooltip.
    const y = e.cy - tooltipContentEl.getBoundingClientRect().height - 8
    tooltipRef.current.style.opacity = "1"
    tooltipRef.current.style.transform = `translate(${x - 60}px, ${y}px)`
    tooltipContentEl.childNodes[0].innerHTML = `Interval ${interval + 1}`
    tooltipContentEl.childNodes[1].innerHTML = formatDate(
      BeaconRewardsHelper.intervalStartOf(interval)
    )

    const formattedRewardAmount = displayAmount(fromTokenUnit(amount))
    tooltipContentEl.childNodes[2].innerHTML = tooltipContentEl.childNodes[2].innerHTML = `${formattedRewardAmount} KEEP`
  }

  const mouseLeaveOnDot = () => {
    if (tooltipRef.current) {
      tooltipRef.current.style.opacity = "0"
    }
  }

  return (
    <>
      <h4>Current Interval</h4>
      <div className="text-caption text-grey-60">Keep per Interval</div>
      <div style={styles.chartWrapper}>
        <LineChart width={500} height={300} data={data}>
          <CartesianGrid
            vertical={false}
            stroke={colors.grey20}
            strokeWidth={2}
          />
          <XAxis
            dataKey="interval"
            // To be able to draw area for a current interval.
            domain={["dataMin" - 0.5, "dataMax" + 0.5]}
            type="number"
            stroke={colors.grey60}
            strokeWidth={2}
            ticks={xAxisTicks}
            tickSize={8}
            tickLine={{ strokeLinejoin: "round" }}
            tickFormatter={(tick) => tick + 1}
          />
          <YAxis
            stroke={colors.grey60}
            strokeWidth={2}
            tickFormatter={(tick) =>
              getNumberWithMetricSuffix(tick).formattedValue
            }
            tickLine={false}
          />
          <Tooltip cursor={false} wrapperStyle={{ display: "none" }} />
          <ReferenceArea
            x1={BeaconRewardsHelper.currentInterval - 0.5}
            x2={BeaconRewardsHelper.currentInterval + 0.5}
            stroke={colors.green30}
            fill={colors.green30}
            fillOpacity={1}
          />
          <Line
            type="monotone"
            dataKey="amount"
            strokeWidth={2}
            stroke={colors.primary}
            dot={styles.dot}
            activeDot={{
              onMouseEnter: mouseEnterOnDot,
              onMouseLeave: mouseLeaveOnDot,
              r: 6,
            }}
          />
        </LineChart>
        <div className="stake-drop-chart__tooltip-wrapper" ref={tooltipRef}>
          <div className="stake-drop-chart__tooltip">
            <p className="stake-drop-chart__tooltip__title" />
            <p className="stake-drop-chart__tooltip__period" />
            <p className="stake-drop-chart__tooltip__amount" />
          </div>
        </div>
      </div>
    </>
  )
}

export default StakeDropChart
