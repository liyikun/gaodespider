package search

import (
	"fmt"
	"testing"

	mathcers "eatinlife.com/gaodespider/matchers"
)

func TestPaginationGet(t *testing.T) {
	tests := []struct {
		name string
		want []mathcers.RestaurantInfo
	}{
		{
			name: "test1",
			want: make([]mathcers.RestaurantInfo, 20),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PaginationGet()
			fmt.Println("hello")
			for _, v := range got {
				fmt.Printf("%s \n", v.Name)
			}
		})
	}
}
