import React from "react"
import CircleSkeleton from "./CircleSkeleton"
import Skeleton from "./Skeleton"

const TokenAmountSkeleton = ({
  wrapperClassName = "flex row center",
  wrapperStyles,
  iconWidth = 35,
  iconHeight = 35,
  textStyles,
  textClassName = "h2 ml-1",
}) => {
  return (
    <div className={wrapperClassName} style={wrapperStyles}>
      <CircleSkeleton shining width={iconWidth} height={iconHeight} />
      <Skeleton
        color="grey-20"
        shining
        className={textClassName}
        styles={textStyles}
      />
    </div>
  )
}

export default React.memo(TokenAmountSkeleton)
