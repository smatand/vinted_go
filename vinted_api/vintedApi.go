// todos: the file should contain functions to convert structure type vinted to vintedapi url
// currency_seller ~ add to user on DC ability to filter among sellers from poland, czechia and svk
package vintedApi

import (
	"fmt"
	"net/http"
	"io"
	"strconv"

	"github.com/smatand/vinted_go/vinted"
)

// global rest api url
const REST_API_ENDPOINT = "https://vinted.com/api/v2/catalog/"
const PAGE = "1"
const ITEMS_PER_PAGE = "16"

// constructs rest api url with default of 1st page and 16 items per page as those specs are not necessary for our purpose
func constructVintedAPIRequest(v vinted.Vinted) string {
	baseURL := REST_API_ENDPOINT + "items?page=" + PAGE + "&per_page=" + ITEMS_PER_PAGE

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
