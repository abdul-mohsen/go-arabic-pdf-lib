// Generates sample PDFs for all invoice types in both Arabic and English.
// Usage: go run cmd/samples/main.go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/invoice"
	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/models"
)

func main() {
	os.MkdirAll("output/samples", 0755)

	gen := invoice.NewGenerator(invoice.WithFontPath("fonts"))

	type sample struct {
		filename string
		inv      models.Invoice
	}

	samples := []sample{
		// ── B2C ──────────────────────────────────────────
		{"b2c_en_thermal", b2cInvoice(false, models.PaperThermal)},
		{"b2c_ar_thermal", b2cInvoice(true, models.PaperThermal)},
		{"b2c_en_a4", b2cInvoice(false, models.PaperA4)},
		{"b2c_ar_a4", b2cInvoice(true, models.PaperA4)},

		// ── B2B (always A4) ─────────────────────────────
		{"b2b_en_a4", b2bInvoice(false)},
		{"b2b_ar_a4", b2bInvoice(true)},

		// ── B2C Credit Note ─────────────────────────────
		{"credit_b2c_en", creditNote(false, models.InvoiceTypeB2CCredit)},
		{"credit_b2c_ar", creditNote(true, models.InvoiceTypeB2CCredit)},

		// ── B2B Credit Note ─────────────────────────────
		{"credit_b2b_en", creditNote(false, models.InvoiceTypeB2BCredit)},
		{"credit_b2b_ar", creditNote(true, models.InvoiceTypeB2BCredit)},

		// ── B2C Debit Note ──────────────────────────────
		{"debit_b2c_en", debitNote(false, models.InvoiceTypeB2CDebit)},
		{"debit_b2c_ar", debitNote(true, models.InvoiceTypeB2CDebit)},

		// ── B2B Debit Note ──────────────────────────────
		{"debit_b2b_en", debitNote(false, models.InvoiceTypeB2BDebit)},
		{"debit_b2b_ar", debitNote(true, models.InvoiceTypeB2BDebit)},
	}

	for _, s := range samples {
		path := fmt.Sprintf("output/samples/%s.pdf", s.filename)
		if err := gen.Generate(s.inv, path); err != nil {
			fmt.Printf("FAIL  %s: %v\n", s.filename, err)
		} else {
			fmt.Printf("OK    %s\n", path)
		}
	}

	fmt.Printf("\nDone — %d samples written to output/samples/\n", len(samples))
}

// ─── Helpers ────────────────────────────────────────────────────────────────

func langBuilder(arabic bool) *invoice.Builder {
	b := invoice.NewBuilder().
		WithDate(time.Date(2026, 3, 15, 14, 30, 0, 0, time.UTC)).
		WithDateFormat("2006-01-02 15:04")
	if arabic {
		b.WithArabic()
	} else {
		b.WithEnglish()
	}
	return b
}

func sellerEN() (string, string, string, string) {
	return "Alpha Trading Co.", "123 King Fahd Road, Riyadh", "300000000000003", "1010123456"
}

func sellerAR() (string, string, string, string) {
	return "شركة ألفا للتجارة", "شارع الملك فهد 123، الرياض", "300000000000003", "1010123456"
}

func buyerEN() (string, string, string, string) {
	return "Beta Industries Ltd.", "456 Tahlia St, Jeddah", "300000000000099", "4030789012"
}

func buyerAR() (string, string, string, string) {
	return "شركة بيتا للصناعات", "شارع التحلية 456، جدة", "300000000000099", "4030789012"
}

func seller(arabic bool) (string, string, string, string) {
	if arabic {
		return sellerAR()
	}
	return sellerEN()
}

func buyer(arabic bool) (string, string, string, string) {
	if arabic {
		return buyerAR()
	}
	return buyerEN()
}

func productName(arabic bool, enName, arName string) string {
	if arabic {
		return arName
	}
	return enName
}

func title(arabic bool, enTitle, arTitle string) string {
	if arabic {
		return arTitle
	}
	return enTitle
}

// ─── Invoice Builders ───────────────────────────────────────────────────────

func b2cInvoice(arabic bool, paper models.PaperSize) models.Invoice {
	b := langBuilder(arabic).
		WithType(models.InvoiceTypeB2C).
		WithPaper(paper).
		WithTitle(title(arabic, "Simplified Tax Invoice", "فاتورة ضريبية مبسطة")).
		WithInvoiceNumber("INV-2026-0042").
		WithSeller(seller(arabic)).
		WithVATPercentage("15.0").
		WithQRCode("AQlhbHBoYSBjbwIKMzAwMDAwMDAwMw==")

	b.AddProduct(productName(arabic, "Laptop Dell XPS 15", "لابتوب ديل XPS 15"), 1, 4500.00, 200.00, 645.00, 4945.00)
	b.AddProduct(productName(arabic, "Wireless Mouse", "ماوس لاسلكي"), 2, 75.00, 0, 22.50, 172.50)
	b.AddProduct(productName(arabic, "USB-C Hub", "موزع USB-C"), 1, 120.00, 10.00, 16.50, 126.50)

	b.WithTotals(210.00, 4510.00, 684.00, 5244.00)

	return b.Build()
}

