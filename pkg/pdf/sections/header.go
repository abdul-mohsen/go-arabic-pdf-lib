package sections

import "github.com/ssda/bill-generator/pkg/textutil"

// DrawHeader draws the invoice title centered.
func DrawHeader(ctx *DrawContext) {
	setFont(ctx.PDF, true, ctx.Layout.TitleSize)
	ctx.PDF.SetTextColor(0, 0, 0)
	textutil.DrawTextCentered(ctx.PDF, ctx.Invoice.Title, ctx.Layout.Margin, ctx.CurrentY+4, ctx.Layout.ContentW, ctx.Invoice.IsRTL)
	ctx.CurrentY += ctx.Layout.TitleSize + 4
}
