import React from "react"
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ReferenceArea,
} from "recharts"
import { toBN } from "web3-utils"
import { colors } from "../constants/colors"
import { displayAmountWithMetricSuffix } from "../utils/token.utils"

const keepAllocationsInInterval = [
  /* eslint-disable*/
    792000,     1520640,    1748736,    1888635,    2077498,    1765874,
    1500993,    1275844,    1084467,    921797,     783528,     665998,
    566099,     481184,     409006,     347655,     295507,     251181,
    213504,     181478,     154257,     131118,     111450,     94733
  /* eslint-enable*/
]

const data = keepAllocationsInInterval.map((amount, index) => {
  const amountInTokenUnit = toBN(amount)
    .mul(toBN(10).pow(toBN(18)))
    .toString()
  const yAxisKey = displayAmountWithMetricSuffix(amountInTokenUnit)
  return { yAxisKey, interval: index + 1, amount, amountInTokenUnit }
})

// We want to display only first and last tick.
const xAxisTicks = [1, data.length]

const styles = {
  dot: {
    r: 6,
    stroke: "white",
    fill: colors.primary,
    strokeWidth: 2,
  },
}

const StakeDropChart = () => {
  return (
    <>
      <h4>Current Interval</h4>
      <span className="text-caption text-grey-60">Keep per Interval</span>
      <LineChart width={500} height={300} data={data}>
        <CartesianGrid
          vertical={false}
          stroke={colors.grey20}
          strokeWidth={2}
        />
        <XAxis
          dataKey="interval"
          interval="preserveStartEnd"
          stroke={colors.grey60}
          strokeWidth={2}
          ticks={xAxisTicks}
          tickSize={8}
          tickLine={{ strokeLinejoin: "round" }}
        />
        <YAxis tickLine={false} stroke={colors.grey60} strokeWidth={2} />
        <Tooltip cursor={false} />
        <Line
          type="monotone"
          dataKey="amount"
          strokeWidth={2}
          stroke={colors.primary}
          dot={styles.dot}
          activeDot={false}
        />
        <ReferenceArea
          x1={2.5}
          x2={3.5}
          y1={0}
          y2={2200000}
          stroke="red"
          strokeOpacity={0.3}
        />
      </LineChart>
    </>
  )
}

export default StakeDropChart
