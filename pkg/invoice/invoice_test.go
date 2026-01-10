package invoice

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	if builder == nil {
		t.Fatal("NewBuilder returned nil")
	}
}

func TestBuilder_BasicFields(t *testing.T) {
	inv := NewBuilder().
		WithTitle("Test Invoice").
		WithInvoiceNumber("INV-001").
		WithStoreName("Test Store").
		WithStoreAddress("123 Main St").
		WithDate("2024/01/15").
		WithVATRegistration("123456789").
		WithQRCode("test-qr").
		Build()

	if inv.Title != "Test Invoice" {
		t.Errorf("Expected title 'Test Invoice', got '%s'", inv.Title)
	}
	if inv.InvoiceNumber != "INV-001" {
		t.Errorf("Expected invoice number 'INV-001', got '%s'", inv.InvoiceNumber)
	}
	if inv.StoreName != "Test Store" {
		t.Errorf("Expected store name 'Test Store', got '%s'", inv.StoreName)
	}
	if inv.StoreAddress != "123 Main St" {
		t.Errorf("Expected store address '123 Main St', got '%s'", inv.StoreAddress)
	}
	if inv.Date != "2024/01/15" {
		t.Errorf("Expected date '2024/01/15', got '%s'", inv.Date)
	}
}

func TestBuilder_VATCalculation(t *testing.T) {
	inv := NewBuilder().
		WithVATPercentage(15.0).
		AddProduct("Product 1", 1, 100.00).
		Build()

	if len(inv.Products) != 1 {
		t.Fatalf("Expected 1 product, got %d", len(inv.Products))
	}

	p := inv.Products[0]
	// 100 * 15% = 15 VAT, Total = 115
	if p.TaxableAmt != 100 {
		t.Errorf("Expected taxable 100, got %.2f", p.TaxableAmt)
	}
	if p.VATAmount != 15 {
		t.Errorf("Expected VAT 15, got %.2f", p.VATAmount)
	}
	if p.TotalWithVAT != 115 {
		t.Errorf("Expected total 115, got %.2f", p.TotalWithVAT)
	}
}

func TestBuilder_Discount(t *testing.T) {
	inv := NewBuilder().
		WithVATPercentage(15.0).
		AddProductWithDiscount("Product 1", 1, 100.00, 10, 0). // 10% discount
		Build()

	p := inv.Products[0]
	// Gross: 100, Discount: 10, Net: 90, VAT: 13.5, Total: 103.5
	if p.GrossAmount != 100 {
		t.Errorf("Expected gross 100, got %.2f", p.GrossAmount)
	}
	if p.DiscountAmount != 10 {
		t.Errorf("Expected discount 10, got %.2f", p.DiscountAmount)
	}
	if p.NetAmount != 90 {
		t.Errorf("Expected net 90, got %.2f", p.NetAmount)
	}
	if p.TotalWithVAT != 103.5 {
		t.Errorf("Expected total 103.5, got %.2f", p.TotalWithVAT)
	}

	if inv.TotalDiscount != 10 {
		t.Errorf("Expected total discount 10, got %.2f", inv.TotalDiscount)
	}
}

func TestBuilder_EnglishLanguage(t *testing.T) {
	inv := NewBuilder().
		WithEnglish().
		Build()

	if inv.Language != "en" {
		t.Errorf("Expected language 'en', got '%s'", inv.Language)
	}
	if inv.IsRTL {
		t.Error("Expected IsRTL to be false for English")
	}
}

func TestBuilder_ArabicLanguage(t *testing.T) {
	inv := NewBuilder().
		WithArabic().
		Build()

	if inv.Language != "ar" {
		t.Errorf("Expected language 'ar', got '%s'", inv.Language)
	}
	if !inv.IsRTL {
		t.Error("Expected IsRTL to be true for Arabic")
	}
}

func TestBuilder_DefaultIsArabic(t *testing.T) {
	inv := NewBuilder().Build()

	if inv.Language != "ar" {
		t.Errorf("Expected default language 'ar', got '%s'", inv.Language)
	}
	if !inv.IsRTL {
		t.Error("Expected IsRTL to be true by default")
	}
}

func TestBuilder_Labels(t *testing.T) {
	labels := Labels{
		InvoiceNumber: "Invoice #:",
		Date:          "Date:",
		Footer:        "Thank you!",
	}

	inv := NewBuilder().
		WithLabels(labels).
		Build()

	if inv.Labels.InvoiceNumber != "Invoice #:" {
		t.Errorf("Expected invoice number label 'Invoice #:', got '%s'", inv.Labels.InvoiceNumber)
	}
	if inv.Labels.Footer != "Thank you!" {
		t.Errorf("Expected footer 'Thank you!', got '%s'", inv.Labels.Footer)
	}
}

func TestDefaultEnglishLabels(t *testing.T) {
	labels := DefaultEnglishLabels()

	if labels.InvoiceNumber == "" {
		t.Error("Invoice number label should not be empty")
	}
	if labels.DiscountColumn == "" {
		t.Error("Discount column label should not be empty")
	}
	if labels.TotalDiscount == "" {
		t.Error("Total discount label should not be empty")
	}
}

func TestDefaultArabicLabels(t *testing.T) {
	labels := DefaultArabicLabels()

	if labels.InvoiceNumber == "" {
		t.Error("Invoice number label should not be empty")
	}
	if labels.DiscountColumn == "" {
		t.Error("Discount column label should not be empty")
	}
	if labels.TotalDiscount == "" {
		t.Error("Total discount label should not be empty")
	}
}

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator()
	if gen == nil {
		t.Fatal("NewGenerator returned nil")
	}
	if gen.fontPath != "/fonts/Amiri-Regular.ttf" {
		t.Errorf("Expected default font path, got '%s'", gen.fontPath)
	}
}

func TestNewGenerator_WithFontPath(t *testing.T) {
	gen := NewGenerator(WithFontPath("/custom/font.ttf"))
	if gen.fontPath != "/custom/font.ttf" {
		t.Errorf("Expected font path '/custom/font.ttf', got '%s'", gen.fontPath)
	}
}

func TestBuilder_MultipleProducts(t *testing.T) {
	inv := NewBuilder().
		WithVATPercentage(15.0).
		AddProduct("Product 1", 2, 50.00).  // 100
		AddProduct("Product 2", 1, 100.00). // 100
		AddProductWithDiscount("Product 3", 1, 80.00, 0, 10). // 70 after discount
		Build()

	if len(inv.Products) != 3 {
		t.Fatalf("Expected 3 products, got %d", len(inv.Products))
	}

	// Total gross: 100 + 100 + 80 = 280
	if inv.TotalGross != 280 {
		t.Errorf("Expected total gross 280, got %.2f", inv.TotalGross)
	}

	// Total discount: 0 + 0 + 10 = 10
	if inv.TotalDiscount != 10 {
		t.Errorf("Expected total discount 10, got %.2f", inv.TotalDiscount)
	}

	// Total taxable (net): 100 + 100 + 70 = 270
	if inv.TotalTaxableAmt != 270 {
		t.Errorf("Expected total taxable 270, got %.2f", inv.TotalTaxableAmt)
	}

	// Total VAT: 270 * 0.15 = 40.5
	if inv.TotalVAT != 40.5 {
		t.Errorf("Expected total VAT 40.5, got %.2f", inv.TotalVAT)
	}

	// Total with VAT: 270 + 40.5 = 310.5
	if inv.TotalWithVAT != 310.5 {
		t.Errorf("Expected total with VAT 310.5, got %.2f", inv.TotalWithVAT)
	}
}
