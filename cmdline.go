package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Command Line View

func commandLineView(controller *Controller) {
	var rootCmd = &cobra.Command{
		Use:   "productscraper",
		Short: "A web scraper for product information",
		Long:  `A web scraper for product information that accepts an Excel file path containing websites locations and scrapping schemas and initiates a scraping operation.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Handle command line arguments here
			// For example:
			excelFilePath, err := cmd.Flags().GetString("excelfilepath")
			if err != nil {
				log.Fatal(err)
			}

			controller.UpdateExcelFilePath(excelFilePath)

			log.Println("Excel file path:", excelFilePath)
			if err := controller.HandleScrape(); err != nil {
				log.Fatal(err)
			}
		},
	}

	log.Println("In commandline mode...")
	rootCmd.Flags().StringP("excelfilepath", "e", "", "Excel file path")
	rootCmd.MarkFlagRequired("excelfilepath")
	rootCmd.Version = "1.0.0"

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
