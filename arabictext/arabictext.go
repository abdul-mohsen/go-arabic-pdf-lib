// Package arabictext provides Arabic text shaping and RTL rendering for PDF generation.
// It handles proper Arabic letter forms (isolated, initial, medial, final),
// bidirectional text with mixed Arabic/English/numbers, and parentheses mirroring.
package arabictext

import (
	"strings"
	"unicode"
)

// FormType represents the position of a letter in a word
type FormType int

const (
	Isolated FormType = 0
	Final    FormType = 1
	Medial   FormType = 2
	Initial  FormType = 3
)

// LetterForms holds the four forms of an Arabic letter
// Order: [Isolated, Final, Medial, Initial]
type LetterForms [4]rune

// arabicForms maps Arabic letters to their contextual forms
var arabicForms = map[rune]LetterForms{
	// Basic Arabic letters
	'ا': {'ﺍ', 'ﺎ', 'ﺎ', 'ﺍ'}, // Alef (non-connecting)
	'ب': {'ﺏ', 'ﺐ', 'ﺒ', 'ﺑ'}, // Ba
	'ت': {'ﺕ', 'ﺖ', 'ﺘ', 'ﺗ'}, // Ta
	'ث': {'ﺙ', 'ﺚ', 'ﺜ', 'ﺛ'}, // Tha
	'ج': {'ﺝ', 'ﺞ', 'ﺠ', 'ﺟ'}, // Jeem
	'ح': {'ﺡ', 'ﺢ', 'ﺤ', 'ﺣ'}, // Ha
	'خ': {'ﺥ', 'ﺦ', 'ﺨ', 'ﺧ'}, // Kha
	'د': {'ﺩ', 'ﺪ', 'ﺪ', 'ﺩ'}, // Dal (non-connecting)
	'ذ': {'ﺫ', 'ﺬ', 'ﺬ', 'ﺫ'}, // Thal (non-connecting)
	'ر': {'ﺭ', 'ﺮ', 'ﺮ', 'ﺭ'}, // Ra (non-connecting)
	'ز': {'ﺯ', 'ﺰ', 'ﺰ', 'ﺯ'}, // Zay (non-connecting)
	'س': {'ﺱ', 'ﺲ', 'ﺴ', 'ﺳ'}, // Seen
	'ش': {'ﺵ', 'ﺶ', 'ﺸ', 'ﺷ'}, // Sheen
	'ص': {'ﺹ', 'ﺺ', 'ﺼ', 'ﺻ'}, // Sad
	'ض': {'ﺽ', 'ﺾ', 'ﻀ', 'ﺿ'}, // Dad
	'ط': {'ﻁ', 'ﻂ', 'ﻄ', 'ﻃ'}, // Ta
	'ظ': {'ﻅ', 'ﻆ', 'ﻈ', 'ﻇ'}, // Za
	'ع': {'ﻉ', 'ﻊ', 'ﻌ', 'ﻋ'}, // Ain
	'غ': {'ﻍ', 'ﻎ', 'ﻐ', 'ﻏ'}, // Ghain
	'ف': {'ﻑ', 'ﻒ', 'ﻔ', 'ﻓ'}, // Fa
	'ق': {'ﻕ', 'ﻖ', 'ﻘ', 'ﻗ'}, // Qaf
	'ك': {'ﻙ', 'ﻚ', 'ﻜ', 'ﻛ'}, // Kaf
	'ل': {'ﻝ', 'ﻞ', 'ﻠ', 'ﻟ'}, // Lam
	'م': {'ﻡ', 'ﻢ', 'ﻤ', 'ﻣ'}, // Meem
	'ن': {'ﻥ', 'ﻦ', 'ﻨ', 'ﻧ'}, // Noon
	'ه': {'ﻩ', 'ﻪ', 'ﻬ', 'ﻫ'}, // Ha
	'و': {'ﻭ', 'ﻮ', 'ﻮ', 'ﻭ'}, // Waw (non-connecting)
	'ي': {'ﻱ', 'ﻲ', 'ﻴ', 'ﻳ'}, // Ya
	'ى': {'ﻯ', 'ﻰ', 'ﻰ', 'ﻯ'}, // Alef Maksura (non-connecting)
	'ة': {'ﺓ', 'ﺔ', 'ﺔ', 'ﺓ'}, // Ta Marbuta (non-connecting)
	'ء': {'ء', 'ء', 'ء', 'ء'}, // Hamza
	'أ': {'ﺃ', 'ﺄ', 'ﺄ', 'ﺃ'}, // Alef with Hamza above (non-connecting)
	'إ': {'ﺇ', 'ﺈ', 'ﺈ', 'ﺇ'}, // Alef with Hamza below (non-connecting)
	'آ': {'ﺁ', 'ﺂ', 'ﺂ', 'ﺁ'}, // Alef with Madda (non-connecting)
	'ؤ': {'ﺅ', 'ﺆ', 'ﺆ', 'ﺅ'}, // Waw with Hamza (non-connecting)
	'ئ': {'ﺉ', 'ﺊ', 'ﺌ', 'ﺋ'}, // Ya with Hamza
	'ـ': {'ـ', 'ـ', 'ـ', 'ـ'}, // Tatweel (kashida)
	'ﻻ': {'ﻻ', 'ﻼ', 'ﻼ', 'ﻻ'}, // Lam-Alef ligature
}

// nonConnectingLetters contains letters that don't connect to the next letter
var nonConnectingLetters = map[rune]bool{
	'ا': true, 'د': true, 'ذ': true, 'ر': true, 'ز': true, 'و': true,
	'أ': true, 'إ': true, 'آ': true, 'ؤ': true, 'ة': true, 'ى': true,
	'ء': true,
}

