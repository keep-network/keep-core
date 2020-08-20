import React, { useState } from "react"
import StepNav from "../components/StepNav"
import {
  CopyStakeStepO,
  CopyStakeStep1,
  CopyStakeStep2,
  CopyStakeStep3,
} from "../components/copy-stake-steps"

const copyStakeSteps = ["balance", "upgrade", "review", "complete"]

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
      <main className="copy-stake__content-wrapper">
        <div className="copy-stake__content">
          {/* <CopyStakeStepO /> */}
          {/* <CopyStakeStep1 /> */}
          {/* <CopyStakeStep2
            amount={"10000000000000000000000000"}
            beneficiary={"0x6a0502bcaC31A40C3519920F6FC8E492DCEf87ca"}
            authorizerAddress={"0x6a0502bcaC31A40C3519920F6FC8E492DCEf87ca"}
            operatorAddress={"0x6a0502bcaC31A40C3519920F6FC8E492DCEf87ca"}
          /> */}
          {/* <CopyStakeStep3 /> */}
        </div>
      </main>
    </div>
  )
}

export default CopyStakePage
