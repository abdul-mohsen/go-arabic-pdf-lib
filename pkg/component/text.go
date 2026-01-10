package component

import (
	"bill-generator/arabictext"

	"github.com/signintech/gopdf"
)

// TextBlock is a simple text component with configurable options.
type TextBlock struct {
	Text    string
	Options Options
}

// NewTextBlock creates a text block with the given text and options.
func NewTextBlock(text string, opts ...OptionFunc) *TextBlock {
	options := DefaultOptions()
	ApplyOptions(&options, opts...)
	return &TextBlock{
		Text:    text,
		Options: options,
	}
}

// Draw renders the text block to the PDF.
func (t *TextBlock) Draw(pdf *gopdf.GoPdf) float64 {
	opts := t.Options

	// Set font
	fontName := opts.Style.FontName
	if opts.Style.Bold {
		fontName = fontName + "Bold"
	}
	if err := pdf.SetFont(fontName, "", int(opts.Style.FontSize)); err != nil {
		pdf.SetFont(opts.Style.FontName, "", int(opts.Style.FontSize))
	}

	// Process text for RTL if needed
	text := t.Text
	if opts.RTL {
		text = arabictext.Process(text)
	}

	// Calculate position based on alignment
	textWidth, _ := pdf.MeasureTextWidth(text)
	x := opts.Position.X + opts.Style.Padding

	switch opts.Alignment {
	case AlignCenter:
		x = opts.Position.X + (opts.Size.Width-textWidth)/2
	case AlignRight:
		x = opts.Position.X + opts.Size.Width - textWidth - opts.Style.Padding
	}

	// Draw border if enabled
	if opts.Border {
		pdf.SetStrokeColor(0, 0, 0)
		pdf.SetLineWidth(0.5)
		pdf.RectFromUpperLeftWithStyle(opts.Position.X, opts.Position.Y, opts.Size.Width, opts.Size.Height, "D")
	}

	// Draw text
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(x, opts.Position.Y+opts.Style.Padding)
	pdf.Cell(nil, text)

	return opts.Style.LineHeight
}

// LabelValuePair renders a label on one side and value on the other.
type LabelValuePair struct {
	Label   string
	Value   string
	Options Options
}

// NewLabelValuePair creates a label-value pair component.
func NewLabelValuePair(label, value string, opts ...OptionFunc) *LabelValuePair {
	options := DefaultOptions()
	ApplyOptions(&options, opts...)
	return &LabelValuePair{
		Label:   label,
		Value:   value,
		Options: options,
	}
}

// Draw renders the label-value pair.
func (lv *LabelValuePair) Draw(pdf *gopdf.GoPdf) float64 {
	opts := lv.Options

	// Set font
	fontName := opts.Style.FontName
	if opts.Style.Bold {
		fontName = fontName + "Bold"
	}
	if err := pdf.SetFont(fontName, "", int(opts.Style.FontSize)); err != nil {
		pdf.SetFont(opts.Style.FontName, "", int(opts.Style.FontSize))
	}

	pdf.SetTextColor(0, 0, 0)

	// Process text
	label := lv.Label
	value := lv.Value
	if opts.RTL {
		label = arabictext.Process(label)
	}

	labelW, _ := pdf.MeasureTextWidth(label)
	valueW, _ := pdf.MeasureTextWidth(value)

	if opts.RTL {
		// Label on right, value on left
		pdf.SetXY(opts.Position.X+opts.Size.Width-labelW-opts.Style.Padding, opts.Position.Y)
		pdf.Cell(nil, label)
		pdf.SetXY(opts.Position.X+opts.Style.Padding, opts.Position.Y)
		pdf.Cell(nil, value)
	} else {
		// Label on left, value on right
		pdf.SetXY(opts.Position.X+opts.Style.Padding, opts.Position.Y)
		pdf.Cell(nil, label)
		pdf.SetXY(opts.Position.X+opts.Size.Width-valueW-opts.Style.Padding, opts.Position.Y)
		pdf.Cell(nil, value)
	}

	return opts.Style.LineHeight
}

// Header is a centered header text component.
type Header struct {
	Text    string
	Options Options
}

// NewHeader creates a header component.
func NewHeader(text string, opts ...OptionFunc) *Header {
	options := DefaultOptions()
	options.Alignment = AlignCenter
	options.Style.Bold = true
	options.Style.FontSize = 14
	ApplyOptions(&options, opts...)
	return &Header{
		Text:    text,
		Options: options,
	}
}

// Draw renders the header.
func (h *Header) Draw(pdf *gopdf.GoPdf) float64 {
	opts := h.Options

	fontName := opts.Style.FontName
	if opts.Style.Bold {
		fontName = fontName + "Bold"
	}
	if err := pdf.SetFont(fontName, "", int(opts.Style.FontSize)); err != nil {
		pdf.SetFont(opts.Style.FontName, "", int(opts.Style.FontSize))
	}

	text := h.Text
	if opts.RTL {
		text = arabictext.Process(text)
	}

	pdf.SetTextColor(0, 0, 0)
	textWidth, _ := pdf.MeasureTextWidth(text)
	x := opts.Position.X + (opts.Size.Width-textWidth)/2
	pdf.SetXY(x, opts.Position.Y+4)
	pdf.Cell(nil, text)

	return opts.Style.LineHeight + 6
}

// WrappedText handles long text that needs to wrap across lines.
type WrappedText struct {
	Text    string
	Options Options
}

// NewWrappedText creates a wrapped text component.
func NewWrappedText(text string, opts ...OptionFunc) *WrappedText {
	options := DefaultOptions()
	options.WrapText = true
	ApplyOptions(&options, opts...)
	return &WrappedText{
		Text:    text,
		Options: options,
	}
}

// Draw renders the wrapped text.
func (w *WrappedText) Draw(pdf *gopdf.GoPdf) float64 {
	opts := w.Options

	fontName := opts.Style.FontName
	if opts.Style.Bold {
		fontName = fontName + "Bold"
	}
	if err := pdf.SetFont(fontName, "", int(opts.Style.FontSize)); err != nil {
		pdf.SetFont(opts.Style.FontName, "", int(opts.Style.FontSize))
	}

	pdf.SetTextColor(0, 0, 0)

	text := w.Text
	if opts.RTL {
		text = arabictext.Process(text)
	}

	// Simple word wrapping
	maxWidth := opts.Size.Width - (2 * opts.Style.Padding)
	words := splitWords(text)
	lines := []string{}
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		lineW, _ := pdf.MeasureTextWidth(testLine)
		if lineW > maxWidth && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	// Draw lines
	y := opts.Position.Y + opts.Style.Padding
	for _, line := range lines {
		lineW, _ := pdf.MeasureTextWidth(line)
		x := opts.Position.X + opts.Style.Padding

		switch opts.Alignment {
		case AlignCenter:
			x = opts.Position.X + (opts.Size.Width-lineW)/2
		case AlignRight:
			x = opts.Position.X + opts.Size.Width - lineW - opts.Style.Padding
		}

		pdf.SetXY(x, y)
		pdf.Cell(nil, line)
		y += opts.Style.LineHeight
	}

	return float64(len(lines)) * opts.Style.LineHeight
}

func splitWords(text string) []string {
	words := []string{}
	current := ""
	for _, r := range text {
		if r == ' ' || r == '\n' || r == '\t' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(r)
		}
	}
	if current != "" {
		words = append(words, current)
	}
	return words
}
