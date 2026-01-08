package rpn

import "testing"

// TestSingleNumber tests the parsing of single numbers.
// A string with a single number should return that number.
//
// Example Inputs:
// - "5" -> 5
// - "-1" -> -1
func TestSingleNumber(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       int
		wantErr    bool
	}{
		{
			name:       "positive number",
			expression: "5",
			want:       5,
			wantErr:    false,
		},
		{
			name:       "negative number",
			expression: "-1",
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
