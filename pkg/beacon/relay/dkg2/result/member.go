package resultsubmission

// ResultSubmittingMember represents a member submitting a DKG result to the
// blockchain along with signatures received from other group members supporting
// the result.
type ResultSubmittingMember struct {
	*ResultSigningMember
	// Predefined step for each submitting window. The value is used to determine
	// eligible submitting member.
	blockStep uint32
}
