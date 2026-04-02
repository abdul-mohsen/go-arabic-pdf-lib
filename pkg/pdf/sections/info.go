package sections

import (
	"github.com/ssda/bill-generator/pkg/models"
	"github.com/ssda/bill-generator/pkg/textutil"
	"github.com/signintech/gopdf"
)

// DrawInvoiceInfo draws invoice number and date.
func DrawInvoiceInfo(ctx *DrawContext) {
	inv := ctx.Invoice
	l := ctx.Layout

	// Invoice Number
	drawLabelValue(ctx, inv.Labels.InvoiceNumber, inv.InvoiceNumber, l.BodySize)

	// Date
	dateFormat := inv.DateFormat
	if dateFormat == "" {
		dateFormat = "2006-01-02 15:04:05"
	}
	dateStr := inv.Date.Format(dateFormat)
	drawLabelValue(ctx, inv.Labels.Date, dateStr, l.BodySize)
}

// DrawSellerInfo draws the seller information in a bordered table.
func DrawSellerInfo(ctx *DrawContext) {
	inv := ctx.Invoice
	l := ctx.Layout

	if l.ContentW > 400 {
		// A4: draw as bordered table
		nameLabel := "Name:"
		if inv.IsRTL {
			nameLabel = "الاسم:"
		}
		drawPartyTable(ctx, inv.Labels.SellerInfo, []tableRow{
			{nameLabel, inv.Seller.Name, true},
			{addressLabel(inv.IsRTL), inv.Seller.Address, false},
			{inv.Labels.VATRegistration, inv.Seller.VATRegistrationNo, false},
			{inv.Labels.CommercialReg, inv.Seller.CommercialRegNo, false},
		})
	} else {
		// Thermal: compact centered layout
		drawPartyCompact(ctx, inv.Labels.SellerInfo, inv.Seller, inv.Labels.VATRegistration, inv.Labels.CommercialReg)
	}
}

// DrawBuyerInfo draws the buyer information in a bordered table (B2B only).
func DrawBuyerInfo(ctx *DrawContext) {
	inv := ctx.Invoice
	l := ctx.Layout

	if l.ContentW > 400 {
		nameLabel := "Name:"
		if inv.IsRTL {
			nameLabel = "الاسم:"
		}
		drawPartyTable(ctx, inv.Labels.BuyerInfo, []tableRow{
			{nameLabel, inv.Buyer.Name, true},
			{addressLabel(inv.IsRTL), inv.Buyer.Address, false},
			{inv.Labels.BuyerVATRegistration, inv.Buyer.VATRegistrationNo, false},
			{inv.Labels.BuyerCommercialReg, inv.Buyer.CommercialRegNo, false},
		})
	} else {
		drawPartyCompact(ctx, inv.Labels.BuyerInfo, inv.Buyer, inv.Labels.BuyerVATRegistration, inv.Labels.BuyerCommercialReg)
	}
}

func addressLabel(isRTL bool) string {
	if isRTL {
		return "العنوان:"
	}
	return "Address:"
}

type tableRow struct {
	Label string
	Value string
	Bold  bool
}

// drawPartyTable draws a 2-column bordered table for party info (A4).
func drawPartyTable(ctx *DrawContext, title string, rows []tableRow) {
	pdf := ctx.PDF
	inv := ctx.Invoice
	l := ctx.Layout

	// Section heading
	if title != "" {
		setFont(pdf, true, l.BodySize)
		pdf.SetTextColor(0, 0, 0)
		textutil.DrawTextCentered(pdf, title, l.Margin, ctx.CurrentY, l.ContentW, inv.IsRTL)
		ctx.CurrentY += l.LineHeight + 2
	}

	tableX := l.Margin
	tableW := l.ContentW
	labelColW := tableW * 0.40
	valueColW := tableW * 0.60
	rowH := l.LineHeight + 6

	pdf.SetStrokeColor(180, 180, 180)
	pdf.SetLineWidth(0.5)

	for _, row := range rows {
		if row.Value == "" {
			continue
		}

		// Calculate row height based on value text wrapping
		if row.Bold {
			setFont(pdf, true, l.BodySize)
		} else {
			pdf.SetFont("Amiri", "", l.SmallSize+1)
		}
		valColInner := valueColW - 8 // padding
		_, valH := textutil.WrapText(pdf, row.Value, valColInner, l.LineHeight, inv.IsRTL)
		actualRowH := valH + 6
		if actualRowH < rowH {
			actualRowH = rowH
		}

		// Draw cell borders
		if inv.IsRTL {
			pdf.RectFromUpperLeftWithStyle(tableX, ctx.CurrentY, valueColW, actualRowH, "D")
			pdf.RectFromUpperLeftWithStyle(tableX+valueColW, ctx.CurrentY, labelColW, actualRowH, "D")
		} else {
			pdf.RectFromUpperLeftWithStyle(tableX, ctx.CurrentY, labelColW, actualRowH, "D")
			pdf.RectFromUpperLeftWithStyle(tableX+labelColW, ctx.CurrentY, valueColW, actualRowH, "D")
		}

		// Draw label
		setFont(pdf, true, l.SmallSize+1)
		pdf.SetTextColor(80, 80, 80)
		drawInfoCell(pdf, row.Label, tableX, ctx.CurrentY, labelColW, valueColW, actualRowH, inv.IsRTL, true)

		// Draw value (wrapped)
		if row.Bold {
			setFont(pdf, true, l.BodySize)
		} else {
			pdf.SetFont("Amiri", "", l.SmallSize+1)
		}
		pdf.SetTextColor(0, 0, 0)
		drawInfoCellWrapped(ctx, row.Value, tableX, ctx.CurrentY, labelColW, valueColW, actualRowH, inv.IsRTL)

		ctx.CurrentY += actualRowH
	}

	pdf.SetStrokeColor(0, 0, 0)
	ctx.CurrentY += 4
}

