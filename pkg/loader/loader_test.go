package loader

import (
	"testing"
)

func TestParseJSON(t *testing.T) {
	// All values are pre-calculated - this library is for visualization only
	jsonData := []byte(`{
		"config": {
			"vatPercentage": 15,
			"currencySymbol": "SAR",
			"dateFormat": "2006/01/02",
			"english": true
		},
		"invoice": {
			"title": "Test Invoice",
			"invoiceNumber": "INV001",
			"storeName": "Test Store",
			"storeAddress": "123 Test St",
			"date": "2024/01/15",
			"vatRegistrationNo": "123456789",
			"qrCodeData": "test-qr-data",
			"totalDiscount": 10.0,
			"totalTaxable": 190.0,
			"totalVat": 28.5,
			"totalWithVat": 218.5
		},
		"products": [
			{"name": "Product 1", "quantity": 2, "unitPrice": 50, "discount": 5.0, "vatAmount": 14.25, "total": 109.25},
			{"name": "Product 2", "quantity": 1, "unitPrice": 100, "discount": 5.0, "vatAmount": 14.25, "total": 109.25}
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
			"discountColumn": "Discount",
			"vatAmountColumn": "VAT",
			"totalColumn": "Total",
			"totalDiscount": "Total Discount",
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

	// Check products are passed through correctly
	if len(invoice.Products) != 2 {
		t.Fatalf("Expected 2 products, got %d", len(invoice.Products))
	}

	// Verify product values are passed through (not calculated)
	p1 := invoice.Products[0]
	if p1.Discount != 5.0 {
		t.Errorf("Expected discount 5.0, got %.2f", p1.Discount)
	}
	if p1.VATAmount != 14.25 {
		t.Errorf("Expected VAT 14.25, got %.2f", p1.VATAmount)
	}
	if p1.Total != 109.25 {
		t.Errorf("Expected total 109.25, got %.2f", p1.Total)
	}

	// Check totals are passed through
	if invoice.TotalDiscount != 10.0 {
		t.Errorf("Expected total discount 10.0, got %.2f", invoice.TotalDiscount)
	}
	if invoice.TotalTaxableAmt != 190.0 {
		t.Errorf("Expected total taxable 190.0, got %.2f", invoice.TotalTaxableAmt)
	}
	if invoice.TotalVAT != 28.5 {
		t.Errorf("Expected total VAT 28.5, got %.2f", invoice.TotalVAT)
	}
	if invoice.TotalWithVAT != 218.5 {
		t.Errorf("Expected total with VAT 218.5, got %.2f", invoice.TotalWithVAT)
	}
}

func TestParseJSON_ArabicRTL(t *testing.T) {
	jsonData := []byte(`{
		"config": {"vatPercentage": 15},
		"invoice": {"title": "فاتورة", "invoiceNumber": "1", "storeName": "متجر", "storeAddress": "عنوان", "date": "2024/01/01", "vatRegistrationNo": "123", "qrCodeData": "qr", "totalTaxable": 100, "totalVat": 15, "totalWithVat": 115},
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
		"invoice": {"title": "Test", "invoiceNumber": "1", "storeName": "Store", "storeAddress": "Addr", "date": "2024/01/01", "vatRegistrationNo": "123", "qrCodeData": "qr", "totalTaxable": 100, "totalVat": 15, "totalWithVat": 115},
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

func TestParseJSON_ProductWithDiscount(t *testing.T) {
	jsonData := []byte(`{
		"config": {"vatPercentage": 15},
		"invoice": {"title": "Test", "invoiceNumber": "1", "storeName": "Store", "storeAddress": "Addr", "date": "2024/01/01", "vatRegistrationNo": "123", "qrCodeData": "qr", "totalDiscount": 10.0, "totalTaxable": 90, "totalVat": 13.5, "totalWithVat": 103.5},
		"products": [
			{"name": "Product", "quantity": 1, "unitPrice": 100, "discount": 10, "vatAmount": 13.5, "total": 103.5}
		],
		"labels": {"invoiceNumber": "", "date": "", "vatRegistration": "", "totalTaxable": "", "totalWithVat": "", "productColumn": "", "quantityColumn": "", "unitPriceColumn": "", "vatAmountColumn": "", "totalColumn": "", "footer": ""}
	}`)

	invoice, err := ParseJSON(jsonData)
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	p := invoice.Products[0]

	// All values should be passed through exactly as provided
	if p.Discount != 10 {
		t.Errorf("Expected discount 10, got %.2f", p.Discount)
	}
	if p.VATAmount != 13.5 {
		t.Errorf("Expected VAT 13.5, got %.2f", p.VATAmount)
	}
	if p.Total != 103.5 {
		t.Errorf("Expected total 103.5, got %.2f", p.Total)
	}

	// Check totals passed through
	if invoice.TotalDiscount != 10.0 {
		t.Errorf("Expected total discount 10.0, got %.2f", invoice.TotalDiscount)
	}
}

func TestParseJSON_NoDiscount(t *testing.T) {
	jsonData := []byte(`{
		"config": {"vatPercentage": 15},
		"invoice": {"title": "Test", "invoiceNumber": "1", "storeName": "Store", "storeAddress": "Addr", "date": "2024/01/01", "vatRegistrationNo": "123", "qrCodeData": "qr", "totalTaxable": 100, "totalVat": 15, "totalWithVat": 115},
		"products": [
			{"name": "Product", "quantity": 1, "unitPrice": 100, "vatAmount": 15, "total": 115}
		],
		"labels": {"invoiceNumber": "", "date": "", "vatRegistration": "", "totalTaxable": "", "totalWithVat": "", "productColumn": "", "quantityColumn": "", "unitPriceColumn": "", "vatAmountColumn": "", "totalColumn": "", "footer": ""}
	}`)

	invoice, err := ParseJSON(jsonData)
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	p := invoice.Products[0]

	// Discount should be 0 when not provided
	if p.Discount != 0 {
		t.Errorf("Expected discount 0, got %.2f", p.Discount)
	}

	// Total discount should be 0
	if invoice.TotalDiscount != 0 {
		t.Errorf("Expected total discount 0, got %.2f", invoice.TotalDiscount)
	}
}
