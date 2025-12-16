package rpn

import "testing"

// TestMultiplication tests the multiplication operator.
// Implement the logic to handle the `*` operator.
//
// Example Input:
// - "4 5 *" -> 20
func TestMultiplication(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       int
		wantErr    bool
	}{
		{
			name:       "simple multiplication",
			expression: "4 5 *",
			want:       20,
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
