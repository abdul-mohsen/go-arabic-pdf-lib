// Package pdf provides PDF generation for invoices.
package pdf

import (
	"fmt"
	"os"

	"bill-generator/pkg/models"
	"bill-generator/pkg/textutil"

	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
)

// Generator handles PDF generation for invoices.
type Generator struct {
	pdf      gopdf.GoPdf
	fontDir  string
	invoice  models.Invoice
	pageW    float64
	margin   float64
	contentW float64
	currentY float64
}

// NewGenerator creates a new PDF generator.
func NewGenerator(fontDir string) *Generator {
	return &Generator{
		fontDir: fontDir,
		pageW:   226.77, // 80mm in points
		margin:  10.0,
	}
}

// Generate creates a PDF from the invoice and saves it to filename.
func (g *Generator) Generate(invoice models.Invoice, filename string) error {
	g.invoice = invoice
	g.contentW = g.pageW - (2 * g.margin)
	g.currentY = 10.0

	// Initialize PDF
	g.pdf = gopdf.GoPdf{}
	g.pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: 226.77, H: 708.66}, // 80mm x 250mm
	})

	// Load fonts
	if err := g.loadFonts(); err != nil {
		return err
	}

	g.pdf.AddPage()

	// Draw invoice sections
	g.drawHeader()
	g.drawInvoiceInfo()
	g.drawProductsTable()
	g.drawTotals()
	g.drawFooter()
	g.drawQRCode()

	return g.pdf.WritePdf(filename)
}

func (g *Generator) loadFonts() error {
	regularPath := g.fontDir + "/Amiri-Regular.ttf"
	if err := g.pdf.AddTTFFont("Amiri", regularPath); err != nil {
		return fmt.Errorf("failed to load regular font: %w", err)
	}

	boldPath := g.fontDir + "/Amiri-Bold.ttf"
	if err := g.pdf.AddTTFFont("AmiriBold", boldPath); err != nil {
		// Fallback to regular
		_ = g.pdf.AddTTFFont("AmiriBold", regularPath)
	}

	return nil
}

func (g *Generator) drawHeader() {
	if err := g.pdf.SetFont("AmiriBold", "", 14); err != nil {
		g.pdf.SetFont("Amiri", "", 14)
	}
	g.pdf.SetTextColor(0, 0, 0)
	textutil.DrawTextCentered(&g.pdf, g.invoice.Title, g.margin, g.currentY+4, g.contentW, g.invoice.IsRTL)
	g.currentY += 18
}

func (g *Generator) drawInvoiceInfo() {
	inv := g.invoice
	isRTL := inv.IsRTL

	// Invoice Number
	g.pdf.SetFont("Amiri", "", 9)
	g.pdf.SetTextColor(0, 0, 0)

	labelText := textutil.ProcessText(inv.Labels.InvoiceNumber, isRTL)
	labelW, _ := g.pdf.MeasureTextWidth(labelText)

	if isRTL {
		g.pdf.SetXY(g.margin+g.contentW-labelW-3, g.currentY)
		g.pdf.Cell(nil, labelText)
		g.pdf.SetXY(g.margin+3, g.currentY)
		g.pdf.Cell(nil, inv.InvoiceNumber)
	} else {
		g.pdf.SetXY(g.margin+3, g.currentY)
		g.pdf.Cell(nil, labelText)
		valueW, _ := g.pdf.MeasureTextWidth(inv.InvoiceNumber)
		g.pdf.SetXY(g.margin+g.contentW-valueW-3, g.currentY)
		g.pdf.Cell(nil, inv.InvoiceNumber)
	}
	g.currentY += 12

	// Store Name
	if err := g.pdf.SetFont("AmiriBold", "", 11); err != nil {
		g.pdf.SetFont("Amiri", "", 11)
	}
	textutil.DrawTextCentered(&g.pdf, inv.StoreName, g.margin, g.currentY, g.contentW, isRTL)
	g.currentY += 14

	// Store Address
	g.pdf.SetFont("Amiri", "", 9)
	textutil.DrawTextCentered(&g.pdf, inv.StoreAddress, g.margin, g.currentY, g.contentW, isRTL)
	g.currentY += 14

	// Date
	dateLabel := textutil.ProcessText(inv.Labels.Date, isRTL)
	dateLabelW, _ := g.pdf.MeasureTextWidth(dateLabel)

	if isRTL {
		g.pdf.SetXY(g.margin+g.contentW-dateLabelW-3, g.currentY)
		g.pdf.Cell(nil, dateLabel)
		g.pdf.SetXY(g.margin+3, g.currentY)
		g.pdf.Cell(nil, inv.Date)
	} else {
		g.pdf.SetXY(g.margin+3, g.currentY)
		g.pdf.Cell(nil, dateLabel)
		dateW, _ := g.pdf.MeasureTextWidth(inv.Date)
		g.pdf.SetXY(g.margin+g.contentW-dateW-3, g.currentY)
		g.pdf.Cell(nil, inv.Date)
	}
	g.currentY += 12

	// VAT Registration
	g.pdf.SetFont("Amiri", "", 8)
	vatLabel := textutil.ProcessText(inv.Labels.VATRegistration, isRTL)
	vatLabelW, _ := g.pdf.MeasureTextWidth(vatLabel)

	if isRTL {
		g.pdf.SetXY(g.margin+g.contentW-vatLabelW, g.currentY)
		g.pdf.Cell(nil, vatLabel)
		g.pdf.SetXY(g.margin, g.currentY)
		g.pdf.Cell(nil, inv.VATRegistrationNo)
	} else {
		g.pdf.SetXY(g.margin, g.currentY)
		g.pdf.Cell(nil, vatLabel)
		vatNoW, _ := g.pdf.MeasureTextWidth(inv.VATRegistrationNo)
		g.pdf.SetXY(g.margin+g.contentW-vatNoW, g.currentY)
		g.pdf.Cell(nil, inv.VATRegistrationNo)
	}
	g.currentY += 14
}

