import React from "react"

const NoData = ({ title, iconComponent, content }) => (
  <div className="no-data">
    <h2>{title}</h2>
    {iconComponent}
    <span>{content}</span>
  </div>
)

export default NoData
