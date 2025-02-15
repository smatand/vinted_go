package vinted

import (
	"reflect"
	"testing"
)

func TestExtractIDs(t *testing.T) {
	const catalogParam = "catalog[]"

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
			urlStr:    "https://www.vinted.sk/cataloga?search_text=",
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

func TestParseFilterParams(t *testing.T) {
	tests := []struct {
		name   string
		urlStr string
		want   FilterParams
	}{
		{
			name:   "all parameters",
			urlStr: "https://www.vinted.sk/catalog?search_text=&brand[]=1&catalog[]=79&catalog[]=80&color[]=1&material[]=1&size[]=1&size[]=2&status[]=1&status[]=2",
			want: FilterParams{
				brandIDs:    []int{1},
				catalogIDs:  []int{79, 80},
				colorIDs:    []int{1},
				materialIDs: []int{1},
				sizeIDs:     []int{1, 2},
				statusIDs:   []int{1, 2},
			},
		},
		{
			name:   "no parameters",
			urlStr: "https://www.vinted.sk/cataloga?search_text=",
			want: FilterParams{
				brandIDs:    nil,
				catalogIDs:  nil,
				colorIDs:    nil,
				materialIDs: nil,
				sizeIDs:     nil,
				statusIDs:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFilterParams(tt.urlStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFilterParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractPrices(t *testing.T) {
	tests := []struct {
		name      string
		urlStr    string
		paramName string
		want      float32
	}{
		{
			name:      "single price",
			urlStr:    "https://www.vinted.sk/catalog?search_text=&price_from=10.0",
			paramName: "price_from",
			want:      10.0,
		},
		{
			name:      "both prices",
			urlStr:    "https://www.vinted.sk/catalog?search_text=&price_from=10.0&price_to=20.0",
			paramName: "price_to",
			want:      20.0,
		},
		{
			name:      "no prices",
			urlStr:    "https://www.vinted.sk/catalog?search_text=",
			paramName: "price_from",
			want:      0.0,
		},
		{
			name:      "multiple prices - error",
			urlStr:    "https://www.vinted.sk/catalog?search_text=&price_from=10.0&price_from=20.0",
			paramName: "price_from",
			want:      0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractPrices(tt.urlStr, tt.paramName); got != tt.want {
				t.Errorf("extractPrices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePrices(t *testing.T) {
	tests := []struct {
		name   string
		urlStr string
		want   PriceParams
	}{
		{
			name:   "both prices",
			urlStr: "https://www.vinted.sk/catalog?search_text=&price_from=10.0&price_to=20.0",
			want: PriceParams{
				priceFrom: 10.0,
				priceTo:   20.0,
			},
		},
		{
			name:   "no prices",
			urlStr: "https://www.vinted.sk/catalog?search_text=",
			want: PriceParams{
				priceFrom: 0.0,
				priceTo:   0.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParsePrices(tt.urlStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePrices() = %v, want %v", got, tt.want)
			}
		})
	}
}
