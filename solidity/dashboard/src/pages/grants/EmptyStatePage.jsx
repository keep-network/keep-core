import React from "react"
import PageWrapper from "../../components/PageWrapper"
import EmptyState from "../../components/empty-states/EmptyState"
import * as Icons from "../../components/Icons"
import {
  SkeletonBox,
  Skeleton,
  SkeletonProgressBar,
} from "../../components/skeletons"

const title = "Track your token grant"
const subtitle =
  "Find information about your token grant unlocking schedule here. Connect a wallet first to view any token grants."

const EmptyStatePage = ({ title: pageTitle }) => {
  return (
    <PageWrapper title={pageTitle}>
      <EmptyState>
        <EmptyState.Skeleton className="empty-page--grants">
          <SkeletonBox>
            <div>
              <Skeleton tag="h3" className="mb-2" />
              <div className="flex row mb-2">
                <Icons.KeepCircle />
                <div className="flex-1 ml-1">
                  <Skeleton
                    tag="h6"
                    color="grey-20"
                    className="mb-1"
                    width="90%"
                  />
                  <Skeleton
                    tag="h6"
                    color="grey-20"
                    className="mb-1"
                    width="50%"
                  />
                </div>
              </div>
            </div>
            <div>
              <div>
                <SkeletonProgressBar fillingInPercentage="65" />
                <Skeleton tag="h5" color="grey-20" className="mt-1" />
              </div>
              <div>
                <SkeletonProgressBar fillingInPercentage="75" />
                <Skeleton tag="h5" color="grey-20" className="mt-1" />
              </div>
            </div>
          </SkeletonBox>
        </EmptyState.Skeleton>
        {/* TODO add tooltip to the title when PR with tooltip will be merged */}
        <EmptyState.Title text={title} />
        <EmptyState.Subtitle text={subtitle} />
        <EmptyState.ConnectWalletBtn />
      </EmptyState>
    </PageWrapper>
  )
}

export { EmptyStatePage }
