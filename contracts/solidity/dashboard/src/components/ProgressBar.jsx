import React from 'react'
import { displayAmount } from '../utils'
import web3Utils from 'web3-utils'

const calculateWidth = (value, total) => {
  const valueInBN = web3Utils.toBN(value || 0)
  const totalInBN = web3Utils.toBN(total || 1)

  return valueInBN.mul(web3Utils.toBN(100)).div(totalInBN).toString()
}

const ProgressBar = ({ total, items, height, withLegend }) => {
  const bars = items
    .map((item) => ({ ...item, width: calculateWidth(item.value, total) }))
    .sort((a, b) => b.width - a.width)

  const renderProgressBar = (item, index) => <ProgressBarItem
    key={index}
    {...item}
    index={index}
    wrapperHeight={height}
  />

  return (
    <React.Fragment>
      <div className="progress-bar-wrapper" style={{ height: `${height}px` }}>
        {bars.map(renderProgressBar)}
      </div>
      {withLegend && items.map(renderProgressBarLegendItem)}
    </React.Fragment>
  )
}

const renderProgressBarLegendItem = (item, index) => <ProgressBarLegendItem key={index} {...item} />

const ProgressBarLegendItem = React.memo(({ value, label, color }) => {
  return (
    <div className="flex flex-row-center">
      <div className="dot" style={{ backgroundColor: color }}/>
      {displayAmount(value)}&nbsp;KEEP&nbsp;<span className="text-small text-grey">{label}</span>
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

export default ProgressBar
