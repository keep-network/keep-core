import React, { useMemo, useRef, useEffect } from "react"

const StepNav = ({ steps, activeStep = 0 }) => {
  const stepNavElement = useRef(null)
  const stepActiveIndicatorElement = useRef(null)
  const numberOfSteps = useRef(steps.length)

  useEffect(() => {
    if (
      !stepNavElement.current ||
      !stepActiveIndicatorElement ||
      activeStep > numberOfSteps.current ||
      activeStep === 0
    ) {
      return
    }

    const stepOffsetHeight =
      activeStep === numberOfSteps.current
        ? "100%"
        : `${
            stepNavElement.current
              .getElementsByClassName("step-nav__step")
              [activeStep - 1].getBoundingClientRect().y
          }px`

    stepActiveIndicatorElement.current.style.height = stepOffsetHeight
  }, [activeStep])

  const stepsComponents = useMemo(() => {
    return steps.map(renderStep)
  }, [steps])

  return (
    <div ref={stepNavElement} className="step-nav">
      <span className="step-nav__indicator" />
      <span
        ref={stepActiveIndicatorElement}
        className="step-nav__indicator--active"
      />
      {stepsComponents}
    </div>
  )
}

const Step = ({ step }) => {
  return (
    <div className="step-nav__step">
      <h5 className="step-nav__step__content">{step}</h5>
    </div>
  )
}

const renderStep = (item, index) => <Step key={index} step={item} />

export default React.memo(StepNav)
