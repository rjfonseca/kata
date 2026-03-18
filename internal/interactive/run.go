package interactive

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/rjfonseca/kata/internal/state"
	"github.com/rjfonseca/kata/internal/taskrunner"
	"github.com/rjfonseca/kata/internal/watcher"
)

type Options struct {
	LoadState  func() (*state.State, error)
	RunCmd     func() *exec.Cmd
	NextCmd    func() *exec.Cmd
	RunTaskCmd func(name string) *exec.Cmd
	ListTasks  func() ([]taskrunner.TaskInfo, error)
	Watch      bool
}

func Run(ctx context.Context, opts Options) error {
	st, err := opts.LoadState()
	if err != nil {
		return err
	}

	m := &model{
		opts: opts,
		st:   st,
	}

	if opts.Watch {
		w, err := watcher.New(".", nil)
		if err != nil {
			return err
		}
		defer w.Close()
		m.watcher = w
	}

	m.refreshForm()

	p := tea.NewProgram(m, tea.WithContext(ctx))
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

type model struct {
	opts    Options
	st      *state.State
	form    *huh.Form
	choice  Action
	watcher *watcher.Watcher

	running bool
	msg     string
	err     error
}

func (m *model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.form.Init()}
	if m.watcher != nil {
		cmds = append(cmds, waitForFileChange(m.watcher.Events()))
	}
	return tea.Batch(cmds...)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case fileChangeMsg:
		if !m.running {
			m.running = true
			m.msg = "File changed. Running tests..."
			return m, tea.Batch(
				runExec(m.opts.RunCmd()),
				waitForFileChange(m.watcher.Events()),
			)
		}
		return m, waitForFileChange(m.watcher.Events())

	case runResultMsg:
		m.running = false
		m.err = msg.err
		m.st, _ = m.opts.LoadState()

		m.msg = statusMessage(m.st)
		// If command failed but state says passing, it might be a real error
		if m.err != nil && m.st.TestPassing {
			m.msg = fmt.Sprintf("Run error: %v", m.err)
		}

		m.refreshForm()
		return m, m.form.Init()

	case actionResultMsg:
		m.running = false
		m.err = msg.err
		m.st, _ = m.opts.LoadState()
		m.refreshForm()
		return m, m.form.Init()
	}

	if !m.running {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
			if m.form.State == huh.StateCompleted {
				return m, m.handleAction(m.choice)
			}
		}
		return m, cmd
	}

	return m, nil
}

func (m *model) View() string {
	if m.running {
		return fmt.Sprintf("\n%s\n", m.msg)
	}
	return m.form.View()
}

func (m *model) handleAction(a Action) tea.Cmd {
	m.running = true
	switch {
	case a == ActionRun:
		m.msg = "Running tests..."
		return runExec(m.opts.RunCmd())
	case a == ActionNext:
		m.msg = "Advancing..."
		return actionExec(m.opts.NextCmd())
	case a == ActionExit:
		return tea.Quit
	case strings.HasPrefix(string(a), "task:"):
		name := strings.TrimPrefix(string(a), "task:")
		m.msg = "Running task " + name + "..."
		return runExec(m.opts.RunTaskCmd(name))
	}
	return nil
}

func (m *model) refreshForm() {
	tasks, _ := m.opts.ListTasks()

	selectField := huh.NewSelect[Action]().
		Title("What would you like to do next?").
		Options(availableActions(m.st, tasks)...).
		Value(&m.choice)

	note := huh.NewNote().Title(statusMessage(m.st))
	if m.msg != "" {
		note = huh.NewNote().Title(m.msg)
	}

	m.form = huh.NewForm(huh.NewGroup(note, selectField))
}

type fileChangeMsg struct{}
type runResultMsg struct{ err error }
type actionResultMsg struct{ err error }

func waitForFileChange(ch <-chan struct{}) tea.Cmd {
	return func() tea.Msg {
		if ch == nil {
			return nil
		}
		<-ch
		return fileChangeMsg{}
	}
}

func runExec(c *exec.Cmd) tea.Cmd {
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return runResultMsg{err}
	})
}

func actionExec(c *exec.Cmd) tea.Cmd {
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return actionResultMsg{err}
	})
}
