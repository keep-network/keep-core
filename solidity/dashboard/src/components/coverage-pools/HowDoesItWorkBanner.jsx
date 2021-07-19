import React from "react"
import Banner from "../Banner"
import * as Icons from "../Icons"
import Tag from "../Tag"
import List from "../List"
import OnlyIf from "../OnlyIf"
import { useHideComponent } from "../../hooks/useHideComponent"

const benefits = [
  { icon: Icons.FeesVector, label: "High APY rewards" },
  { icon: Icons.Decentralize, label: "Secure the network" },
]

const CheckListBanner = () => {
  const [isBannerVisible, hide] = useHideComponent(false)

  return (
    <OnlyIf condition={isBannerVisible}>
      <Banner className="coverage-pool__checklist">
        <Banner.CloseIcon onClick={hide} />
        <Banner.Title>How does it work?</Banner.Title>
        <div className="checklist-wrapper">
          <List className="checklist__section checklist__section--docs">
            <List.Title className="text-grey-60 mb-1">Overview</List.Title>
            <List.Content>
              <List.Item>
                KEEP holders (the underwriters) deposit into the coverage pool
                to cover risk. Deposit KEEP to ensure the safety of the peg in
                exchange for rewards.
              </List.Item>
              <List.Item>
                <Tag
                  text={
                    <a
                      href="https://example.com"
                      rel="noopener noreferrer"
                      target="_blank"
                      className="text-black"
                    >
                      Read Documentation
                    </a>
                  }
                />
              </List.Item>
            </List.Content>
          </List>

          <List
            className="checklist__section checklist__section--needed"
            items={benefits}
          >
            <List.Title className="text-grey-60 mb-1">Overview</List.Title>
            <List.Content />
          </List>
        </div>
      </Banner>
    </OnlyIf>
  )
}

export default CheckListBanner
