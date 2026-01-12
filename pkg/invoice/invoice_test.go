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
		WithTotals(10.0, 190.0, 28.5, 218.5).
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

func TestBuilder_Totals(t *testing.T) {
	inv := NewBuilder().
		WithTotals(10.0, 190.0, 28.5, 218.5).
		Build()

	if inv.TotalDiscount != 10.0 {
		t.Errorf("Expected total discount 10.0, got %.2f", inv.TotalDiscount)
	}
	if inv.TotalTaxableAmt != 190.0 {
		t.Errorf("Expected total taxable 190.0, got %.2f", inv.TotalTaxableAmt)
	}
	if inv.TotalVAT != 28.5 {
		t.Errorf("Expected total VAT 28.5, got %.2f", inv.TotalVAT)
	}
	if inv.TotalWithVAT != 218.5 {
		t.Errorf("Expected total with VAT 218.5, got %.2f", inv.TotalWithVAT)
	}
}

func TestBuilder_AddProduct(t *testing.T) {
	inv := NewBuilder().
		AddProduct("Product 1", 2, 50.00, 5.0, 14.25, 109.25).
		Build()

	if len(inv.Products) != 1 {
		t.Fatalf("Expected 1 product, got %d", len(inv.Products))
	}

	p := inv.Products[0]
	if p.Name != "Product 1" {
		t.Errorf("Expected name 'Product 1', got '%s'", p.Name)
	}
	if p.Quantity != 2 {
		t.Errorf("Expected quantity 2, got %.2f", p.Quantity)
	}
	if p.UnitPrice != 50.00 {
		t.Errorf("Expected unit price 50, got %.2f", p.UnitPrice)
	}
	if p.Discount != 5.0 {
		t.Errorf("Expected discount 5.0, got %.2f", p.Discount)
	}
	if p.VATAmount != 14.25 {
		t.Errorf("Expected VAT 14.25, got %.2f", p.VATAmount)
	}
	if p.Total != 109.25 {
		t.Errorf("Expected total 109.25, got %.2f", p.Total)
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
		AddProduct("Product 1", 2, 50.00, 0, 15.0, 115.0).
		AddProduct("Product 2", 1, 100.00, 5.0, 14.25, 109.25).
		AddProduct("Product 3", 1, 80.00, 10.0, 10.5, 80.5).
		WithTotals(15.0, 215.0, 39.75, 304.75).
		Build()

	if len(inv.Products) != 3 {
		t.Fatalf("Expected 3 products, got %d", len(inv.Products))
	}

	// Check totals are passed through
	if inv.TotalDiscount != 15.0 {
		t.Errorf("Expected total discount 15.0, got %.2f", inv.TotalDiscount)
	}
	if inv.TotalWithVAT != 304.75 {
		t.Errorf("Expected total with VAT 304.75, got %.2f", inv.TotalWithVAT)
	}
}
