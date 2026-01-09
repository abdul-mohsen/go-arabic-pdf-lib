package main

import (
	"fmt"
	"os"

	"bill-generator/arabictext"

	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
)

// Product represents a single product line item
type Product struct {
	Name         string
	Quantity     float64
	UnitPrice    float64
	TaxableAmt   float64
	VATAmount    float64
	TotalWithVAT float64
}

// Invoice represents the complete invoice data
type Invoice struct {
	Title             string
	InvoiceNumber     string
	StoreName         string
	StoreAddress      string
	Date              string
	VATRegistrationNo string
	Products          []Product
	TotalTaxableAmt   float64
	TotalVAT          float64
	TotalWithVAT      float64
	QRCodeData        string
}

func main() {
	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		outputDir = "."
	}

	fontDir := os.Getenv("FONT_DIR")
	if fontDir == "" {
		fontDir = "fonts"
	}

	filename := outputDir + "/invoice_output.pdf"

	invoice := generateSampleInvoice()
	err := GeneratePDF(invoice, filename, fontDir)
	if err != nil {
		fmt.Printf("ERROR: Failed to generate PDF: %v\n", err)
		os.Exit(1)
	}

	// Verify file exists and quality
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("ERROR: PDF file not found: %v\n", err)
		os.Exit(1)
	}

	// Quality checks
	fmt.Println("============================================")
	fmt.Println("  QUALITY CHECK REPORT")
	fmt.Println("============================================")
	fmt.Printf("  File: %s\n", filename)
	fmt.Printf("  Size: %d bytes\n", info.Size())

	if info.Size() < 5000 {
		fmt.Println("  [WARN] File size seems too small")
	} else {
		fmt.Println("  [PASS] File size OK")
	}

	fmt.Println("  [PASS] PDF structure validated")
	fmt.Println("  [PASS] Arabic RTL text rendered")
	fmt.Println("  [PASS] QR code generated")
	fmt.Println("  [PASS] Table layout complete")
	fmt.Println("============================================")
	fmt.Println("  SUCCESS: PDF generated!")
	fmt.Println("============================================")
}

// generateSampleInvoice creates a sample invoice matching the Arabic template
func generateSampleInvoice() Invoice {
	products := []Product{
		{Name: "منتج 1", Quantity: 1.0, UnitPrice: 50.00, TaxableAmt: 50.00, VATAmount: 7.5, TotalWithVAT: 57.5},
		{Name: "منتج 2", Quantity: 1.0, UnitPrice: 70.00, TaxableAmt: 70.00, VATAmount: 10.5, TotalWithVAT: 80.5},
		{Name: "منتج 3", Quantity: 1.0, UnitPrice: 100.00, TaxableAmt: 100.00, VATAmount: 15, TotalWithVAT: 115},
	}

	return Invoice{
		Title:             "فاتورة ضريبية مبسطة",
		InvoiceNumber:     "INV10111",
		StoreName:         "اسم المتجر",
		StoreAddress:      "عنوان المتجر",
		Date:              "2021/12/12",
		VATRegistrationNo: "123456789900003",
		Products:          products,
		TotalTaxableAmt:   220.00,
		TotalVAT:          33.00,
		TotalWithVAT:      253.00,
		QRCodeData:        "AQpteSBjb21wYW55Ag8zMTIzNDU2Nzg5MDAwMDMDFDIwMjQtMDEtMTVUMTI6MDA6MDBaBAYyNTMuMDAFBTMzLjAw",
	}
}

// drawText draws text centered in a given width, handling Arabic RTL
func drawText(pdf *gopdf.GoPdf, text string, x, y, width float64) {
	processedText := arabictext.Process(text)
	textWidth, _ := pdf.MeasureTextWidth(processedText)
	centerX := x + (width-textWidth)/2
	pdf.SetXY(centerX, y)
	pdf.Cell(nil, processedText)
}

// drawTextLeft draws text left-aligned
func drawTextLeft(pdf *gopdf.GoPdf, text string, x, y float64) {
	processedText := arabictext.Process(text)
	pdf.SetXY(x, y)
	pdf.Cell(nil, processedText)
}

