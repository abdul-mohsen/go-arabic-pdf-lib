// Package pdf provides PDF generation for invoices.
package pdf

import "github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/models"

// Layout holds all dimensions and font sizes for a given paper size.
type Layout struct {
	// Page
	PageW    float64
	PageH    float64
	Margin   float64
	ContentW float64

	// Fonts
	TitleSize     float64
	HeadingSize   float64
	BodySize      float64
	SmallSize     float64
	TableHeadSize float64
	FooterSize    float64

	// Spacing
	SectionGap      float64
	LineHeight      float64
	RowMinHeight    float64
	HeaderRowH      float64
	TotalsRowH      float64
	TotalsFinalRowH float64

	// QR
	QRSize float64

	// Table columns for B2C: Product, Qty, Price, Discount, VAT, Total
	B2CColWidths []float64
	// Table columns for B2B: Product, Qty, Price, Discount, Subtotal, VAT, Total
	B2BColWidths []float64
}

// ThermalLayout returns the layout for 80mm thermal receipt paper.
func ThermalLayout() Layout {
	pageW := 226.77 // 80mm
	margin := 10.0
	contentW := pageW - (2 * margin)
	return Layout{
		PageW:    pageW,
		PageH:    708.66, // 250mm
		Margin:   margin,
		ContentW: contentW,

		TitleSize:     14,
		HeadingSize:   11,
		BodySize:      9,
		SmallSize:     7,
		TableHeadSize: 6,
		FooterSize:    7,

		SectionGap:      14,
		LineHeight:      12,
		RowMinHeight:    18,
		HeaderRowH:      32,
		TotalsRowH:      16,
		TotalsFinalRowH: 18,

		QRSize: 55,

		// B2C: Product(70), Qty(20), Price(30), Discount(26), VAT(30), Total(30) ≈ 206
		B2CColWidths: []float64{70, 20, 30, 26, 30, 30},
		// B2B not supported on thermal — use A4
		B2BColWidths: nil,
	}
}

// A4Layout returns the layout for A4 paper.
func A4Layout() Layout {
	pageW := 595.28 // 210mm
	margin := 40.0
	contentW := pageW - (2 * margin)
	return Layout{
		PageW:    pageW,
		PageH:    841.89, // 297mm
		Margin:   margin,
		ContentW: contentW,

		TitleSize:     20,
		HeadingSize:   13,
		BodySize:      10,
		SmallSize:     8,
		TableHeadSize: 8,
		FooterSize:    8,

		SectionGap:      16,
		LineHeight:      14,
		RowMinHeight:    22,
		HeaderRowH:      32,
		TotalsRowH:      20,
		TotalsFinalRowH: 24,

		QRSize: 80,

		// B2C on A4: Product(180), Qty(45), Price(65), Discount(60), VAT(65), Total(100) = 515
		B2CColWidths: []float64{180, 45, 65, 60, 65, 100},
		// B2B on A4: Product(120), Qty(35), Price(60), Discount(55), Subtotal(80), VAT(65), Total(100) = 515
		B2BColWidths: []float64{120, 35, 60, 55, 80, 65, 100},
	}
}

// LayoutForInvoice returns the appropriate layout based on the invoice config.
func LayoutForInvoice(inv models.Invoice) Layout {
	if inv.PaperSize == models.PaperA4 {
		return A4Layout()
	}
	return ThermalLayout()
}

// ColWidths returns the column widths appropriate for the invoice type.
// For RTL, the returned widths are reversed (right-to-left visual order).
func (l Layout) ColWidths(inv models.Invoice) []float64 {
	var widths []float64
	if inv.Type == models.InvoiceTypeB2B || inv.Type == models.InvoiceTypeB2BCredit || inv.Type == models.InvoiceTypeB2BDebit {
		widths = make([]float64, len(l.B2BColWidths))
		copy(widths, l.B2BColWidths)
	} else {
		widths = make([]float64, len(l.B2CColWidths))
		copy(widths, l.B2CColWidths)
	}
	if inv.IsRTL {
		// Reverse column order for RTL
		for i, j := 0, len(widths)-1; i < j; i, j = i+1, j-1 {
			widths[i], widths[j] = widths[j], widths[i]
		}
	}
	return widths
}

// IsB2B returns true if the invoice type is a B2B variant.
func IsB2B(inv models.Invoice) bool {
	return inv.Type == models.InvoiceTypeB2B || inv.Type == models.InvoiceTypeB2BCredit || inv.Type == models.InvoiceTypeB2BDebit
}

// IsCreditOrDebit returns true if the invoice is a credit or debit note.
func IsCreditOrDebit(inv models.Invoice) bool {
	switch inv.Type {
	case models.InvoiceTypeB2CCredit, models.InvoiceTypeB2BCredit, models.InvoiceTypeB2CDebit, models.InvoiceTypeB2BDebit:
		return true
	}
	return false
}
