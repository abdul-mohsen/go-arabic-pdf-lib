// Package invoice provides a clean API for generating PDF invoices.
// This is the main entry point for using the github.com/abdul-mohsen/go-arabic-pdf-lib library.
//
// Example usage:
//
//	import "github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/invoice"
//
//	// Generate from JSON file
//	err := invoice.GenerateFromFile("invoice.json", "output.pdf")
//
//	// Generate from JSON data
//	jsonData := []byte(`{"invoice": {...}, "products": [...], ...}`)
//	err := invoice.GenerateFromJSON(jsonData, "output.pdf")
//
//	// Generate from Invoice struct
//	inv := invoice.Invoice{...}
//	err := invoice.Generate(inv, "output.pdf")
package invoice

import (
	"fmt"

	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/loader"
	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/models"
	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/pdf"
)

// Invoice represents a complete invoice with all calculated values.
// This is the main data structure for invoice generation.
type Invoice = models.Invoice

// InvoiceData represents the raw input data for an invoice (from JSON).
type InvoiceData = models.InvoiceData

// Product represents a calculated product line item.
type Product = models.Product

// ProductInput represents raw product input data.
type ProductInput = models.ProductInput

// Labels holds all text labels for i18n support.
type Labels = models.Labels

// Config holds invoice configuration options.
type Config = models.Config

// Generator creates PDF invoices.
type Generator struct {
	fontPath string
}

// Option configures a Generator.
type Option func(*Generator)

// WithFontPath sets the font path for the generator.
func WithFontPath(path string) Option {
	return func(g *Generator) {
		g.fontPath = path
	}
}