func (g *Generator) drawProductsTable() {
	inv := g.invoice
	isRTL := inv.IsRTL

	// Column widths (order depends on RTL)
	// For RTL: Total, VAT, Price, Qty, Product (right to left visually)
	// For LTR: Product, Qty, Price, VAT, Total (left to right visually)
	var colWidths []float64
	if isRTL {
		colWidths = []float64{30, 30, 30, 20, 96}
	} else {
		colWidths = []float64{96, 20, 30, 30, 30}
	}

	tableWidth := 0.0
	for _, w := range colWidths {
		tableWidth += w
	}
	tableX := g.margin

	// Draw header
	g.drawTableHeader(tableX, colWidths, isRTL)

	// Draw rows
	g.drawTableRows(tableX, colWidths, isRTL)
}

func (g *Generator) drawTableHeader(tableX float64, colWidths []float64, isRTL bool) {
	inv := g.invoice

	g.pdf.SetStrokeColor(0, 0, 0)
	g.pdf.SetLineWidth(0.5)

	headerHeight := 28.0
	xPos := tableX

	// Draw header cell borders
	for i := range colWidths {
		g.pdf.RectFromUpperLeftWithStyle(xPos, g.currentY, colWidths[i], headerHeight, "D")
		xPos += colWidths[i]
	}

	// Header text
	if err := g.pdf.SetFont("AmiriBold", "", 7); err != nil {
		g.pdf.SetFont("Amiri", "", 7)
	}
	g.pdf.SetTextColor(0, 0, 0)

	var headers [][]string
	if isRTL {
		headers = [][]string{
			{"السعر شامل", "الضريبة"},
			{"ضريبة القيمة", "المضافة"},
			{"سعر", "الوحدة"},
			{"", "الكمية"},
			{"", "المنتجات"},
		}
	} else {
		headers = [][]string{
			{"", "Product"},
			{"", "Qty"},
			{"Unit", "Price"},
			{"VAT", "Amount"},
			{"Total", "(inc. VAT)"},
		}
	}

	xPos = tableX
	for i, header := range headers {
		// Line 1
		if header[0] != "" {
			h1Text := textutil.ProcessText(header[0], isRTL)
			h1w, _ := g.pdf.MeasureTextWidth(h1Text)
			g.pdf.SetXY(xPos+(colWidths[i]-h1w)/2, g.currentY+4)
			g.pdf.Cell(nil, h1Text)
		}
		// Line 2
		h2Text := textutil.ProcessText(header[1], isRTL)
		h2w, _ := g.pdf.MeasureTextWidth(h2Text)
		g.pdf.SetXY(xPos+(colWidths[i]-h2w)/2, g.currentY+14)
		g.pdf.Cell(nil, h2Text)
		xPos += colWidths[i]
	}

	g.currentY += headerHeight
}

