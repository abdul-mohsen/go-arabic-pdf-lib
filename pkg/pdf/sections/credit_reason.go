package sections

import "github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/textutil"

// DrawCreditDebitReason draws the reason for a credit or debit note.
func DrawCreditDebitReason(ctx *DrawContext) {
	inv := ctx.Invoice
	l := ctx.Layout
	pdf := ctx.PDF

	if inv.NoteReason == "" {
		return
	}

	label := inv.Labels.NoteReason
	if label == "" {
		if inv.IsRTL {
			label = "السبب:"
		} else {
			label = "Reason:"
		}
	}

	setFont(pdf, true, l.BodySize)
	pdf.SetTextColor(0, 0, 0)
	labelText := textutil.ProcessText(label, inv.IsRTL)
	labelW, _ := pdf.MeasureTextWidth(labelText)

	if inv.IsRTL {
		pdf.SetXY(l.Margin+l.ContentW-labelW-3, ctx.CurrentY)
	} else {
		pdf.SetXY(l.Margin+3, ctx.CurrentY)
	}
	pdf.Cell(nil, labelText)
	ctx.CurrentY += l.LineHeight

	// Reason text (auto-wrap)
	pdf.SetFont("Amiri", "", l.BodySize)
	h := textutil.DrawWrappedText(pdf, inv.NoteReason, l.Margin, ctx.CurrentY, l.ContentW, l.LineHeight, inv.IsRTL)
	ctx.CurrentY += h + 2
}
