package sections

import (
	"fmt"

	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/models"
	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/textutil"

	"github.com/signintech/gopdf"
)

// DrawProductsTable draws the products table with header and rows.
func DrawProductsTable(ctx *DrawContext) {
	inv := ctx.Invoice
	l := ctx.Layout
	isRTL := inv.IsRTL
	colWidths := l.ColWidths
	isB2B := inv.Type == models.InvoiceTypeB2B || inv.Type == models.InvoiceTypeB2BCredit || inv.Type == models.InvoiceTypeB2BDebit
	tableX := l.Margin

	drawTableHeader(ctx, tableX, colWidths, isRTL, isB2B)
	drawTableRows(ctx, tableX, colWidths, isRTL, isB2B)
}

func drawTableHeader(ctx *DrawContext, tableX float64, colWidths []float64, isRTL, isB2B bool) {
	pdf := ctx.PDF
	l := ctx.Layout

	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.5)

	headerHeight := l.HeaderRowH
	xPos := tableX

	// Draw header cell borders
	for i := range colWidths {
		pdf.RectFromUpperLeftWithStyle(xPos, ctx.CurrentY, colWidths[i], headerHeight, "D")
		xPos += colWidths[i]
	}

	setFont(pdf, true, l.TableHeadSize)
	pdf.SetTextColor(0, 0, 0)

	headers := buildHeaders(ctx.Invoice, isB2B)

	xPos = tableX
	for i, header := range headers {
		text := header[0]
		drawHeaderCellWrapped(pdf, text, xPos, ctx.CurrentY, colWidths[i], headerHeight, l.TableHeadSize, isRTL)
		xPos += colWidths[i]
	}

	ctx.CurrentY += headerHeight
}

