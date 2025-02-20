// todos: the file should contain functions to convert structure type vinted to vintedapi url
// currency_seller ~ add to user on DC ability to filter among sellers from poland, czechia and svk
package vintedApi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/smatand/vinted_go/vinted"
)

// for exponential backoff ~ waitExponential()
var retryCount int

// global rest api url
const accessTokenCookieName = "access_token_web"
const refreshTokenWebName = "refresh_token_web"
const restAPIEndpoint = "https://vinted.com/api/v2/catalog/"
const pageNth = "1"
const itemsPerPage = "16"
const maxExponentialWait = 14400 // 4 hours

type cookies struct {
	accessTokenWeb  string
	refreshTokenWeb string
}

// constructs rest api url with default of 1st page and 16 items per page as those specs are not necessary for our purpose
func constructVintedAPIRequest(v vinted.Vinted) string {
	baseURL := restAPIEndpoint + "items?page=" + pageNth + "&per_page=" + itemsPerPage

	baseURL += constructPriceParams(v.PriceParams)
	baseURL += constructFilterParams(v.FilterParams)
	baseURL += constructMiscParams(v.MiscParams)

	return baseURL
}

func constructPriceParams(p vinted.PriceParams) string {
	toRet := ""
	if p.PriceFrom != 0.0 {
		toRet += "&price_from=" + fmt.Sprintf("%.2f", p.PriceFrom)
	}
	if p.PriceTo != 0.0 {
		toRet += "&price_to=" + fmt.Sprintf("%.2f", p.PriceTo)
	}

	return toRet
}

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

func constructParamString(paramName string, ids []int) string {
	toRet := ""
	for _, id := range ids {
		toRet += "&" + paramName + "=" + strconv.Itoa(id)
	}

	return toRet
}

// wait exponential time, maximum is 4 hours
func waitExponential() {
	delaySecs := 1 << retryCount
	if delaySecs > maxExponentialWait {
		delaySecs = maxExponentialWait
	}

	time.Sleep(time.Duration(delaySecs) * time.Second)

	retryCount++
}

// from a given cookie string, return the value after '=' sign and before the first ';'
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

func FetchVintedCookies(host string) (cookies, error) {
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

		cookieData, err = FetchVintedCookies(host)
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

		retryCount = 0
	}

	return cookieData, nil
}