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
import moment from "moment"
import { colors } from "../constants/colors"
import {
  displayAmount,
  getNumberWithMetricSuffix,
  fromTokenUnit,
} from "../utils/token.utils"
import { formatDate } from "../utils/general.utils"

const keepAllocationsInInterval = [
  /* eslint-disable*/
    792000,     1520640,    1748736,    1888635,    2077498,    1765874,
    1500993,    1275844,    1084467,    921797,     783528,     665998,
    566099,     481184,     409006,     347655,     295507,     251181,
    213504,     181478,     154257,     131118,     111450,     94733
  /* eslint-enable*/
]

// Beacon genesis date, 2020-09-24, is the first interval start.
// https://etherscan.io/tx/0xe2e8ab5631473a3d7d8122ce4853c38f5cc7d3dcbfab3607f6b27a7ef3b86da2
const beaconFirstIntervalStart = 1600905600

// Each interval is 30 days long.
const beaconTermLength = moment.duration(30, "days").asSeconds()

const data = keepAllocationsInInterval.map((amount, index) => {
  return {
    interval: index,
    amount,
  }
})

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

const intervalStartOf = (interval) => {
  return moment
    .unix(beaconFirstIntervalStart)
    .add(interval * beaconTermLength, "seconds")
}

const currentInterval = Math.floor(
  (moment().unix() - beaconFirstIntervalStart) / beaconTermLength
)

const StakeDropChart = () => {
  const tooltipRef = useRef(null)

  useEffect(() => {
    const tooltipElement = tooltipRef.current

    return () => {
      if (tooltipElement) {
        tooltipElement.style.display = "none"
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
      intervalStartOf(interval)
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
            x1={currentInterval - 0.5}
            x2={currentInterval + 0.5}
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
