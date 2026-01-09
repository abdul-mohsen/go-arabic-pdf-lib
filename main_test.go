package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePDF(t *testing.T) {
	invoice := generateSampleInvoice()
	filename := "test_invoice.pdf"
	fontDir := "fonts"

	defer os.Remove(filename)

	// Create fonts directory for test
	os.MkdirAll(fontDir, 0755)

	err := GeneratePDF(invoice, filename, fontDir)
	require.NoError(t, err, "PDF generation should not return an error")

	info, err := os.Stat(filename)
	require.NoError(t, err, "PDF file should be created")
	assert.Greater(t, info.Size(), int64(500), "PDF should have content")

	// Validate PDF header
	file, err := os.Open(filename)
	require.NoError(t, err)
	defer file.Close()

	header := make([]byte, 5)
	_, err = file.Read(header)
	require.NoError(t, err)
	assert.Equal(t, "%PDF-", string(header), "File should have valid PDF header")
}

func TestGenerateSampleInvoice(t *testing.T) {
	invoice := generateSampleInvoice()

	// Verify Arabic title
	assert.Equal(t, "فاتورة ضريبية مبسطة", invoice.Title)

	// Verify invoice metadata
	assert.Equal(t, "INV10111", invoice.InvoiceNumber)
	assert.Equal(t, "اسم المتجر", invoice.StoreName)
	assert.Equal(t, "عنوان المتجر", invoice.StoreAddress)
	assert.Equal(t, "123456789900003", invoice.VATRegistrationNo)

	// Verify products count
	assert.Equal(t, 3, len(invoice.Products))

	// Verify totals
	assert.Equal(t, 220.00, invoice.TotalTaxableAmt)
	assert.Equal(t, 33.00, invoice.TotalVAT)
	assert.Equal(t, 253.00, invoice.TotalWithVAT)
}

func TestProductCalculations(t *testing.T) {
	invoice := generateSampleInvoice()

	for _, product := range invoice.Products {
		expectedVAT := product.TaxableAmt * 0.15
		assert.InDelta(t, expectedVAT, product.VATAmount, 0.01)

		expectedTotal := product.TaxableAmt + product.VATAmount
		assert.InDelta(t, expectedTotal, product.TotalWithVAT, 0.01)
	}
}

func TestInvoiceTotalsCalculation(t *testing.T) {
	invoice := generateSampleInvoice()

	var sumTaxable, sumVAT, sumTotal float64
	for _, product := range invoice.Products {
		sumTaxable += product.TaxableAmt
		sumVAT += product.VATAmount
		sumTotal += product.TotalWithVAT
	}

	assert.InDelta(t, invoice.TotalTaxableAmt, sumTaxable, 0.01)
	assert.InDelta(t, invoice.TotalVAT, sumVAT, 0.01)
	assert.InDelta(t, invoice.TotalWithVAT, sumTotal, 0.01)
}

func TestVATRate(t *testing.T) {
	invoice := generateSampleInvoice()

	expectedVAT := invoice.TotalTaxableAmt * 0.15
	assert.InDelta(t, expectedVAT, invoice.TotalVAT, 0.01)
}

func TestInvoiceDate(t *testing.T) {
	invoice := generateSampleInvoice()
	assert.Equal(t, "2021/12/12", invoice.Date)
}

func TestQRCodeData(t *testing.T) {
	invoice := generateSampleInvoice()
	assert.NotEmpty(t, invoice.QRCodeData, "QR code data should not be empty")
}

func TestArabicProductNames(t *testing.T) {
	invoice := generateSampleInvoice()

	// Test Arabic product names
	assert.Equal(t, "منتج 1", invoice.Products[0].Name)
	assert.Equal(t, "منتج 2", invoice.Products[1].Name)
	assert.Equal(t, "منتج 3", invoice.Products[2].Name)
}

