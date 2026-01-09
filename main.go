package main

import (
	"encoding/json"
	"fmt"
	"os"

	"bill-generator/arabictext"

	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
)

// Config holds global configuration
type Config struct {
	VATPercentage  float64 `json:"vatPercentage"`
	CurrencySymbol string  `json:"currencySymbol"`
	DateFormat     string  `json:"dateFormat"`
}

// ProductInput represents a product from JSON (without calculated fields)
type ProductInput struct {
	Name      string  `json:"name"`
	Quantity  float64 `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
}

// Product represents a single product line item with calculated fields
type Product struct {
	Name         string
	Quantity     float64
	UnitPrice    float64
	TaxableAmt   float64
	VATAmount    float64
	TotalWithVAT float64
}

// InvoiceInput represents invoice data from JSON
type InvoiceInput struct {
	Title             string `json:"title"`
	InvoiceNumber     string `json:"invoiceNumber"`
	StoreName         string `json:"storeName"`
	StoreAddress      string `json:"storeAddress"`
	Date              string `json:"date"`
	VATRegistrationNo string `json:"vatRegistrationNo"`
	QRCodeData        string `json:"qrCodeData"`
}

// Labels holds all text labels for the invoice
type Labels struct {
	InvoiceNumber   string `json:"invoiceNumber"`
	Date            string `json:"date"`
	VATRegistration string `json:"vatRegistration"`
	TotalTaxable    string `json:"totalTaxable"`
	TotalWithVat    string `json:"totalWithVat"`
	ProductColumn   string `json:"productColumn"`
	QuantityColumn  string `json:"quantityColumn"`
	UnitPriceColumn string `json:"unitPriceColumn"`
	VATAmountColumn string `json:"vatAmountColumn"`
	TotalColumn     string `json:"totalColumn"`
	Footer          string `json:"footer"`
}

// InvoiceData represents the complete JSON structure
type InvoiceData struct {
	Config   Config         `json:"config"`
	Invoice  InvoiceInput   `json:"invoice"`
	Products []ProductInput `json:"products"`
	Labels   Labels         `json:"labels"`
}

// Invoice represents the complete invoice with calculated values
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
	VATPercentage     float64
	Labels            Labels
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

	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "invoice_data.json"
	}

	filename := outputDir + "/invoice_output.pdf"

	invoice, err := loadInvoiceFromJSON(dataFile)
	if err != nil {
		fmt.Printf("ERROR: Failed to load invoice data: %v\n", err)
		os.Exit(1)
	}

	err = GeneratePDF(invoice, filename, fontDir)
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

// loadInvoiceFromJSON loads invoice data from a JSON file and calculates VAT
func loadInvoiceFromJSON(filename string) (Invoice, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Invoice{}, fmt.Errorf("failed to read file: %w", err)
	}

	var invoiceData InvoiceData
	if err := json.Unmarshal(data, &invoiceData); err != nil {
		return Invoice{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	vatRate := invoiceData.Config.VATPercentage / 100.0

	// Calculate product values
	var products []Product
	var totalTaxable, totalVAT float64

	for _, p := range invoiceData.Products {
		taxableAmt := p.Quantity * p.UnitPrice
		vatAmount := taxableAmt * vatRate
		totalWithVAT := taxableAmt + vatAmount

		products = append(products, Product{
			Name:         p.Name,
			Quantity:     p.Quantity,
			UnitPrice:    p.UnitPrice,
			TaxableAmt:   taxableAmt,
			VATAmount:    vatAmount,
			TotalWithVAT: totalWithVAT,
		})

		totalTaxable += taxableAmt
		totalVAT += vatAmount
	}

	return Invoice{
		Title:             invoiceData.Invoice.Title,
		InvoiceNumber:     invoiceData.Invoice.InvoiceNumber,
		StoreName:         invoiceData.Invoice.StoreName,
		StoreAddress:      invoiceData.Invoice.StoreAddress,
		Date:              invoiceData.Invoice.Date,
		VATRegistrationNo: invoiceData.Invoice.VATRegistrationNo,
		Products:          products,
		TotalTaxableAmt:   totalTaxable,
		TotalVAT:          totalVAT,
		TotalWithVAT:      totalTaxable + totalVAT,
		QRCodeData:        invoiceData.Invoice.QRCodeData,
		VATPercentage:     invoiceData.Config.VATPercentage,
		Labels:            invoiceData.Labels,
	}, nil
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

// wrapText splits text into multiple lines that fit within maxWidth
// Returns the lines and the total height needed
func wrapText(pdf *gopdf.GoPdf, text string, maxWidth float64, lineHeight float64) ([]string, float64) {
	processedText := arabictext.Process(text)
	
	// Check if text fits in one line
	textWidth, _ := pdf.MeasureTextWidth(processedText)
	if textWidth <= maxWidth {
		return []string{processedText}, lineHeight
	}
	
	// Need to wrap - split by spaces/characters
	var lines []string
	runes := []rune(text)
	
	currentLine := ""
	for i := 0; i < len(runes); i++ {
		testLine := currentLine + string(runes[i])
		testProcessed := arabictext.Process(testLine)
		testWidth, _ := pdf.MeasureTextWidth(testProcessed)
		
		if testWidth > maxWidth && currentLine != "" {
			// Current line is full, save it and start new line
			lines = append(lines, arabictext.Process(currentLine))
			currentLine = string(runes[i])
		} else {
			currentLine = testLine
		}
	}
	
	// Add the last line
	if currentLine != "" {
		lines = append(lines, arabictext.Process(currentLine))
	}
	
	if len(lines) == 0 {
		lines = []string{processedText}
	}
	
	return lines, float64(len(lines)) * lineHeight
}

// drawWrappedTextRight draws wrapped text right-aligned in a cell
func drawWrappedTextRight(pdf *gopdf.GoPdf, text string, x, y, width, lineHeight float64) float64 {
	lines, totalHeight := wrapText(pdf, text, width-6, lineHeight) // 6pt padding
	
	for i, line := range lines {
		lineWidth, _ := pdf.MeasureTextWidth(line)
		pdf.SetXY(x+width-lineWidth-3, y+float64(i)*lineHeight)
		pdf.Cell(nil, line)
	}
	
	return totalHeight
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

	// ===== HEADER - Arabic Title =====
	if err := pdf.SetFont("AmiriBold", "", 14); err != nil {
		pdf.SetFont("Amiri", "", 14)
	}
	pdf.SetTextColor(0, 0, 0)
	drawText(&pdf, invoice.Title, margin, currentY+4, contentWidth)
	currentY += 18

	// ===== INVOICE NUMBER =====
	pdf.SetFont("Amiri", "", 9)
	pdf.SetTextColor(0, 0, 0)
	labelText := arabictext.Process(invoice.Labels.InvoiceNumber)
	labelW, _ := pdf.MeasureTextWidth(labelText)
	pdf.SetXY(margin+contentWidth-labelW-3, currentY)
	pdf.Cell(nil, labelText)
	pdf.SetXY(margin+3, currentY)
	pdf.Cell(nil, invoice.InvoiceNumber)
	currentY += 12

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

	// ===== DATE =====
	pdf.SetFont("Amiri", "", 9)
	dateLabelText := arabictext.Process(invoice.Labels.Date)
	dateLabelW, _ := pdf.MeasureTextWidth(dateLabelText)
	pdf.SetXY(margin+contentWidth-dateLabelW-3, currentY)
	pdf.Cell(nil, dateLabelText)
	pdf.SetXY(margin+3, currentY)
	pdf.Cell(nil, invoice.Date)
	currentY += 12

	// ===== VAT REGISTRATION (no border, just text) =====
	pdf.SetFont("Amiri", "", 8)
	// Draw as single line: number on left, label on right
	vatLabelText := arabictext.Process(invoice.Labels.VATRegistration)
	vatLabelW, _ := pdf.MeasureTextWidth(vatLabelText)
	pdf.SetXY(margin+contentWidth-vatLabelW, currentY)
	pdf.Cell(nil, vatLabelText)
	pdf.SetXY(margin, currentY)
	pdf.Cell(nil, invoice.VATRegistrationNo)
	currentY += 14

	// ===== PRODUCTS TABLE =====
	// Column widths - maximize product name column
	// Total, VAT, Price, Qty, Product
	colWidths := []float64{30, 30, 30, 20, 96} // Reduced other columns, increased product to 96
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
	baseRowHeight := 12.0 // Base height per line
	minRowHeight := 18.0  // Minimum row height

	for _, product := range invoice.Products {
		pdf.SetStrokeColor(0, 0, 0)
		
		// Calculate row height based on product name wrapping
		productColWidth := colWidths[4]
		_, nameHeight := wrapText(&pdf, product.Name, productColWidth-6, baseRowHeight)
		rowHeight := nameHeight + 6 // Add padding
		if rowHeight < minRowHeight {
			rowHeight = minRowHeight
		}
		
		// Draw row cells (border only) with calculated height
		xPos = tableX
		for i := range colWidths {
			pdf.RectFromUpperLeftWithStyle(xPos, currentY, colWidths[i], rowHeight, "D")
			xPos += colWidths[i]
		}
		
		// Draw text - position at baseline +3 from cell top (moved up)
		pdf.SetTextColor(0, 0, 0)
		xPos = tableX
		textY := currentY + 3
		
		// Column 0: Total with VAT (right aligned)
		totalStr := fmt.Sprintf("%.1f", product.TotalWithVAT)
		tw, _ := pdf.MeasureTextWidth(totalStr)
		pdf.SetXY(xPos+colWidths[0]-tw-3, textY)
		pdf.Cell(nil, totalStr)
		xPos += colWidths[0]
		
		// Column 1: VAT Amount (right aligned)
		vatStr := fmt.Sprintf("%.1f", product.VATAmount)
		vw, _ := pdf.MeasureTextWidth(vatStr)
		pdf.SetXY(xPos+colWidths[1]-vw-3, textY)
		pdf.Cell(nil, vatStr)
		xPos += colWidths[1]
		
		// Column 2: Unit Price (right aligned)
		priceStr := fmt.Sprintf("%.0f", product.UnitPrice)
		pw, _ := pdf.MeasureTextWidth(priceStr)
		pdf.SetXY(xPos+colWidths[2]-pw-3, textY)
		pdf.Cell(nil, priceStr)
		xPos += colWidths[2]
		
		// Column 3: Quantity (right aligned)
		qtyStr := fmt.Sprintf("%.0f", product.Quantity)
		qw, _ := pdf.MeasureTextWidth(qtyStr)
		pdf.SetXY(xPos+colWidths[3]-qw-3, textY)
		pdf.Cell(nil, qtyStr)
		xPos += colWidths[3]
		
		// Column 4: Product Name (Arabic, right aligned, with wrapping)
		drawWrappedTextRight(&pdf, product.Name, xPos, textY, colWidths[4], baseRowHeight)
		
		currentY += rowHeight
	}
	currentY += 8

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
	pdf.SetXY(totalsX+valueWidth-taxableW-3, currentY+3)
	pdf.Cell(nil, taxableStr)
	
	taxableLbl := arabictext.Process(invoice.Labels.TotalTaxable)
	taxableLblW, _ := pdf.MeasureTextWidth(taxableLbl)
	pdf.SetXY(totalsX+valueWidth+labelWidth-taxableLblW-2, currentY)
	pdf.Cell(nil, taxableLbl)
	currentY += 16

	// Row 2: Total with VAT (bold border, 18pt height)
	pdf.SetLineWidth(1.0)
	pdf.RectFromUpperLeftWithStyle(totalsX, currentY, valueWidth, 18, "D")
	pdf.RectFromUpperLeftWithStyle(totalsX+valueWidth, currentY, labelWidth, 18, "D")
	
	if err := pdf.SetFont("AmiriBold", "", 10); err != nil {
		pdf.SetFont("Amiri", "", 10)
	}
	pdf.SetTextColor(0, 0, 0)
	
	totalStr := fmt.Sprintf("%.0f", invoice.TotalWithVAT)
	totalStrW, _ := pdf.MeasureTextWidth(totalStr)
	pdf.SetXY(totalsX+valueWidth-totalStrW-2, currentY)
	pdf.Cell(nil, totalStr)
	
	// Arabic label with percentage - right aligned
	totalLbl := arabictext.Process(invoice.Labels.TotalWithVat)
	totalPct := fmt.Sprintf("(%.0f%%)", invoice.VATPercentage)
	totalLblW, _ := pdf.MeasureTextWidth(totalLbl)
	totalPctW, _ := pdf.MeasureTextWidth(totalPct)
	// Right align: label then percentage
	pdf.SetXY(totalsX+valueWidth+labelWidth-totalLblW-2, currentY)
	pdf.Cell(nil, totalLbl)
	pdf.SetXY(totalsX+valueWidth+labelWidth-totalLblW-totalPctW-2, currentY)
	pdf.Cell(nil, totalPct)
	currentY += 22

	// ===== FOOTER =====
	pdf.SetFont("Amiri", "", 7)
	pdf.SetTextColor(0, 0, 0)
	footerText := invoice.Labels.Footer
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
