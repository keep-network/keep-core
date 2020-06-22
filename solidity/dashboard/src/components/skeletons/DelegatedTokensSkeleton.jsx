import React from "react"

const DelegatedTokensSkeleton = () => {
  return (
    <section className="flex flex-1 row wrap">
      <section className="tile delegation-overview">
        <div className="skeleton h2 mb-1" />
        <div className="skeleton h1 mb-1" />
        <div
          className="skeleton h6"
          style={{ width: "70%", marginTop: "0.5rem" }}
        />
        <div
          className="skeleton h6"
          style={{ width: "40%", marginTop: "0.5rem" }}
        />
        <div
          className="skeleton h6"
          style={{ width: "55%", marginTop: "0.5rem" }}
        />
      </section>
      <section className="tile flex column undelegation-section">
        <div className="skeleton h4 mb-1" />
        <div
          className="skeleton"
          style={{ marginBottom: "auto", width: "75%" }}
        />
        <div className="skeleton" style={{ width: "40%", padding: "2rem" }} />
      </section>
    </section>
  )
}

export default DelegatedTokensSkeleton
