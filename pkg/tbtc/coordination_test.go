package tbtc

import (
	"github.com/keep-network/keep-core/internal/testutils"
	"testing"
)

func TestCoordinationWindow_ActivePhaseEndBlock(t *testing.T) {
	window := newCoordinationWindow(900)

	testutils.AssertIntsEqual(
		t,
		"active phase end block",
		980,
		int(window.activePhaseEndBlock()),
	)
}

func TestCoordinationWindow_EndBlock(t *testing.T) {
	window := newCoordinationWindow(900)

	testutils.AssertIntsEqual(
		t,
		"end block",
		1000,
		int(window.endBlock()),
	)
}

func TestCoordinationWindow_IsAfterActivePhase(t *testing.T) {
	window := newCoordinationWindow(1800)

	previousWindow := newCoordinationWindow(900)
	sameWindow := newCoordinationWindow(1800)
	nextWindow := newCoordinationWindow(2700)

	testutils.AssertBoolsEqual(
		t,
		"result for nil",
		true,
		window.isAfter(nil),
	)
	testutils.AssertBoolsEqual(
		t,
		"result for previous window",
		true,
		window.isAfter(previousWindow),
	)
	testutils.AssertBoolsEqual(
		t,
		"result for same window",
		false,
		window.isAfter(sameWindow),
	)
	testutils.AssertBoolsEqual(
		t,
		"result for next window",
		false,
		window.isAfter(nextWindow),
	)
}