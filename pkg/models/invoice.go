// Package models contains data structures for invoice generation.
package models

import (
	"time"
)

// InvoiceType represents the type of document being generated.
type InvoiceType string

const (
	InvoiceTypeB2C       InvoiceType = "b2c"
	InvoiceTypeB2B       InvoiceType = "b2b"
	InvoiceTypeB2CCredit InvoiceType = "b2c-credit"
	InvoiceTypeB2BCredit InvoiceType = "b2b-credit"
	InvoiceTypeB2CDebit  InvoiceType = "b2c-debit"
	InvoiceTypeB2BDebit  InvoiceType = "b2b-debit"
)

// PaperSize represents the output paper format.
type PaperSize string

const (
	PaperThermal PaperSize = "thermal" // 80mm receipt
	PaperA4      PaperSize = "a4"      // A4 (210mm x 297mm)
)

// PartyInfo represents a business entity (seller or buyer).
type PartyInfo struct {
	Name              string `json:"name"`
	Address           string `json:"address"`
	VATRegistrationNo string `json:"vatRegistrationNo"`
	CommercialRegNo   string `json:"commercialRegNo,omitempty"`
}

// Config holds global configuration for invoice generation.
type Config struct {
	VATPercentage  string `json:"vatPercentage"`
	CurrencySymbol string `json:"currencySymbol"`
	DateFormat     string `json:"dateFormat"`
	English        bool   `json:"english"` // false (default) = Arabic RTL, true = English LTR
}

// ProductInput represents a product from JSON input.
// All values are pre-calculated - this library only visualizes, no calculations.
type ProductInput struct {
	Name            string `json:"name"`
	Quantity        string `json:"quantity"`
	UnitPrice       string `json:"unitPrice"`
	Discount        string `json:"discount,omitempty"`        // Pre-calculated discount amount
	SubtotalExclVAT string `json:"subtotalExclVat,omitempty"` // Pre-calculated subtotal before VAT
	VATAmount       string `json:"vatAmount"`                 // Pre-calculated VAT
	Total           string `json:"total"`                     // Pre-calculated total (inc. VAT)
}

// Product represents a single product line item for rendering.
// All values are pre-calculated and passed directly from input.
type Product struct {
	Name            string
	Quantity        string
	UnitPrice       string
	Discount        string // Pre-calculated discount amount
	SubtotalExclVAT string // Pre-calculated subtotal before VAT
	VATAmount       string // Pre-calculated VAT
	Total           string // Pre-calculated total (inc. VAT)
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
	CommercialRegNo   string `json:"commercialRegNo,omitempty"`
	QRCodeData        string `json:"qrCodeData"`
	TotalDiscount     string `json:"totalDiscount,omitempty"` // Pre-calculated total discount
	TotalTaxable      string `json:"totalTaxable"`            // Pre-calculated taxable amount
	TotalVAT          string `json:"totalVat"`                // Pre-calculated total VAT
	TotalWithVAT      string `json:"totalWithVat"`            // Pre-calculated grand total
}

// BuyerInput represents buyer information from JSON input (B2B only).
type BuyerInput struct {
	Name              string `json:"name"`
	Address           string `json:"address"`
	VATRegistrationNo string `json:"vatRegistrationNo"`
	CommercialRegNo   string `json:"commercialRegNo,omitempty"`
}

// Labels holds all text labels for the invoice (supports i18n).
type Labels struct {
	// Document
	InvoiceNumber string `json:"invoiceNumber"`
	Date          string `json:"date"`
	Footer        string `json:"footer"`

	// Seller section
	SellerInfo      string `json:"sellerInfo,omitempty"`
	VATRegistration string `json:"vatRegistration"`
	CommercialReg   string `json:"commercialReg,omitempty"`

	// Buyer section (B2B)
	BuyerInfo            string `json:"buyerInfo,omitempty"`
	BuyerVATRegistration string `json:"buyerVatRegistration,omitempty"`
	BuyerCommercialReg   string `json:"buyerCommercialReg,omitempty"`

	// Credit/Debit note
	NoteReason string `json:"noteReason,omitempty"`

	// Table columns
	ProductColumn         string `json:"productColumn"`
	QuantityColumn        string `json:"quantityColumn"`
	UnitPriceColumn       string `json:"unitPriceColumn"`
	DiscountColumn        string `json:"discountColumn,omitempty"`
	SubtotalExclVATColumn string `json:"subtotalExclVatColumn,omitempty"`
	VATAmountColumn       string `json:"vatAmountColumn"`
	TotalColumn           string `json:"totalColumn"`

	// Totals
	TotalDiscount string `json:"totalDiscount,omitempty"`
	TotalTaxable  string `json:"totalTaxable"`
	TotalVAT      string `json:"totalVat,omitempty"`
	TotalWithVat  string `json:"totalWithVat"`
}

// InvoiceData represents the complete JSON input structure.
type InvoiceData struct {
	Config      Config         `json:"config"`
	Invoice     InvoiceInput   `json:"invoice"`
	Buyer       BuyerInput     `json:"buyer,omitempty"`
	Products    []ProductInput `json:"products"`
	Labels      Labels         `json:"labels"`
	InvoiceType string         `json:"invoiceType,omitempty"` // "b2c", "b2b", "b2c-credit", etc.
	PaperSize   string         `json:"paperSize,omitempty"`   // "thermal", "a4"
	NoteReason  string         `json:"noteReason,omitempty"`  // For credit/debit notes
}

// Invoice represents a fully processed invoice ready for PDF generation.
type Invoice struct {
	// Document type and format
	Type      InvoiceType
	PaperSize PaperSize

	// Header
	Title         string
	InvoiceNumber string
	Date          time.Time
	DateFormat    string

	// Seller info (all invoice types)
	Seller PartyInfo

	// Buyer info (B2B only)
	Buyer PartyInfo

	// Credit/Debit note
	NoteReason string

	// Products
	Products []Product

	// Totals (all pre-calculated)
	TotalGross      string
	TotalDiscount   string
	TotalTaxableAmt string
	TotalVAT        string
	TotalWithVAT    string

	// Metadata
	QRCodeData    string
	VATPercentage string
	Labels        Labels
	Language      string // "ar" or "en"
	IsRTL         bool
}
