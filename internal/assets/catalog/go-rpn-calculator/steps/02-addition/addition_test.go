package rpn

import "testing"

// TestAddition tests the addition operator.
// Implement the logic to handle the `+` operator for two numbers.
//
// Example Input:
// - "1 2 +" -> 3
func TestAddition(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       int
		wantErr    bool
	}{
		{
			name:       "simple addition",
			expression: "1 2 +",
			want:       3,
			wantErr:    false,
		},
		{
			name:       "addition with negative",
			expression: "1 -2 +",
			want:       -1,
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
