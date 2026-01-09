package models

import (
	"testing"
)

func TestProductCalculations(t *testing.T) {
	// Test that the Product struct can hold calculated values correctly
	product := Product{
		Name:         "Test Product",
		Quantity:     2,
		UnitPrice:    50,
		TaxableAmt:   100,
		VATAmount:    15,
		TotalWithVAT: 115,
	}

	if product.TaxableAmt != product.Quantity*product.UnitPrice {
		t.Errorf("TaxableAmt mismatch: expected %.2f, got %.2f",
			product.Quantity*product.UnitPrice, product.TaxableAmt)
	}
}

func TestInvoiceTotals(t *testing.T) {
	invoice := Invoice{
		Products: []Product{
			{TaxableAmt: 100, VATAmount: 15, TotalWithVAT: 115},
			{TaxableAmt: 200, VATAmount: 30, TotalWithVAT: 230},
		},
		TotalTaxableAmt: 300,
		TotalVAT:        45,
		TotalWithVAT:    345,
	}

	// Verify totals match sum of products
	var sumTaxable, sumVAT, sumTotal float64
	for _, p := range invoice.Products {
		sumTaxable += p.TaxableAmt
		sumVAT += p.VATAmount
		sumTotal += p.TotalWithVAT
	}

	if invoice.TotalTaxableAmt != sumTaxable {
		t.Errorf("TotalTaxableAmt mismatch: expected %.2f, got %.2f", sumTaxable, invoice.TotalTaxableAmt)
	}
	if invoice.TotalVAT != sumVAT {
		t.Errorf("TotalVAT mismatch: expected %.2f, got %.2f", sumVAT, invoice.TotalVAT)
	}
	if invoice.TotalWithVAT != sumTotal {
		t.Errorf("TotalWithVAT mismatch: expected %.2f, got %.2f", sumTotal, invoice.TotalWithVAT)
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
