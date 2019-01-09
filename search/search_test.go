package search

import "testing"

func TestRunSpider(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "testdb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunSpider()
		})
	}
}
