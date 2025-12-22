package rpn

import "testing"

// TestEmptyExpression tests the handling of an empty expression.
// An empty string should return 0.
//
// Example Input:
// - "" -> 0
func TestEmptyExpression(t *testing.T) {
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
