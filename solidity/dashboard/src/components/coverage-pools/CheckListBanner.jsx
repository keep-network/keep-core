import React from "react"
import Banner from "../Banner"
import * as Icons from "../Icons"
import Tag from "../Tag"
import List from "../List"

const needed = [
  { icon: Icons.BrowserWindow, label: "Web3-compatible browser" },
  { icon: Icons.KeepOutline, label: "KEEP tokens" },
  { icon: Icons.ETH, label: "Ethereum wallet" },
]

const notes = [
  { icon: Icons.Warning, label: "Cooldown periods apply" },
  { icon: Icons.Warning, label: "Risk warning" },
]

const CheckListBanner = () => {
  return (
    <Banner className="coverage-pool__checklist">
      <Banner.CloseIcon />
      <Banner.Title>Checklist</Banner.Title>
      <div className="checklist-wrapper">
        <List
          className="checklist__section checklist__section--needed"
          items={needed}
        >
          <List.Title className="text-grey-60 mb-1">
            What You&apos;ll Need
          </List.Title>
          <List.Content />
        </List>

        <List className="checklist__section checklist__section--docs">
          <List.Title className="text-grey-60 mb-1">Documentation</List.Title>
          <List.Content>
            <List.Item>
              Read the documentation to learn more about participating in the
              coverage pool.
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
          className="checklist__section checklist__section--notes"
          items={notes}
        >
          <List.Title className="text-grey-60 mb-1">Take note</List.Title>
          <List.Content />
        </List>
      </div>
    </Banner>
  )
}

export default CheckListBanner