// drawHeaderCellWrapped draws header text, auto-wrapping to multiple lines if it overflows the cell.
func drawHeaderCellWrapped(pdf *gopdf.GoPdf, text string, x, y, colW, cellH, fontSize float64, isRTL bool) {
	processed := textutil.ProcessText(text, isRTL)
	tw, _ := pdf.MeasureTextWidth(processed)
	pad := 2.0
	availW := colW - 2*pad

	if tw <= availW {
		// Single line — vertically center
		pdf.SetXY(x+(colW-tw)/2, y+(cellH-fontSize)/2)
		pdf.Cell(nil, processed)
		return
	}

	// Word-wrap the text
	words := splitWords(text)
	var lines []string
	currentLine := ""
	for _, word := range words {
		candidate := currentLine
		if candidate != "" {
			candidate += " "
		}
		candidate += word
		procCand := textutil.ProcessText(candidate, isRTL)
		cw, _ := pdf.MeasureTextWidth(procCand)
		if cw > availW && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = candidate
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	lineH := fontSize + 2
	totalH := float64(len(lines)) * lineH
	startY := y + (cellH-totalH)/2

	for _, line := range lines {
		procLine := textutil.ProcessText(line, isRTL)
		lw, _ := pdf.MeasureTextWidth(procLine)
		pdf.SetXY(x+(colW-lw)/2, startY)
		pdf.Cell(nil, procLine)
		startY += lineH
	}
}

// splitWords splits text by spaces.
func splitWords(s string) []string {
	var words []string
	current := ""
	for _, r := range s {
		if r == ' ' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		words = append(words, current)
	}
	return words
}

func buildHeaders(inv models.Invoice, isB2B bool) [][]string {
	labels := inv.Labels
	if inv.IsRTL {
		return buildHeadersRTL(labels, isB2B)
	}
	return buildHeadersLTR(labels, isB2B)
}

func buildHeadersLTR(labels models.Labels, isB2B bool) [][]string {
	headers := [][]string{
		{or(labels.ProductColumn, "Product")},
		{or(labels.QuantityColumn, "Qty")},
		{or(labels.UnitPriceColumn, "Unit Price")},
		{or(labels.DiscountColumn, "Discount")},
	}
	if isB2B {
		headers = append(headers, []string{or(labels.SubtotalExclVATColumn, "Subtotal (excl. VAT)")})
	}
	headers = append(headers,
		[]string{or(labels.VATAmountColumn, "VAT")},
		[]string{or(labels.TotalColumn, "Total")},
	)
	return headers
}

func buildHeadersRTL(labels models.Labels, isB2B bool) [][]string {
	// RTL: columns are reversed (rightmost = first column visually)
	headers := [][]string{
		{or(labels.TotalColumn, "السعر شامل الضريبة")},
		{or(labels.VATAmountColumn, "ضريبة القيمة المضافة")},
	}
	if isB2B {
		headers = append(headers, []string{or(labels.SubtotalExclVATColumn, "المجموع الفرعي بدون الضريبة")})
	}
	headers = append(headers,
		[]string{or(labels.DiscountColumn, "الخصم")},
		[]string{or(labels.UnitPriceColumn, "سعر الوحدة")},
		[]string{or(labels.QuantityColumn, "الكمية")},
		[]string{or(labels.ProductColumn, "المنتجات")},
	)
	return headers
}

func or(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func drawTableRows(ctx *DrawContext, tableX float64, colWidths []float64, isRTL, isB2B bool) {
	inv := ctx.Invoice
	l := ctx.Layout
	pdf := ctx.PDF

	pdf.SetFont("Amiri", "", l.BodySize)
	baseRowHeight := l.LineHeight
	minRowHeight := l.RowMinHeight

	for _, product := range inv.Products {
		pdf.SetStrokeColor(0, 0, 0)

		// Product column index
		var productColIdx int
		if isRTL {
			productColIdx = len(colWidths) - 1
		} else {
			productColIdx = 0
		}
		_, nameHeight := textutil.WrapText(pdf, product.Name, colWidths[productColIdx]-6, baseRowHeight, isRTL)
		rowHeight := nameHeight + 6
		if rowHeight < minRowHeight {
			rowHeight = minRowHeight
		}

		// Draw row cell borders
		xPos := tableX
		for i := range colWidths {
			pdf.RectFromUpperLeftWithStyle(xPos, ctx.CurrentY, colWidths[i], rowHeight, "D")
			xPos += colWidths[i]
		}

		// Draw cell values
		pdf.SetTextColor(0, 0, 0)
		textY := ctx.CurrentY + 3

		if isRTL {
			drawRowCellsRTL(ctx, tableX, colWidths, textY, baseRowHeight, product, isB2B)
		} else {
			drawRowCellsLTR(ctx, tableX, colWidths, textY, baseRowHeight, product, isB2B)
		}

		ctx.CurrentY += rowHeight
	}

	ctx.CurrentY += 8
}

func drawRowCellsLTR(ctx *DrawContext, tableX float64, colWidths []float64, textY, lineHeight float64, product models.Product, isB2B bool) {
	pdf := ctx.PDF
	fs := ctx.Layout.BodySize
	xPos := tableX
	col := 0

	// Product Name
	textutil.DrawWrappedText(pdf, product.Name, xPos, textY, colWidths[col], lineHeight, false)
	xPos += colWidths[col]
	col++

	// Quantity
	drawCellNumberSized(pdf, product.Quantity, xPos, colWidths[col], textY, true, fs)
	xPos += colWidths[col]
	col++

	// Unit Price
	drawCellNumberSized(pdf, product.UnitPrice, xPos, colWidths[col], textY, false, fs)
	xPos += colWidths[col]
	col++

	// Discount
	// if product.Discount > 0 {
	// }
	drawCellNumberSized(pdf, product.Discount, xPos, colWidths[col], textY, false, fs)
	xPos += colWidths[col]
	col++

	// Subtotal excl VAT (B2B only)
	if isB2B {
		drawCellNumberSized(pdf, product.SubtotalExclVAT, xPos, colWidths[col], textY, false, fs)
		xPos += colWidths[col]
		col++
	}

	// VAT Amount
	drawCellNumberSized(pdf, product.VATAmount, xPos, colWidths[col], textY, false, fs)
	xPos += colWidths[col]
	col++

	// Total
	drawCellNumberSized(pdf, product.Total, xPos, colWidths[col], textY, false, fs)
}

func drawRowCellsRTL(ctx *DrawContext, tableX float64, colWidths []float64, textY, lineHeight float64, product models.Product, isB2B bool) {
	pdf := ctx.PDF
	fs := ctx.Layout.BodySize
	xPos := tableX
	col := 0

	// Column 0: Total (leftmost in RTL visual)
	drawCellNumberSized(pdf, product.Total, xPos, colWidths[col], textY, false, fs)
	xPos += colWidths[col]
	col++

	// Column 1: VAT Amount
	drawCellNumberSized(pdf, product.VATAmount, xPos, colWidths[col], textY, false, fs)
	xPos += colWidths[col]
	col++

	// Subtotal excl VAT (B2B only)
	if isB2B {
		drawCellNumberSized(pdf, product.SubtotalExclVAT, xPos, colWidths[col], textY, false, fs)
		xPos += colWidths[col]
		col++
	}

	// Discount
	// if product.Discount > 0 {
	// 	drawCellNumberSized(pdf, fmt.Sprintf("%.2f", product.Discount), xPos, colWidths[col], textY, false, fs)
	// }
	drawCellNumberSized(pdf, product.Discount, xPos, colWidths[col], textY, false, fs)
	xPos += colWidths[col]
	col++

	// Unit Price
	drawCellNumberSized(pdf, product.UnitPrice, xPos, colWidths[col], textY, false, fs)
	xPos += colWidths[col]
	col++

	// Quantity
	drawCellNumberSized(pdf, product.Quantity, xPos, colWidths[col], textY, false, fs)
	xPos += colWidths[col]
	col++

	// Product Name (rightmost in RTL)
	textutil.DrawWrappedText(pdf, product.Name, xPos, textY, colWidths[col], lineHeight, true)
}

// drawCellNumber draws a right-aligned number in a cell. If center is true, it centers.
// Auto-scales the font down if the text exceeds the available cell width.
func drawCellNumber(pdf *gopdf.GoPdf, text string, xPos, colWidth, textY float64, center bool) {
	drawCellNumberSized(pdf, text, xPos, colWidth, textY, center, 0)
}

// drawCellNumberSized draws a number with auto-scale-down. fontSize=0 means use current size.
func drawCellNumberSized(pdf *gopdf.GoPdf, text string, xPos, colWidth, textY float64, center bool, fontSize float64) {
	pad := 3.0
	availW := colWidth - 2*pad

	w, _ := pdf.MeasureTextWidth(text)

	// Auto-scale down if text overflows the cell
	if w > availW && fontSize > 0 {
		for trySize := fontSize - 1; trySize >= 5; trySize-- {
			pdf.SetFontSize(trySize)
			w, _ = pdf.MeasureTextWidth(text)
			if w <= availW {
				break
			}
		}
		defer pdf.SetFontSize(fontSize)
	}

	if center {
		pdf.SetXY(xPos+(colWidth-w)/2, textY)
	} else {
		pdf.SetXY(xPos+colWidth-w-pad, textY)
	}
	pdf.Cell(nil, text)
}
