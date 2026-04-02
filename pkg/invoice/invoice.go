// Package invoice provides a clean API for generating PDF invoices.
// This is the main entry point for using the bill-generator library.
//
// Example usage:
//
//	import "github.com/ssda/bill-generator/pkg/invoice"
//
//	// Generate from JSON file
//	err := invoice.GenerateFromFile("invoice.json", "output.pdf")
//
//	// Generate from Invoice struct
//	inv := invoice.Invoice{...}
//	err := invoice.Generate(inv, "output.pdf")
package invoice

import (
	"fmt"
	"time"

	"github.com/ssda/bill-generator/pkg/loader"
	"github.com/ssda/bill-generator/pkg/models"
	"github.com/ssda/bill-generator/pkg/pdf"
)

// Type aliases for convenience.
type (
	Invoice      = models.Invoice
	InvoiceData  = models.InvoiceData
	Product      = models.Product
	ProductInput = models.ProductInput
	Labels       = models.Labels
	Config       = models.Config
	PartyInfo    = models.PartyInfo
	InvoiceType  = models.InvoiceType
	PaperSize    = models.PaperSize
)

// Re-export constants.
const (
	TypeB2C       = models.InvoiceTypeB2C
	TypeB2B       = models.InvoiceTypeB2B
	TypeB2CCredit = models.InvoiceTypeB2CCredit
	TypeB2BCredit = models.InvoiceTypeB2BCredit
	TypeB2CDebit  = models.InvoiceTypeB2CDebit
	TypeB2BDebit  = models.InvoiceTypeB2BDebit
	PaperThermal  = models.PaperThermal
	PaperA4       = models.PaperA4
)

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
		fontPath: "fonts",
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
	FillDefaultLabels(&inv)
	return g.Generate(inv, outputPath)
}

// GenerateFromJSON parses JSON data and generates a PDF invoice.
func (g *Generator) GenerateFromJSON(jsonData []byte, outputPath string) error {
	inv, err := loader.ParseJSON(jsonData)
	if err != nil {
		return fmt.Errorf("failed to parse invoice: %w", err)
	}
	FillDefaultLabels(&inv)
	return g.Generate(inv, outputPath)
}

// Generate creates a PDF invoice from an Invoice struct.
func (g *Generator) Generate(inv Invoice, outputPath string) error {
	FillDefaultLabels(&inv)
	return pdf.GenerateInvoice(inv, outputPath, g.fontPath)
}

// GenerateBytes creates a PDF invoice and returns it as bytes.
func (g *Generator) GenerateBytes(inv Invoice) ([]byte, error) {
	FillDefaultLabels(&inv)
	return pdf.GenerateInvoiceBytes(inv, g.fontPath)
}

// --- Package-level convenience functions ---

var defaultGenerator = NewGenerator()

// SetDefaultFontPath sets the font path for the default generator.
func SetDefaultFontPath(path string) {
	defaultGenerator.fontPath = path
}

// GenerateFromFile reads a JSON file and generates a PDF invoice.
func GenerateFromFile(jsonPath, outputPath string) error {
	return defaultGenerator.GenerateFromFile(jsonPath, outputPath)
}

// GenerateFromJSON parses JSON data and generates a PDF invoice.
func GenerateFromJSON(jsonData []byte, outputPath string) error {
	return defaultGenerator.GenerateFromJSON(jsonData, outputPath)
}

// Generate creates a PDF invoice from an Invoice struct.
func Generate(inv Invoice, outputPath string) error {
	return defaultGenerator.Generate(inv, outputPath)
}

// GenerateBytes creates a PDF invoice and returns it as bytes.
func GenerateBytes(inv Invoice) ([]byte, error) {
	return defaultGenerator.GenerateBytes(inv)
}

// --- Builder Pattern ---

// Builder provides a fluent API for building invoices programmatically.
type Builder struct {
	inv models.Invoice
}

// NewBuilder creates a new invoice builder with sensible defaults.
func NewBuilder() *Builder {
	return &Builder{
		inv: models.Invoice{
			Type:          models.InvoiceTypeB2C,
			PaperSize:     models.PaperThermal,
			VATPercentage: 15.0,
			Date:          time.Now(),
			DateFormat:    "2006-01-02 15:04:05",
			Language:      "ar",
			IsRTL:         true,
		},
	}
}

