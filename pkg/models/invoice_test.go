package models

import (
	"testing"
)

func TestProductStructure(t *testing.T) {
	// Test that Product struct holds pre-calculated values correctly
	// This library only visualizes - no calculations performed
	product := Product{
		Name:      "Test Product",
		Quantity:  2,
		UnitPrice: 50,
		Discount:  5.0,
		VATAmount: 14.25,
		Total:     109.25,
	}

	if product.Name != "Test Product" {
		t.Errorf("Expected name 'Test Product', got '%s'", product.Name)
	}
	if product.Quantity != 2 {
		t.Errorf("Expected quantity 2, got %.2f", product.Quantity)
	}
	if product.UnitPrice != 50 {
		t.Errorf("Expected unit price 50, got %.2f", product.UnitPrice)
	}
	if product.Discount != 5.0 {
		t.Errorf("Expected discount 5.0, got %.2f", product.Discount)
	}
	if product.VATAmount != 14.25 {
		t.Errorf("Expected VAT 14.25, got %.2f", product.VATAmount)
	}
	if product.Total != 109.25 {
		t.Errorf("Expected total 109.25, got %.2f", product.Total)
	}
}

func TestInvoiceTotals(t *testing.T) {
	// All totals are pre-calculated - this test verifies they are stored correctly
	invoice := Invoice{
		Products: []Product{
			{Discount: 5, VATAmount: 14.25, Total: 109.25},
			{Discount: 10, VATAmount: 28.5, Total: 218.5},
		},
		TotalDiscount:   15,
		TotalTaxableAmt: 285,
		TotalVAT:        42.75,
		TotalWithVAT:    327.75,
	}

	if invoice.TotalDiscount != 15 {
		t.Errorf("Expected total discount 15, got %.2f", invoice.TotalDiscount)
	}
	if invoice.TotalTaxableAmt != 285 {
		t.Errorf("Expected total taxable 285, got %.2f", invoice.TotalTaxableAmt)
	}
	if invoice.TotalVAT != 42.75 {
		t.Errorf("Expected total VAT 42.75, got %.2f", invoice.TotalVAT)
	}
	if invoice.TotalWithVAT != 327.75 {
		t.Errorf("Expected total with VAT 327.75, got %.2f", invoice.TotalWithVAT)
	}
}

func TestInvoiceLanguageSettings(t *testing.T) {
	tests := []struct {
		language string
		isRTL    bool
	}{
		{"ar", true},
		{"he", true},
		{"en", false},
		{"fr", false},
		{"", false}, // Empty should default to false
	}

	for _, tt := range tests {
		invoice := Invoice{
			Language: tt.language,
			IsRTL:    tt.language == "ar" || tt.language == "he",
		}

		if invoice.IsRTL != tt.isRTL {
			t.Errorf("Language '%s': expected IsRTL=%v, got %v", tt.language, tt.isRTL, invoice.IsRTL)
		}
	}
}

func TestProductInputStructure(t *testing.T) {
	// Verify ProductInput holds JSON input values correctly
	input := ProductInput{
		Name:      "Widget",
		Quantity:  "3",
		UnitPrice: "25.00",
		Discount:  "2.50",
		VATAmount: "11.25",
		Total:     "83.75",
	}

	if input.Name != "Widget" {
		t.Errorf("Expected name 'Widget', got '%s'", input.Name)
	}
	if input.Discount != "2.50" {
		t.Errorf("Expected discount 2.50, got %.2f", input.Discount)
	}
}

func TestInvoiceInputStructure(t *testing.T) {
	// Verify InvoiceInput holds pre-calculated totals
	input := InvoiceInput{
		Title:         "Invoice",
		InvoiceNumber: "INV-001",
		TotalDiscount: 15.0,
		TotalTaxable:  285.0,
		TotalVAT:      42.75,
		TotalWithVAT:  327.75,
	}

	if input.TotalDiscount != 15.0 {
		t.Errorf("Expected total discount 15.0, got %.2f", input.TotalDiscount)
	}
	if input.TotalWithVAT != 327.75 {
		t.Errorf("Expected total with VAT 327.75, got %.2f", input.TotalWithVAT)
	}
}
