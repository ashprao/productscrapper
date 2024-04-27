package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

// SiteSchemaStruct contains the data for a retailer
type SiteSchemaStruct struct {
	Name             string
	CategoriesSchema map[Category]CategorySchemaStruct
}

type CategorySchemaStruct struct {
	URL                    string
	ProductPageElement     string
	NextProductPageElement string
	ProductElements        ProductElementsSchemaMap
}

// ProductElementsSchemaMap contains the HTML elements for a product
type ProductElementsSchemaMap map[string]string

func (m ProductElementsSchemaMap) findKey(value string) (string, error) {
	// find the key associated with the value.
	for k, v := range m {
		if v == value {
			return k, nil
		}
	}
	return "", fmt.Errorf("key not found")
}

// getSitesSchema returns the data for a retailer
func getSitesSchema(excelFile string) (sitesSchemaList []SiteSchemaStruct, err error) {

	// Open the excel file
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		log.Printf("Error opening file: %s\n; %v", excelFile, err)
		return sitesSchemaList, fmt.Errorf("error opening file: %s; %w", excelFile, err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Get all sheet names
	sheetNames := []string{}
	for _, name := range f.GetSheetMap() {
		sheetNames = append(sheetNames, name)
	}

	// Slice to store RetailerData for each sheet
	sitesSchemaList = make([]SiteSchemaStruct, len(sheetNames))

	for i, sheetName := range sheetNames {
		sitesSchemaList[i].Name = sheetName

		// Get rows from the sheet
		rows, err := f.GetRows(sheetName)
		if err != nil {
			log.Printf("Error getting data from sheet %s: %v", sheetName, err)
			continue
		}

		// Extract retailer data from the sheet
		sitesSchemaList[i].CategoriesSchema = extractSiteData(rows)
	}

	// Print the extracted retailer data
	for _, data := range sitesSchemaList {
		fmt.Println("Name:", data.Name)
		for category, categoryData := range data.CategoriesSchema {
			fmt.Println("Category:", category)
			fmt.Println("URL:", categoryData.URL)
			fmt.Println("Product Page Element:", categoryData.ProductPageElement)
			fmt.Println("Next Product Page Element:", categoryData.NextProductPageElement)
			fmt.Println("Product Elements:", categoryData.ProductElements)
			fmt.Println("-------")
		}
	}
	return
}

// extractSiteData extracts retailer data from a sheet
func extractSiteData(rows [][]string) map[Category]CategorySchemaStruct {
	var categoriesSchema map[Category]CategorySchemaStruct = make(map[Category]CategorySchemaStruct)

	type rowData struct {
		Category     string
		CategoryData CategorySchemaStruct
	}

	// Define expected column names
	category := "Category"
	urlCol := "URL"
	productPageCol := "Product Page"
	nextProductCol := "Next Catalog"

	// Loop through rows starting from the second row (assuming headers are in the first row)
	for i := 1; i < len(rows); i++ {
		// Loop through rows starting from the second row (assuming headers are in the first row)

		rowData := rowData{}
		rowData.CategoryData.ProductElements = make(ProductElementsSchemaMap)

		for j, cell := range rows[i] {
			switch j {
			// case 0:
			// 	// Sheet name is the retailer name
			// 	data.Name = cell
			case findColIndex(rows[0], category):
				rowData.Category = cell
			case findColIndex(rows[0], urlCol):
				rowData.CategoryData.URL = cell
			case findColIndex(rows[0], productPageCol):
				rowData.CategoryData.ProductPageElement = cell
			case findColIndex(rows[0], nextProductCol):
				rowData.CategoryData.NextProductPageElement = cell
			default:
				// Other columns are product attributes
				rowData.CategoryData.ProductElements[rows[0][j]] = cell
			}
		}

		categoriesSchema[Category(rowData.Category)] = CategorySchemaStruct{
			URL:                    rowData.CategoryData.URL,
			ProductPageElement:     rowData.CategoryData.ProductPageElement,
			NextProductPageElement: rowData.CategoryData.NextProductPageElement,
			ProductElements:        rowData.CategoryData.ProductElements,
		}
	}

	return categoriesSchema
}
