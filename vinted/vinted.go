// vinted.go source file contains abstraction over the Vinted API parameters
package vinted

import (
	"log"
	"net/url"
	"strconv"
	"strings"
)

type Vinted struct {
	PriceParams
	FilterParams
	MiscParams
}

// todo: consider string type of elems
type PriceParams struct {
	PriceFrom float32
	PriceTo   float32
}

type FilterParams struct {
	BrandIDs    []int
	CatalogIDs  []int
	ColorIDs    []int
	MaterialIDs []int
	SizeIDs     []int
	StatusIDs   []int
}

type MiscParams struct {
	SearchText string
	Currency   string
	Order      string
}

func ParsePrices(urlStr string) PriceParams {
	PriceFrom := extractPrices(urlStr, "price_from")
	PriceTo := extractPrices(urlStr, "price_to")

	return PriceParams{
		PriceFrom: PriceFrom,
		PriceTo:   PriceTo,
	}
}

func ParseFilterParams(urlStr string) FilterParams {
	BrandIDs := extractIDs(urlStr, "brand[]")
	CatalogIDs := extractIDs(urlStr, "catalog[]")
	ColorIDs := extractIDs(urlStr, "color[]")
	MaterialIDs := extractIDs(urlStr, "material[]")
	SizeIDs := extractIDs(urlStr, "size[]")
	StatusIDs := extractIDs(urlStr, "status[]")

	return FilterParams{
		BrandIDs:    BrandIDs,
		CatalogIDs:  CatalogIDs,
		ColorIDs:    ColorIDs,
		MaterialIDs: MaterialIDs,
		SizeIDs:     SizeIDs,
		StatusIDs:   StatusIDs,
	}
}

func ParseMiscParams(urlStr string) MiscParams {
	SearchText := extractMiscParams(urlStr, "search_text")
	Currency := extractMiscParams(urlStr, "Currency")
	Order := extractMiscParams(urlStr, "Order")

	return MiscParams{
		SearchText: SearchText,
		Currency:   Currency,
		Order:      Order,
	}
}

func parseQueryParams(urlStr string) (url.Values, error) {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("Could not parse url %v: %v", urlStr, err)
		return nil, err
	}

	queryParams, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return nil, err
	}

	return queryParams, nil
}

func extractPrices(urlStr string, paramName string) float32 {
	queryParams, err := parseQueryParams(urlStr)
	if err != nil {
		return 0.0
	}

	numericValues := queryParams[paramName]
	if len(numericValues) != 1 {
		log.Printf("Expected only one value of %v, got %v", paramName, numericValues)
		return 0.0
	}

	result, err := strconv.ParseFloat(queryParams[paramName][0], 32)
	if err != nil {
		return 0.0
	}

	if result < 0 {
		log.Printf("Expected positive value of %v, got %v, continuing with 0.0", paramName, result)
		return 0.0
	}

	return float32(result)
}

// extracts IDs from string like 'catalog[]=79&catalog[]=80' -> [79, 80]
func extractIDs(urlStr string, paramName string) []int {
	if paramName == "catalog[]" {
		// could be url like /catalog/2050-clothing, handle this separately
		pathSegments := strings.Split(urlStr, "/")

		// https://www.../catalog/2050-clothing -> 5 segments, 4th is catalogID
		if len(pathSegments) > 4 {
			catalogSegment := pathSegments[4]

			// "2050-clothing"[0] -> 2050
			catalogIDStr := strings.Split(catalogSegment, "-")[0]
			catalogID, err := strconv.Atoi(catalogIDStr)
			if err == nil {
				return []int{catalogID}
			}

			// unhandled cases for /catalog/...
			return nil
		}
	}

	queryParams, err := parseQueryParams(urlStr)
	if err != nil {
		return nil
	}

	CatalogIDs := queryParams[paramName]
	lenCatalogIDs := len(CatalogIDs)
	if lenCatalogIDs == 0 {
		return nil
	}

	result := make([]int, 0, lenCatalogIDs)
	// if there's invalid catalogID like catalog[]=281a, ignore it
	for _, catalogID := range CatalogIDs {
		id, err := strconv.Atoi(catalogID)
		if err != nil {
			continue
		}

		result = append(result, id)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func extractMiscParams(urlStr string, paramName string) string {
	queryParams, err := parseQueryParams(urlStr)
	if err != nil {
		return ""
	}

	values := queryParams[paramName]
	if len(values) != 1 {
		return ""
	}

	return values[0]
}