func (b *Builder) WithType(t models.InvoiceType) *Builder { b.inv.Type = t; return b }
func (b *Builder) WithPaper(p models.PaperSize) *Builder  { b.inv.PaperSize = p; return b }
func (b *Builder) WithTitle(title string) *Builder        { b.inv.Title = title; return b }
func (b *Builder) WithInvoiceNumber(n string) *Builder    { b.inv.InvoiceNumber = n; return b }
func (b *Builder) WithDate(date time.Time) *Builder       { b.inv.Date = date; return b }
func (b *Builder) WithDateFormat(f string) *Builder       { b.inv.DateFormat = f; return b }
func (b *Builder) WithQRCode(data string) *Builder        { b.inv.QRCodeData = data; return b }
func (b *Builder) WithVATPercentage(p float64) *Builder   { b.inv.VATPercentage = p; return b }
func (b *Builder) WithNoteReason(reason string) *Builder  { b.inv.NoteReason = reason; return b }
func (b *Builder) WithLabels(labels Labels) *Builder      { b.inv.Labels = labels; return b }

// WithSeller sets the seller (store) information.
func (b *Builder) WithSeller(name, address, vatNo, commercialReg string) *Builder {
	b.inv.Seller = models.PartyInfo{
		Name: name, Address: address,
		VATRegistrationNo: vatNo, CommercialRegNo: commercialReg,
	}
	return b
}

// WithBuyer sets the buyer information (B2B only).
func (b *Builder) WithBuyer(name, address, vatNo, commercialReg string) *Builder {
	b.inv.Buyer = models.PartyInfo{
		Name: name, Address: address,
		VATRegistrationNo: vatNo, CommercialRegNo: commercialReg,
	}
	return b
}

// Backward compatible aliases
func (b *Builder) WithStoreName(name string) *Builder    { b.inv.Seller.Name = name; return b }
func (b *Builder) WithStoreAddress(addr string) *Builder { b.inv.Seller.Address = addr; return b }
func (b *Builder) WithVATRegistration(vatNo string) *Builder {
	b.inv.Seller.VATRegistrationNo = vatNo
	return b
}

// WithEnglish sets the invoice language to English.
func (b *Builder) WithEnglish() *Builder {
	b.inv.Language = "en"
	b.inv.IsRTL = false
	return b
}

// WithArabic sets the invoice language to Arabic (default).
func (b *Builder) WithArabic() *Builder {
	b.inv.Language = "ar"
	b.inv.IsRTL = true
	return b
}

