// Package sections provides composable PDF drawing blocks for invoice generation.
package sections

import (
	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/models"
	"github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/textutil"

	"github.com/signintech/gopdf"
)

// DrawContext holds shared state for all section draw calls.
type DrawContext struct {
	PDF      *gopdf.GoPdf
	Invoice  models.Invoice
	Layout   LayoutInfo
	CurrentY float64
}

// LayoutInfo is the subset of layout data sections need.
type LayoutInfo struct {
	PageW           float64
	PageH           float64
	Margin          float64
	ContentW        float64
	TitleSize       float64
	HeadingSize     float64
	BodySize        float64
	SmallSize       float64
	TableHeadSize   float64
	FooterSize      float64
	SectionGap      float64
	LineHeight      float64
	RowMinHeight    float64
	HeaderRowH      float64
	TotalsRowH      float64
	TotalsFinalRowH float64
	QRSize          float64
	ColWidths       []float64 // Already direction-aware
}

// setFont sets font with bold fallback.
func setFont(pdf *gopdf.GoPdf, bold bool, size float64) {
	if bold {
		if err := pdf.SetFont("AmiriBold", "", size); err != nil {
			pdf.SetFont("Amiri", "", size)
		}
	} else {
		pdf.SetFont("Amiri", "", size)
	}
}

// drawLabelValue draws a label on one side and value on the other.
func drawLabelValue(ctx *DrawContext, label, value string, fontSize float64) {
	pdf := ctx.PDF
	inv := ctx.Invoice
	l := ctx.Layout

	pdf.SetFont("Amiri", "", fontSize)
	pdf.SetTextColor(0, 0, 0)

	labelText := textutil.ProcessText(label, inv.IsRTL)
	labelW, _ := pdf.MeasureTextWidth(labelText)

	if inv.IsRTL {
		pdf.SetXY(l.Margin+l.ContentW-labelW-3, ctx.CurrentY)
		pdf.Cell(nil, labelText)
		pdf.SetXY(l.Margin+3, ctx.CurrentY)
		pdf.Cell(nil, value)
	} else {
		pdf.SetXY(l.Margin+3, ctx.CurrentY)
		pdf.Cell(nil, labelText)
		valueW, _ := pdf.MeasureTextWidth(value)
		pdf.SetXY(l.Margin+l.ContentW-valueW-3, ctx.CurrentY)
		pdf.Cell(nil, value)
	}
	ctx.CurrentY += l.LineHeight
}