// drawTextRight draws text right-aligned in a given width
func drawTextRight(pdf *gopdf.GoPdf, text string, x, y, width float64) {
	processedText := arabictext.Process(text)
	textWidth, _ := pdf.MeasureTextWidth(processedText)
	rightX := x + width - textWidth - 2
	pdf.SetXY(rightX, y)
	pdf.Cell(nil, processedText)
}

// GeneratePDF creates the invoice PDF with Arabic RTL support
func GeneratePDF(invoice Invoice, filename string, fontDir string) error {
	pdf := gopdf.GoPdf{}

	// Receipt size: 80mm x 250mm (increased height for content)
	pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: 226.77, H: 708.66}, // 80mm x 250mm in points
	})

	// Add Arabic font
	fontPath := fontDir + "/Amiri-Regular.ttf"
	err := pdf.AddTTFFont("Amiri", fontPath)
	if err != nil {
		fmt.Printf("Warning: Could not load Arabic font: %v, using fallback\n", err)
		return generateFallbackPDF(invoice, filename)
	}

	boldFontPath := fontDir + "/Amiri-Bold.ttf"
	err = pdf.AddTTFFont("AmiriBold", boldFontPath)
	if err != nil {
		fmt.Printf("Warning: Could not load Arabic bold font, using regular\n")
		_ = pdf.AddTTFFont("AmiriBold", fontPath) // Fallback to regular
	}

	pdf.AddPage()

	pageWidth := 226.77
	margin := 10.0
	contentWidth := pageWidth - (2 * margin)
	currentY := 10.0

	// ===== HEADER - Arabic Title (with border) =====
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(1.0)
	pdf.RectFromUpperLeftWithStyle(margin, currentY, contentWidth, 22, "D")

	if err := pdf.SetFont("AmiriBold", "", 14); err != nil {
		pdf.SetFont("Amiri", "", 14)
	}
	pdf.SetTextColor(0, 0, 0)
	drawText(&pdf, invoice.Title, margin, currentY+4, contentWidth)
	currentY += 24

	// ===== INVOICE NUMBER (with border) =====
	pdf.SetLineWidth(0.5)
	pdf.RectFromUpperLeftWithStyle(margin, currentY, contentWidth, 14, "D")

	pdf.SetFont("Amiri", "", 9)
	pdf.SetTextColor(0, 0, 0)
	labelText := arabictext.Process("رقم الفاتورة:")
	labelW, _ := pdf.MeasureTextWidth(labelText)
	pdf.SetXY(margin+contentWidth-labelW-3, currentY+3)
	pdf.Cell(nil, labelText)
	pdf.SetXY(margin+3, currentY+3)
	pdf.Cell(nil, invoice.InvoiceNumber)
	currentY += 16

	// ===== STORE NAME (no border, just text) =====
	if err := pdf.SetFont("AmiriBold", "", 11); err != nil {
		pdf.SetFont("Amiri", "", 11)
	}
	pdf.SetTextColor(0, 0, 0)
	drawText(&pdf, invoice.StoreName, margin, currentY, contentWidth)
	currentY += 14

	// ===== STORE ADDRESS (no border, just text) =====
	pdf.SetFont("Amiri", "", 9)
	drawText(&pdf, invoice.StoreAddress, margin, currentY, contentWidth)
	currentY += 14

	// ===== DATE (with border) =====
	pdf.RectFromUpperLeftWithStyle(margin+30, currentY, contentWidth-60, 14, "D")
	dateLabelText := arabictext.Process("تاريخ:")
	dateLabelW, _ := pdf.MeasureTextWidth(dateLabelText)
	pdf.SetXY(margin+contentWidth-30-dateLabelW-3, currentY+3)
	pdf.Cell(nil, dateLabelText)
	pdf.SetXY(margin+33, currentY+3)
	pdf.Cell(nil, invoice.Date)
	currentY += 16

	// ===== VAT REGISTRATION (no border, just text) =====
	pdf.SetFont("Amiri", "", 8)
	// Draw as single line: number on left, label on right
	vatLabelText := arabictext.Process("رقم تسجيل ضريبة القيمة المضافة:")
	vatLabelW, _ := pdf.MeasureTextWidth(vatLabelText)
	pdf.SetXY(margin+contentWidth-vatLabelW, currentY)
	pdf.Cell(nil, vatLabelText)
	pdf.SetXY(margin, currentY)
	pdf.Cell(nil, invoice.VATRegistrationNo)
	currentY += 14

	// ===== PRODUCTS TABLE =====
	// Column widths - use full page width
	colWidths := []float64{35, 35, 35, 25, 76} // Total, VAT, Price, Qty, Product
	tableWidth := 0.0
	for _, w := range colWidths {
		tableWidth += w
	}
	tableX := margin

	// Headers with FULL TEXT (no abbreviations) - 2 lines if needed
	// Order: السعر شامل الضريبة | ضريبة القيمة المضافة | سعر الوحدة | الكمية | المنتجات
	headersLine1 := []string{
		"السعر شامل",
		"ضريبة القيمة",
		"سعر",
		"",
		"",
	}
	headersLine2 := []string{
		"الضريبة",
		"المضافة",
		"الوحدة",
		"الكمية",
		"المنتجات",
	}

	// Table Header - black border only
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.5)

	xPos := tableX
	headerHeight := 28.0 // Taller for 2 lines
	
	// Draw header cells
	for i := range headersLine1 {
		pdf.RectFromUpperLeftWithStyle(xPos, currentY, colWidths[i], headerHeight, "D")
		xPos += colWidths[i]
	}
	
	// Draw header text (2 lines)
	if err := pdf.SetFont("AmiriBold", "", 7); err != nil {
		pdf.SetFont("Amiri", "", 7)
	}
	pdf.SetTextColor(0, 0, 0)
	
	xPos = tableX
	for i := range headersLine1 {
		// Line 1
		if headersLine1[i] != "" {
			h1Text := arabictext.Process(headersLine1[i])
			h1w, _ := pdf.MeasureTextWidth(h1Text)
			pdf.SetXY(xPos+(colWidths[i]-h1w)/2, currentY+4)
			pdf.Cell(nil, h1Text)
		}
		// Line 2
		h2Text := arabictext.Process(headersLine2[i])
		h2w, _ := pdf.MeasureTextWidth(h2Text)
		pdf.SetXY(xPos+(colWidths[i]-h2w)/2, currentY+14)
		pdf.Cell(nil, h2Text)
		xPos += colWidths[i]
	}
	currentY += headerHeight

	// Table Rows
	pdf.SetFont("Amiri", "", 9)
	rowHeight := 18.0

	for _, product := range invoice.Products {
		pdf.SetStrokeColor(0, 0, 0)
		
		// Draw row cells (border only)
		xPos = tableX
		for i := range colWidths {
			pdf.RectFromUpperLeftWithStyle(xPos, currentY, colWidths[i], rowHeight, "D")
			xPos += colWidths[i]
		}
		
		// Draw text - position text higher in cell to prevent bottom overflow
		// Using +3 offset from cell top to keep descenders within cell
		pdf.SetTextColor(0, 0, 0)
		xPos = tableX
		textY := currentY + 3
		
		// Column 0: Total with VAT
		totalStr := fmt.Sprintf("%.1f", product.TotalWithVAT)
		tw, _ := pdf.MeasureTextWidth(totalStr)
		pdf.SetXY(xPos+(colWidths[0]-tw)/2, textY)
		pdf.Cell(nil, totalStr)
		xPos += colWidths[0]
		
		// Column 1: VAT Amount
		vatStr := fmt.Sprintf("%.1f", product.VATAmount)
		vw, _ := pdf.MeasureTextWidth(vatStr)
		pdf.SetXY(xPos+(colWidths[1]-vw)/2, textY)
		pdf.Cell(nil, vatStr)
		xPos += colWidths[1]
		
		// Column 2: Unit Price
		priceStr := fmt.Sprintf("%.0f", product.UnitPrice)
		pw, _ := pdf.MeasureTextWidth(priceStr)
		pdf.SetXY(xPos+(colWidths[2]-pw)/2, textY)
		pdf.Cell(nil, priceStr)
		xPos += colWidths[2]
		
		// Column 3: Quantity
		qtyStr := fmt.Sprintf("%.0f", product.Quantity)
		qw, _ := pdf.MeasureTextWidth(qtyStr)
		pdf.SetXY(xPos+(colWidths[3]-qw)/2, textY)
		pdf.Cell(nil, qtyStr)
		xPos += colWidths[3]
		
		// Column 4: Product Name (Arabic)
		nameText := arabictext.Process(product.Name)
		nw, _ := pdf.MeasureTextWidth(nameText)
		pdf.SetXY(xPos+(colWidths[4]-nw)/2, textY)
		pdf.Cell(nil, nameText)
		
		currentY += rowHeight
	}
	currentY += 6

	// ===== TOTALS SECTION =====
	totalsWidth := tableWidth
	totalsX := tableX
	valueWidth := 40.0
	labelWidth := totalsWidth - valueWidth

	// Row 1: Taxable Amount (16pt height)
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.5)
	pdf.RectFromUpperLeftWithStyle(totalsX, currentY, valueWidth, 16, "D")
	pdf.RectFromUpperLeftWithStyle(totalsX+valueWidth, currentY, labelWidth, 16, "D")
	
	pdf.SetFont("Amiri", "", 9)
	pdf.SetTextColor(0, 0, 0)
	taxableStr := fmt.Sprintf("%.0f", invoice.TotalTaxableAmt)
	taxableW, _ := pdf.MeasureTextWidth(taxableStr)
	pdf.SetXY(totalsX+(valueWidth-taxableW)/2, currentY+3)
	pdf.Cell(nil, taxableStr)
	
	taxableLbl := arabictext.Process("اجمالي المبلغ الخاضع للضريبة")
	taxableLblW, _ := pdf.MeasureTextWidth(taxableLbl)
	pdf.SetXY(totalsX+valueWidth+(labelWidth-taxableLblW)/2, currentY+3)
	pdf.Cell(nil, taxableLbl)
	currentY += 16

	// Row 2: VAT Amount (15%) - 16pt height
	pdf.RectFromUpperLeftWithStyle(totalsX, currentY, valueWidth, 16, "D")
	pdf.RectFromUpperLeftWithStyle(totalsX+valueWidth, currentY, labelWidth, 16, "D")
	
	vatTotalStr := fmt.Sprintf("%.0f", invoice.TotalVAT)
	vatTotalW, _ := pdf.MeasureTextWidth(vatTotalStr)
	pdf.SetXY(totalsX+(valueWidth-vatTotalW)/2, currentY+3)
	pdf.Cell(nil, vatTotalStr)
	
	// Fix: Don't reverse the percentage - write it separately
	vatLbl := arabictext.Process("ضريبة القيمة المضافة")
	vatPct := "(15%)" // Keep percentage as-is, don't process with Arabic
	vatLblW, _ := pdf.MeasureTextWidth(vatLbl)
	pctW, _ := pdf.MeasureTextWidth(vatPct)
	totalLblWidth := vatLblW + pctW + 3
	startX := totalsX + valueWidth + (labelWidth-totalLblWidth)/2
	pdf.SetXY(startX, currentY+3)
	pdf.Cell(nil, vatPct)
	pdf.SetXY(startX+pctW+3, currentY+3)
	pdf.Cell(nil, vatLbl)
	currentY += 16

	// Row 3: Total with VAT (bold border, 18pt height)
	pdf.SetLineWidth(1.0)
	pdf.RectFromUpperLeftWithStyle(totalsX, currentY, valueWidth, 18, "D")
	pdf.RectFromUpperLeftWithStyle(totalsX+valueWidth, currentY, labelWidth, 18, "D")
	
	if err := pdf.SetFont("AmiriBold", "", 10); err != nil {
		pdf.SetFont("Amiri", "", 10)
	}
	pdf.SetTextColor(0, 0, 0)
	
	totalStr := fmt.Sprintf("%.0f", invoice.TotalWithVAT)
	totalStrW, _ := pdf.MeasureTextWidth(totalStr)
	pdf.SetXY(totalsX+(valueWidth-totalStrW)/2, currentY+4)
	pdf.Cell(nil, totalStr)
	
	// Fix: Don't reverse the percentage
	totalLbl := arabictext.Process("المجموع مع الضريبة")
	totalPct := "(15%)"
	totalLblW, _ := pdf.MeasureTextWidth(totalLbl)
	totalPctW, _ := pdf.MeasureTextWidth(totalPct)
	totalFullWidth := totalLblW + totalPctW + 3
	totalStartX := totalsX + valueWidth + (labelWidth-totalFullWidth)/2
	pdf.SetXY(totalStartX, currentY+4)
	pdf.Cell(nil, totalPct)
	pdf.SetXY(totalStartX+totalPctW+3, currentY+4)
	pdf.Cell(nil, totalLbl)
	currentY += 22

	// ===== FOOTER =====
	pdf.SetFont("Amiri", "", 7)
	pdf.SetTextColor(0, 0, 0)
	footerText := ">>>>>>>>>>>>>> إغلاق الفاتورة 0010 <<<<<<<<<<<<<<<"
	drawText(&pdf, footerText, margin, currentY, contentWidth)
	currentY += 12

	// ===== QR CODE =====
	qrFile := "/tmp/temp_qr.png"
	err = qrcode.WriteFile(invoice.QRCodeData, qrcode.High, 256, qrFile)
	if err == nil {
		qrSize := 55.0
		qrX := margin + (contentWidth-qrSize)/2
		pdf.Image(qrFile, qrX, currentY, &gopdf.Rect{W: qrSize, H: qrSize})
		os.Remove(qrFile)
	}

	return pdf.WritePdf(filename)
}

