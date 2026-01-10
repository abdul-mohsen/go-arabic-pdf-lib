// Package models contains data structures for invoice generation.
package models

// Config holds global configuration for invoice generation.
type Config struct {
	VATPercentage  float64 `json:"vatPercentage"`
	CurrencySymbol string  `json:"currencySymbol"`
	DateFormat     string  `json:"dateFormat"`
	English        bool    `json:"english"` // false (default) = Arabic RTL, true = English LTR
}

// ProductInput represents a product from JSON input (without calculated fields).
type ProductInput struct {
	Name            string  `json:"name"`
	Quantity        float64 `json:"quantity"`
	UnitPrice       float64 `json:"unitPrice"`
	DiscountPercent float64 `json:"discountPercent,omitempty"` // Discount percentage (0-100)
	DiscountAmount  float64 `json:"discountAmount,omitempty"`  // Fixed discount amount
}

// Product represents a single product line item with calculated fields.
type Product struct {
	Name            string
	Quantity        float64
	UnitPrice       float64
	DiscountPercent float64
	DiscountAmount  float64
	GrossAmount     float64 // Quantity * UnitPrice (before discount)
	NetAmount       float64 // After discount, before VAT
	TaxableAmt      float64 // Same as NetAmount
	VATAmount       float64
	TotalWithVAT    float64
}

// InvoiceInput represents invoice header data from JSON input.
type InvoiceInput struct {
	Title             string `json:"title"`
	InvoiceNumber     string `json:"invoiceNumber"`
	StoreName         string `json:"storeName"`
	StoreAddress      string `json:"storeAddress"`
	Date              string `json:"date"`
	VATRegistrationNo string `json:"vatRegistrationNo"`
	QRCodeData        string `json:"qrCodeData"`
}

// Labels holds all text labels for the invoice (supports i18n).
type Labels struct {
	InvoiceNumber   string `json:"invoiceNumber"`
	Date            string `json:"date"`
	VATRegistration string `json:"vatRegistration"`
	TotalTaxable    string `json:"totalTaxable"`
	TotalWithVat    string `json:"totalWithVat"`
	ProductColumn   string `json:"productColumn"`
	QuantityColumn  string `json:"quantityColumn"`
	UnitPriceColumn string `json:"unitPriceColumn"`
	DiscountColumn  string `json:"discountColumn,omitempty"`
	VATAmountColumn string `json:"vatAmountColumn"`
	TotalColumn     string `json:"totalColumn"`
	TotalDiscount   string `json:"totalDiscount,omitempty"`
	Footer          string `json:"footer"`
}

// InvoiceData represents the complete JSON input structure.
type InvoiceData struct {
	Config   Config         `json:"config"`
	Invoice  InvoiceInput   `json:"invoice"`
	Products []ProductInput `json:"products"`
	Labels   Labels         `json:"labels"`
}

// Invoice represents a fully processed invoice ready for PDF generation.
type Invoice struct {
	Title             string
	InvoiceNumber     string
	StoreName         string
	StoreAddress      string
	Date              string
	VATRegistrationNo string
	Products          []Product
	TotalGross        float64 // Sum of all gross amounts (before discounts)
	TotalDiscount     float64 // Sum of all discounts
	TotalTaxableAmt   float64 // Sum of net amounts (after discounts)
	TotalVAT          float64
	TotalWithVAT      float64
	QRCodeData        string
	VATPercentage     float64
	Labels            Labels
	Language          string // "ar" or "en"
	IsRTL             bool   // true for Arabic/Hebrew, false for English
}
