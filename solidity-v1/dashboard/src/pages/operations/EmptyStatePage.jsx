import React from "react"
import EmptyState from "../../components/empty-states/EmptyState"
import * as Icons from "../../components/Icons"
import { SkeletonBox, Skeleton } from "../../components/skeletons"
import PageWrapper from "../../components/PageWrapper"

const title = "Check your staking balance"
const subtitle =
  "This page tracks delegated tokens to your address. Connect a wallet to view your current balance and slashing events."

const EmptyStatePage = (props) => {
  return (
    <PageWrapper {...props}>
      <EmptyState>
        <EmptyState.Skeleton className="empty-page--operatorations">
          <SkeletonBox>
            <div className="flex row">
              <div className="flex-1">
                <Skeleton tag="h3" className="mb-2" width="75%" />
                <div className="flex row center mb-2">
                  <Icons.KeepCircle />
                  <Skeleton
                    tag="h3"
                    color="grey-20"
                    className="ml-1"
                    width="50%"
                  />
                </div>
                <Skeleton
                  tag="h5"
                  className="mb-1"
                  color="grey-20"
                  width="50%"
                />
                <Skeleton
                  tag="h5"
                  className="mb-1"
                  color="grey-20"
                  width="45%"
                />
                <Skeleton
                  tag="h5"
                  className="mb-1"
                  color="grey-20"
                  width="80%"
                />
              </div>
              <div className="flex column flex-2">
                <Skeleton tag="h3" className="mb-1" width="35%" />
                <Skeleton
                  tag="h5"
                  color="grey-20"
                  className="mb-2"
                  width="75%"
                />
                <Skeleton
                  tag="h1"
                  color="grey-20"
                  className="mt-a"
                  width="45%"
                />
              </div>
            </div>
          </SkeletonBox>
        </EmptyState.Skeleton>
        <EmptyState.Title text={title} />
        <EmptyState.Subtitle text={subtitle} />
        <EmptyState.ConnectWalletBtn />
      </EmptyState>
    </PageWrapper>
  )
}

export default EmptyStatePage
