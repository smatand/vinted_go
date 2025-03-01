package vintedApi

import (
	"testing"

	"github.com/smatand/vinted_go/vinted"
)

func TestFetchCookie(t *testing.T) {
	const testURL = "https://www.vinted.sk"
	tests := []struct {
		name string
	}{
		{
			name: "fetch cookie",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fetchVintedCookies(testURL)
			if err != nil {
				t.Logf("API is not available: %v", err)
			}
		})
	}
}

func TestGetVintedItems(t *testing.T) {
	// just test whether it works, but use parseParam fro mvinted to get the params
	params := vinted.Vinted{}
	params.ParseParams("https://www.vinted.sk/catalog?search_text=&catalog[]=1806&price_from=0.01&price_to=9999.0&currency=EUR&search_id=21069023846&order=newest_first")

	_, err := GetVintedItems(params)
	if err != nil {
		t.Errorf("error: %v", err)
	}
}
