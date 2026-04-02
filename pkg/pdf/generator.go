// Package pdf provides PDF generation for invoices.
package pdf

import (
	"fmt"

	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/models"
	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/pdf/sections"

	"github.com/signintech/gopdf"
)

// Generator handles PDF generation for invoices.
type Generator struct {
	fontDir string
}

// NewGenerator creates a new PDF generator.
func NewGenerator(fontDir string) *Generator {
	return &Generator{fontDir: fontDir}
}

// Generate creates a PDF from the invoice and saves it to filename.
func (g *Generator) Generate(invoice models.Invoice, filename string) error {
	pdfDoc, err := g.buildPDF(invoice)
	if err != nil {
		return err
	}
	return pdfDoc.WritePdf(filename)
}

// GenerateBytes creates a PDF from the invoice and returns it as bytes.
func (g *Generator) GenerateBytes(invoice models.Invoice) ([]byte, error) {
	pdfDoc, err := g.buildPDF(invoice)
	if err != nil {
		return nil, err
	}
	return pdfDoc.GetBytesPdf(), nil
}

// buildPDF constructs the full PDF document by composing sections.
func (g *Generator) buildPDF(invoice models.Invoice) (*gopdf.GoPdf, error) {
	layout := LayoutForInvoice(invoice)
	colWidths := layout.ColWidths(invoice)

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: layout.PageW, H: layout.PageH},
	})

	// Load fonts
	if err := loadFonts(pdf, g.fontDir); err != nil {
		return nil, err
	}

	pdf.AddPage()

	// Build draw context
	ctx := &sections.DrawContext{
		PDF:     pdf,
		Invoice: invoice,
		Layout: sections.LayoutInfo{
			PageW:           layout.PageW,
			PageH:           layout.PageH,
			Margin:          layout.Margin,
			ContentW:        layout.ContentW,
			TitleSize:       layout.TitleSize,
			HeadingSize:     layout.HeadingSize,
			BodySize:        layout.BodySize,
			SmallSize:       layout.SmallSize,
			TableHeadSize:   layout.TableHeadSize,
			FooterSize:      layout.FooterSize,
			SectionGap:      layout.SectionGap,
			LineHeight:      layout.LineHeight,
			RowMinHeight:    layout.RowMinHeight,
			HeaderRowH:      layout.HeaderRowH,
			TotalsRowH:      layout.TotalsRowH,
			TotalsFinalRowH: layout.TotalsFinalRowH,
			QRSize:          layout.QRSize,
			ColWidths:       colWidths,
		},
		CurrentY: layout.Margin,
	}

	// Compose sections based on invoice type
	sections.DrawHeader(ctx)
	sections.DrawInvoiceInfo(ctx)
	sections.DrawSellerInfo(ctx)

	if IsB2B(invoice) {
		sections.DrawBuyerInfo(ctx)
	}

	if IsCreditOrDebit(invoice) {
		sections.DrawCreditDebitReason(ctx)
	}

	sections.DrawProductsTable(ctx)
	sections.DrawTotals(ctx)
	sections.DrawFooter(ctx)
	sections.DrawQRCode(ctx)

	return pdf, nil
}

func loadFonts(pdf *gopdf.GoPdf, fontDir string) error {
	regularPath := fontDir + "/Amiri-Regular.ttf"
	if err := pdf.AddTTFFont("Amiri", regularPath); err != nil {
		return fmt.Errorf("failed to load regular font: %w", err)
	}

	boldPath := fontDir + "/Amiri-Bold.ttf"
	if err := pdf.AddTTFFont("AmiriBold", boldPath); err != nil {
		// Fallback to regular
		_ = pdf.AddTTFFont("AmiriBold", regularPath)
	}

	return nil
}

// --- Package-level convenience functions ---

// GenerateInvoice creates a PDF invoice and saves it to a file.
func GenerateInvoice(invoice models.Invoice, filename, fontDir string) error {
	gen := NewGenerator(fontDir)
	return gen.Generate(invoice, filename)
}

// GenerateInvoiceBytes creates a PDF invoice and returns it as bytes.
func GenerateInvoiceBytes(invoice models.Invoice, fontDir string) ([]byte, error) {
	gen := NewGenerator(fontDir)
	return gen.GenerateBytes(invoice)
}
