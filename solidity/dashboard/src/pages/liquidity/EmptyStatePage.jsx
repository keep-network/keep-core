import React from "react"
import PageWrapper from "../../components/PageWrapper"
import EmptyState from "../../components/empty-states/EmptyState"
import { Skeleton, SkeletonBox } from "../../components/skeletons"
import * as Icons from "../../components/Icons"

const emptyPageTitle = "Check your liquidity rewards balance"
const emptyPageSubtitle =
  "This page tracks your liqudity rewards. Connect a wallet to view your available rewards."

const EmptyStatePage = (props) => {
  const { headerTitle, ...restProps } = props
  return (
    <PageWrapper {...restProps} title={headerTitle} newPage={true}>
      <EmptyState>
        <EmptyState.Skeleton className="empty-page--liquidity-page">
          {Array.from(Array(3)).map(rednerLiquidityCardSkeleton)}
        </EmptyState.Skeleton>
        <EmptyState.Title text={emptyPageTitle} />
        <EmptyState.Subtitle text={emptyPageSubtitle} />
        <EmptyState.ConnectWalletBtn />
      </EmptyState>
    </PageWrapper>
  )
}

const rednerLiquidityCardSkeleton = (_, index) => (
  <LiquidityCardSkeleton key={index} />
)

const styles = {
  liquidityInfoSkeleton: {
    height: "75px",
    borderRadius: "8px",
    marginRight: "0.5rem",
  },
  btnSkeleton: { borderRadius: "0" },
}

const LiquidityCardSkeleton = () => (
  <SkeletonBox>
    <Skeleton tag="h3" className="mb-1" />
    <Skeleton tag="h4" width="70%" color="grey-20" />
    <div className="flex row mt-2">
      <Skeleton
        className="flex-1"
        styles={styles.liquidityInfoSkeleton}
        color="grey-10"
      />
      <Skeleton
        className="flex-1 m-1"
        styles={styles.liquidityInfoSkeleton}
        color="grey-10"
      />
    </div>
    <div className="flex row center mt-2">
      <Icons.KeepCircle />
      <Skeleton tag="h3" color="grey-20" className="ml-1" />
    </div>
    <Skeleton
      tag="h2"
      color="grey-20"
      className="mt-1"
      styles={styles.btnSkeleton}
    />
    <Skeleton
      tag="h2"
      color="grey-10"
      className="mt-1"
      styles={styles.btnSkeleton}
    />
  </SkeletonBox>
)

export default EmptyStatePage
