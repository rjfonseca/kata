package rpn

import "testing"

// TestDivision tests the division operator.
// Implement the logic to handle the `/` operator. Division by zero should return an error.
//
// Example Inputs:
// - "10 2 /" -> 5
// - "10 0 /" -> error
func TestDivision(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       int
		wantErr    bool
	}{
		{
			name:       "simple division",
			expression: "10 2 /",
			want:       5,
			wantErr:    false,
		},
		{
			name:       "division by zero",
			expression: "10 0 /",
			want:       0,
			wantErr:    true,
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
