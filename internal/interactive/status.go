package interactive

import "github.com/rjfonseca/kata/internal/state"

func statusMessage(st *state.State) string {
	switch {
	case !st.TestPassing:
		return "Tests are failing ❌"

	case st.TestPassing && !st.KataFinished:
		return "Tests are passing ✅"

	default:
		return "Kata completed 🎉 You can keep refactoring and re-running tests."
	}
}
