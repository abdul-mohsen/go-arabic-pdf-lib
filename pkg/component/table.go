package component

import (
	"fmt"
	"github.com/abdul-mohsen/go-arabic-pdf-lib/arabictext"

	"github.com/signintech/gopdf"
)

// TableColumn defines a column in a table.
type TableColumn struct {
	Header []string // Multi-line header text
	Width  float64
	Align  Alignment
}

// TableRow represents a row of data.
type TableRow struct {
	Cells     []string
	Heights   []float64 // Optional per-cell height hints
	WrapCells []bool    // Which cells should wrap
}

// Table is a component for rendering tabular data.
type Table struct {
	Columns []TableColumn
	Rows    []TableRow
	Options Options
}

// NewTable creates a new table component.
func NewTable(columns []TableColumn, opts ...OptionFunc) *Table {
	options := DefaultOptions()
	ApplyOptions(&options, opts...)
	return &Table{
		Columns: columns,
		Options: options,
	}
}

// AddRow adds a row to the table.
func (t *Table) AddRow(cells []string, wrapCells []bool) {
	row := TableRow{
		Cells:     cells,
		WrapCells: wrapCells,
	}
	t.Rows = append(t.Rows, row)
}

// Draw renders the table.
func (t *Table) Draw(pdf *gopdf.GoPdf) float64 {
	opts := t.Options
	totalHeight := 0.0

	// Draw header
	headerHeight := t.drawHeader(pdf)
	totalHeight += headerHeight

	// Draw rows
	for _, row := range t.Rows {
		rowHeight := t.drawRow(pdf, row, opts.Position.Y+totalHeight)
		totalHeight += rowHeight
	}

	return totalHeight
}

func (t *Table) drawHeader(pdf *gopdf.GoPdf) float64 {
	opts := t.Options
	headerHeight := 28.0

	fontName := opts.Style.FontName + "Bold"
	if err := pdf.SetFont(fontName, "", 7); err != nil {
		pdf.SetFont(opts.Style.FontName, "", 7)
	}

	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.5)
	pdf.SetTextColor(0, 0, 0)

	xPos := opts.Position.X
	for i, col := range t.Columns {
		// Draw cell border
		pdf.RectFromUpperLeftWithStyle(xPos, opts.Position.Y, col.Width, headerHeight, "D")

		// Draw header text (multi-line)
		for j, line := range col.Header {
			text := line
			if opts.RTL {
				text = arabictext.Process(line)
			}
			if text == "" {
				continue
			}

			textW, _ := pdf.MeasureTextWidth(text)
			textX := xPos + (col.Width-textW)/2
			textY := opts.Position.Y + 4 + float64(j)*10

			pdf.SetXY(textX, textY)
			pdf.Cell(nil, text)
		}

		_ = i
		xPos += col.Width
	}

	return headerHeight
}

func (t *Table) drawRow(pdf *gopdf.GoPdf, row TableRow, y float64) float64 {
	opts := t.Options

	if err := pdf.SetFont(opts.Style.FontName, "", int(opts.Style.FontSize)); err != nil {
		return 0
	}

	// Calculate row height (find max cell height)
	baseRowHeight := 12.0
	minRowHeight := 18.0
	maxHeight := minRowHeight

	for i, cell := range row.Cells {
		if i < len(row.WrapCells) && row.WrapCells[i] && i < len(t.Columns) {
			text := cell
			if opts.RTL {
				text = arabictext.Process(cell)
			}
			_, cellHeight := measureWrappedHeight(pdf, text, t.Columns[i].Width-6, baseRowHeight)
			if cellHeight+6 > maxHeight {
				maxHeight = cellHeight + 6
			}
		}
	}

	rowHeight := maxHeight

	// Draw cell borders and content
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetTextColor(0, 0, 0)

	xPos := opts.Position.X
	for i, col := range t.Columns {
		// Draw border
		pdf.RectFromUpperLeftWithStyle(xPos, y, col.Width, rowHeight, "D")

		// Draw cell content
		if i < len(row.Cells) {
			cellText := row.Cells[i]
			if opts.RTL && i < len(row.WrapCells) && row.WrapCells[i] {
				cellText = arabictext.Process(cellText)
			}

			shouldWrap := i < len(row.WrapCells) && row.WrapCells[i]
			if shouldWrap {
				drawWrappedCell(pdf, cellText, xPos, y+3, col.Width, baseRowHeight, col.Align, opts.RTL)
			} else {
				textW, _ := pdf.MeasureTextWidth(cellText)
				textX := xPos + opts.Style.Padding
				switch col.Align {
				case AlignCenter:
					textX = xPos + (col.Width-textW)/2
				case AlignRight:
					textX = xPos + col.Width - textW - opts.Style.Padding
				}
				pdf.SetXY(textX, y+3)
				pdf.Cell(nil, cellText)
			}
		}

		xPos += col.Width
	}

	return rowHeight
}

