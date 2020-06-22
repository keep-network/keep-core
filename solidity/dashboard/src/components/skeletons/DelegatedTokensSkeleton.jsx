import React from "react"

const DelegatedTokensSkeleton = () => {
  return (
    <section className="flex flex-1 row wrap">
      <section className="tile mb-0 delegation-overview">
        <div className="skeleton h2 mb-1" />
        <div className="skeleton h1 mb-1" />
        <div
          className="skeleton h6"
          style={{ width: "70%", marginTop: "2rem" }}
        />
        <div
          className="skeleton h6"
          style={{ width: "40%", marginTop: "1rem" }}
        />
        <div
          className="skeleton h6"
          style={{ width: "55%", marginTop: "1rem" }}
        />
      </section>
      <section className="tile mb-0 flex column undelegation-section">
        <div className="skeleton h4 mb-1" />
        <div
          className="skeleton mt-1"
          style={{ marginBottom: "auto", width: "75%" }}
        />
        <div
          className="skeleton self-start"
          style={{ width: "40%", padding: "2rem 4rem" }}
        />
      </section>
    </section>
  )
}

export default DelegatedTokensSkeleton
