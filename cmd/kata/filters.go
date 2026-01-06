package main

import "strings"

func filterTasks(tasks []string) []string {
	var filtered []string
	for _, t := range tasks {
		if t == "test" {
			continue
		}
		if strings.HasSuffix(t, "_hook") {
			continue
		}
		filtered = append(filtered, t)
	}
	return filtered
}
