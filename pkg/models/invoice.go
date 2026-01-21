// Package models contains data structures for invoice generation.
package models

// Config holds global configuration for invoice generation.
type Config struct {
	VATPercentage  string
	CurrencySymbol string `json:"currencySymbol"`
	DateFormat     string `json:"dateFormat"`
	English        bool   `json:"english"` // false (default) = Arabic RTL, true = English LTR
}

// ProductInput represents a product from JSON input.
// All values are pre-calculated - this library only visualizes, no calculations.
type ProductInput struct {
	Name      string `json:"name"`
	Quantity  string `json:"quantity"`
	UnitPrice string `json:"unitPrice"`
	Discount  string `json:"discount,omitempty"` // Pre-calculated discount amount
	VATAmount string `json:"vatAmount"`          // Pre-calculated VAT
	Total     string `json:"total"`              // Pre-calculated total (inc. VAT)
}

// Product represents a single product line item for rendering.
// All values are pre-calculated and passed directly from input.
type Product struct {
	Name      string
	Quantity  string
	UnitPrice string
	Discount  string
	VATAmount string
	Total     string
}

// InvoiceInput represents invoice header data from JSON input.
// All totals are pre-calculated - this library only visualizes.
type InvoiceInput struct {
	Title             string `json:"title"`
	InvoiceNumber     string `json:"invoiceNumber"`
	StoreName         string `json:"storeName"`
	StoreAddress      string `json:"storeAddress"`
	Date              string `json:"date"`
	VATRegistrationNo string `json:"vatRegistrationNo"`
	QRCodeData        string `json:"qrCodeData"`
	TotalDiscount     string
	TotalTaxable      string
	TotalVAT          string
	TotalWithVAT      string
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
	TotalGross        string
	TotalDiscount     string
	TotalTaxableAmt   string
	TotalVAT          string
	TotalWithVAT      string
	QRCodeData        string
	VATPercentage     string
	Labels            Labels
	Language          string // "ar" or "en"
	IsRTL             bool   // true for Arabic/Hebrew, false for English
}
