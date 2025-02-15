// vinted.go source file contains abstraction over the Vinted API parameters
package vinted

import (
	"log"
	"net/url"
	"strconv"
	"strings"
)

type PriceParams struct {
	priceFrom float32
	priceTo   float32
}

type FilterParams struct {
	brandIDs    []int
	catalogIDs  []int
	colorIDs    []int
	materialIDs []int
	sizeIDs     []int
	statusIDs   []int
}

func ParsePrices(urlStr string) PriceParams {
	priceFrom := extractPrices(urlStr, "price_from")
	priceTo := extractPrices(urlStr, "price_to")

	return PriceParams{
		priceFrom: priceFrom,
		priceTo:   priceTo,
	}
}

func ParseFilterParams(urlStr string) FilterParams {
	brandIDs := extractIDs(urlStr, "brand[]")
	catalogIDs := extractIDs(urlStr, "catalog[]")
	colorIDs := extractIDs(urlStr, "color[]")
	materialIDs := extractIDs(urlStr, "material[]")
	sizeIDs := extractIDs(urlStr, "size[]")
	statusIDs := extractIDs(urlStr, "status[]")

	return FilterParams{
		brandIDs:    brandIDs,
		catalogIDs:  catalogIDs,
		colorIDs:    colorIDs,
		materialIDs: materialIDs,
		sizeIDs:     sizeIDs,
		statusIDs:   statusIDs,
	}
}

func extractPrices(urlStr string, paramName string) float32 {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("Could not parse url %v: %v", urlStr, err)
		return 0.0
	}

	queryParams, err := url.ParseQuery(parsedUrl.RawQuery)
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

	return float32(result)
}

// extracts IDs from string like 'catalog[]=79&catalog[]=80' -> [79, 80]
func extractIDs(urlStr string, paramName string) []int {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("Could not parse url %v: %v", urlStr, err)
		return nil
	}

	if paramName == "catalog[]" {
		// could be url like /catalog/2050-clothing, handle this separately
		pathSegments := strings.Split(parsedUrl.Path, "/")
		if len(pathSegments) > 2 {
			catalogSegment := pathSegments[2]
			catalogIDStr := strings.Split(catalogSegment, "-")[0]
			catalogID, err := strconv.Atoi(catalogIDStr)
			if err == nil {
				return []int{catalogID}
			}

			// unhandled cases for /catalog/...
			return nil
		}
	}

	queryParams, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return nil
	}

	catalogIDs := queryParams[paramName]
	lenCatalogIDs := len(catalogIDs)
	if lenCatalogIDs == 0 {
		return nil
	}

	result := make([]int, 0, lenCatalogIDs)
	// if there's invalid catalogID like catalog[]=281a, ignore it
	for _, catalogID := range catalogIDs {
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
