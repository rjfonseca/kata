package rpn

import "testing"

// TestComplexExpressions tests handling of complex expressions.
// Handle more complex expressions with multiple operators and also negative numbers.
//
// Example Inputs:
// - "15 7 1 1 + - / 3 * 2 1 1 + + -" -> 5
// - "5 1 2 + 4 * + 3 -" -> 14
func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       int
		wantErr    bool
	}{
		{
			name:       "complex expression 1",
			expression: "15 7 1 1 + - / 3 * 2 1 1 + + -",
			want:       5,
			wantErr:    false,
		},
		{
			name:       "complex expression 2",
			expression: "5 1 2 + 4 * + 3 -",
			want:       14,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate(%s) error = %v, wantErr %v", tt.expression, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluate(%s) = %v, want %v", tt.expression, got, tt.want)
			}
		})
	}
}
