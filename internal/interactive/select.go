package interactive

import (
	"context"

	"github.com/charmbracelet/huh"
)

// SelectKata prompts the user to select a kata from the given list.
func SelectKata(ctx context.Context, names []string) (string, error) {
	var selected string

	opts := make([]huh.Option[string], 0, len(names))
	for _, name := range names {
		opts = append(opts, huh.NewOption(name, name))
	}

	selectField := huh.NewSelect[string]().
		Title("Which kata do you want to start?").
		Options(opts...).
		Height(10).
		Value(&selected)

	form := huh.NewForm(
		huh.NewGroup(selectField),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}
