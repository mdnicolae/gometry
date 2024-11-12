package logging

import (
	"reflect"
	"testing"
)

func TestNormalizeAttributes(t *testing.T) {
	tests := []struct {
		name       string
		attributes []map[string]interface{}
		expected   map[string]interface{}
	}{
		{
			name:       "nil attributes",
			attributes: nil,
			expected:   map[string]interface{}{},
		},
		{
			name:       "empty attributes",
			attributes: []map[string]interface{}{},
			expected:   map[string]interface{}{},
		},
		{
			name:       "non-nil attributes",
			attributes: []map[string]interface{}{{"key": "value"}},
			expected:   map[string]interface{}{"key": "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeAttributes(tt.attributes)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
