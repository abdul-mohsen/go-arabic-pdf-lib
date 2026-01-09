// Package loader handles loading and parsing invoice data from JSON files.
package loader

import (
	"encoding/json"
	"fmt"
	"os"

	"bill-generator/pkg/models"
)

// LoadFromJSON reads a JSON file and returns a fully calculated Invoice.
func LoadFromJSON(filename string) (models.Invoice, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return models.Invoice{}, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return ParseJSON(data)
}

// ParseJSON parses JSON data and returns a fully calculated Invoice.
func ParseJSON(data []byte) (models.Invoice, error) {
	var invoiceData models.InvoiceData
	if err := json.Unmarshal(data, &invoiceData); err != nil {
		return models.Invoice{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return BuildInvoice(invoiceData), nil
}

// BuildInvoice creates a complete Invoice from InvoiceData, calculating all derived values.
func BuildInvoice(data models.InvoiceData) models.Invoice {
	vatRate := data.Config.VATPercentage / 100.0
	language := data.Config.Language
	if language == "" {
		language = "ar" // Default to Arabic
	}

	// Calculate product values
	products := make([]models.Product, 0, len(data.Products))
	var totalTaxable, totalVAT float64

	for _, p := range data.Products {
		taxableAmt := p.Quantity * p.UnitPrice
		vatAmount := taxableAmt * vatRate
		totalWithVAT := taxableAmt + vatAmount

		products = append(products, models.Product{
			Name:         p.Name,
			Quantity:     p.Quantity,
			UnitPrice:    p.UnitPrice,
			TaxableAmt:   taxableAmt,
			VATAmount:    vatAmount,
			TotalWithVAT: totalWithVAT,
		})

		totalTaxable += taxableAmt
		totalVAT += vatAmount
	}

	return models.Invoice{
		Title:             data.Invoice.Title,
		InvoiceNumber:     data.Invoice.InvoiceNumber,
		StoreName:         data.Invoice.StoreName,
		StoreAddress:      data.Invoice.StoreAddress,
		Date:              data.Invoice.Date,
		VATRegistrationNo: data.Invoice.VATRegistrationNo,
		Products:          products,
		TotalTaxableAmt:   totalTaxable,
		TotalVAT:          totalVAT,
		TotalWithVAT:      totalTaxable + totalVAT,
		QRCodeData:        data.Invoice.QRCodeData,
		VATPercentage:     data.Config.VATPercentage,
		Labels:            data.Labels,
		Language:          language,
		IsRTL:             language == "ar" || language == "he",
	}
}
