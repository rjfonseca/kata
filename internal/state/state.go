package state

import "time"

// State represents the execution state of a kata.
// It models the TDD cycle explicitly: tests failing â tests passing.
type State struct {
	KataName string `json:"kata_name"`

	Runner           string   `json:"runner"`
	Steps            []string `json:"steps"`
	CurrentStepIndex int      `json:"current_step_index"`

	// TestPassing represents whether the last test run passed.
	// false == Red, true == Green
	TestPassing bool `json:"test_passing"`

	// KataFinished is set to true the first time the user
	// attempts to advance past the last step.
	KataFinished bool `json:"kata_finished"`

	StartedAt time.Time `json:"started_at"`
	LastRunAt time.Time `json:"last_run_at,omitzero"`
}
