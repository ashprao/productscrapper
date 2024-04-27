package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func printSiteData(categoriesProducts CategoryProductsMap, retailerSchema SiteSchemaStruct) {
	for category, products := range categoriesProducts {
		fmt.Println(strings.ToTitle(string(category)) + " has " + strconv.Itoa(len(products)) + " products")
		// get all attribute keys from ProductDataMap

		var attributeKeys []string
		attributeKeys = append(attributeKeys, "URL")
		for productAttributes := range retailerSchema.CategoriesSchema[category].ProductElements {
			attributeKeys = append(attributeKeys, string(productAttributes))
		}

		for _, product := range products {
			for _, attributeKey := range attributeKeys {
				fmt.Printf("%10s: %v", attributeKey, "")
				for i, attributeValue := range product[ProductAttribute(attributeKey)] {
					if i == 0 {
						fmt.Printf("%v", attributeValue)
						continue
					}
					fmt.Printf("\n%10s  %s", "", string(attributeValue))

				}
				fmt.Println()
			}
			fmt.Println("------------")
		}
	}
}

func exportToExcel(retailerProducts CategoryProductsMap, retailerSchema SiteSchemaStruct) {
	// Create a new Excel file
	f := excelize.NewFile()

	// For each category, create a new sheet and add the products
	for category, products := range retailerProducts {
		// Create a new sheet for the category
		sheet := string(category)
		f.NewSheet(sheet)

		// Get all attribute keys from ProductDataMap
		var attributeKeys []string
		attributeKeys = append(attributeKeys, "URL")
		for productAttributes := range retailerSchema.CategoriesSchema[category].ProductElements {
			attributeKeys = append(attributeKeys, string(productAttributes))
		}

		// Write the header row
		for i, attributeKey := range attributeKeys {
			f.SetCellValue(sheet, string('A'+i)+"1", attributeKey)
		}

		// Write the product data
		row := 2
		for _, product := range products {
			maxAttributeRow := 1
			for i, attributeKey := range attributeKeys {
				attributeRow := row

				// Get the product attribute values
				if values, ok := product[ProductAttribute(attributeKey)]; ok {
					// Write each product attribute value to a new cell
					for _, attributeValue := range values {
						f.SetCellValue(sheet, string('A'+i)+strconv.Itoa(attributeRow), string(attributeValue))
						attributeRow++
					}
					if attributeRow > maxAttributeRow {
						maxAttributeRow = attributeRow
					}
				}
			}
			row = maxAttributeRow
		}
	}

	// Save the Excel file
	err := f.SaveAs(retailerSchema.Name + ".xlsx")
	if err != nil {
		log.Println("Error saving Excel file:", err)
	}
	log.Printf("Created '%s.xlsx' file\n", retailerSchema.Name)
	for category, products := range retailerProducts {
		log.Printf("Extracted %d products from category %s", len(products), category)
	}

}

// Controller
type Controller struct {
	model *Model
}

// NewController creates a new controller
func NewController(m *Model) *Controller {
	return &Controller{
		model: m,
	}
}

func (c *Controller) UpdateExcelFilePath(filePath string) {
	c.model.UpdateExcelFilePath(filePath)
}

// HandleFileSelect handles the file select event
// HandleScrape handles the Scrape button click
func (c *Controller) HandleScrape() error {
	if c.model.ExcelFilePath == "" {
		return fmt.Errorf("excel filepath not provided")
	}

	// Get the retailers schema from the Excel file
	retailersSchema, err := getSitesSchema(c.model.ExcelFilePath)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	for _, retailerSchema := range retailersSchema {
		retailerProducts := make(CategoryProductsMap)

		fmt.Println("Scraping", retailerSchema.Name)
		for category, categorySchema := range retailerSchema.CategoriesSchema {
			fmt.Println("Scraping", category)
			retailerProducts[category] = scrapeCategoryCatalogPage(categorySchema)
		}

		printSiteData(retailerProducts, retailerSchema)
		exportToExcel(retailerProducts, retailerSchema)
		// fmt.Println(retailerProducts)
	}

	return nil
}
