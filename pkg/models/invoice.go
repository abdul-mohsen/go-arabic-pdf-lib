// Package models contains data structures for invoice generation.
package models

// Config holds global configuration for invoice generation.
type Config struct {
	VATPercentage  float64 `json:"vatPercentage"`
	CurrencySymbol string  `json:"currencySymbol"`
	DateFormat     string  `json:"dateFormat"`
	English        bool    `json:"english"` // false (default) = Arabic RTL, true = English LTR
}

// ProductInput represents a product from JSON input.
// All values are pre-calculated - this library only visualizes, no calculations.
type ProductInput struct {
	Name      string  `json:"name"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
	Discount  float64 `json:"discount,omitempty"` // Pre-calculated discount amount
	VATAmount float64 `json:"vatAmount"`          // Pre-calculated VAT
	Total     float64 `json:"total"`              // Pre-calculated total (inc. VAT)
}

// Product represents a single product line item for rendering.
// All values are pre-calculated and passed directly from input.
type Product struct {
	Name      string
	Quantity  float64
	UnitPrice float64
	Discount  float64 // Pre-calculated discount amount
	VATAmount float64 // Pre-calculated VAT
	Total     float64 // Pre-calculated total (inc. VAT)
}

// InvoiceInput represents invoice header data from JSON input.
// All totals are pre-calculated - this library only visualizes.
type InvoiceInput struct {
	Title             string  `json:"title"`
	InvoiceNumber     string  `json:"invoiceNumber"`
	StoreName         string  `json:"storeName"`
	StoreAddress      string  `json:"storeAddress"`
	Date              string  `json:"date"`
	VATRegistrationNo string  `json:"vatRegistrationNo"`
	QRCodeData        string  `json:"qrCodeData"`
	TotalDiscount     float64 `json:"totalDiscount,omitempty"`  // Pre-calculated total discount
	TotalTaxable      float64 `json:"totalTaxable"`             // Pre-calculated taxable amount
	TotalVAT          float64 `json:"totalVat"`                 // Pre-calculated total VAT
	TotalWithVAT      float64 `json:"totalWithVat"`             // Pre-calculated grand total
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
