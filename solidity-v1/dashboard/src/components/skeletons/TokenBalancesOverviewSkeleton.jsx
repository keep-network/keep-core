import React from "react"
import CardSkeleton from "./CardSkeleton"
import Skeleton from "./Skeleton"
import CircleSkeleton from "./CircleSkeleton"

const TokenOverviewBalancesSkeleton = () => {
  return (
    <section className="tile" id="token-overview-balance">
      <div style={{ marginRight: "auto", width: "35%" }}>
        <Skeleton className="h2" styles={{ width: "45%", padding: "0.75em" }} />
        <div className="flex row center mt-2" styles={{ width: "50%" }}>
          <CircleSkeleton width={60} height={60} />
          <Skeleton className="h1 ml-1" styles={{ width: "50%" }} />
        </div>
        <Skeleton className="h3 mt-1" styles={{ width: "30%" }} />
      </div>
      <CardSkeleton />
      <CardSkeleton />
    </section>
  )
}

export default React.memo(TokenOverviewBalancesSkeleton)
