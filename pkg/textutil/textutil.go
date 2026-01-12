// Package textutil provides text processing utilities for PDF generation.
package textutil

import (
	"bill-generator/arabictext"

	"github.com/signintech/gopdf"
)

// ProcessText processes text based on language settings.
// For Arabic (RTL), it applies reshaping and reversal.
// For English (LTR), it returns the text as-is.
func ProcessText(text string, isRTL bool) string {
	if isRTL {
		return arabictext.Process(text)
	}
	return text
}

// WrapText splits text into multiple lines that fit within maxWidth.
// Returns the lines and the total height needed.
func WrapText(pdf *gopdf.GoPdf, text string, maxWidth float64, lineHeight float64, isRTL bool) ([]string, float64) {
	processedText := ProcessText(text, isRTL)

	// Check if text fits in one line
	textWidth, _ := pdf.MeasureTextWidth(processedText)
	if textWidth <= maxWidth {
		return []string{processedText}, lineHeight
	}

	// Need to wrap - split by characters
	var lines []string
	runes := []rune(text)

	currentLine := ""
	for i := 0; i < len(runes); i++ {
		testLine := currentLine + string(runes[i])
		testProcessed := ProcessText(testLine, isRTL)
		testWidth, _ := pdf.MeasureTextWidth(testProcessed)

		if testWidth > maxWidth && currentLine != "" {
			// Current line is full, save it and start new line
			lines = append(lines, ProcessText(currentLine, isRTL))
			currentLine = string(runes[i])
		} else {
			currentLine = testLine
		}
	}

	// Add the last line
	if currentLine != "" {
		lines = append(lines, ProcessText(currentLine, isRTL))
	}

	if len(lines) == 0 {
		lines = []string{processedText}
	}

	return lines, float64(len(lines)) * lineHeight
}

// DrawTextCentered draws text centered within a given width.
func DrawTextCentered(pdf *gopdf.GoPdf, text string, x, y, width float64, isRTL bool) {
	processedText := ProcessText(text, isRTL)
	textWidth, _ := pdf.MeasureTextWidth(processedText)
	centerX := x + (width-textWidth)/2
	pdf.SetXY(centerX, y)
	pdf.Cell(nil, processedText)
}

// DrawTextRight draws text right-aligned within a given width.
func DrawTextRight(pdf *gopdf.GoPdf, text string, x, y, width float64, isRTL bool) {
	processedText := ProcessText(text, isRTL)
	textWidth, _ := pdf.MeasureTextWidth(processedText)
	rightX := x + width - textWidth - 2
	pdf.SetXY(rightX, y)
	pdf.Cell(nil, processedText)
}

// DrawTextLeft draws text left-aligned.
func DrawTextLeft(pdf *gopdf.GoPdf, text string, x, y float64, isRTL bool) {
	processedText := ProcessText(text, isRTL)
	pdf.SetXY(x, y)
	pdf.Cell(nil, processedText)
}

// DrawWrappedText draws wrapped text aligned appropriately.
// For RTL languages, text is right-aligned; for LTR, left-aligned.
func DrawWrappedText(pdf *gopdf.GoPdf, text string, x, y, width, lineHeight float64, isRTL bool) float64 {
	lines, totalHeight := WrapText(pdf, text, width-6, lineHeight, isRTL)

	for i, line := range lines {
		lineWidth, _ := pdf.MeasureTextWidth(line)
		var drawX float64
		if isRTL {
			drawX = x + width - lineWidth - 3
		} else {
			drawX = x + 3
		}
		pdf.SetXY(drawX, y+float64(i)*lineHeight)
		pdf.Cell(nil, line)
	}

	return totalHeight
}
