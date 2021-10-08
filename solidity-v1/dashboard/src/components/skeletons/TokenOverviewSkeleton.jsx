import React from "react"
import CircleSkeleton from "./CircleSkeleton"
import Skeleton from "./Skeleton"

const TokenOverviewSkeleton = ({ items = 1 }) => {
  return (
    <>
      <Header />
      {Array.from(Array(items)).map((_, index) => (
        <CircularSumUpSkeleton key={index} />
      ))}
    </>
  )
}

export default TokenOverviewSkeleton

const Header = () => (
  <>
    <Skeleton shining color="grey-20" className="h4 mb-2" />
    <Skeleton
      shining
      className="h5"
      color="grey-20"
      styles={{ width: "50%", marginBottom: "0.5rem" }}
    />
    <Skeleton
      shining
      className="h5"
      color="grey-20"
      styles={{ width: "65%", marginBottom: "0.5rem" }}
    />
    <Skeleton
      shining
      color="grey-20"
      className="h5"
      styles={{ width: "40%", marginBottom: "0.5rem" }}
    />
  </>
)

const CircularSumUpSkeleton = ({ wrapperClassName = "" }) => (
  <div className={`${wrapperClassName} flex row mb-1`}>
    <CircleSkeleton shining width={110} height={110} />
    <div className="flex-1 ml-2 mt-1">
      <Skeleton shining color="grey-20" className="h5 mb-1" />
      <Skeleton shining color="grey-20" className="h5 mb-1" />
    </div>
  </div>
)

export const TokenGrantSkeletonOverview = () => (
  <section className="tile token-grant-overview">
    <div className="grant-amount">
      <Header />
    </div>
    <div className="unlocking-details flex-1">
      <CircularSumUpSkeleton wrapperClassName="flex-1" />
    </div>
    <div className="staked-details">
      <CircularSumUpSkeleton wrapperClassName="flex-1" />
    </div>
  </section>
)
