# Product Scraper

This is a web scraper built with Golang and Colly that scrapes product data from ecommerce websites. It accepts an Excel file containing the scraping schema and outputs the results back to Excel.

## Usage

`go run main.go -f <excel_file>`

The Excel file should contain the schema with a sheet per retailer. Each sheet has columns for category, URL, product page element, product attributes etc. See the example file in the schemas directory.

## Code Structure

- `main.go` - Program entrypoint 
- `controller.go` - Handles app logic and flows
- `model.go` - Data structures 
- `scrapper.go` - Web scraping logic with Colly 
- `siteschema.go` - Parses schema from Excel
- `desktop.go` - GUI implementation
- `cmdline.go` - Command line implementation

## Features

- Accepts Excel file as input scraping schema
- Scrapes category pages to find products
- Follows links to scrape individual product pages
- Extracts attributes like title, price, images as defined in the excel mapping file.  
- Supports multiple retailers in one scraping run
- Respects robots.txt and noindex meta directives
- Outputs scraped data back to Excel 

## Testing

Run unit tests:

`go test ./...`

## Contributing

Pull requests are welcome! Please follow conventions in the existing code.

## License

MIT

## Credits

- Uses [Colly](https://go-colly.org/) for web scraping.
- Uses [Fyne](https://fyne.io/) for desktop UI.
- Uses [Cobra](https://cobra.dev) for command line UI.

