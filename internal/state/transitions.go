package state

import (
	"errors"
	"time"
)

var (
	ErrTestsFailing = errors.New("tests are failing")
	ErrNoNextStep   = errors.New("no next step available")
)

// CanRun always returns true for an initialized kata.
func (s *State) CanRun() bool {
	return true
}

// MarkRunResult updates the test result after a test execution.
func (s *State) MarkRunResult(passed bool) {
	s.TestPassing = passed
	s.LastRunAt = time.Now()
}

// HasNextStep returns true if there is another step after the current one.
func (s *State) HasNextStep() bool {
	return s.CurrentStepIndex < len(s.Steps)-1
}

// CanNext returns true if the kata can advance to the next step.
func (s *State) CanNext() bool {
	return s.TestPassing && s.HasNextStep()
}

// Next advances the kata to the next step.
// It resets TestPassing to false (Red).
func (s *State) Next() error {
	if !s.TestPassing {
		return ErrTestsFailing
	}

	if !s.HasNextStep() {
		return ErrNoNextStep
	}

	s.CurrentStepIndex++
	s.TestPassing = false
	return nil
}
