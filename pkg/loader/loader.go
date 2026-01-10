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

	// Default is Arabic (RTL), English requires explicit flag
	language := "ar"
	isRTL := true
	if data.Config.English {
		language = "en"
		isRTL = false
	}

	// Calculate product values with discount support
	products := make([]models.Product, 0, len(data.Products))
	var totalGross, totalDiscount, totalTaxable, totalVAT float64

	for _, p := range data.Products {
		// Calculate gross amount (before discount)
		grossAmount := p.Quantity * p.UnitPrice

		// Calculate discount: percentage discount + fixed discount
		discountAmt := (grossAmount * p.DiscountPercent / 100.0) + p.DiscountAmount

		// Net amount after discount is the taxable amount
		netAmount := grossAmount - discountAmt
		if netAmount < 0 {
			netAmount = 0 // Prevent negative amounts
		}

		// VAT is calculated on net amount (after discount)
		vatAmount := netAmount * vatRate
		totalWithVAT := netAmount + vatAmount

		products = append(products, models.Product{
			Name:            p.Name,
			Quantity:        p.Quantity,
			UnitPrice:       p.UnitPrice,
			DiscountPercent: p.DiscountPercent,
			DiscountAmount:  discountAmt,
			GrossAmount:     grossAmount,
			NetAmount:       netAmount,
			TaxableAmt:      netAmount,
			VATAmount:       vatAmount,
			TotalWithVAT:    totalWithVAT,
		})

		totalGross += grossAmount
		totalDiscount += discountAmt
		totalTaxable += netAmount
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
		TotalGross:        totalGross,
		TotalDiscount:     totalDiscount,
		TotalTaxableAmt:   totalTaxable,
		TotalVAT:          totalVAT,
		TotalWithVAT:      totalTaxable + totalVAT,
		QRCodeData:        data.Invoice.QRCodeData,
		VATPercentage:     data.Config.VATPercentage,
		Labels:            data.Labels,
		Language:          language,
		IsRTL:             isRTL,
	}
}
