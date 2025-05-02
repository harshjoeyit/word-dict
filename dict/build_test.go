package dict

import "testing"

func TestCalcIndexSize(t *testing.T) {
	tests := []struct {
		name     string
		index    []IndexEntry
		expected int64
	}{
		{
			name:     "empty index",
			index:    []IndexEntry{},
			expected: int64(8),
		},
		{
			name: "single entry",
			index: []IndexEntry{
				{Word: "hello", Offset: 0, DefSize: 5},
			},
			expected: int64(8) + int64(len("hello")+1) + int64(8+1) + int64(2+1),
		},
		{
			name: "multiple entries",
			index: []IndexEntry{
				{Word: "hello", Offset: 0, DefSize: 5},
				{Word: "john", Offset: 10, DefSize: 5},
			},
			expected: int64(8) + int64(len("hello")+1) + int64(8+1) + int64(2+1) + int64(len("john")+1) + int64(8+1) + int64(2+1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calcIndexSize(tt.index)
			if result != tt.expected {
				t.Errorf("calcIndexSize() = %v, want %v", result, tt.expected)
			}
		})
	}
}
