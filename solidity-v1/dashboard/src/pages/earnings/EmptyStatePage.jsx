import React from "react"
import EmptyState from "../../components/empty-states/EmptyState"
import * as Icons from "../../components/Icons"
import { SkeletonBox, Skeleton } from "../../components/skeletons"

const title = "Check your earnings balance"
const subtitle = "Connect a wallet to view your current earnings balance."

const EmptyStatePage = (props) => {
  return (
    <EmptyState>
      <EmptyState.Skeleton className="empty-page--operatorations">
        <SkeletonBox>
          <Skeleton tag="h3" className="mb-2" width="75%" />
          <div className="flex row center mb-2">
            <Icons.KeepCircle />
            <Skeleton tag="h3" color="grey-20" className="ml-1" width="50%" />
          </div>
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="50%" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="45%" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="80%" />
        </SkeletonBox>
      </EmptyState.Skeleton>
      <EmptyState.Title text={title} />
      <EmptyState.Subtitle text={subtitle} />
      <EmptyState.ConnectWalletBtn />
    </EmptyState>
  )
}

export default EmptyStatePage
