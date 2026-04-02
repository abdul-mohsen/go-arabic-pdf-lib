package sections

import (
	"fmt"

	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/textutil"
)

// DrawTotals draws the totals section (discount, taxable, VAT, grand total).
func DrawTotals(ctx *DrawContext) {
	inv := ctx.Invoice
	l := ctx.Layout
	pdf := ctx.PDF
	isRTL := inv.IsRTL

	tableWidth := l.ContentW
	totalsX := l.Margin
	valueWidth := 50.0
	if l.ContentW > 400 { // A4
		valueWidth = 80.0
	}
	labelWidth := tableWidth - valueWidth

	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetFont("Amiri", "", l.BodySize)
	pdf.SetTextColor(0, 0, 0)

	// Row: Total Discount (only if > 0)
	// if inv.TotalDiscount > 0 {
	// 	drawTotalsRow(ctx, totalsX, valueWidth, labelWidth, l.TotalsRowH,
	// 		inv.Labels.TotalDiscount, fmt.Sprintf("%.2f", inv.TotalDiscount), isRTL, false)
	// }
	drawTotalsRow(ctx, totalsX, valueWidth, labelWidth, l.TotalsRowH,
		inv.Labels.TotalDiscount, fmt.Sprintf("%.2f", inv.TotalDiscount), isRTL, false)

	// Row: Taxable Amount
	drawTotalsRow(ctx, totalsX, valueWidth, labelWidth, l.TotalsRowH,
		inv.Labels.TotalTaxable, fmt.Sprintf("%.2f", inv.TotalTaxableAmt), isRTL, false)

	// Row: VAT Amount — draw label and percentage separately to avoid RTL mangling
	vatLabel := or(inv.Labels.TotalVAT, "")
	if vatLabel == "" {
		if isRTL {
			vatLabel = "ضريبة القيمة المضافة"
		} else {
			vatLabel = "VAT Amount"
		}
	}
	pctSuffix := fmt.Sprintf(" (%.0f%%)", inv.VATPercentage)
	drawTotalsRowWithSuffix(ctx, totalsX, valueWidth, labelWidth, l.TotalsRowH,
		vatLabel, pctSuffix, fmt.Sprintf("%.2f", inv.TotalVAT), isRTL, false)

	// Row: Grand Total (bold, thicker border)
	pdf.SetLineWidth(1.0)
	setFont(pdf, true, l.BodySize+1)

	totalLabel := inv.Labels.TotalWithVat
	drawTotalsRowWithSuffix(ctx, totalsX, valueWidth, labelWidth, l.TotalsFinalRowH,
		totalLabel, pctSuffix, fmt.Sprintf("%.2f", inv.TotalWithVAT), isRTL, true)

	ctx.CurrentY += 4
}

// drawTotalsRowWithSuffix draws a totals row with a label (RTL-processed) and an LTR suffix (e.g. "(15%)").
func drawTotalsRowWithSuffix(ctx *DrawContext, x, valueW, labelW, rowH float64, label, suffix, value string, isRTL, bold bool) {
	pdf := ctx.PDF

	pdf.SetLineWidth(0.5)

	if isRTL {
		pdf.RectFromUpperLeftWithStyle(x, ctx.CurrentY, valueW, rowH, "D")
		pdf.RectFromUpperLeftWithStyle(x+valueW, ctx.CurrentY, labelW, rowH, "D")
	} else {
		pdf.RectFromUpperLeftWithStyle(x, ctx.CurrentY, labelW, rowH, "D")
		pdf.RectFromUpperLeftWithStyle(x+labelW, ctx.CurrentY, valueW, rowH, "D")
	}

	// Draw label (RTL processed)
	labelText := textutil.ProcessText(label, isRTL)
	labelTextW, _ := pdf.MeasureTextWidth(labelText)
	suffixW, _ := pdf.MeasureTextWidth(suffix)

	if isRTL {
		// Arabic: label right-aligned, suffix (LTR) drawn to the left of label
		pdf.SetXY(x+valueW+labelW-labelTextW-4, ctx.CurrentY+3)
		pdf.Cell(nil, labelText)
		pdf.SetXY(x+valueW+labelW-labelTextW-suffixW-5, ctx.CurrentY+3)
		pdf.Cell(nil, suffix)
	} else {
		// English: label left-aligned, suffix right after
		pdf.SetXY(x+4, ctx.CurrentY+3)
		pdf.Cell(nil, labelText)
		pdf.SetXY(x+4+labelTextW+1, ctx.CurrentY+3)
		pdf.Cell(nil, suffix)
	}

	// Draw value
	vw, _ := pdf.MeasureTextWidth(value)
	if isRTL {
		pdf.SetXY(x+valueW-vw-3, ctx.CurrentY+3)
	} else {
		pdf.SetXY(x+labelW+valueW-vw-3, ctx.CurrentY+3)
	}
	pdf.Cell(nil, value)

	ctx.CurrentY += rowH
}

func drawTotalsRow(ctx *DrawContext, x, valueW, labelW, rowH float64, label, value string, isRTL, bold bool) {
	pdf := ctx.PDF

	if bold {
		setFont(pdf, true, ctx.Layout.BodySize+1)
	}

	pdf.SetLineWidth(0.5)
	if bold {
		pdf.SetLineWidth(1.0)
	}

	if isRTL {
		// RTL: value on left, label on right
		pdf.RectFromUpperLeftWithStyle(x, ctx.CurrentY, valueW, rowH, "D")
		pdf.RectFromUpperLeftWithStyle(x+valueW, ctx.CurrentY, labelW, rowH, "D")
	} else {
		// LTR: label on left, value on right
		pdf.RectFromUpperLeftWithStyle(x, ctx.CurrentY, labelW, rowH, "D")
		pdf.RectFromUpperLeftWithStyle(x+labelW, ctx.CurrentY, valueW, rowH, "D")
	}

	// Draw label
	labelText := textutil.ProcessText(label, isRTL)
	labelTextW, _ := pdf.MeasureTextWidth(labelText)

	if isRTL {
		pdf.SetXY(x+valueW+labelW-labelTextW-4, ctx.CurrentY+3)
	} else {
		pdf.SetXY(x+4, ctx.CurrentY+3)
	}
	pdf.Cell(nil, labelText)

	// Draw value (right-aligned in value cell)
	vw, _ := pdf.MeasureTextWidth(value)
	if isRTL {
		pdf.SetXY(x+valueW-vw-3, ctx.CurrentY+3)
	} else {
		pdf.SetXY(x+labelW+valueW-vw-3, ctx.CurrentY+3)
	}
	pdf.Cell(nil, value)

	ctx.CurrentY += rowH

	if bold {
		pdf.SetFont("Amiri", "", ctx.Layout.BodySize)
		pdf.SetLineWidth(0.5)
	}
}
