import React from "react"
import PageWrapper from "../components/PageWrapper"
import Tile from "../components/Tile"
import * as Icons from "../components/Icons"
import ApplicationBox, {
  EmptyApplicationBox,
} from "../components/ApplicationBox"

const applications = [
  {
    icon: <Icons.TBTC />,
    name: "tBTC",
    websiteUrl: "https://tbtc.network/",
    websiteName: "tBTC Website",
    description:
      "Stake KEEP tokens and bond ETH to earn TBTC through ECDSA keeps.",
    btnLink: "/applications/tbtc",
  },
  {
    icon: <Icons.KeepBlackGreen />,
    name: "Random Beacon",
    websiteUrl: "https://keep.network/",
    websiteName: "Keep Website",
    description:
      "Stake KEEP tokens and earn ETH through Random Beacon signing groups.",
    btnLink: "/applications/random-beacon",
  },
]

const ApplicationOverviewPage = () => {
  return (
    <PageWrapper title="Applications">
      <Tile title="Select an application below to manage authorization contracts for operators.">
        <div className="flex row space-evenly wrap mt-1">
          {applications.map(renderApplicationBox)}
          <EmptyApplicationBox />
        </div>
      </Tile>
    </PageWrapper>
  )
}

const renderApplicationBox = (application) => (
  <ApplicationBox key={application.name} {...application} />
)

export default ApplicationOverviewPage
