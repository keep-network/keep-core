import React, { useState } from "react"
import StepNav from "../components/StepNav"

const copyStakeSteps = ["stake", "review", "complete"]

const CopyStakePage = () => {
  const [step, setStep] = useState(0)

  return (
    <div className="copy-stake__layout">
      <nav className="copy-stake__nav">
        <div
          className="copy-stake__nav__indicator"
          onClick={() => setStep((prevState) => prevState + 1)}
        >
          <StepNav steps={copyStakeSteps} activeStep={step} />
        </div>
      </nav>
      <main className="copy-stake__content"></main>
    </div>
  )
}

export default CopyStakePage
