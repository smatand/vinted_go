// vinted.go source file contains abstraction over the Vinted API parameters
package vinted

import (
	"log"
	"net/url"
	"strconv"
	"strings"
)

// extracts IDs from string like 'catalog[]=79&catalog[]=80' -> [79, 80]
func extractCatalogIDs(urlStr string) []int {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("Could not parse url %v: %e", urlStr, err)
		return nil
	}

	// could be url like /catalog/2050-clothing, handle this separately
	pathSegments := strings.Split(parsedUrl.Path, "/")
	if len(pathSegments) > 2 {
		catalogSegment := pathSegments[2]
		catalogIDStr := strings.Split(catalogSegment, "-")[0]
		catalogID, err := strconv.Atoi(catalogIDStr)
		if err == nil {
			return []int{catalogID}
		}

		// unhandled cases
		return nil
	}

	queryParams, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return nil
	}

	catalogIDs := queryParams["catalog[]"]
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

	return result
}
