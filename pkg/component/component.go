// Package component provides a component-based system for building PDF documents.
// Components are reusable building blocks with configurable options for positioning,
// alignment, styling, and RTL support.
package component

import (
	"github.com/signintech/gopdf"
)

// Alignment specifies text alignment within a component.
type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

// Position represents x,y coordinates on the page.
type Position struct {
	X float64
	Y float64
}

// Size represents width and height dimensions.
type Size struct {
	Width  float64
	Height float64
}

// Style holds styling options for components.
type Style struct {
	FontName   string
	FontSize   float64
	Bold       bool
	LineHeight float64
	Padding    float64
}

// DefaultStyle returns sensible default styling.
func DefaultStyle() Style {
	return Style{
		FontName:   "Amiri",
		FontSize:   9,
		Bold:       false,
		LineHeight: 12,
		Padding:    3,
	}
}

// Options configures component behavior and appearance.
type Options struct {
	Position  Position
	Size      Size
	Style     Style
	Alignment Alignment
	RTL       bool
	Border    bool
	WrapText  bool
}

// DefaultOptions returns sensible default options.
func DefaultOptions() Options {
	return Options{
		Style:     DefaultStyle(),
		Alignment: AlignLeft,
		RTL:       false,
		Border:    false,
		WrapText:  false,
	}
}

// Component is the interface for all drawable PDF components.
type Component interface {
	// Draw renders the component to the PDF at the current position.
	// Returns the height consumed by the component.
	Draw(pdf *gopdf.GoPdf) float64
}

// Container manages a collection of components and their layout.
type Container interface {
	// Add appends a component to the container.
	Add(c Component)

	// Render draws all components in order.
	Render(pdf *gopdf.GoPdf) float64
}

// OptionFunc is a function that modifies Options.
type OptionFunc func(*Options)

// WithPosition sets the component position.
func WithPosition(x, y float64) OptionFunc {
	return func(o *Options) {
		o.Position = Position{X: x, Y: y}
	}
}

// WithSize sets the component size.
func WithSize(w, h float64) OptionFunc {
	return func(o *Options) {
		o.Size = Size{Width: w, Height: h}
	}
}

// WithAlignment sets text alignment.
func WithAlignment(a Alignment) OptionFunc {
	return func(o *Options) {
		o.Alignment = a
	}
}

// WithRTL enables right-to-left text direction.
func WithRTL(rtl bool) OptionFunc {
	return func(o *Options) {
		o.RTL = rtl
	}
}

// WithBorder enables border drawing.
func WithBorder(border bool) OptionFunc {
	return func(o *Options) {
		o.Border = border
	}
}

// WithWrapText enables text wrapping.
func WithWrapText(wrap bool) OptionFunc {
	return func(o *Options) {
		o.WrapText = wrap
	}
}

// WithStyle sets the component style.
func WithStyle(s Style) OptionFunc {
	return func(o *Options) {
		o.Style = s
	}
}

// WithFontSize sets just the font size.
func WithFontSize(size float64) OptionFunc {
	return func(o *Options) {
		o.Style.FontSize = size
	}
}

// WithBold enables bold text.
func WithBold(bold bool) OptionFunc {
	return func(o *Options) {
		o.Style.Bold = bold
	}
}

// WithPadding sets the internal padding.
func WithPadding(p float64) OptionFunc {
	return func(o *Options) {
		o.Style.Padding = p
	}
}

// ApplyOptions applies a list of option functions to Options.
func ApplyOptions(opts *Options, fns ...OptionFunc) {
	for _, fn := range fns {
		fn(opts)
	}
}
