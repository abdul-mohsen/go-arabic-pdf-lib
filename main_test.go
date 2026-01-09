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

	defer os.Remove(filename)

	err := GeneratePDF(invoice, filename)
	require.NoError(t, err, "PDF generation should not return an error")

	info, err := os.Stat(filename)
	require.NoError(t, err, "PDF file should be created")
	assert.Greater(t, info.Size(), int64(1000), "PDF should have substantial content")

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

	// Verify title
	assert.Equal(t, "Simplified Tax Invoice", invoice.Title)

	// Verify invoice metadata
	assert.Equal(t, "INV10111", invoice.InvoiceNumber)
	assert.Equal(t, "Store Name", invoice.StoreName)
	assert.Equal(t, "Store Address", invoice.StoreAddress)
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

func TestPDFQuality(t *testing.T) {
	invoice := generateSampleInvoice()
	filename := "test_quality.pdf"

	defer os.Remove(filename)

	err := GeneratePDF(invoice, filename)
	require.NoError(t, err)

	info, err := os.Stat(filename)
	require.NoError(t, err)

	// Quality checks - receipt style PDF is smaller
	assert.Greater(t, info.Size(), int64(2000), "PDF should be at least 2KB")
	assert.Less(t, info.Size(), int64(1000000), "PDF should be less than 1MB")
}

func TestInvoiceDate(t *testing.T) {
	invoice := generateSampleInvoice()
	assert.Equal(t, "2021/12/12", invoice.Date)
}

func TestQRCodeData(t *testing.T) {
	invoice := generateSampleInvoice()
	assert.NotEmpty(t, invoice.QRCodeData, "QR code data should not be empty")
}

func TestProductDetails(t *testing.T) {
	invoice := generateSampleInvoice()

	// Test first product
	assert.Equal(t, "Product 1", invoice.Products[0].Name)
	assert.Equal(t, 1.0, invoice.Products[0].Quantity)
	assert.Equal(t, 50.00, invoice.Products[0].UnitPrice)
	assert.Equal(t, 7.5, invoice.Products[0].VATAmount)
	assert.Equal(t, 57.5, invoice.Products[0].TotalWithVAT)

	// Test second product
	assert.Equal(t, "Product 2", invoice.Products[1].Name)
	assert.Equal(t, 70.00, invoice.Products[1].UnitPrice)
	assert.Equal(t, 10.5, invoice.Products[1].VATAmount)
	assert.Equal(t, 80.5, invoice.Products[1].TotalWithVAT)

	// Test third product
	assert.Equal(t, "Product 3", invoice.Products[2].Name)
	assert.Equal(t, 100.00, invoice.Products[2].UnitPrice)
	assert.Equal(t, 15.0, invoice.Products[2].VATAmount)
	assert.Equal(t, 115.0, invoice.Products[2].TotalWithVAT)
}
