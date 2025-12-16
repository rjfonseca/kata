package interactive

import "github.com/charmbracelet/huh"

// Select prompts the user to select an option from the given list.
func Select(title string, options []string) (string, error) {
	if len(options) == 0 {
		return "", nil
	}

	var selected string

	selectField := huh.NewSelect[string]().
		Title(title).
		Options(huh.NewOptions(options...)...).
		Value(&selected)

	form := huh.NewForm(
		huh.NewGroup(selectField),
	)

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
}
