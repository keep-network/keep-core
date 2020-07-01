import React from "react"
import Skeleton from "./Skeleton"

const DelegatedTokensSkeleton = () => {
  return (
    <section className="flex flex-1 row wrap">
      <section className="tile delegation-overview">
        <Skeleton className="h2 mb-1" />
        <Skeleton className="h1 mb-1" />
        <Skeleton className="h6" styles={{ width: "70%", marginTop: "2rem" }} />
        <Skeleton className="h6" styles={{ width: "40%", marginTop: "1rem" }} />
        <Skeleton className="h6" styles={{ width: "55%", marginTop: "1rem" }} />
      </section>
      <section className="tile flex column undelegation-section">
        <Skeleton className="h4 mb-1" />
        <Skeleton
          className="mt-1"
          styles={{ marginBottom: "auto", width: "75%" }}
        />
        <Skeleton
          className="skeleton self-start"
          styles={{ width: "40%", padding: "2rem 4rem" }}
        />
      </section>
    </section>
  )
}

export default DelegatedTokensSkeleton
