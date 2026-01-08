package rpn

import "testing"

// TestComplexExpressions tests handling of complex expressions.
// Handle more complex expressions with multiple operators and also negative numbers.
//
// Example Inputs:
// - "15 7 1 1 + - / 3 * 2 1 1 + + -" -> 4
// - "-5" -> -5
// - "1 -5 +" -> -4
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
			want:       4,
			wantErr:    false,
		},
		{
			name:       "single negative number",
			expression: "-5",
			want:       -5,
			wantErr:    false,
		},
		{
			name:       "addition with negative number",
			expression: "1 -5 +",
			want:       -4,
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