func b2bInvoice(arabic bool) models.Invoice {
	b := langBuilder(arabic).
		WithType(models.InvoiceTypeB2B).
		WithTitle(title(arabic, "Tax Invoice", "فاتورة ضريبية")).
		WithInvoiceNumber("INV-2026-B2B-0018").
		WithSeller(seller(arabic)).
		WithBuyer(buyer(arabic)).
		WithVATPercentage(15).
		WithQRCode("AQlhbHBoYSBjbwIKMzAwMDAwMDAwMw==")

	// B2B needs SubtotalExclVAT — set directly on products
	inv := b.Build()
	inv.Products = []models.Product{
		{
			Name:     productName(arabic, "Enterprise Server Rack", "خادم مؤسسي"),
			Quantity: 2, UnitPrice: 12000.00, Discount: 1000.00,
			SubtotalExclVAT: 23000.00, VATAmount: 3450.00, Total: 26450.00,
		},
		{
			Name:     productName(arabic, "Annual Support Contract", "عقد دعم سنوي"),
			Quantity: 1, UnitPrice: 8000.00, Discount: 500.00,
			SubtotalExclVAT: 7500.00, VATAmount: 1125.00, Total: 8625.00,
		},
		{
			Name:     productName(arabic, "Network Switch 48-Port", "سويتش شبكة 48 منفذ"),
			Quantity: 4, UnitPrice: 1500.00, Discount: 0,
			SubtotalExclVAT: 6000.00, VATAmount: 900.00, Total: 6900.00,
		},
	}
	inv.TotalDiscount = 1500.00
	inv.TotalTaxableAmt = 36500.00
	inv.TotalVAT = 5475.00
	inv.TotalWithVAT = 41975.00
	return inv
}

func creditNote(arabic bool, invType models.InvoiceType) models.Invoice {
	b := langBuilder(arabic).
		WithType(invType).
		WithTitle(title(arabic, "Credit Note", "إشعار دائن")).
		WithInvoiceNumber("CN-2026-0007").
		WithSeller(seller(arabic)).
		WithNoteReason(title(arabic,
			"Returned goods — defective items from INV-2026-0042",
			"إرجاع بضاعة — أصناف معيبة من الفاتورة INV-2026-0042")).
		WithVATPercentage(15).
		WithQRCode("AQlhbHBoYSBjbwIKMzAwMDAwMDAwMw==")

	if invType == models.InvoiceTypeB2BCredit {
		b.WithBuyer(buyer(arabic))
		inv := b.Build()
		inv.Products = []models.Product{
			{
				Name:     productName(arabic, "Enterprise Server Rack", "خادم مؤسسي"),
				Quantity: 1, UnitPrice: 12000.00, Discount: 500.00,
				SubtotalExclVAT: 11500.00, VATAmount: 1725.00, Total: 13225.00,
			},
		}
		inv.TotalDiscount = 500.00
		inv.TotalTaxableAmt = 11500.00
		inv.TotalVAT = 1725.00
		inv.TotalWithVAT = 13225.00
		return inv
	}

	b.AddProduct(productName(arabic, "Wireless Mouse", "ماوس لاسلكي"), 2, 75.00, 0, 22.50, 172.50)
	b.WithTotals(0, 150.00, 22.50, 172.50)
	return b.Build()
}

func debitNote(arabic bool, invType models.InvoiceType) models.Invoice {
	b := langBuilder(arabic).
		WithType(invType).
		WithTitle(title(arabic, "Debit Note", "إشعار مدين")).
		WithInvoiceNumber("DN-2026-0003").
		WithSeller(seller(arabic)).
		WithNoteReason(title(arabic,
			"Price correction — additional charges for INV-2026-B2B-0018",
			"تصحيح سعر — رسوم إضافية للفاتورة INV-2026-B2B-0018")).
		WithVATPercentage(15).
		WithQRCode("AQlhbHBoYSBjbwIKMzAwMDAwMDAwMw==")

	if invType == models.InvoiceTypeB2BDebit {
		b.WithBuyer(buyer(arabic))
		inv := b.Build()
		inv.Products = []models.Product{
			{
				Name:     productName(arabic, "Shipping Surcharge", "رسوم شحن إضافية"),
				Quantity: 1, UnitPrice: 500.00, Discount: 0,
				SubtotalExclVAT: 500.00, VATAmount: 75.00, Total: 575.00,
			},
		}
		inv.TotalDiscount = 0
		inv.TotalTaxableAmt = 500.00
		inv.TotalVAT = 75.00
		inv.TotalWithVAT = 575.00
		return inv
	}

	b.AddProduct(productName(arabic, "Delivery Surcharge", "رسوم توصيل إضافية"), 1, 300.00, 0, 45.00, 345.00)
	b.WithTotals(0, 300.00, 45.00, 345.00)
	return b.Build()
}
