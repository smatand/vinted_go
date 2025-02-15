package vinted

import (
	"reflect"
	"testing"
)

func TestExtractCatalogIDs(t *testing.T) {
	tests := []struct {
		name   string
		params string
		want   []int
	}{
		{
			name:   "single catalog ID",
			params: "https://www.vinted.sk/catalog?search_text=&catalog[]=79",
			want:   []int{79},
		},
		{
			name:   "multiple catalog IDs",
			params: "https://www.vinted.sk/catalog?search_text=&catalog[]=79&catalog[]=80",
			want:   []int{79, 80},
		},
		{
			name:   "no catalog IDs",
			params: "https://www.vinted.sk/catalog?search_text=",
			want:   nil,
		},
		{
			name:   "empty string",
			params: "",
			want:   nil,
		},
		{
			name:   "invalid catalogID",
			params: "https://www.vinted.sk/catalog?search_text=&catalog[]=281a",
			want:   []int{},
		},
		{
			name:   "invalid query parameters - missing catalog?",
			params: "https://www.vinted.sk/catalog[]=5",
			want:   nil,
		},
		{
			name:   "specific case: /catalog/2050-clothing",
			params: "https://www.vinted.sk/catalog/2050-clothing?search_text=",
			want:   []int{2050},
		},
		{
			name:   "specific case: /catalog/20a50-clothing, but return nil",
			params: "https://www.vinted.sk/catalog/20a50-clothing?search_text=",
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractCatalogIDs(tt.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractCatalogIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