// AddProduct adds a product to the invoice with pre-calculated values.
func (b *Builder) AddProduct(name string, quantity, unitPrice, discount, vatAmount, total float64) *Builder {
	b.inv.Products = append(b.inv.Products, models.Product{
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
	b.inv.TotalDiscount = totalDiscount
	b.inv.TotalTaxableAmt = totalTaxable
	b.inv.TotalVAT = totalVAT
	b.inv.TotalWithVAT = totalWithVAT
	return b
}

// Build creates the Invoice, filling in default labels for any empty fields.
func (b *Builder) Build() Invoice {
	inv := b.inv
	// Enforce B2B = A4
	if inv.Type == models.InvoiceTypeB2B || inv.Type == models.InvoiceTypeB2BCredit || inv.Type == models.InvoiceTypeB2BDebit {
		inv.PaperSize = models.PaperA4
	}
	FillDefaultLabels(&inv)
	return inv
}

// --- Default Labels ---

// FillDefaultLabels fills any empty label fields with defaults based on language.
func FillDefaultLabels(inv *models.Invoice) {
	var defaults models.Labels
	if inv.IsRTL {
		defaults = DefaultArabicLabels()
	} else {
		defaults = DefaultEnglishLabels()
	}

	l := &inv.Labels
	if l.InvoiceNumber == "" {
		l.InvoiceNumber = defaults.InvoiceNumber
	}
	if l.Date == "" {
		l.Date = defaults.Date
	}
	if l.Footer == "" {
		l.Footer = defaults.Footer
	}
	if l.SellerInfo == "" {
		l.SellerInfo = defaults.SellerInfo
	}
	if l.VATRegistration == "" {
		l.VATRegistration = defaults.VATRegistration
	}
	if l.CommercialReg == "" {
		l.CommercialReg = defaults.CommercialReg
	}
	if l.BuyerInfo == "" {
		l.BuyerInfo = defaults.BuyerInfo
	}
	if l.BuyerVATRegistration == "" {
		l.BuyerVATRegistration = defaults.BuyerVATRegistration
	}
	if l.BuyerCommercialReg == "" {
		l.BuyerCommercialReg = defaults.BuyerCommercialReg
	}
	if l.NoteReason == "" {
		l.NoteReason = defaults.NoteReason
	}
	if l.ProductColumn == "" {
		l.ProductColumn = defaults.ProductColumn
	}
	if l.QuantityColumn == "" {
		l.QuantityColumn = defaults.QuantityColumn
	}
	if l.UnitPriceColumn == "" {
		l.UnitPriceColumn = defaults.UnitPriceColumn
	}
	if l.DiscountColumn == "" {
		l.DiscountColumn = defaults.DiscountColumn
	}
	if l.SubtotalExclVATColumn == "" {
		l.SubtotalExclVATColumn = defaults.SubtotalExclVATColumn
	}
	if l.VATAmountColumn == "" {
		l.VATAmountColumn = defaults.VATAmountColumn
	}
	if l.TotalColumn == "" {
		l.TotalColumn = defaults.TotalColumn
	}
	if l.TotalDiscount == "" {
		l.TotalDiscount = defaults.TotalDiscount
	}
	if l.TotalTaxable == "" {
		l.TotalTaxable = defaults.TotalTaxable
	}
	if l.TotalVAT == "" {
		l.TotalVAT = defaults.TotalVAT
	}
	if l.TotalWithVat == "" {
		l.TotalWithVat = defaults.TotalWithVat
	}
}

// DefaultEnglishLabels returns the default English labels.
func DefaultEnglishLabels() Labels {
	return Labels{
		InvoiceNumber:         "Invoice No:",
		Date:                  "Date:",
		Footer:                "Thank you for your business",
		SellerInfo:            "Seller Information",
		VATRegistration:       "VAT Registration No:",
		CommercialReg:         "Commercial Reg:",
		BuyerInfo:             "Buyer Information",
		BuyerVATRegistration:  "Buyer VAT Reg:",
		BuyerCommercialReg:    "Buyer Commercial Reg:",
		NoteReason:            "Reason:",
		ProductColumn:         "Product",
		QuantityColumn:        "Qty",
		UnitPriceColumn:       "Unit Price",
		DiscountColumn:        "Discount",
		SubtotalExclVATColumn: "Subtotal (excl. VAT)",
		VATAmountColumn:       "VAT",
		TotalColumn:           "Total",
		TotalDiscount:         "Total Discount:",
		TotalTaxable:          "Total Taxable Amount:",
		TotalVAT:              "VAT Amount:",
		TotalWithVat:          "Total with VAT:",
	}
}

// DefaultArabicLabels returns the default Arabic labels.
func DefaultArabicLabels() Labels {
	return Labels{
		InvoiceNumber:         "رقم الفاتورة:",
		Date:                  "التاريخ:",
		Footer:                "شكراً لتعاملكم معنا",
		SellerInfo:            "معلومات البائع",
		VATRegistration:       "رقم تسجيل ضريبة القيمة المضافة:",
		CommercialReg:         "رقم السجل التجاري:",
		BuyerInfo:             "معلومات المشتري",
		BuyerVATRegistration:  "رقم تسجيل ضريبة القيمة المضافة للمشتري:",
		BuyerCommercialReg:    "رقم السجل التجاري للمشتري:",
		NoteReason:            "السبب:",
		ProductColumn:         "المنتج",
		QuantityColumn:        "الكمية",
		UnitPriceColumn:       "سعر الوحدة",
		DiscountColumn:        "الخصم",
		SubtotalExclVATColumn: "المجموع الفرعي بدون الضريبة",
		VATAmountColumn:       "ضريبة القيمة المضافة",
		TotalColumn:           "السعر شامل الضريبة",
		TotalDiscount:         "إجمالي الخصم:",
		TotalTaxable:          "المجموع",
		TotalVAT:              "ضريبة القيمة المضافة",
		TotalWithVat:          "المجموع شامل ضريبة القيمة المضافة",
	}
}