func measureWrappedHeight(pdf *gopdf.GoPdf, text string, maxWidth, lineHeight float64) ([]string, float64) {
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

	return lines, float64(len(lines)) * lineHeight
}

func drawWrappedCell(pdf *gopdf.GoPdf, text string, x, y, width, lineHeight float64, align Alignment, rtl bool) {
	lines, _ := measureWrappedHeight(pdf, text, width-6, lineHeight)

	for _, line := range lines {
		lineW, _ := pdf.MeasureTextWidth(line)
		textX := x + 3
		switch align {
		case AlignCenter:
			textX = x + (width-lineW)/2
		case AlignRight:
			textX = x + width - lineW - 3
		}
		pdf.SetXY(textX, y)
		pdf.Cell(nil, line)
		y += lineHeight
	}
}

// TotalsRow represents a row in the totals section.
type TotalsRow struct {
	Label string
	Value string
	Bold  bool
	Thick bool // Thicker border
}

// TotalsTable displays summary totals.
type TotalsTable struct {
	Rows       []TotalsRow
	LabelWidth float64
	ValueWidth float64
	Options    Options
}

// NewTotalsTable creates a totals table.
func NewTotalsTable(labelWidth, valueWidth float64, opts ...OptionFunc) *TotalsTable {
	options := DefaultOptions()
	ApplyOptions(&options, opts...)
	return &TotalsTable{
		LabelWidth: labelWidth,
		ValueWidth: valueWidth,
		Options:    options,
	}
}

// AddRow adds a row to the totals table.
func (t *TotalsTable) AddRow(label, value string, bold, thick bool) {
	t.Rows = append(t.Rows, TotalsRow{
		Label: label,
		Value: value,
		Bold:  bold,
		Thick: thick,
	})
}

// Draw renders the totals table.
func (t *TotalsTable) Draw(pdf *gopdf.GoPdf) float64 {
	opts := t.Options
	totalHeight := 0.0

	y := opts.Position.Y

	for _, row := range t.Rows {
		rowHeight := 16.0
		if row.Bold {
			rowHeight = 18.0
		}

		// Draw borders
		pdf.SetStrokeColor(0, 0, 0)
		if row.Thick {
			pdf.SetLineWidth(1.0)
		} else {
			pdf.SetLineWidth(0.5)
		}

		pdf.RectFromUpperLeftWithStyle(opts.Position.X, y, t.ValueWidth, rowHeight, "D")
		pdf.RectFromUpperLeftWithStyle(opts.Position.X+t.ValueWidth, y, t.LabelWidth, rowHeight, "D")

		// Set font
		fontName := opts.Style.FontName
		fontSize := opts.Style.FontSize
		if row.Bold {
			fontName += "Bold"
			fontSize = 10
		}
		if err := pdf.SetFont(fontName, "", int(fontSize)); err != nil {
			pdf.SetFont(opts.Style.FontName, "", int(fontSize))
		}

		pdf.SetTextColor(0, 0, 0)

		// Draw value (right-aligned in value cell)
		valueW, _ := pdf.MeasureTextWidth(row.Value)
		pdf.SetXY(opts.Position.X+t.ValueWidth-valueW-3, y+3)
		pdf.Cell(nil, row.Value)

		// Draw label
		label := row.Label
		if opts.RTL {
			label = arabictext.Process(label)
		}
		labelW, _ := pdf.MeasureTextWidth(label)

		if opts.RTL {
			pdf.SetXY(opts.Position.X+t.ValueWidth+t.LabelWidth-labelW-2, y+3)
		} else {
			pdf.SetXY(opts.Position.X+t.ValueWidth+3, y+3)
		}
		pdf.Cell(nil, label)

		y += rowHeight
		totalHeight += rowHeight
	}

	return totalHeight
}

// FormatNumber formats a number for display.
func FormatNumber(n float64, decimals int) string {
	format := fmt.Sprintf("%%.%df", decimals)
	return fmt.Sprintf(format, n)
}
