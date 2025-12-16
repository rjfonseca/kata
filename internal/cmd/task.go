package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"

	"github.com/rjfonseca/kata/internal/i18n"
	"github.com/rjfonseca/kata/internal/taskrunner"
)

// Task executes a task or lists available tasks.
func Task(root string, args []string, translator i18n.Translator) error {
	runner, err := taskrunner.New(root)
	if err != nil {
		return fmt.Errorf("%s: %w", translator.T("task.error_init_runner"), err)
	}

	if len(args) > 0 {
		taskName := args[0]
		return runner.Run(taskName)
	}

	tasks, err := runner.ListTasks()
	if err != nil {
		return fmt.Errorf("%s: %w", translator.T("task.error_list_tasks"), err)
	}

	if len(tasks) == 0 {
		slog.Info(translator.T("task.log_no_tasks"))
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, t := range tasks {
		if t.Description == "" {
			_, _ = fmt.Fprintf(w, "%s\n", t.Name)
		} else {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", t.Name, t.Description)
		}
	}
	return w.Flush()
}
