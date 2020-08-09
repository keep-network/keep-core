import React from "react"

const Tag = ({ IconComponent, text, type, size }) => {
  return (
    <div
      className={`tag${type ? `--${type}` : ""} ${size ? `tag--${size}` : ""}`}
    >
      <div className="flex row center">
        <div className="tag__icon">
          <IconComponent />
        </div>
        <div className="tag__text">{text}</div>
      </div>
    </div>
  )
}

export default React.memo(Tag)
