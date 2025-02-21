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
		{
			name: "real example",
			urlStr: "https://www.vinted.sk/catalog?search_text=&catalog[]=79&price_from=2.1&price_to=2.5&currency=EUR&color_ids[]=10&color_ids[]=16&color_ids[]=28&size_ids[]=209&size_ids[]=210&brand_ids[]=90804&brand_ids[]=2319&brand_ids[]=255056&brand_ids[]=17161&brand_ids[]=319730&search_id=20007005180&order=newest_first&time=1740144460",
			paramName: "color_ids[]",
			want: []int{10, 16, 28},
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
				BrandIDs:    []int{1},
				CatalogIDs:  []int{79, 80},
				ColorIDs:    []int{1},
				MaterialIDs: []int{1},
				SizeIDs:     []int{1, 2},
				StatusIDs:   []int{1, 2},
			},
		},
		{
			name:   "no parameters",
			urlStr: "https://www.vinted.sk/cataloga?search_text=",
			want: FilterParams{
				BrandIDs:    nil,
				CatalogIDs:  nil,
				ColorIDs:    nil,
				MaterialIDs: nil,
				SizeIDs:     nil,
				StatusIDs:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFilterParams(tt.urlStr); !reflect.DeepEqual(got, tt.want) {
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
		{
			name:      "negative price",
			urlStr:    "https://www.vinted.sk/catalog?search_text=&price_from=-10.0",
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
				PriceFrom: 10.0,
				PriceTo:   20.0,
			},
		},
		{
			name:   "no prices",
			urlStr: "https://www.vinted.sk/catalog?search_text=",
			want: PriceParams{
				PriceFrom: 0.0,
				PriceTo:   0.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePrices(tt.urlStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePrices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractMiscParams(t *testing.T) {
	tests := []struct {
		name      string
		urlStr    string
		paramName string
		want      string
	}{
		{
			name:      "search_text",
			urlStr:    "https://www.vinted.sk/catalog?search_text=hello",
			paramName: "search_text",
			want:      "hello",
		},
		{
			name:      "currency",
			urlStr:    "https://www.vinted.sk/catalog?currency=EUR",
			paramName: "currency",
			want:      "EUR",
		},
		{
			name:      "order",
			urlStr:    "https://www.vinted.sk/catalog?order=price_asc",
			paramName: "order",
			want:      "price_asc",
		},
		{
			name:      "empty param",
			urlStr:    "https://www.vinted.sk/catalog?search_text=",
			paramName: "search_text",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractMiscParams(tt.urlStr, tt.paramName); got != tt.want {
				t.Errorf("extractMiscParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseMiscParams(t *testing.T) {
	tests := []struct {
		name   string
		urlStr string
		want   MiscParams
	}{
		{
			name:   "all parameters",
			urlStr: "https://www.vinted.sk/catalog?search_text=hello&currency=EUR&order=price_asc",
			want: MiscParams{
				SearchText: "hello",
				Currency:   "EUR",
				Order:      "price_asc",
			},
		},
		{
			name:   "no parameters",
			urlStr: "https://www.vinted.sk/catalog?search_text=",
			want: MiscParams{
				SearchText: "",
				Currency:   "",
				Order:      "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseMiscParams(tt.urlStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMiscParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
