// Package loader handles loading and parsing invoice data from JSON files.
package loader

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/models"
)

// LoadFromJSON reads a JSON file and returns an Invoice for visualization.
func LoadFromJSON(filename string) (models.Invoice, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return models.Invoice{}, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return ParseJSON(data)
}

// ParseJSON parses JSON data and returns an Invoice for visualization.
func ParseJSON(data []byte) (models.Invoice, error) {
	var invoiceData models.InvoiceData
	if err := json.Unmarshal(data, &invoiceData); err != nil {
		return models.Invoice{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return BuildInvoice(invoiceData), nil
}

// BuildInvoice creates an Invoice from InvoiceData.
// This library is for visualization only - all values are pre-calculated in the input.
func BuildInvoice(data models.InvoiceData) models.Invoice {
	// Default is Arabic (RTL), English requires explicit flag
	language := "ar"
	isRTL := true
	if data.Config.English {
		language = "en"
		isRTL = false
	}

	// Map product inputs directly to products (no calculations)
	products := make([]models.Product, 0, len(data.Products))
	for _, p := range data.Products {
		products = append(products, models.Product{
			Name:      p.Name,
			Quantity:  p.Quantity,
			UnitPrice: p.UnitPrice,
			Discount:  p.Discount,
			VATAmount: p.VATAmount,
			Total:     p.Total,
		})
	}

	return models.Invoice{
		Title:             data.Invoice.Title,
		InvoiceNumber:     data.Invoice.InvoiceNumber,
		StoreName:         data.Invoice.StoreName,
		StoreAddress:      data.Invoice.StoreAddress,
		Date:              data.Invoice.Date,
		VATRegistrationNo: data.Invoice.VATRegistrationNo,
		Products:          products,
		TotalDiscount:     data.Invoice.TotalDiscount,
		TotalTaxableAmt:   data.Invoice.TotalTaxable,
		TotalVAT:          data.Invoice.TotalVAT,
		TotalWithVAT:      data.Invoice.TotalWithVAT,
		QRCodeData:        data.Invoice.QRCodeData,
		VATPercentage:     data.Config.VATPercentage,
		Labels:            data.Labels,
		Language:          language,
		IsRTL:             isRTL,
	}
}
