import { useFormikContext } from "formik"

export const useCustomOnSubmitFormik = (onSubmitAction) => {
  const {
    values,
    setSubmitting,
    setTouched,
    resetForm,
    validateForm,
    setFormikState,
  } = useFormikContext()

  const onSubmit = async (awaitingPromise) => {
    // Pre-submit
    const touched = {}
    Object.keys(values).forEach((name) => {
      touched[name] = true
    })
    setTouched(touched, false)
    setSubmitting(true)
    setFormikState((prevState) => ({
      ...prevState,
      submitCount: prevState.submitCount + 1,
    }))

    // Validation
    const errors = await validateForm(values)
    if (Object.keys(errors).length > 0) {
      setSubmitting(false)
      throw new Error("Invalid form")
    }

    // Submission
    try {
      await onSubmitAction(values, awaitingPromise)
      await awaitingPromise.promise
      setSubmitting(false)
      resetForm()
    } catch (error) {
      setSubmitting(false)
      throw error
    }
  }

  return onSubmit
}
