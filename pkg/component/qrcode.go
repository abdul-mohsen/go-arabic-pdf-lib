package component

import (
	"os"

	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
)

// QRCode renders a QR code image.
type QRCode struct {
	Data    string
	Options Options
}

// NewQRCode creates a QR code component.
func NewQRCode(data string, opts ...OptionFunc) *QRCode {
	options := DefaultOptions()
	options.Size = Size{Width: 55, Height: 55}
	ApplyOptions(&options, opts...)
	return &QRCode{
		Data:    data,
		Options: options,
	}
}

// Draw renders the QR code.
func (q *QRCode) Draw(pdf *gopdf.GoPdf) float64 {
	opts := q.Options

	// Generate QR code to temp file
	tmpFile := "/tmp/invoice_qr.png"
	err := qrcode.WriteFile(q.Data, qrcode.High, 256, tmpFile)
	if err != nil {
		return 0
	}
	defer os.Remove(tmpFile)

	// Calculate centered position if width is provided
	x := opts.Position.X
	if opts.Alignment == AlignCenter && opts.Size.Width > 0 {
		// Position.X is the container left edge, calculate center
		// But we need to know container width - use Size.Width as QR size
	}

	pdf.Image(tmpFile, x, opts.Position.Y, &gopdf.Rect{
		W: opts.Size.Width,
		H: opts.Size.Height,
	})

	return opts.Size.Height
}

// CenteredQRCode creates a QR code centered within a container.
type CenteredQRCode struct {
	Data           string
	QRSize         float64
	ContainerWidth float64
	Options        Options
}

// NewCenteredQRCode creates a centered QR code.
func NewCenteredQRCode(data string, qrSize, containerWidth float64, opts ...OptionFunc) *CenteredQRCode {
	options := DefaultOptions()
	ApplyOptions(&options, opts...)
	return &CenteredQRCode{
		Data:           data,
		QRSize:         qrSize,
		ContainerWidth: containerWidth,
		Options:        options,
	}
}

// Draw renders the centered QR code.
func (q *CenteredQRCode) Draw(pdf *gopdf.GoPdf) float64 {
	opts := q.Options

	// Generate QR code
	tmpFile := "/tmp/invoice_qr.png"
	err := qrcode.WriteFile(q.Data, qrcode.High, 256, tmpFile)
	if err != nil {
		return 0
	}
	defer os.Remove(tmpFile)

	// Calculate centered X position
	x := opts.Position.X + (q.ContainerWidth-q.QRSize)/2

	pdf.Image(tmpFile, x, opts.Position.Y, &gopdf.Rect{
		W: q.QRSize,
		H: q.QRSize,
	})

	return q.QRSize
}
