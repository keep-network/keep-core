import React from "react"
import Banner from "../Banner"
import * as Icons from "../Icons"

const infoBannerTitle = "The cooldown period is 21 days."

const infoBannerDescription =
  "A withdrawn deposit will be available to claim after 21 days. During cooldown, your funds will accumulate rewards but are also subject to risk to cover for a hit."

const styles = {
  banner: { minWidth: "100%" },
  iconWrapper: { flex: "0 0 auto" },
}

const CooldownPeriodBanner = () => {
  return (
    <Banner
      className="banner--info flex row center mt-2 mb-2"
      style={styles.banner}
    >
      <div style={styles.iconWrapper}>
        <Banner.Icon
          icon={Icons.Tooltip}
          className="mr-1"
          backgroundColor="transparent"
          color="black"
        />
      </div>
      <div>
        <Banner.Title>{infoBannerTitle}</Banner.Title>
        <Banner.Description>{infoBannerDescription}</Banner.Description>
      </div>
    </Banner>
  )
}

export default CooldownPeriodBanner