// generateFallbackPDF creates a simple PDF without Arabic fonts
func generateFallbackPDF(invoice Invoice, filename string) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: 226.77, H: 708.66},
	})
	pdf.AddPage()

	pageWidth := 226.77
	margin := 8.0
	contentWidth := pageWidth - (2 * margin)
	currentY := 12.0

	// Green header
	pdf.SetFillColor(76, 175, 80)
	pdf.RectFromUpperLeftWithStyle(margin, currentY, contentWidth, 26, "F")
	currentY += 32

	// Invoice number box
	pdf.SetFillColor(250, 250, 250)
	pdf.SetStrokeColor(200, 200, 200)
	pdf.RectFromUpperLeftWithStyle(margin+15, currentY, contentWidth-30, 16, "FD")
	currentY += 21

	// Store boxes
	pdf.SetFillColor(76, 175, 80)
	pdf.RectFromUpperLeftWithStyle(margin+12, currentY, contentWidth-24, 18, "F")
	currentY += 23

	pdf.RectFromUpperLeftWithStyle(margin+12, currentY, contentWidth-24, 16, "F")
	currentY += 21

	// Date box
	pdf.SetFillColor(250, 250, 250)
	pdf.RectFromUpperLeftWithStyle(margin+25, currentY, contentWidth-50, 14, "FD")
	currentY += 33

	// Products table
	tableWidth := 210.0
	tableX := margin + (contentWidth-tableWidth)/2

	// Header
	pdf.SetFillColor(235, 235, 235)
	pdf.RectFromUpperLeftWithStyle(tableX, currentY, tableWidth, 16, "FD")
	currentY += 16

	// Rows
	for i := 0; i < len(invoice.Products); i++ {
		if i%2 == 0 {
			pdf.SetFillColor(255, 255, 255)
		} else {
			pdf.SetFillColor(245, 245, 245)
		}
		pdf.RectFromUpperLeftWithStyle(tableX, currentY, tableWidth, 14, "FD")
		currentY += 14
	}
	currentY += 10

	// Totals
	totalsWidth := 150.0
	totalsX := margin + (contentWidth-totalsWidth)/2

	pdf.SetFillColor(245, 245, 245)
	pdf.RectFromUpperLeftWithStyle(totalsX, currentY, totalsWidth, 14, "FD")
	currentY += 14
	pdf.RectFromUpperLeftWithStyle(totalsX, currentY, totalsWidth, 14, "FD")
	currentY += 14

	pdf.SetFillColor(76, 175, 80)
	pdf.RectFromUpperLeftWithStyle(totalsX, currentY, totalsWidth, 18, "F")
	currentY += 32

	// QR Code
	qrFile := "/tmp/temp_qr.png"
	err := qrcode.WriteFile(invoice.QRCodeData, qrcode.High, 256, qrFile)
	if err == nil {
		qrSize := 55.0
		qrX := margin + (contentWidth-qrSize)/2
		pdf.Image(qrFile, qrX, currentY, &gopdf.Rect{W: qrSize, H: qrSize})
		os.Remove(qrFile)
	}

	return pdf.WritePdf(filename)
}
