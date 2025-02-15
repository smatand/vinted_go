package vinted

import (
	"reflect"
	"testing"
)

const catalogParam = "catalog[]"

func TestExtractCatalogIDs(t *testing.T) {
	tests := []struct {
		name      string
		urlStr    string
		paramName string
		want      []int
	}{
		{
			name:      "single catalog ID",
			urlStr:    "https://www.vinted.sk/catalog?search_text=&catalog[]=79",
			paramName: catalogParam,
			want:      []int{79},
		},
		{
			name:      "multiple catalog IDs",
			urlStr:    "https://www.vinted.sk/catalog?search_text=&catalog[]=79&catalog[]=80",
			paramName: catalogParam,
			want:      []int{79, 80},
		},
		{
			name:      "no catalog IDs",
			urlStr:    "https://www.vinted.sk/catalog?search_text=",
			paramName: catalogParam,
			want:      nil,
		},
		{
			name:      "empty string",
			urlStr:    "",
			paramName: catalogParam,
			want:      nil,
		},
		{
			name:      "invalid catalogID",
			urlStr:    "https://www.vinted.sk/catalog?search_text=&catalog[]=281a",
			paramName: catalogParam,
			want:      nil,
		},
		{
			name:      "invalid query parameters - missing catalog?",
			urlStr:    "https://www.vinted.sk/catalog[]=5",
			paramName: catalogParam,
			want:      nil,
		},
		{
			name:      "specific case: /catalog/2050-clothing",
			urlStr:    "https://www.vinted.sk/catalog/2050-clothing?search_text=",
			paramName: catalogParam,
			want:      []int{2050},
		},
		{
			name:      "specific case: /catalog/20a50-clothing, but return nil",
			urlStr:    "https://www.vinted.sk/catalog/20a50-clothing?search_text=",
			paramName: catalogParam,
			want:      nil,
		},
		{
			name:      "size_ids[] parameter",
			urlStr:    "https://www.vinted.sk/catalog?search_text=&size_ids[]=1&size_ids[]=2",
			paramName: "size_ids[]",
			want:      []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractIDs(tt.urlStr, tt.paramName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractCatalogIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
