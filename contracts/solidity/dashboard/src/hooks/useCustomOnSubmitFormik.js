import { useFormikContext } from 'formik'

export const useCustomOnSubmitFormik = (onSubmitAction) => {
  const {
    values,
    setSubmitting,
    setTouched,
    resetForm,
    validateForm,
    setFormikState,
  } = useFormikContext()

  const onSubmit = async (onTransactionHashCallback, openInfoMessage, setIsFetchng) => {
    // Pre-submit
    const touched = {}
    Object.keys(values).forEach((name) => {
      touched[name] = true
    })
    setTouched(touched, false)
    setSubmitting(true)
    setFormikState((prevState) => ({ ...prevState, submitCount: prevState.submitCount + 1 }))

    // Validation
    const errors = await validateForm(values)
    console.log('errors', errors)
    if (errors) {
      setSubmitting(false)
      return
    }

    // Submission
    try {
      openInfoMessage()
      setIsFetchng()
      await onSubmitAction(values, onTransactionHashCallback)
      setSubmitting(false)
      resetForm()
    } catch (error) {
      setSubmitting(false)
    }
  }

  return onSubmit
}
