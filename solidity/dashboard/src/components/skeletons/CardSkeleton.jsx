import React from "react"
import CircleSkeleton from "./CircleSkeleton"
import Skeleton from "./Skeleton"
import { colors } from "../../constants/colors"

const cardWrapperDefaultStyle = {
  border: `1px solid ${colors.grey20}`,
  padding: "2rem",
  minWidth: "325px",
  minHeight: "345px",
  margin: "2rem",
}

const CardSkeleton = ({ cardWrapperStyle = cardWrapperDefaultStyle }) => {
  return (
    <div
      style={cardWrapperStyle}
      className="flex column center space-between mt-1"
    >
      <CircleSkeleton width={60} height={60} />
      <Skeleton className="h4" styles={{ width: "50%", marginTop: "-2rem" }} />
      <Skeleton
        className="text-small"
        styles={{ width: "80%", marginTop: "-2rem" }}
      />
      <Skeleton className="skeleton" styles={{ padding: "1.5rem 5rem" }} />
    </div>
  )
}

export default React.memo(CardSkeleton)