// drawInfoCell draws text in the label or value column of the party table (single line).
func drawInfoCell(pdf *gopdf.GoPdf, text string, tableX, y, labelW, valueW, rowH float64, isRTL, isLabel bool) {
	processed := textutil.ProcessText(text, isRTL)
	tw, _ := pdf.MeasureTextWidth(processed)
	pad := 4.0
	textY := y + (rowH-10)/2 + 1

	if isRTL {
		if isLabel {
			pdf.SetXY(tableX+valueW+labelW-tw-pad, textY)
		} else {
			pdf.SetXY(tableX+pad, textY)
		}
	} else {
		if isLabel {
			pdf.SetXY(tableX+pad, textY)
		} else {
			pdf.SetXY(tableX+labelW+pad, textY)
		}
	}
	pdf.Cell(nil, processed)
	_ = tw
}

// drawInfoCellWrapped draws value text with word wrapping in the value column.
func drawInfoCellWrapped(ctx *DrawContext, text string, tableX, y, labelW, valueW, rowH float64, isRTL bool) {
	pdf := ctx.PDF
	l := ctx.Layout
	pad := 4.0
	textY := y + 3

	var cellX, cellW float64
	if isRTL {
		cellX = tableX + pad
		cellW = valueW - 2*pad
	} else {
		cellX = tableX + labelW + pad
		cellW = valueW - 2*pad
	}

	lines, _ := textutil.WrapText(pdf, text, cellW, l.LineHeight, isRTL)
	for i, line := range lines {
		if isRTL {
			tw, _ := pdf.MeasureTextWidth(line)
			pdf.SetXY(cellX+cellW-tw, textY+float64(i)*l.LineHeight)
		} else {
			pdf.SetXY(cellX, textY+float64(i)*l.LineHeight)
		}
		pdf.Cell(nil, line)
	}
}

// drawPartyCompact draws party info in compact centered form (thermal).
func drawPartyCompact(ctx *DrawContext, sectionLabel string, party models.PartyInfo, vatLabel, crLabel string) {
	pdf := ctx.PDF
	inv := ctx.Invoice
	l := ctx.Layout

	// Section label
	if sectionLabel != "" {
		setFont(pdf, true, l.BodySize)
		pdf.SetTextColor(0, 0, 0)
		textutil.DrawTextCentered(pdf, sectionLabel, l.Margin, ctx.CurrentY, l.ContentW, inv.IsRTL)
		ctx.CurrentY += l.LineHeight
	}

	// Name (wrapped)
	setFont(pdf, true, l.HeadingSize)
	nameH := textutil.DrawTextCenteredWrapped(pdf, party.Name, l.Margin, ctx.CurrentY, l.ContentW, l.LineHeight, inv.IsRTL)
	ctx.CurrentY += nameH + 2

	// Address (wrapped)
	if party.Address != "" {
		pdf.SetFont("Amiri", "", l.BodySize)
		addrH := textutil.DrawTextCenteredWrapped(pdf, party.Address, l.Margin, ctx.CurrentY, l.ContentW, l.LineHeight, inv.IsRTL)
		ctx.CurrentY += addrH + 2
	}

	// VAT
	if party.VATRegistrationNo != "" && vatLabel != "" {
		drawLabelValue(ctx, vatLabel, party.VATRegistrationNo, l.SmallSize+1)
	}

	// Commercial Reg
	if party.CommercialRegNo != "" && crLabel != "" {
		drawLabelValue(ctx, crLabel, party.CommercialRegNo, l.SmallSize+1)
	}

	ctx.CurrentY += 2
}