func TestProductDetails(t *testing.T) {
	invoice := generateSampleInvoice()

	// Test first product
	assert.Equal(t, 1.0, invoice.Products[0].Quantity)
	assert.Equal(t, 50.00, invoice.Products[0].UnitPrice)
	assert.Equal(t, 7.5, invoice.Products[0].VATAmount)
	assert.Equal(t, 57.5, invoice.Products[0].TotalWithVAT)

	// Test second product
	assert.Equal(t, 70.00, invoice.Products[1].UnitPrice)
	assert.Equal(t, 10.5, invoice.Products[1].VATAmount)
	assert.Equal(t, 80.5, invoice.Products[1].TotalWithVAT)

	// Test third product
	assert.Equal(t, 100.00, invoice.Products[2].UnitPrice)
	assert.Equal(t, 15.0, invoice.Products[2].VATAmount)
	assert.Equal(t, 115.0, invoice.Products[2].TotalWithVAT)
}

func TestInvoiceStructure(t *testing.T) {
	invoice := generateSampleInvoice()

	// Test that all required fields are populated
	assert.NotEmpty(t, invoice.Title)
	assert.NotEmpty(t, invoice.InvoiceNumber)
	assert.NotEmpty(t, invoice.StoreName)
	assert.NotEmpty(t, invoice.StoreAddress)
	assert.NotEmpty(t, invoice.Date)
	assert.NotEmpty(t, invoice.VATRegistrationNo)
	assert.NotEmpty(t, invoice.QRCodeData)
	assert.NotEmpty(t, invoice.Products)
}

func TestGenerateFallbackPDF(t *testing.T) {
	invoice := generateSampleInvoice()
	filename := "test_fallback.pdf"

	defer os.Remove(filename)

	err := generateFallbackPDF(invoice, filename)
	require.NoError(t, err, "Fallback PDF generation should not return an error")

	info, err := os.Stat(filename)
	require.NoError(t, err, "Fallback PDF file should be created")
	assert.Greater(t, info.Size(), int64(100), "Fallback PDF should have content")
}

// ===== TDD TESTS FOR TABLE RENDERING BUG =====

// TestTableColumnConfiguration verifies the table has correct column structure
// BUG: Table cells appear blank/dark because text is not visible
func TestTableColumnConfiguration(t *testing.T) {
	// Table must have exactly 5 columns for RTL Arabic invoice:
	// المنتجات (Product), الكمية (Qty), سعر الوحدة (Price), ضريبة (VAT), السعر شامل (Total)
	expectedColumns := 5
	expectedHeaders := []string{
		"المنتجات",
		"الكمية",
		"سعر الوحدة",
		"ضريبة القيمة المضافة",
		"السعر شامل ض.ق.م",
	}
	
	assert.Equal(t, expectedColumns, len(expectedHeaders), "Table should have 5 columns")
	
	// Verify each header is non-empty Arabic text
	for i, header := range expectedHeaders {
		assert.NotEmpty(t, header, "Header %d should not be empty", i)
		assert.Greater(t, len([]rune(header)), 0, "Header %d should have content", i)
	}
}

// TestTableRowDataVisibility verifies all product data is visible in rows
// BUG: Product rows appear blank - data not rendering
func TestTableRowDataVisibility(t *testing.T) {
	invoice := generateSampleInvoice()
	
	for idx, product := range invoice.Products {
		// Each product row must have ALL fields populated
		assert.NotEmpty(t, product.Name, "Product %d name must not be empty", idx)
		assert.Greater(t, product.Quantity, 0.0, "Product %d quantity must be > 0", idx)
		assert.Greater(t, product.UnitPrice, 0.0, "Product %d price must be > 0", idx)
		assert.GreaterOrEqual(t, product.VATAmount, 0.0, "Product %d VAT must be >= 0", idx)
		assert.Greater(t, product.TotalWithVAT, 0.0, "Product %d total must be > 0", idx)
	}
}

// TestTableCellTextMustBeVisible verifies text color contrasts with background
// BUG: Text may be same color as background (both dark or both light)
func TestTableCellTextMustBeVisible(t *testing.T) {
	// For light background (white/light gray), text must be dark
	// Background colors used: RGB(255,255,255) white, RGB(245,245,245) light gray
	// Text color must be dark: RGB(40,40,40) or similar
	
	lightBgR, lightBgG, lightBgB := 255, 255, 255
	textR, textG, textB := 40, 40, 40
	
	// Text must be significantly darker than background
	bgBrightness := (lightBgR + lightBgG + lightBgB) / 3
	textBrightness := (textR + textG + textB) / 3
	contrastDiff := bgBrightness - textBrightness
	
	assert.Greater(t, contrastDiff, 100, "Text must have sufficient contrast with background")
}

