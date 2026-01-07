package rpn

import "testing"

// Step 1: Basic Cases
//
// Handle the basic cases:
// - An empty string should return `0`.
// - A string with a single number should return that number.
//
// Example Inputs:
// - `""` -> `0`
// - `"5"` -> `5`

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       int
		wantErr    bool
	}{
		{
			name:       "empty expression",
			expression: "",
			want:       0,
			wantErr:    false,
		},
		{
			name:       "single number",
			expression: "5",
			want:       5,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
