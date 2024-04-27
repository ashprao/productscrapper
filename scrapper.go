package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

type ProductAttribute string
type ProductAttributeValue string
type Category string

// ProductDataMap stores the product values of all required product attributes
type ProductDataMap map[ProductAttribute][]ProductAttributeValue

// CategoryProductsMap stores the product data per category
type CategoryProductsMap map[Category][]ProductDataMap

// Define multiple disallowed robots tags (modify as needed)
var disallowedTags []string = []string{"noindex", "nofollow"}

// findColIndex finds the index of a column by name in the header row
func findColIndex(header []string, name string) int {
	for i, colName := range header {
		if colName == name {
			return i
		}
	}
	return -1
}

func scrapeCategoryCatalogPage(categoryData CategorySchemaStruct) []ProductDataMap {
	products := []ProductDataMap{}

	// Instantiate default collector
	catalogCollector := colly.NewCollector()
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
	catalogCollector.UserAgent = userAgent

	catalogCollector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting product catalog page:", r.URL)
	})

	catalogCollector.OnError(func(r *colly.Response, err error) {
		log.Println("Error scraping product catalog page:", err)
	})

	allowed := true
	catalogCollector.OnResponse(func(r *colly.Response) {
		log.Println("Received response from catalog page:", r.Request.URL)
		robotsTag := r.Headers.Get("X-Robots-Tag")
		if robotsTag != "" {
			for _, disallowedTag := range disallowedTags {
				if strings.Contains(strings.ToLower(robotsTag), disallowedTag) {
					fmt.Println("Found disallowed robots tag:", robotsTag)
					// Ignore the response and cancel visit
					allowed = false
					return
				}
			}
		}
		// Proceed with normal processing if no disallowed tag found
		log.Println("Processing response from:", r.Request.URL)
	})

	catalogCollector.OnHTML(categoryData.ProductPageElement, func(h *colly.HTMLElement) {
		if !allowed {
			return
		}

		visitLink := h.Attr("href")
		log.Println("Found link:", visitLink)
		if !strings.HasPrefix(visitLink, "http") {
			visitLink = h.Request.AbsoluteURL(visitLink)
			if !strings.HasPrefix(visitLink, "http") {
				// HTML may have the "base href" set to "/", so we need to account for that
				visitLink = h.Request.URL.Scheme + "://" + h.Request.URL.Host + visitLink
			}
		}
		products = append(products, scrapeProductLink(visitLink, categoryData.ProductElements))
	})

	catalogCollector.OnHTML(categoryData.NextProductPageElement, func(h *colly.HTMLElement) {
		if !allowed {
			return
		}

		visitLink := h.Attr("href")
		log.Println("\nFound Next link:", visitLink)
		if !strings.HasPrefix(visitLink, "http") {
			visitLink = h.Request.AbsoluteURL(visitLink)
			if !strings.HasPrefix(visitLink, "http") {
				// HTML may have the "base href" set to "/", so we need to account for that
				visitLink = h.Request.URL.Scheme + "://" + h.Request.URL.Host + visitLink
			}
		}
		catalogCollector.Visit(visitLink)
	})

	catalogCollector.OnScraped(func(r *colly.Response) {
		log.Println("Scraped product catalog page:", r.Request.URL)
	})

	catalogCollector.Visit(categoryData.URL)
	return products
}

// Function to check if a URL is absolute
func isAbsoluteURL(urlString string) bool {
	parsedUrl, err := url.Parse(urlString)
	return err == nil && parsedUrl.Scheme != ""
}

// Function to clean and construct absolute URL
func cleanUrl(baseUrl *url.URL, relUrl string) (string, error) {
	// Handle absolute URLs directly
	if isAbsoluteURL(relUrl) {
		return relUrl, nil
	}

	// Parse the relative URL
	rel, err := url.Parse(relUrl)
	if err != nil {
		return "", err
	}

	// Use the url package's ResolveReference function to resolve the relative URL against the base URL
	resolvedUrl := baseUrl.ResolveReference(rel)

	return resolvedUrl.String(), nil
}

// removeQueryPart is a function to remove the query part from a urlString
func removeQueryPart(urlString string) (string, error) {
	url, err := url.Parse(urlString)
	if err != nil {
		log.Println("Error encountered in parsing URL")
		return "", err
	}

	url.RawQuery = ""

	return url.String(), nil
}

func scrapeProductLink(productLink string, productElements ProductElementsSchemaMap) (product ProductDataMap) {

	localProductElements := productElements

	product = make(ProductDataMap)

	// Create a collector for scraping product pages
	collector := colly.NewCollector()
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/000000000 Safari/537.36"
	collector.UserAgent = userAgent

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting product page:", r.URL)
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Error scraping product page:", err)
	})

	allowed := true
	// Check for "X-Robots-Tag" header in OnResponse
	collector.OnResponse(func(r *colly.Response) {
		log.Println("Received response from product page:", r.Request.URL)
		robotsTag := r.Headers.Get("X-Robots-Tag")
		if robotsTag != "" {
			for _, disallowedTag := range disallowedTags {
				if strings.Contains(strings.ToLower(robotsTag), disallowedTag) {
					log.Println("Found disallowed robots tag:", robotsTag)
					// Ignore the response and cancel visit
					allowed = false
					return
				}
			}
		}
		// Proceed with normal processing if no disallowed tag found
		log.Println("Processing response from:", r.Request.URL)
	})

	for _, productElement := range localProductElements {
		localProductElement := productElement
		collector.OnHTML(localProductElement, func(h *colly.HTMLElement) {

			if !allowed {
				return
			}

			// log.Printf("Found product element: %s\n", productElement)
			attribute, err := localProductElements.findKey(localProductElement)
			localProductAttribute := ProductAttribute(attribute)
			if err != nil {
				log.Println("Error: Could not find key for product element:", productElement)
				return
			}

			// if productElement represents an image then get the src attribute
			// other wise get the text content
			switch {
			case strings.Contains(strings.ToLower(localProductElement), "img"):
				if ProductAttributeValue(h.Attr("src")) != "" {
					if !isAbsoluteURL(h.Attr("src")) {
						// get base url
						// Get base URL and handle non-absolute URLs
						baseUrl, err := cleanUrl(h.Request.URL, h.Attr("src"))
						if err != nil || baseUrl == "" {
							log.Println("Error: Could not construct absolute URL")
							return
						}

						baseUrl, err = removeQueryPart(baseUrl)
						if err != nil {
							log.Println("Error: Could not remove query part from URL")
							return
						}
						product[localProductAttribute] = append(product[localProductAttribute], ProductAttributeValue(baseUrl))
					} else {
						baseUrl, err := removeQueryPart(h.Attr("src"))
						if err != nil {
							log.Println("Error: Could not remove query part from URL")
							return
						}
						product[localProductAttribute] = append(product[localProductAttribute], ProductAttributeValue(baseUrl))
					}
				}
			default:
				// log.Printf("Found product element value: %s for product attribute: %s\n", h.Text, localProductAttribute)
				if ProductAttributeValue(h.Text) != "" {
					product[localProductAttribute] = append(product[localProductAttribute], ProductAttributeValue(h.Text))
				}
			}
			// log.Println("Found product element:", productElement)
		})
	}

	collector.OnScraped(func(r *colly.Response) {
		log.Println("Scraped product page:", r.Request.URL)
		product["URL"] = append(product["URL"], ProductAttributeValue(r.Request.URL.String()))
	})

	collector.Visit(productLink)

	return
}
