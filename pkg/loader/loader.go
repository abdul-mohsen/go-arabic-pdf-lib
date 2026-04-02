// Package loader handles loading and parsing invoice data from JSON files.
package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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

	// Determine invoice type
	invoiceType := models.InvoiceType(data.InvoiceType)
	if invoiceType == "" {
		invoiceType = models.InvoiceTypeB2C
	}

	// Determine paper size
	paperSize := models.PaperSize(data.PaperSize)
	if paperSize == "" {
		paperSize = models.PaperThermal
	}
	// B2B types must be A4
	if invoiceType == models.InvoiceTypeB2B || invoiceType == models.InvoiceTypeB2BCredit || invoiceType == models.InvoiceTypeB2BDebit {
		paperSize = models.PaperA4
	}

	// Map product inputs directly to products (no calculations)
	products := make([]models.Product, 0, len(data.Products))
	for _, p := range data.Products {
		products = append(products, models.Product{
			Name:            p.Name,
			Quantity:        p.Quantity,
			UnitPrice:       p.UnitPrice,
			Discount:        p.Discount,
			SubtotalExclVAT: p.SubtotalExclVAT,
			VATAmount:       p.VATAmount,
			Total:           p.Total,
		})
	}

	// Parse date string - try common formats
	parsedDate := parseDate(data.Invoice.Date)

	return models.Invoice{
		Type:          invoiceType,
		PaperSize:     paperSize,
		Title:         data.Invoice.Title,
		InvoiceNumber: data.Invoice.InvoiceNumber,
		Date:          parsedDate,
		DateFormat:    data.Config.DateFormat,
		Seller: models.PartyInfo{
			Name:              data.Invoice.StoreName,
			Address:           data.Invoice.StoreAddress,
			VATRegistrationNo: data.Invoice.VATRegistrationNo,
			CommercialRegNo:   data.Invoice.CommercialRegNo,
		},
		Buyer: models.PartyInfo{
			Name:              data.Buyer.Name,
			Address:           data.Buyer.Address,
			VATRegistrationNo: data.Buyer.VATRegistrationNo,
			CommercialRegNo:   data.Buyer.CommercialRegNo,
		},
		NoteReason:      data.NoteReason,
		Products:        products,
		TotalDiscount:   data.Invoice.TotalDiscount,
		TotalTaxableAmt: data.Invoice.TotalTaxable,
		TotalVAT:        data.Invoice.TotalVAT,
		TotalWithVAT:    data.Invoice.TotalWithVAT,
		QRCodeData:      data.Invoice.QRCodeData,
		VATPercentage:   data.Config.VATPercentage,
		Labels:          data.Labels,
		Language:        language,
		IsRTL:           isRTL,
	}
}

// parseDate tries common date formats and returns the parsed time, or time.Now() on failure.
func parseDate(s string) time.Time {
	formats := []string{
		time.RFC3339,
		"2006-01-02",
		"2006/01/02",
		"02/01/2006",
		"02-01-2006",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Now()
}
