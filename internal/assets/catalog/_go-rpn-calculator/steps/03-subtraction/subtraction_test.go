package rpn

import "testing"

// TestSubtraction tests the subtraction operator.
// Implement the logic to handle the `-` operator. Remember that order matters.
//
// Example Input:
// - "5 3 -" -> 2
func TestSubtraction(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       int
		wantErr    bool
	}{
		{
			name:       "simple subtraction",
			expression: "5 3 -",
			want:       2,
			wantErr:    false,
		},
		{
			name:       "subtraction resulting in negative",
			expression: "3 5 -",
			want:       -2,
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