// NewGenerator creates a new invoice generator with optional configuration.
func NewGenerator(opts ...Option) *Generator {
	g := &Generator{
		fontPath: "/fonts/Amiri-Regular.ttf", // default font path
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// GenerateFromFile reads a JSON file and generates a PDF invoice.
func (g *Generator) GenerateFromFile(jsonPath, outputPath string) error {
	inv, err := loader.LoadFromJSON(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to load invoice: %w", err)
	}
	return g.Generate(inv, outputPath)
}

// GenerateFromJSON parses JSON data and generates a PDF invoice.
func (g *Generator) GenerateFromJSON(jsonData []byte, outputPath string) error {
	inv, err := loader.ParseJSON(jsonData)
	if err != nil {
		return fmt.Errorf("failed to parse invoice: %w", err)
	}
	return g.Generate(inv, outputPath)
}

// Generate creates a PDF invoice from an Invoice struct.
func (g *Generator) Generate(inv Invoice, outputPath string) error {
	return pdf.GenerateInvoice(inv, outputPath, g.fontPath)
}

// GenerateBytes creates a PDF invoice and returns it as bytes.
func (g *Generator) GenerateBytes(inv Invoice) ([]byte, error) {
	return pdf.GenerateInvoiceBytes(inv, g.fontPath)
}

// --- Package-level convenience functions ---

// defaultGenerator is used by package-level functions.
var defaultGenerator = NewGenerator()

// SetDefaultFontPath sets the font path for the default generator.
func SetDefaultFontPath(path string) {
	defaultGenerator.fontPath = path
}

// GenerateFromFile reads a JSON file and generates a PDF invoice.
// Uses the default generator configuration.
func GenerateFromFile(jsonPath, outputPath string) error {
	return defaultGenerator.GenerateFromFile(jsonPath, outputPath)
}

// GenerateFromJSON parses JSON data and generates a PDF invoice.
// Uses the default generator configuration.
func GenerateFromJSON(jsonData []byte, outputPath string) error {
	return defaultGenerator.GenerateFromJSON(jsonData, outputPath)
}

// Generate creates a PDF invoice from an Invoice struct.
// Uses the default generator configuration.
func Generate(inv Invoice, outputPath string) error {
	return defaultGenerator.Generate(inv, outputPath)
}

// GenerateBytes creates a PDF invoice and returns it as bytes.
// Uses the default generator configuration.
func GenerateBytes(inv Invoice) ([]byte, error) {
	return defaultGenerator.GenerateBytes(inv)
}

// --- Builder Pattern for programmatic invoice creation ---

// Builder provides a fluent API for building invoices programmatically.
type Builder struct {
	data models.InvoiceData
}

// NewBuilder creates a new invoice builder.
func NewBuilder() *Builder {
	return &Builder{
		data: models.InvoiceData{
			Config: models.Config{
				VATPercentage: 15.0, // Default VAT
			},
		},
	}
}

// WithTitle sets the invoice title.
func (b *Builder) WithTitle(title string) *Builder {
	b.data.Invoice.Title = title
	return b
}

// WithInvoiceNumber sets the invoice number.
func (b *Builder) WithInvoiceNumber(number string) *Builder {
	b.data.Invoice.InvoiceNumber = number
	return b
}

// WithStoreName sets the store name.
func (b *Builder) WithStoreName(name string) *Builder {
	b.data.Invoice.StoreName = name
	return b
}

// WithStoreAddress sets the store address.
func (b *Builder) WithStoreAddress(address string) *Builder {
	b.data.Invoice.StoreAddress = address
	return b
}

// WithDate sets the invoice date.
func (b *Builder) WithDate(date string) *Builder {
	b.data.Invoice.Date = date
	return b
}

// WithVATRegistration sets the VAT registration number.
func (b *Builder) WithVATRegistration(vatNo string) *Builder {
	b.data.Invoice.VATRegistrationNo = vatNo
	return b
}

// WithQRCode sets the QR code data.
func (b *Builder) WithQRCode(data string) *Builder {
	b.data.Invoice.QRCodeData = data
	return b
}

// WithVATPercentage sets the VAT percentage.
func (b *Builder) WithVATPercentage(percentage float64) *Builder {
	b.data.Config.VATPercentage = percentage
	return b
}

// WithEnglish sets the invoice language to English.
func (b *Builder) WithEnglish() *Builder {
	b.data.Config.English = true
	return b
}

// WithArabic sets the invoice language to Arabic (default).
func (b *Builder) WithArabic() *Builder {
	b.data.Config.English = false
	return b
}

// WithLabels sets the labels for the invoice.
func (b *Builder) WithLabels(labels Labels) *Builder {
	b.data.Labels = labels
	return b
}

// AddProduct adds a product to the invoice with pre-calculated values.
// This library is for visualization only - all values must be pre-calculated.
func (b *Builder) AddProduct(name string, quantity, unitPrice, discount, vatAmount, total float64) *Builder {
	b.data.Products = append(b.data.Products, models.ProductInput{
		Name:      name,
		Quantity:  quantity,
		UnitPrice: unitPrice,
		Discount:  discount,
		VATAmount: vatAmount,
		Total:     total,
	})
	return b
}

// WithTotals sets the pre-calculated totals for the invoice.
func (b *Builder) WithTotals(totalDiscount, totalTaxable, totalVAT, totalWithVAT float64) *Builder {
	b.data.Invoice.TotalDiscount = totalDiscount
	b.data.Invoice.TotalTaxable = totalTaxable
	b.data.Invoice.TotalVAT = totalVAT
	b.data.Invoice.TotalWithVAT = totalWithVAT
	return b
}

// Build creates the Invoice from the builder data.
func (b *Builder) Build() Invoice {
	return loader.BuildInvoice(b.data)
}

// DefaultEnglishLabels returns the default English labels.
func DefaultEnglishLabels() Labels {
	return Labels{
		InvoiceNumber:   "Invoice No:",
		Date:            "Date:",
		VATRegistration: "VAT Registration No:",
		TotalTaxable:    "Total Taxable Amount:",
		TotalWithVat:    "Total with VAT:",
		ProductColumn:   "Product",
		QuantityColumn:  "Qty",
		UnitPriceColumn: "Unit Price",
		DiscountColumn:  "Discount",
		VATAmountColumn: "VAT",
		TotalColumn:     "Total",
		TotalDiscount:   "Total Discount:",
		Footer:          "Thank you for your business",
	}
}

// DefaultArabicLabels returns the default Arabic labels.
func DefaultArabicLabels() Labels {
	return Labels{
		InvoiceNumber:   "رقم الفاتورة:",
		Date:            "التاريخ:",
		VATRegistration: "رقم التسجيل الضريبي:",
		TotalTaxable:    "إجمالي المبلغ الخاضع للضريبة:",
		TotalWithVat:    "الإجمالي شامل الضريبة:",
		ProductColumn:   "المنتج",
		QuantityColumn:  "الكمية",
		UnitPriceColumn: "سعر الوحدة",
		DiscountColumn:  "الخصم",
		VATAmountColumn: "الضريبة",
		TotalColumn:     "الإجمالي",
		TotalDiscount:   "إجمالي الخصم:",
		Footer:          "شكراً لتعاملكم معنا",
	}
}
