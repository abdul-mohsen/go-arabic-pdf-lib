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
			"english": true
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
		"config": {"vatPercentage": 15},
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

func TestParseJSON_DiscountPercent(t *testing.T) {
	jsonData := []byte(`{
		"config": {"vatPercentage": 15},
		"invoice": {"title": "Test", "invoiceNumber": "1", "storeName": "Store", "storeAddress": "Addr", "date": "2024/01/01", "vatRegistrationNo": "123", "qrCodeData": "qr"},
		"products": [
			{"name": "Product", "quantity": 1, "unitPrice": 100, "discountPercent": 10}
		],
		"labels": {"invoiceNumber": "", "date": "", "vatRegistration": "", "totalTaxable": "", "totalWithVat": "", "productColumn": "", "quantityColumn": "", "unitPriceColumn": "", "vatAmountColumn": "", "totalColumn": "", "footer": ""}
	}`)

	invoice, err := ParseJSON(jsonData)
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	p := invoice.Products[0]

	// Gross: 100, Discount: 10% = 10, Net: 90, VAT: 13.5, Total: 103.5
	if p.GrossAmount != 100 {
		t.Errorf("Expected gross 100, got %.2f", p.GrossAmount)
	}
	if p.DiscountAmount != 10 {
		t.Errorf("Expected discount 10, got %.2f", p.DiscountAmount)
	}
	if p.NetAmount != 90 {
		t.Errorf("Expected net 90, got %.2f", p.NetAmount)
	}
	if p.VATAmount != 13.5 {
		t.Errorf("Expected VAT 13.5, got %.2f", p.VATAmount)
	}
	if p.TotalWithVAT != 103.5 {
		t.Errorf("Expected total 103.5, got %.2f", p.TotalWithVAT)
	}

	// Check totals
	if invoice.TotalGross != 100 {
		t.Errorf("Expected total gross 100, got %.2f", invoice.TotalGross)
	}
	if invoice.TotalDiscount != 10 {
		t.Errorf("Expected total discount 10, got %.2f", invoice.TotalDiscount)
	}
}

func TestParseJSON_DiscountAmount(t *testing.T) {
	jsonData := []byte(`{
		"config": {"vatPercentage": 15},
		"invoice": {"title": "Test", "invoiceNumber": "1", "storeName": "Store", "storeAddress": "Addr", "date": "2024/01/01", "vatRegistrationNo": "123", "qrCodeData": "qr"},
		"products": [
			{"name": "Product", "quantity": 1, "unitPrice": 100, "discountAmount": 20}
		],
		"labels": {"invoiceNumber": "", "date": "", "vatRegistration": "", "totalTaxable": "", "totalWithVat": "", "productColumn": "", "quantityColumn": "", "unitPriceColumn": "", "vatAmountColumn": "", "totalColumn": "", "footer": ""}
	}`)

	invoice, err := ParseJSON(jsonData)
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	p := invoice.Products[0]

	// Gross: 100, Discount: 20, Net: 80, VAT: 12, Total: 92
	if p.DiscountAmount != 20 {
		t.Errorf("Expected discount 20, got %.2f", p.DiscountAmount)
	}
	if p.NetAmount != 80 {
		t.Errorf("Expected net 80, got %.2f", p.NetAmount)
	}
	if p.VATAmount != 12 {
		t.Errorf("Expected VAT 12, got %.2f", p.VATAmount)
	}
	if p.TotalWithVAT != 92 {
		t.Errorf("Expected total 92, got %.2f", p.TotalWithVAT)
	}
}

func TestParseJSON_CombinedDiscount(t *testing.T) {
	jsonData := []byte(`{
		"config": {"vatPercentage": 15},
		"invoice": {"title": "Test", "invoiceNumber": "1", "storeName": "Store", "storeAddress": "Addr", "date": "2024/01/01", "vatRegistrationNo": "123", "qrCodeData": "qr"},
		"products": [
			{"name": "Product", "quantity": 1, "unitPrice": 100, "discountPercent": 10, "discountAmount": 5}
		],
		"labels": {"invoiceNumber": "", "date": "", "vatRegistration": "", "totalTaxable": "", "totalWithVat": "", "productColumn": "", "quantityColumn": "", "unitPriceColumn": "", "vatAmountColumn": "", "totalColumn": "", "footer": ""}
	}`)

	invoice, err := ParseJSON(jsonData)
	if err != nil {
		t.Fatalf("ParseJSON failed: %v", err)
	}

	p := invoice.Products[0]

	// Gross: 100, Discount: 10% (10) + 5 = 15, Net: 85, VAT: 12.75, Total: 97.75
	if p.DiscountAmount != 15 {
		t.Errorf("Expected discount 15, got %.2f", p.DiscountAmount)
	}
	if p.NetAmount != 85 {
		t.Errorf("Expected net 85, got %.2f", p.NetAmount)
	}
}

