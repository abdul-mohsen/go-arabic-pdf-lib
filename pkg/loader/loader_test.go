package loader

import (
	"testing"
)

func TestParseJSON(t *testing.T) {
	jsonData := []byte(`{
		"config": {
			"vatPercentage": 15,
			"currencySymbol": "SAR",
			"dateFormat": "2006/01/02",
			"language": "en"
		},
		"invoice": {
			"title": "Test Invoice",
			"invoiceNumber": "INV001",
			"storeName": "Test Store",
			"storeAddress": "123 Test St",
			"date": "2024/01/15",
			"vatRegistrationNo": "123456789",
			"qrCodeData": "test-qr-data"
		},
		"products": [
			{"name": "Product 1", "quantity": 2, "unitPrice": 50},
			{"name": "Product 2", "quantity": 1, "unitPrice": 100}
		],
		"labels": {
			"invoiceNumber": "Invoice #:",
			"date": "Date:",
			"vatRegistration": "VAT Reg:",
			"totalTaxable": "Subtotal",
			"totalWithVat": "Total",
			"productColumn": "Product",
			"quantityColumn": "Qty",
			"unitPriceColumn": "Price",
			"vatAmountColumn": "VAT",
			"totalColumn": "Total",
			"footer": "Thank you!"
		}
	}`)

	invoice, err := ParseJSON(jsonData)
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	// Check basic fields
	if invoice.Title != "Test Invoice" {
		t.Errorf("Expected title 'Test Invoice', got '%s'", invoice.Title)
	}

	if invoice.Language != "en" {
		t.Errorf("Expected language 'en', got '%s'", invoice.Language)
	}

	if invoice.IsRTL {
		t.Error("Expected IsRTL to be false for English")
	}

	// Check product calculations
	if len(invoice.Products) != 2 {
		t.Fatalf("Expected 2 products, got %d", len(invoice.Products))
	}

	// Product 1: 2 * 50 = 100, VAT = 15, Total = 115
	if invoice.Products[0].TaxableAmt != 100 {
		t.Errorf("Expected taxable 100, got %.2f", invoice.Products[0].TaxableAmt)
	}
	if invoice.Products[0].VATAmount != 15 {
		t.Errorf("Expected VAT 15, got %.2f", invoice.Products[0].VATAmount)
	}
	if invoice.Products[0].TotalWithVAT != 115 {
		t.Errorf("Expected total 115, got %.2f", invoice.Products[0].TotalWithVAT)
	}

	// Product 2: 1 * 100 = 100, VAT = 15, Total = 115
	if invoice.Products[1].TotalWithVAT != 115 {
		t.Errorf("Expected total 115, got %.2f", invoice.Products[1].TotalWithVAT)
	}

	// Check totals: 200 taxable, 30 VAT, 230 total
	if invoice.TotalTaxableAmt != 200 {
		t.Errorf("Expected total taxable 200, got %.2f", invoice.TotalTaxableAmt)
	}
	if invoice.TotalVAT != 30 {
		t.Errorf("Expected total VAT 30, got %.2f", invoice.TotalVAT)
	}
	if invoice.TotalWithVAT != 230 {
		t.Errorf("Expected total with VAT 230, got %.2f", invoice.TotalWithVAT)
	}
}

func TestParseJSON_ArabicRTL(t *testing.T) {
	jsonData := []byte(`{
		"config": {"vatPercentage": 15, "language": "ar"},
		"invoice": {"title": "فاتورة", "invoiceNumber": "1", "storeName": "متجر", "storeAddress": "عنوان", "date": "2024/01/01", "vatRegistrationNo": "123", "qrCodeData": "qr"},
		"products": [],
		"labels": {"invoiceNumber": "", "date": "", "vatRegistration": "", "totalTaxable": "", "totalWithVat": "", "productColumn": "", "quantityColumn": "", "unitPriceColumn": "", "vatAmountColumn": "", "totalColumn": "", "footer": ""}
	}`)

	invoice, err := ParseJSON(jsonData)
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	if invoice.Language != "ar" {
		t.Errorf("Expected language 'ar', got '%s'", invoice.Language)
	}

	if !invoice.IsRTL {
		t.Error("Expected IsRTL to be true for Arabic")
	}
}

func TestParseJSON_DefaultLanguage(t *testing.T) {
	jsonData := []byte(`{
		"config": {"vatPercentage": 15},
		"invoice": {"title": "Test", "invoiceNumber": "1", "storeName": "Store", "storeAddress": "Addr", "date": "2024/01/01", "vatRegistrationNo": "123", "qrCodeData": "qr"},
		"products": [],
		"labels": {"invoiceNumber": "", "date": "", "vatRegistration": "", "totalTaxable": "", "totalWithVat": "", "productColumn": "", "quantityColumn": "", "unitPriceColumn": "", "vatAmountColumn": "", "totalColumn": "", "footer": ""}
	}`)

	invoice, err := ParseJSON(jsonData)
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	// Default should be Arabic
	if invoice.Language != "ar" {
		t.Errorf("Expected default language 'ar', got '%s'", invoice.Language)
	}

	if !invoice.IsRTL {
		t.Error("Expected IsRTL to be true for default Arabic")
	}
}

func TestParseJSON_InvalidJSON(t *testing.T) {
	jsonData := []byte(`{invalid json}`)

	_, err := ParseJSON(jsonData)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}
