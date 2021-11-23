import React, { useMemo } from "react"
import CircleSkeleton from "./CircleSkeleton"
import Skeleton from "./Skeleton"

const TokenAmountSkeleton = ({
  wrapperClassName = "flex row center",
  wrapperStyles,
  iconWidth = 35,
  iconHeight = 35,
  textStyles,
  textClassName = "h2 ml-1",
  icon: IconComponent = null,
  iconClassName = "",
}) => {
  const renderIcon = useMemo(() => {
    if (IconComponent) {
      return (
        <IconComponent
          className={`${iconClassName}`}
          width={iconWidth}
          height={iconHeight}
        />
      )
    }

    return <CircleSkeleton shining width={iconWidth} height={iconHeight} />
  }, [IconComponent, iconClassName, iconWidth, iconHeight])

  return (
    <div className={wrapperClassName} style={wrapperStyles}>
      {renderIcon}
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