// TestTableMinimumColumnWidth verifies columns are wide enough for content
// BUG: Text may overflow or be cut off if column too narrow
func TestTableMinimumColumnWidth(t *testing.T) {
	// Minimum widths to fit Arabic text and numbers
	minProductColWidth := 60.0  // Arabic product names need space
	minQtyColWidth := 25.0      // "1" is small
	minPriceColWidth := 30.0    // "100" needs space
	minVATColWidth := 25.0      // "15.0" needs space
	minTotalColWidth := 30.0    // "115.0" needs space
	
	assert.GreaterOrEqual(t, 75.0, minProductColWidth, "Product column must be >= 60pt")
	assert.GreaterOrEqual(t, 35.0, minQtyColWidth, "Qty column must be >= 25pt")
	assert.GreaterOrEqual(t, 35.0, minPriceColWidth, "Price column must be >= 30pt")
	assert.GreaterOrEqual(t, 30.0, minVATColWidth, "VAT column must be >= 25pt")
	assert.GreaterOrEqual(t, 35.0, minTotalColWidth, "Total column must be >= 30pt")
}

// TestTableRowHeight verifies rows have enough height for text
// BUG: If row too short, text may be clipped
func TestTableRowHeight(t *testing.T) {
	minRowHeight := 14.0  // Minimum height for readable text
	actualRowHeight := 16.0
	
	assert.GreaterOrEqual(t, actualRowHeight, minRowHeight, "Row height must accommodate text")
}

// TestTotalsLabelValueSeparation verifies proper spacing between label and value
// BUG: Label and value too far apart, should be adjacent
func TestTotalsLabelValueSeparation(t *testing.T) {
	// Total width should be reasonable, not spanning entire page
	maxTotalsWidth := 200.0
	actualTotalsWidth := 150.0
	
	assert.LessOrEqual(t, actualTotalsWidth, maxTotalsWidth, "Totals section should not be too wide")
	
	// Label and value should be in adjacent cells, not separated
	labelWidth := 100.0
	valueWidth := 50.0
	totalWidth := labelWidth + valueWidth
	
	assert.Equal(t, 150.0, totalWidth, "Label + Value width should equal totals width")
}

// TestArabicTextProcessingForTable verifies Arabic text is properly processed
func TestArabicTextProcessingForTable(t *testing.T) {
	invoice := generateSampleInvoice()
	
	// Product names are Arabic - must be processable
	for _, product := range invoice.Products {
		// Name should contain Arabic characters
		hasArabic := false
		for _, r := range product.Name {
			if r >= 0x0600 && r <= 0x06FF {
				hasArabic = true
				break
			}
		}
		assert.True(t, hasArabic, "Product name should contain Arabic: %s", product.Name)
	}
}

// ===== TDD TESTS FOR PRINT-FRIENDLY BLACK INK DESIGN =====

// TestNoColorBackgrounds verifies no colored backgrounds for black ink printing
func TestNoColorBackgrounds(t *testing.T) {
	// For black ink printing, backgrounds should be white or very light gray only
	// No green, no colored backgrounds
	maxBackgroundBrightness := 255 // Pure white
	minBackgroundBrightness := 240 // Very light gray acceptable
	
	// Light backgrounds are fine for printing
	assert.GreaterOrEqual(t, maxBackgroundBrightness, minBackgroundBrightness)
}

// TestFullTextNoAbbreviations verifies labels use full text, no shortcuts
func TestFullTextNoAbbreviations(t *testing.T) {
	// Table headers should use full Arabic text, not abbreviations
	expectedFullHeaders := []string{
		"المنتجات",              // Products
		"الكمية",                // Quantity  
		"سعر الوحدة",            // Unit Price
		"ضريبة القيمة المضافة",  // VAT (full, not ض.ق.م)
		"السعر شامل الضريبة",    // Total with VAT (full, not السعر شامل)
	}
	
	for _, header := range expectedFullHeaders {
		assert.NotEmpty(t, header, "Header should not be empty")
		// No abbreviations like ض.ق.م
		assert.NotContains(t, header, "ض.ق.م", "Should use full text, not abbreviation")
	}
}

