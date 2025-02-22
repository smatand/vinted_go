// todos: the file should contain functions to convert structure type vinted to vintedapi url
// currency_seller ~ add to user on DC ability to filter among sellers from poland, czechia and svk
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
	refreshTokenWebName   = "refresh_token_web"
	restAPIEndpoint       = "https://vinted.sk/api/v2/catalog/"
	pageNth               = "1"
	itemsPerPage          = "16"
	maxExponentialWait    = 14400 // 4 hours
)

// Keeps the accessTokenWeb for authentification in API and refreshTokenWeb for refreshing the accessTokenWeb after it expires.
// todo: the use of refreshTokenWeb is not yet tested/implemented, I might consider deleting it
type cookies struct {
	accessTokenWeb  string
	refreshTokenWeb string
}

// Structure of response from Vinted API. The struct contains only the necessary fields.
type VintedItemsResp struct {
	Items []struct {
		ID         int              `json:"id"`
		Title      string           `json:"title"`
		Price      VintedPrice      `json:"price"`
		BrandTitle string           `json:"brand_title"`
		Url        string           `json:"url"`
		Conversion VintedConversion `json:"conversion"`
		Photo      VintedPhoto      `json:"photo"`
	} `json:"items"`
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
func constructVintedAPIRequest(v vinted.Vinted) string {
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
	toRet += constructParamString("brand_ids", f.BrandIDs)
	toRet += constructParamString("catalog_ids", f.CatalogIDs)
	toRet += constructParamString("color_ids", f.ColorIDs)
	toRet += constructParamString("material_ids", f.MaterialIDs)
	toRet += constructParamString("size_ids", f.SizeIDs)
	toRet += constructParamString("status_ids", f.StatusIDs)

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

// From a given cookie string, the function returns the value after '=' sign and before the first ';'.
func getToken(cookie string) (string, error) {
	if !strings.Contains(cookie, "=") {
		return "", fmt.Errorf("cookie does not contain = operator")
	}

	token := strings.SplitN(cookie, "=", 2)[1]
	randomizedBytes := strings.SplitN(token, ";", 2)[0]

	// let's return error if there's another = operator in token
	if strings.Contains(randomizedBytes, "=") {
		return "", fmt.Errorf("token contains more than 1 = operator")
	}

	return randomizedBytes, nil
}

// Fetches cookie access_token_web and refresh_token_web from the given host.
func fetchVintedCookies(host string) (cookies, error) {
	var cookieData cookies
	client := http.Client{}

	req, err := http.NewRequest("GET", host, nil)
	if err != nil {
		return cookies{}, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return cookies{}, fmt.Errorf("could not create client: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		waitExponential()

		cookieData, err = fetchVintedCookies(host)
		if err != nil {
			return cookies{}, err
		}
	} else {
		sessionCookies := resp.Header["Set-Cookie"]
		for _, cookie := range sessionCookies {
			if strings.Contains(cookie, accessTokenCookieName+"=") {
				// access_token_web contains some randomized data + basic info about Domain, Max-Age divided by semicolons, we only want the 1st one
				cookieData.accessTokenWeb, err = getToken(cookie)
			} else if strings.Contains(cookie, refreshTokenWebName+"=") {
				// same goes for refresh_token_web, the bytes before first ; are the token
				cookieData.refreshTokenWeb, err = getToken(cookie)
			}

			if err != nil {
				return cookies{}, fmt.Errorf("could not get token: %v", err)
			}
		}

		retryCountExp = 0
	}

	return cookieData, nil
}

func extractHost(URL string) string {
	return strings.Split(URL, "/api")[0]
}

// Retrieves items from Vinted API based on the given parameters from vinted.Vinted structure
// The data are json unmarshalled into VintedItemsResp structure.
func GetVintedItems(v vinted.Vinted) (VintedItemsResp, error) {
	requestURL := constructVintedAPIRequest(v)

	host := extractHost(requestURL)

	// todo: do not fetch it always
	cookies, err := fetchVintedCookies(host)
	if err != nil {
		return VintedItemsResp{}, err
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return VintedItemsResp{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Cookie", accessTokenCookieName+"="+cookies.accessTokenWeb)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return VintedItemsResp{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return VintedItemsResp{}, fmt.Errorf("status code from url %v : %v", requestURL, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return VintedItemsResp{}, fmt.Errorf("failed to read response body: %v", err)
	}

	var vintedResp VintedItemsResp
	err = json.Unmarshal(body, &vintedResp)
	if err != nil {
		return VintedItemsResp{}, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return vintedResp, nil
}
