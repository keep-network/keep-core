import { useFormikContext } from 'formik'

export const useCustomOnSubmitFormik = (onSubmitAction) => {
  const {
    values,
    setSubmitting,
    setFormikState,
    setTouched,
    resetForm,
    validateForm,
  } = useFormikContext()

  const onSubmit = async (onTransactionHashCallback, openInfoMessage, setIsFetchng) => {
    // Pre-submit
    setSubmitting(true)
    setFormikState((prevState) => ({ submitCount: prevState.submitCount + 1 }))

    // Validation
    const errors = await validateForm(values)
    console.log('errors', errors)
    if (errors) {
      const touched = {}
      Object.keys(errors).forEach((name) => {
        touched[name] = true
      })
      setTouched(touched, false)
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