// mirroredBrackets maps brackets to their mirrored versions for RTL
var mirroredBrackets = map[rune]rune{
	'(': ')',
	')': '(',
	'[': ']',
	']': '[',
	'{': '}',
	'}': '{',
	'<': '>',
	'>': '<',
	'«': '»',
	'»': '«',
}

// IsArabic checks if a rune is an Arabic character
func IsArabic(r rune) bool {
	return unicode.Is(unicode.Arabic, r)
}

// IsArabicPresentationForm checks if a rune is in Arabic Presentation Forms
func IsArabicPresentationForm(r rune) bool {
	return (r >= 0xFB50 && r <= 0xFDFF) || (r >= 0xFE70 && r <= 0xFEFF)
}

// HasArabic checks if a string contains any Arabic characters
func HasArabic(s string) bool {
	for _, r := range s {
		if IsArabic(r) || IsArabicPresentationForm(r) {
			return true
		}
	}
	return false
}

// IsNonConnecting checks if a letter doesn't connect to the following letter
func IsNonConnecting(r rune) bool {
	return nonConnectingLetters[r]
}

// GetLetterForm returns the appropriate form of an Arabic letter based on its position
func GetLetterForm(letter rune, formType FormType) rune {
	forms, ok := arabicForms[letter]
	if !ok {
		return letter
	}
	return forms[formType]
}

// Reshape transforms Arabic text by applying contextual letter forms
func Reshape(text string) string {
	runes := []rune(text)
	if len(runes) == 0 {
		return text
	}

	result := make([]rune, 0, len(runes))

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Skip non-Arabic characters
		if !IsArabic(r) {
			result = append(result, r)
			continue
		}

		// Check if this letter has forms defined
		forms, hasForms := arabicForms[r]
		if !hasForms {
			result = append(result, r)
			continue
		}

		// Determine connections
		prevConnects := false
		if i > 0 {
			prev := runes[i-1]
			prevConnects = IsArabic(prev) && !IsNonConnecting(prev)
		}

		nextConnects := false
		if i < len(runes)-1 {
			next := runes[i+1]
			nextConnects = IsArabic(next)
		}

		// Determine the appropriate form
		var form rune
		isNonConnecting := IsNonConnecting(r)

		if isNonConnecting {
			// Non-connecting letters only have isolated and final forms
			if prevConnects {
				form = forms[Final]
			} else {
				form = forms[Isolated]
			}
		} else {
			// Regular connecting letters
			if prevConnects && nextConnects {
				form = forms[Medial]
			} else if prevConnects {
				form = forms[Final]
			} else if nextConnects {
				form = forms[Initial]
			} else {
				form = forms[Isolated]
			}
		}

		result = append(result, form)
	}

	return string(result)
}

// Reverse reverses a string (for RTL display)
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// MirrorBrackets swaps brackets for RTL display
func MirrorBrackets(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if mirrored, ok := mirroredBrackets[r]; ok {
			runes[i] = mirrored
		}
	}
	return string(runes)
}

// segment represents a text segment with its direction
type segment struct {
	text  string
	isRTL bool
}

// segmentText splits text into RTL (Arabic) and LTR (Latin/numbers) segments
func segmentText(text string) []segment {
	if len(text) == 0 {
		return nil
	}

	var segments []segment
	var current strings.Builder
	var currentIsRTL bool
	first := true

	for _, r := range text {
		isArabic := IsArabic(r) || IsArabicPresentationForm(r)
		
		// Neutral characters (spaces, punctuation) inherit direction
		isNeutral := unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsDigit(r)
		
		if first {
			currentIsRTL = isArabic
			first = false
		}

		// If not neutral and direction changes, start new segment
		if !isNeutral && isArabic != currentIsRTL {
			if current.Len() > 0 {
				segments = append(segments, segment{text: current.String(), isRTL: currentIsRTL})
				current.Reset()
			}
			currentIsRTL = isArabic
		}

		current.WriteRune(r)
	}

	if current.Len() > 0 {
		segments = append(segments, segment{text: current.String(), isRTL: currentIsRTL})
	}

	return segments
}

// Process prepares text for RTL PDF rendering
// It handles mixed Arabic/English text, numbers, and brackets correctly
func Process(text string) string {
	if len(text) == 0 {
		return text
	}

	// If no Arabic, return as-is
	if !HasArabic(text) {
		return text
	}

	// Segment the text
	segments := segmentText(text)
	
	// Process each segment
	var result strings.Builder
	
	// For RTL layout, we need to reverse the order of segments
	// and handle each segment appropriately
	for i := len(segments) - 1; i >= 0; i-- {
		seg := segments[i]
		
		if seg.isRTL {
			// Arabic text: reshape and reverse
			reshaped := Reshape(seg.text)
			reversed := Reverse(reshaped)
			mirrored := MirrorBrackets(reversed)
			result.WriteString(mirrored)
		} else {
			// LTR text (English/numbers): keep as-is but mirror brackets
			result.WriteString(seg.text)
		}
	}

	return result.String()
}

// ProcessSimple is a simpler version that just reshapes and reverses
// Use this for pure Arabic text without mixed content
func ProcessSimple(text string) string {
	if !HasArabic(text) {
		return text
	}
	reshaped := Reshape(text)
	reversed := Reverse(reshaped)
	return MirrorBrackets(reversed)
}

// FormatNumber formats a number for display in Arabic context
// Numbers are kept LTR but the whole string is arranged for RTL layout
func FormatNumber(num string) string {
	return num // Numbers stay as-is in Unicode
}

// ProcessWithNumbers handles text with embedded numbers
// Arabic text is reshaped/reversed, numbers stay LTR
func ProcessWithNumbers(text string) string {
	return Process(text)
}
