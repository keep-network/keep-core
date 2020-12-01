import React from "react"
import EmptyState from "../../components/empty-states/EmptyState"
import * as Icons from "../../components/Icons"
import { SkeletonBox, Skeleton } from "../../components/skeletons"

const title = "Check your rewards balance"
const subtitle =
  "This page tracks your rewards. Connect a wallet to view your available rewards and history."

const EmptyStatePage = () => {
  return (
    <EmptyState>
      <EmptyState.Skeleton className="empty-page--rewards-overview">
        <SkeletonBox>
          <Skeleton tag="h3" className="mb-2" width="75%" />
          <div className="flex row center mb-2">
            <Icons.KeepCircle />
            <Skeleton tag="h3" color="grey-20" className="ml-1" width="50%" />
          </div>
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="65%" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="65%" />
        </SkeletonBox>

        <SkeletonBox>
          <Skeleton tag="h3" className="mb-2" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="65%" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="45%" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="85%" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="40%" />
        </SkeletonBox>

        <SkeletonBox>
          <Skeleton tag="h3" className="mb-2" width="75%" />
          <div className="flex row center mb-2">
            <Icons.KeepCircle />
            <Skeleton tag="h3" color="grey-20" className="ml-1" width="50%" />
          </div>
          <Skeleton tag="h3" color="grey-30" width="50%" />
        </SkeletonBox>

        <SkeletonBox>
          <Skeleton tag="h3" className="mb-2" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="100%" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="85%" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="75%" />
        </SkeletonBox>
      </EmptyState.Skeleton>
      <EmptyState.Title text={title} />
      <EmptyState.Subtitle text={subtitle} />
      <EmptyState.ConnectWalletBtn />
    </EmptyState>
  )
}

export default EmptyStatePage