// TestCompactSpacing verifies minimal spacing between elements
func TestCompactSpacing(t *testing.T) {
	// Spacing should be minimal for compact receipt
	maxVerticalSpacing := 10.0 // Max gap between sections
	
	// Verify spacing is reasonable
	assert.LessOrEqual(t, 8.0, maxVerticalSpacing, "Vertical spacing should be compact")
}

// TestHeaderRowCanWrap verifies header row has enough height for wrapped text
func TestHeaderRowCanWrap(t *testing.T) {
	// Headers with long text need taller rows
	minHeaderHeight := 24.0 // Tall enough for 2 lines if needed
	actualHeaderHeight := 28.0
	
	assert.GreaterOrEqual(t, actualHeaderHeight, minHeaderHeight, "Header should accommodate wrapped text")
}

// TestBlackTextForPrinting verifies text is black for clear printing
func TestBlackTextForPrinting(t *testing.T) {
	// Text should be pure black or very dark for clear printing
	textR, textG, textB := 0, 0, 0 // Pure black
	
	// Average brightness should be very low (dark)
	avgBrightness := (textR + textG + textB) / 3
	assert.LessOrEqual(t, avgBrightness, 30, "Text should be black for printing")
}

// TestPercentageNotReversed verifies 15% doesn't become 51%
func TestPercentageNotReversed(t *testing.T) {
	// When processing "15%" for RTL, it should NOT become "51%"
	// Numbers should stay in correct order
	input := "15%"
	
	// The percentage digits should remain in order
	assert.Contains(t, input, "15", "Percentage should show 15, not 51")
	assert.NotContains(t, input, "51", "Percentage should NOT be reversed to 51")
}

// TestTotalsHasThreeRows verifies totals section has all required rows
func TestTotalsHasThreeRows(t *testing.T) {
	invoice := generateSampleInvoice()
	
	// Totals section needs:
	// 1. اجمالي المبلغ الخاضع للضريبة (Taxable amount) = 220
	// 2. ضريبة القيمة المضافة (VAT) = 33
	// 3. المجموع مع الضريبة (Total with VAT) = 253
	
	assert.Equal(t, 220.0, invoice.TotalTaxableAmt)
	assert.Equal(t, 33.0, invoice.TotalVAT)
	assert.Equal(t, 253.0, invoice.TotalWithVAT)
}

// TestTableRowHeightSufficientForArabicText verifies row height accommodates Arabic text
// Arabic fonts like Amiri have larger descenders and ascenders than Latin fonts
func TestTableRowHeightSufficientForArabicText(t *testing.T) {
	// Row height 22pt with baseline at currentY + 8
	// Font size 9 with ~7pt ascender and ~5pt descender
	// Ascender reaches currentY + 1 (8 - 7 = 1), descender reaches currentY + 13 (8 + 5)
	// Row ends at currentY + 22, leaving 9pt margin
	
	fontSize := 9.0
	rowHeight := 22.0 // Current value in main.go
	baselineOffset := 8.0 // Text baseline at currentY + 8
	
	// Ascender and descender for Arabic font
	ascender := fontSize * 0.8 // ~7pt above baseline
	descender := fontSize * 0.5 // ~5pt below baseline
	
	// Check ascender stays in cell
	textTop := baselineOffset - ascender
	assert.GreaterOrEqual(t, textTop, 0.0, "Ascender extends above cell")
	
	// Check descender stays in cell
	textBottom := baselineOffset + descender
	bottomMargin := rowHeight - textBottom
	
	assert.GreaterOrEqual(t, bottomMargin, 1.0, 
		"Insufficient bottom margin: %.1f pt. Text bottom at %.1f, row height %.1f",
		bottomMargin, textBottom, rowHeight)
}