func (g *Generator) drawTableRows(tableX float64, colWidths []float64, isRTL bool) {
	inv := g.invoice

	g.pdf.SetFont("Amiri", "", 9)
	baseRowHeight := 12.0
	minRowHeight := 18.0

	for _, product := range inv.Products {
		g.pdf.SetStrokeColor(0, 0, 0)

		// Calculate row height based on product name wrapping
		var productColIdx int
		if isRTL {
			productColIdx = 4
		} else {
			productColIdx = 0
		}
		_, nameHeight := textutil.WrapText(&g.pdf, product.Name, colWidths[productColIdx]-6, baseRowHeight, isRTL)
		rowHeight := nameHeight + 6
		if rowHeight < minRowHeight {
			rowHeight = minRowHeight
		}

		// Draw row cell borders
		xPos := tableX
		for i := range colWidths {
			g.pdf.RectFromUpperLeftWithStyle(xPos, g.currentY, colWidths[i], rowHeight, "D")
			xPos += colWidths[i]
		}

		// Draw cell values
		g.pdf.SetTextColor(0, 0, 0)
		textY := g.currentY + 3

		if isRTL {
			g.drawRowCellsRTL(tableX, colWidths, textY, baseRowHeight, product)
		} else {
			g.drawRowCellsLTR(tableX, colWidths, textY, baseRowHeight, product)
		}

		g.currentY += rowHeight
	}

	g.currentY += 8
}

func (g *Generator) drawRowCellsRTL(tableX float64, colWidths []float64, textY, lineHeight float64, product models.Product) {
	xPos := tableX

	// Column 0: Total with VAT
	totalStr := fmt.Sprintf("%.1f", product.TotalWithVAT)
	tw, _ := g.pdf.MeasureTextWidth(totalStr)
	g.pdf.SetXY(xPos+colWidths[0]-tw-3, textY)
	g.pdf.Cell(nil, totalStr)
	xPos += colWidths[0]

	// Column 1: VAT Amount
	vatStr := fmt.Sprintf("%.1f", product.VATAmount)
	vw, _ := g.pdf.MeasureTextWidth(vatStr)
	g.pdf.SetXY(xPos+colWidths[1]-vw-3, textY)
	g.pdf.Cell(nil, vatStr)
	xPos += colWidths[1]

	// Column 2: Unit Price
	priceStr := fmt.Sprintf("%.0f", product.UnitPrice)
	pw, _ := g.pdf.MeasureTextWidth(priceStr)
	g.pdf.SetXY(xPos+colWidths[2]-pw-3, textY)
	g.pdf.Cell(nil, priceStr)
	xPos += colWidths[2]

	// Column 3: Quantity
	qtyStr := fmt.Sprintf("%.0f", product.Quantity)
	qw, _ := g.pdf.MeasureTextWidth(qtyStr)
	g.pdf.SetXY(xPos+colWidths[3]-qw-3, textY)
	g.pdf.Cell(nil, qtyStr)
	xPos += colWidths[3]

	// Column 4: Product Name
	textutil.DrawWrappedText(&g.pdf, product.Name, xPos, textY, colWidths[4], lineHeight, true)
}

func (g *Generator) drawRowCellsLTR(tableX float64, colWidths []float64, textY, lineHeight float64, product models.Product) {
	xPos := tableX

	// Column 0: Product Name
	textutil.DrawWrappedText(&g.pdf, product.Name, xPos, textY, colWidths[0], lineHeight, false)
	xPos += colWidths[0]

	// Column 1: Quantity
	qtyStr := fmt.Sprintf("%.0f", product.Quantity)
	qw, _ := g.pdf.MeasureTextWidth(qtyStr)
	g.pdf.SetXY(xPos+(colWidths[1]-qw)/2, textY)
	g.pdf.Cell(nil, qtyStr)
	xPos += colWidths[1]

	// Column 2: Unit Price
	priceStr := fmt.Sprintf("%.0f", product.UnitPrice)
	pw, _ := g.pdf.MeasureTextWidth(priceStr)
	g.pdf.SetXY(xPos+colWidths[2]-pw-3, textY)
	g.pdf.Cell(nil, priceStr)
	xPos += colWidths[2]

	// Column 3: VAT Amount
	vatStr := fmt.Sprintf("%.1f", product.VATAmount)
	vw, _ := g.pdf.MeasureTextWidth(vatStr)
	g.pdf.SetXY(xPos+colWidths[3]-vw-3, textY)
	g.pdf.Cell(nil, vatStr)
	xPos += colWidths[3]

	// Column 4: Total with VAT
	totalStr := fmt.Sprintf("%.1f", product.TotalWithVAT)
	tw, _ := g.pdf.MeasureTextWidth(totalStr)
	g.pdf.SetXY(xPos+colWidths[4]-tw-3, textY)
	g.pdf.Cell(nil, totalStr)
}

