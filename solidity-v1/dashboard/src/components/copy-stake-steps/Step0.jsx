import React from "react"
import Button from "../Button"
import * as Icons from "../Icons"

const styles = {
  subtitle: { textAlign: "justify" },
  newFeaturesSection: { borderRadius: "10px", width: "100%" },
}

const newFeaturesItems = [
  {
    title: "faster dashboard performance",
    icon: Icons.CarDashboardSpeed,
  },
  {
    title: "add tokens to existing stake",
    icon: Icons.Fees,
  },
  {
    title: "more user friendly staking",
    icon: Icons.UserFriendly,
  },
]

const CopyStakeStepO = ({ incrementStep }) => {
  return (
    <>
      <h1 className="mb-1">
        Move your stake to Keep’s upgraded staking contract!
      </h1>
      <h3 className="mb-3 text-grey-70" style={styles.subtitle}>
        To continue running smoothly on the Keep network, any stake that was
        delegated to the previous contract version will need to move to the
        upgraded staking contract.
      </h3>
      <div className="tile" style={styles.newFeaturesSection}>
        <h3 className="mb-2">New Features</h3>
        <div className="flex row center space-between">
          {newFeaturesItems.map(renderNewFeatureItem)}
        </div>
      </div>
      <Button className="btn btn-primary btn-lg mt-2" onClick={incrementStep}>
        start upgrade
      </Button>
    </>
  )
}

const renderNewFeatureItem = (item) => (
  <NewFeatureItem key={item.title} {...item} />
)

const NewFeatureItem = ({ title, icon: IconComponent }) => {
  return (
    <div className="flex column center">
      <IconComponent />
      <h5 className="text-grey-50 text-center">{title}</h5>
    </div>
  )
}

export default CopyStakeStepO
