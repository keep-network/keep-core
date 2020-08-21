import React from "react"
import * as Icons from "../Icons"
import Button from "../Button"

const CopyStakeStep4 = () => {
  return (
    <>
      <Icons.KEEPTower />
      <h2 className="mb-2 mt-2">
        Your stake balance is successfully copied and redelegated.
      </h2>
      <section className="tile">
        <h4 className="text-grey-40">
          Once your former stake fully undelegates, you’ll need to initiate the
          recovery process in the dashboard. You’ll see a notification in the
          dashboard when it’s time to do this.
        </h4>
      </section>
      <Button className="btn btn-primary btn-lg">cancel</Button>
    </>
  )
}

export default CopyStakeStep4
