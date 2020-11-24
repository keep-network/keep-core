import React from "react"
import EmptyState from "../../components/empty-states/EmptyState"
import * as Icons from "../../components/Icons"
import {
  SkeletonBox,
  Skeleton,
  SkeletonProgressBar,
} from "../../components/skeletons"

const title = "Earn with your KEEP"
const subtitle =
  "This page tracks your staked KEEP tokens. Connect a wallet that contains KEEP tokens to view your staked tokens."

const EmptyStateComponent = () => {
  return (
    <EmptyState>
      <EmptyState.Skeleton className="empty-page--delegation">
        {/* row 1 column 1 */}
        <SkeletonBox>
          <Skeleton tag="h3" className="mb-2" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="40%" />
          <Skeleton tag="h5" className="mb-1" color="grey-20" width="80%" />
        </SkeletonBox>

        {/* row 1-2 column 2 */}
        <SkeletonBox>
          <Skeleton tag="h3" className="mb-2" width="75%" />
          <div className="flex row mb-2">
            <Icons.KeepCircle />
            <div className="flex-1 ml-1">
              <Skeleton tag="h6" color="grey-20" className="mb-1" width="50%" />
              <Skeleton tag="h6" color="grey-20" className="mb-1" width="90%" />
            </div>
          </div>

          <Skeleton tag="h4" color="grey-10" className="mb-1" />
          <Skeleton tag="h4" color="grey-10" className="mb-1" />
          <Skeleton tag="h4" color="grey-10" className="mb-1" />
          <Skeleton tag="h3" color="grey-30" />
        </SkeletonBox>

        {/* row 1 column 3 */}
        <SkeletonBox>
          <Skeleton />
          <Skeleton
            tag="h5"
            color="grey-10"
            className="mt-1 mb-2"
            width="70%"
          />
        </SkeletonBox>

        {/* row 2 column 1 */}
        <SkeletonBox>
          <SkeletonProgressBar fillingInPercentage="75" />
          <Skeleton tag="h5" color="grey-20" className="mt-1" />
        </SkeletonBox>

        {/* row 2 column 3 */}
        <SkeletonBox>
          <SkeletonProgressBar />
          <Skeleton tag="h5" />
        </SkeletonBox>
      </EmptyState.Skeleton>
      <EmptyState.Title text={title} />
      <EmptyState.Subtitle text={subtitle} />
      <EmptyState.ConnectWalletBtn />
    </EmptyState>
  )
}

export default EmptyStateComponent
