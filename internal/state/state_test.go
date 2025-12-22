package state

import "testing"

func newTestState() *State {
	return &State{
		KataName:         "hello-world",
		Steps:            []string{"01", "02"},
		CurrentStepIndex: 0,
		TestPassing:      false,
	}
}

func TestRunTransitions(t *testing.T) {
	s := newTestState()

	s.MarkRunResult(false)
	if s.TestPassing != false {
		t.Fatalf("expected TestsNotPassing")
	}

	s.MarkRunResult(true)
	if s.TestPassing != true {
		t.Fatalf("expected TestsPassing")
	}
}

func TestCannotAdvanceWhenTestsFailing(t *testing.T) {
	s := newTestState()

	if s.CanNext() {
		t.Fatalf("should not allow next when tests are failing")
	}

	if err := s.Next(); err == nil {
		t.Fatalf("expected error when advancing in red")
	}
}

func TestAdvanceResetsToRed(t *testing.T) {
	s := newTestState()
	s.MarkRunResult(true)

	if err := s.Next(); err != nil {
		t.Fatal(err)
	}

	if s.CurrentStepIndex != 1 {
		t.Fatalf("expected step index 1")
	}

	if s.TestPassing != false {
		t.Fatalf("expected TestsNotPassing after advancing")
	}
}

func TestCannotAdvancePastLastStep(t *testing.T) {
	s := &State{
		KataName:         "hello-world",
		Steps:            []string{"01"},
		CurrentStepIndex: 0,
		TestPassing:      true,
	}

	if s.HasNextStep() {
		t.Fatalf("should not have next step")
	}

	if s.CanNext() {
		t.Fatalf("should not allow next on last step")
	}
}
