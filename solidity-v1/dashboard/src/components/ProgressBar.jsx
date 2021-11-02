import React, { useMemo, useContext } from "react"
import BigNumber from "bignumber.js"
import CircularProgressBar from "./CircularProgressBar"
import { percentageOf, sub, lt } from "../utils/arithmetics.utils"
import { formatValue } from "../utils/general.utils"
import OnlyIf from "./OnlyIf"

const defaultValue = 0
const totalDefaultValue = 1
const defaultSecondaryValue = 0

const calculateWidth = (value, total, secondaryValue = null) => {
  const valueInBn = new BigNumber(value || defaultValue)
  const totalInBn = new BigNumber(total || totalDefaultValue)
  const secondaryValueInBn = new BigNumber(
    secondaryValue || defaultSecondaryValue
  )

  const finalValueInBn = valueInBn.plus(secondaryValueInBn)

  return finalValueInBn.multipliedBy(100).div(totalInBn).toFixed(2).toString()
}

const ProgressBarContext = React.createContext({ value: 0, total: 0 })

const useProgressBarContext = () => {
  const context = useContext(ProgressBarContext)

  if (!context) {
    throw new Error("ProgressBarContext used outside of ProgressBar component")
  }

  return context
}

const ProgressBar = ({
  value,
  total,
  color,
  bgColor,
  secondaryValue = null,
  secondaryColor = null,
  children,
}) => {
  return (
    <ProgressBarContext.Provider
      value={{ value, total, color, bgColor, secondaryValue, secondaryColor }}
    >
      {children}
    </ProgressBarContext.Provider>
  )
}

const ProgressBarInline = ({ height = 10, className = "" }) => {
  const { value, total, color, bgColor, secondaryValue, secondaryColor } =
    useProgressBarContext()

  const isDoubleColored = !!secondaryColor && !!secondaryValue

  const barWidth = useMemo(
    () => calculateWidth(value, total, secondaryValue),
    [value, total, secondaryValue]
  )
  const secondaryBarWidth = useMemo(
    () => calculateWidth(secondaryValue, total),
    [secondaryValue, total]
  )

  return (
    <div
      className={`progress-bar-wrapper ${className}`}
      style={{
        height: `${height}px`,
        backgroundColor: bgColor,
      }}
    >
      <div
        className={`progress-bar ${
          isDoubleColored ? "progress-bar--main-color" : ""
        }`}
        style={{
          width: `${barWidth}%`,
          backgroundColor: color,
        }}
      />
      <OnlyIf condition={isDoubleColored}>
        <div
          className="progress-bar progress-bar--secondary-color"
          style={{
            width: `${secondaryBarWidth}%`,
            backgroundColor: secondaryColor,
          }}
        />
      </OnlyIf>
    </div>
  )
}

const defaultDisplayLegendValuFn = (value) => value.toString()
export const ProgressBarLegendContext = React.createContext({
  renderValuePattern: defaultDisplayLegendValuFn,
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
  renderValuePattern = defaultDisplayLegendValuFn,
}) => {
  const { value, total, color, bgColor } = useProgressBarContext()
  const leftValue = useMemo(() => {
    const left = sub(total, value).toString()
    return lt(left, 0) ? 0 : left
  }, [value, total])

  return (
    <ProgressBarLegendContext.Provider value={{ renderValuePattern }}>
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
  const { renderValuePattern } = useProgressBarLegendContext()
  const renderedValue = React.isValidElement(renderValuePattern)
    ? React.cloneElement(renderValuePattern, { amount: value })
    : renderValuePattern(value)

  return (
    <div className="progress-bar-legend__item">
      <div className="legend__item__dot" style={{ backgroundColor: color }} />
      <span className="legend__item__value">
        {renderedValue}
        &nbsp;
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
    () => formatValue(percentageOf(value, total)),
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

// const ProgressBarLegend2 = () => {
//   return (
//     <div className="progress-bar__legend">
//       <ProgressBarLegendItem
//         value={leftValue}
//         label={leftValueLabel}
//         color={bgColor}
//       />
//       <ProgressBarLegendItem value={value} label={valueLabel} color={color} />
//     </div>
//   )
// }

ProgressBar.Inline = ProgressBarInline
ProgressBar.Legend = ProgressBarLegend
ProgressBar.LegendItem = ProgressBarLegendItem
ProgressBar.Circular = ProgressBarCircular
ProgressBar.PercentageLabel = PercentageLabel

export default ProgressBar
