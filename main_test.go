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

func TestReverseString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "olleh"},
		{"world", "dlrow"},
		{"123", "321"},
		{"", ""},
		{"a", "a"},
	}

	for _, tt := range tests {
		result := reverseString(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

func TestIsArabic(t *testing.T) {
	// Test Arabic characters
	assert.True(t, isArabic('ا'))
	assert.True(t, isArabic('ب'))
	assert.True(t, isArabic('ت'))
	assert.True(t, isArabic('م'))
	
	// Test non-Arabic characters
	assert.False(t, isArabic('a'))
	assert.False(t, isArabic('1'))
	assert.False(t, isArabic(' '))
	assert.False(t, isArabic('!'))
}

func TestHasArabic(t *testing.T) {
	assert.True(t, hasArabic("مرحبا"))
	assert.True(t, hasArabic("Hello مرحبا World"))
	assert.False(t, hasArabic("Hello World"))
	assert.False(t, hasArabic("12345"))
	assert.False(t, hasArabic(""))
}

func TestReshapeArabic(t *testing.T) {
	// Test that reshaping produces output
	reshaped := reshapeArabic("مرحبا")
	assert.NotEmpty(t, reshaped)
	
	// Test that non-Arabic text is unchanged
	assert.Equal(t, "hello", reshapeArabic("hello"))
	assert.Equal(t, "123", reshapeArabic("123"))
}

func TestPrepareText(t *testing.T) {
	// Non-Arabic text should remain unchanged
	assert.Equal(t, "hello", prepareText("hello"))
	assert.Equal(t, "123", prepareText("123"))
	
	// Arabic text should be transformed
	arabicResult := prepareText("مرحبا")
	assert.NotEqual(t, "مرحبا", arabicResult) // Should be reshaped and reversed
}

func TestPrepareArabicText(t *testing.T) {
	result := prepareArabicText("اسم")
	assert.NotEmpty(t, result)
	// The result should be reshaped and reversed
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

func TestArabicFormsMap(t *testing.T) {
	// Verify common Arabic letters have forms defined
	commonLetters := []rune{'ا', 'ب', 'ت', 'م', 'ن', 'ل'}
	for _, letter := range commonLetters {
		forms, exists := arabicForms[letter]
		assert.True(t, exists, "Forms should exist for letter: %c", letter)
		assert.Equal(t, 4, len(forms), "Each letter should have 4 forms")
	}
}

func TestNonConnectingLetters(t *testing.T) {
	// Verify non-connecting letters are defined
	assert.True(t, nonConnectingLetters['ا'])
	assert.True(t, nonConnectingLetters['د'])
	assert.True(t, nonConnectingLetters['ر'])
	assert.True(t, nonConnectingLetters['و'])
	
	// Verify connecting letters are not in the map
	assert.False(t, nonConnectingLetters['ب'])
	assert.False(t, nonConnectingLetters['ت'])
	assert.False(t, nonConnectingLetters['م'])
}
