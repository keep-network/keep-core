import React from "react"
import Skeleton from "./Skeleton"

const DelegatedTokensSkeleton = () => {
  return (
    <section className="flex flex-1 row wrap">
      <section className="tile delegation-overview">
        <Skeleton shining color="grey-20" className="h2 mb-1" />
        <Skeleton shining color="grey-20" className="h1 mb-1" />
        <Skeleton
          shining
          color="grey-20"
          className="h6"
          styles={{ width: "70%", marginTop: "2rem" }}
        />
        <Skeleton
          shining
          color="grey-20"
          className="h6"
          styles={{ width: "40%", marginTop: "1rem" }}
        />
        <Skeleton
          shining
          color="grey-20"
          className="h6"
          styles={{ width: "55%", marginTop: "1rem" }}
        />
      </section>
      <section className="tile flex column undelegation-section">
        <Skeleton shining color="grey-20" className="h4 mb-1" />
        <Skeleton
          shining
          color="grey-20"
          className="mt-1"
          styles={{ marginBottom: "auto", width: "75%" }}
        />
        <Skeleton
          shining
          color="grey-20"
          className="skeleton shining self-start"
          styles={{ width: "40%", padding: "2rem 4rem" }}
        />
      </section>
    </section>
  )
}

export default DelegatedTokensSkeleton
