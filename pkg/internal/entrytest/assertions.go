package entrytest

import "testing"

// AssertEntryPublished checks if relay entry has been published to the chain.
// It does not inspect the entry.
func AssertEntryPublished(t *testing.T, testResult *Result) {
	if testResult.entry == nil {
		t.Fatal("expected relay entry to be published")
	}
}

// AssertEntryNotPublished checks if no relay entry has been published to
// the chain.
func AssertEntryNotPublished(t *testing.T, testResult *Result) {
	if testResult.entry != nil {
		t.Fatal("expected relay entry not to be published")
	}
}

// AssertNoSignerFailures checks there were no signer failures during the
// protocol execution.
func AssertNoSignerFailures(
	t *testing.T,
	testResult *Result,
) {
	if len(testResult.signerFailures) != 0 {
		t.Errorf(
			"expected no signer failures; has [%v]",
			len(testResult.signerFailures),
		)
	}
}

// AssertSignerFailuresCount checks the number of signers who failed the
// protocol execution. It does not check which particular signers failed.
func AssertSignerFailuresCount(
	t *testing.T,
	testResult *Result,
	expectedCount int,
) {
	if len(testResult.signerFailures) != expectedCount {
		t.Errorf(
			"unexpected number of signer failures\nexpected: [%v]\nactual:   [%v]",
			expectedCount,
			len(testResult.signerFailures),
		)
	}
}
