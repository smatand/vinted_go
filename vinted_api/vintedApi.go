package vintedApi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/smatand/vinted_go/vinted"
)

// For exponential backoff ~ waitExponential().
var retryCountExp int

// Names of tokens, endpoint and page number
const (
	accessTokenCookieName = "access_token_web"
	RefreshTokenWebName   = "refresh_token_web"
	restAPIEndpoint       = "https://vinted.sk/api/v2/catalog/"
	pageNth               = "1"
	itemsPerPage          = "16"
	maxExponentialWait    = 14400 // 4 hours
	cookiesFilePath       = "cookies.json"
)

// Keeps the AccessTokenWeb for authentification in API and RefreshTokenWeb for refreshing the AccessTokenWeb after it expires.
// todo: the use of RefreshTokenWeb is not yet tested/implemented
type cookies struct {
	AccessTokenWeb  string `json:"access_token_web"`
	RefreshTokenWeb string `json:"refresh_token_web"`
}

type VintedItemsResp struct {
	Items []VintedItemResp `json:"items"`
}

// Structure of response from Vinted API. The struct contains only the necessary fields.
type VintedItemResp struct {
	ID         int              `json:"id"`
	Title      string           `json:"title"`
	Price      VintedPrice      `json:"price"`
	BrandTitle string           `json:"brand_title"`
	Url        string           `json:"url"`
	Conversion VintedConversion `json:"conversion"`
	Photo      VintedPhoto      `json:"photo"`
}

// Structure of json price in response from Vinted API.
type VintedPrice struct {
	Amount string `json:"amount"`
}

// Structure which helps to decide the country of the seller.
type VintedConversion struct {
	SellerCurrency string `json:"seller_currency"`
}

// Structure to hold thumbnail of item photo.
type VintedPhoto struct {
	Url string `json:"url"`
}

// Constructs rest API URL which by default retrieves 1st page with 16 items. The function then adds
// other parameters to the URL based on the vinted.Vinted structure.
// The returned value can be pasted to the URL for the API request.
func ConstructVintedAPIRequest(v vinted.Vinted) string {
	baseURL := restAPIEndpoint + "items?page=" + pageNth + "&per_page=" + itemsPerPage

	baseURL += constructPriceParams(v.PriceParams)
	baseURL += constructFilterParams(v.FilterParams)
	baseURL += constructMiscParams(v.MiscParams)

	return baseURL
}

// Constructs price parameters for the API URL.
// The returned value can be pasted to the URL for the API request.
func constructPriceParams(p vinted.PriceParams) string {
	var toRet string

	if p.PriceFrom != 0.0 {
		toRet += "&price_from=" + fmt.Sprintf("%.2f", p.PriceFrom)
	}
	if p.PriceTo != 0.0 {
		toRet += "&price_to=" + fmt.Sprintf("%.2f", p.PriceTo)
	}

	return toRet
}

// Constructs filter parameters (cathegorical parameters) for the API URL.
// The returned value can be pasted to the URL for the API request.
func constructFilterParams(f vinted.FilterParams) string {
	toRet := ""
	toRet += constructParamString("brand_ids[]", f.BrandIDs)
	toRet += constructParamString("catalog_ids[]", f.CatalogIDs)
	toRet += constructParamString("color_ids[]", f.ColorIDs)
	toRet += constructParamString("material_ids[]", f.MaterialIDs)
	toRet += constructParamString("size_ids[]", f.SizeIDs)
	toRet += constructParamString("status_ids[]", f.StatusIDs)

	return toRet
}

// Constructs miscallenous parameters ~ search_text or currency.
// The returned value can be pasted to the URL for the API request.
func constructMiscParams(m vinted.MiscParams) string {
	toRet := ""
	if m.SearchText != "" {
		toRet += "&search_text=" + m.SearchText
	}
	if m.Currency != "" {
		toRet += "&currency=" + m.Currency
	}
	if m.Order != "" {
		toRet += "&order=" + m.Order
	}

	return toRet
}

// Constructs a string from a given parameter name and a slice of integers.
// The returned value can be pasted to the URL for the API request.
func constructParamString(paramName string, ids []int) string {
	toRet := ""
	for _, id := range ids {
		toRet += "&" + paramName + "=" + strconv.Itoa(id)
	}

	return toRet
}

// Waits exponential time, maximum is 4 hours.
func waitExponential() {
	delaySecs := 1 << retryCountExp
	if delaySecs > maxExponentialWait {
		delaySecs = maxExponentialWait
	}

	time.Sleep(time.Duration(delaySecs) * time.Second)

	retryCountExp++
}

// Fetches cookie access_token_web and refresh_token_web from the given host.
func fetchVintedCookies(host string) (*cookies, error) {
	cookieData := &cookies{}

	client := http.Client{}
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			waitExponential()
		}

		req, err := http.NewRequest("GET", host, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("could not create client: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			continue
		}

		for _, cookie := range resp.Cookies() {
			switch cookie.Name {
			case accessTokenCookieName:
				cookieData.AccessTokenWeb = cookie.Value
			case RefreshTokenWebName:
				cookieData.RefreshTokenWeb = cookie.Value
			}
		}

		// Reset the exponentialWait() counter when both cookies are retrieved
		if cookieData.AccessTokenWeb != "" && cookieData.RefreshTokenWeb != "" {
			retryCountExp = 0

			return cookieData, nil
		}
	}

	return nil, fmt.Errorf("could not retrieve cookies")
}

func extractHost(URL string) string {
	return strings.Split(URL, "/api")[0]
}

// Retrieves items from Vinted API based on the given parameters from vinted.Vinted structure
// The data are json unmarshalled into VintedItemsResp structure.
func GetVintedItems(requestURL string) (*VintedItemsResp, error) {
	host := extractHost(requestURL)

	cookies, err := fetchVintedCookies(host)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Cookie", accessTokenCookieName+"="+cookies.AccessTokenWeb)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code from url %v : %v", requestURL, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	vintedResp := &VintedItemsResp{}
	err = json.Unmarshal(body, &vintedResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return vintedResp, nil
}
