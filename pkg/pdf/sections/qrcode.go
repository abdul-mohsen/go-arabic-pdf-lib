package sections

import (
	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
)

// DrawQRCode generates and draws a QR code centered.
func DrawQRCode(ctx *DrawContext) {
	data := ctx.Invoice.QRCodeData
	if data == "" {
		return
	}

	png, err := qrcode.Encode(data, qrcode.High, 256)
	if err != nil {
		return
	}

	imgHolder, err := gopdf.ImageHolderByBytes(png)
	if err != nil {
		return
	}

	l := ctx.Layout
	qrX := l.Margin + (l.ContentW-l.QRSize)/2
	ctx.PDF.ImageByHolder(imgHolder, qrX, ctx.CurrentY, &gopdf.Rect{W: l.QRSize, H: l.QRSize})
	ctx.CurrentY += l.QRSize + 4
}
