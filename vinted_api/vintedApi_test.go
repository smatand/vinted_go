package vintedApi

import (
	"testing"

	"github.com/smatand/vinted_go/vinted"
)

func TestConstructVintedAPIRequest(t *testing.T) {
	const baseURL = restAPIEndpoint + "items?page=" + pageNth + "&per_page=" + itemsPerPage
	tests := []struct {
		name   string
		vinted vinted.Vinted
		want   string
	}{
		{
			name:   "empty vinted",
			vinted: vinted.Vinted{},
			want:   baseURL,
		},
		{
			name: "order",
			vinted: vinted.Vinted{
				MiscParams: vinted.MiscParams{
					Order: "newest_first",
				},
			},
			want: baseURL + "&order=newest_first",
		},
		{
			name: "price",
			vinted: vinted.Vinted{
				PriceParams: vinted.PriceParams{
					PriceFrom: 10.0,
					PriceTo:   20.0,
				},
			},
			want: baseURL + "&price_from=10.00&price_to=20.00",
		},
		{
			name: "filter",
			vinted: vinted.Vinted{
				FilterParams: vinted.FilterParams{
					BrandIDs: []int{1, 2, 3},
				},
			},
			want: baseURL + "&brand_ids=1&brand_ids=2&brand_ids=3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := constructVintedAPIRequest(tt.vinted); got != tt.want {
				t.Errorf("ConstructVintedAPIRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
