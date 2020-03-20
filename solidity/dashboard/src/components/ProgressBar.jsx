import React, { useMemo } from 'react'
import { displayAmount } from '../utils/general.utils'
import BigNumber from 'bignumber.js'

const defaultValue = 0
const totalDefaultValue = 1

const calculateWidth = (value, total) => {
  const valueInBn = new BigNumber(value || defaultValue)
  const totalInBn = new BigNumber(total || totalDefaultValue)

  return valueInBn.multipliedBy(100).div(totalInBn).toFixed(2).toString()
}

const ProgressBar = ({ total, items, height, withLegend }) => {
  const bars = useMemo(() => {
    return items
      .map((item) => ({ ...item, width: calculateWidth(item.value, total) }))
      .sort((a, b) => b.width - a.width)
      .map((item, index) =>
        <ProgressBarItem
          key={index}
          {...item}
          index={index}
          wrapperHeight={height}
        />)
  }, [total, items])

  return (
    <React.Fragment>
      <div className="progress-bar-wrapper" style={{ height: `${height}px` }}>
        {bars}
      </div>
      {withLegend && items.map(renderProgressBarLegendItem)}
    </React.Fragment>
  )
}

export const renderProgressBarLegendItem = (item, index) => <ProgressBarLegendItem key={index} {...item} />

export const ProgressBarLegendItem = React.memo(({ value, label, color }) => {
  return (
    <div className="flex row center">
      <div className="dot" style={{ backgroundColor: color }}/>
      <span className="text-small">{displayAmount(value)}&nbsp;KEEP&nbsp;</span><span className="text-smaller text-grey-40">{label}</span>
    </div>
  )
})

const ProgressBarItem = React.memo(({ width, color, wrapperHeight, index }) => (
  <div
    className="progress-bar"
    style={{
      width: `${width}%`,
      zIndex: index + 1,
      backgroundColor: color,
      height: `${index === 0 ? wrapperHeight : wrapperHeight - index - 1 }px`,
    }}
  />
))

ProgressBar.defaultProps = {
  height: '10',
}

export default React.memo(ProgressBar)
