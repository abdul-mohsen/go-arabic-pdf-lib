package sections

import "github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/textutil"

// DrawFooter draws the footer text centered (with wrapping).
func DrawFooter(ctx *DrawContext) {
	pdf := ctx.PDF
	l := ctx.Layout

	pdf.SetFont("Amiri", "", l.FooterSize)
	pdf.SetTextColor(0, 0, 0)
	h := textutil.DrawTextCenteredWrapped(pdf, ctx.Invoice.Labels.Footer, l.Margin, ctx.CurrentY, l.ContentW, l.LineHeight, ctx.Invoice.IsRTL)
	ctx.CurrentY += h
}