func (g *Generator) drawTotals() {
	inv := g.invoice
	isRTL := inv.IsRTL

	tableWidth := 206.0
	totalsX := g.margin
	valueWidth := 40.0
	labelWidth := tableWidth - valueWidth

	// Row 1: Taxable Amount
	g.pdf.SetStrokeColor(0, 0, 0)
	g.pdf.SetLineWidth(0.5)
	g.pdf.RectFromUpperLeftWithStyle(totalsX, g.currentY, valueWidth, 16, "D")
	g.pdf.RectFromUpperLeftWithStyle(totalsX+valueWidth, g.currentY, labelWidth, 16, "D")

	g.pdf.SetFont("Amiri", "", 9)
	g.pdf.SetTextColor(0, 0, 0)

	taxableStr := fmt.Sprintf("%.0f", inv.TotalTaxableAmt)
	taxableW, _ := g.pdf.MeasureTextWidth(taxableStr)
	g.pdf.SetXY(totalsX+valueWidth-taxableW-3, g.currentY+3)
	g.pdf.Cell(nil, taxableStr)

	taxableLbl := textutil.ProcessText(inv.Labels.TotalTaxable, isRTL)
	taxableLblW, _ := g.pdf.MeasureTextWidth(taxableLbl)

	if isRTL {
		g.pdf.SetXY(totalsX+valueWidth+labelWidth-taxableLblW-2, g.currentY)
	} else {
		g.pdf.SetXY(totalsX+valueWidth+3, g.currentY)
	}
	g.pdf.Cell(nil, taxableLbl)
	g.currentY += 16

	// Row 2: Total with VAT
	g.pdf.SetLineWidth(1.0)
	g.pdf.RectFromUpperLeftWithStyle(totalsX, g.currentY, valueWidth, 18, "D")
	g.pdf.RectFromUpperLeftWithStyle(totalsX+valueWidth, g.currentY, labelWidth, 18, "D")

	if err := g.pdf.SetFont("AmiriBold", "", 10); err != nil {
		g.pdf.SetFont("Amiri", "", 10)
	}

	totalStr := fmt.Sprintf("%.0f", inv.TotalWithVAT)
	totalStrW, _ := g.pdf.MeasureTextWidth(totalStr)
	g.pdf.SetXY(totalsX+valueWidth-totalStrW-2, g.currentY)
	g.pdf.Cell(nil, totalStr)

	totalLbl := textutil.ProcessText(inv.Labels.TotalWithVat, isRTL)
	totalPct := fmt.Sprintf("(%.0f%%)", inv.VATPercentage)
	totalLblW, _ := g.pdf.MeasureTextWidth(totalLbl)
	totalPctW, _ := g.pdf.MeasureTextWidth(totalPct)

	if isRTL {
		g.pdf.SetXY(totalsX+valueWidth+labelWidth-totalLblW-2, g.currentY)
		g.pdf.Cell(nil, totalLbl)
		g.pdf.SetXY(totalsX+valueWidth+labelWidth-totalLblW-totalPctW-2, g.currentY)
		g.pdf.Cell(nil, totalPct)
	} else {
		g.pdf.SetXY(totalsX+valueWidth+3, g.currentY)
		g.pdf.Cell(nil, totalLbl)
		g.pdf.SetXY(totalsX+valueWidth+totalLblW+6, g.currentY)
		g.pdf.Cell(nil, totalPct)
	}

	g.currentY += 22
}

func (g *Generator) drawFooter() {
	g.pdf.SetFont("Amiri", "", 7)
	g.pdf.SetTextColor(0, 0, 0)
	textutil.DrawTextCentered(&g.pdf, g.invoice.Labels.Footer, g.margin, g.currentY, g.contentW, g.invoice.IsRTL)
	g.currentY += 12
}

func (g *Generator) drawQRCode() {
	qrFile := "/tmp/temp_qr.png"
	err := qrcode.WriteFile(g.invoice.QRCodeData, qrcode.High, 256, qrFile)
	if err == nil {
		qrSize := 55.0
		qrX := g.margin + (g.contentW-qrSize)/2
		g.pdf.Image(qrFile, qrX, g.currentY, &gopdf.Rect{W: qrSize, H: qrSize})
		os.Remove(qrFile)
	}
}
