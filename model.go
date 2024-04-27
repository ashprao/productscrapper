package main

// Model
type Model struct {
	ExcelFilePath string
}

// UpdateExcelFilePath updates the Excel file path in the model
func (m *Model) UpdateExcelFilePath(path string) {
	m.ExcelFilePath = path
}
