package interactive

import (
	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/state"
)

func statusMessage(st *state.State, translator i18n.Translator) string {
	switch {
	case !st.TestPassing:
		return translator.T("interactive.status_failing")

	case st.TestPassing && !st.KataFinished:
		return translator.T("interactive.status_passing")

	default:
		return translator.T("interactive.status_completed")
	}
}
