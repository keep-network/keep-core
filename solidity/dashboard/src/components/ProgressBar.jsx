import React, { useMemo, useContext } from "react"
import BigNumber from "bignumber.js"
import CircularProgressBar from "./CircularProgressBar"
import { percentageOf, sub } from "../utils/arithmetics.utils"

const defaultValue = 0
const totalDefaultValue = 1

const calculateWidth = (value, total) => {
  const valueInBn = new BigNumber(value || defaultValue)
  const totalInBn = new BigNumber(total || totalDefaultValue)

  return valueInBn.multipliedBy(100).div(totalInBn).toFixed(2).toString()
}

const ProgressBarContext = React.createContext({ value: 0, total: 0 })

const useProgressBarContext = () => {
  const context = useContext(ProgressBarContext)

  if (!context) {
    throw new Error("ProgressBarContext used outside of ProgressBar component")
  }

  return context
}

const ProgressBar = ({ value, total, color, bgColor, children }) => {
  return (
    <ProgressBarContext.Provider value={{ value, total, color, bgColor }}>
      {children}
    </ProgressBarContext.Provider>
  )
}

const ProgressBarInline = ({ height = 10, className = "" }) => {
  const { value, total, color, bgColor } = useProgressBarContext()

  const barWidth = useMemo(() => calculateWidth(value, total), [value, total])

  return (
    <div
      className={`progress-bar-wrapper ${className}`}
      style={{
        height: `${height}px`,
        backgroundColor: bgColor,
      }}
    >
      <div
        className="progress-bar"
        style={{
          width: `${barWidth}%`,
          backgroundColor: color,
        }}
      />
    </div>
  )
}

const defaultDisplayLegendValuFn = (value) => value.toString()
export const ProgressBarLegendContext = React.createContext({
  displayLegendValuFn: defaultDisplayLegendValuFn,
})

const useProgressBarLegendContext = () => {
  const context = useContext(ProgressBarLegendContext)

  if (!context) {
    throw new Error(
      "ProgressBarLegendContext used outside of ProgressBar component"
    )
  }

  return context
}

const ProgressBarLegend = ({
  valueLabel,
  leftValueLabel,
  displayLegendValuFn = defaultDisplayLegendValuFn,
}) => {
  const { value, total, color, bgColor } = useProgressBarContext()
  const leftValue = useMemo(() => sub(total, value).toString(), [value, total])

  return (
    <ProgressBarLegendContext.Provider value={{ displayLegendValuFn }}>
      <div className="progress-bar__legend">
        <ProgressBarLegendItem
          value={leftValue}
          label={leftValueLabel}
          color={bgColor}
        />
        <ProgressBarLegendItem value={value} label={valueLabel} color={color} />
      </div>
    </ProgressBarLegendContext.Provider>
  )
}

export const ProgressBarLegendItem = React.memo(({ value, label, color }) => {
  const { displayLegendValuFn } = useProgressBarLegendContext()

  return (
    <div className="progress-bar-legend__item">
      <div className="legend__item__dot" style={{ backgroundColor: color }} />
      <span className="legend__item__value">
        {displayLegendValuFn(value)}&nbsp;
      </span>
      <span className="legend__item__label">{label}</span>
    </div>
  )
})

const ProgressBarCircular = (props) => {
  const { value, total, color, bgColor } = useProgressBarContext()

  return (
    <CircularProgressBar
      {...props}
      color={color}
      backgroundStroke={bgColor}
      total={total}
      value={value}
    />
  )
}

export const renderProgressBarLegendItem = (item, index) => (
  <ProgressBarLegendItem key={index} {...item} />
)

const PercentageLabel = ({ text, className = "" }) => {
  const { value, total } = useProgressBarContext()
  const percentageOfValue = useMemo(
    () => percentageOf(value, total).toString(),
    [value, total]
  )

  return (
    <span className={`progress-bar__percentage-label ${className}`}>
      <span className="progress-bar__percentage-label__value">
        {percentageOfValue}%
      </span>
      <span className="progress-bar__percentage-label__text">{text}</span>
    </span>
  )
}

ProgressBar.Inline = ProgressBarInline
ProgressBar.Legend = ProgressBarLegend
ProgressBar.Circular = ProgressBarCircular
ProgressBar.PercentageLabel = PercentageLabel

export default ProgressBar
