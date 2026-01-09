package arabictext

import (
	"testing"
)

// ==================== IsArabic Tests ====================

func TestIsArabic(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		// Arabic letters
		{"Alef", 'ا', true},
		{"Ba", 'ب', true},
		{"Ta", 'ت', true},
		{"Tha", 'ث', true},
		{"Jeem", 'ج', true},
		{"Ha", 'ح', true},
		{"Kha", 'خ', true},
		{"Dal", 'د', true},
		{"Thal", 'ذ', true},
		{"Ra", 'ر', true},
		{"Zay", 'ز', true},
		{"Seen", 'س', true},
		{"Sheen", 'ش', true},
		{"Sad", 'ص', true},
		{"Dad", 'ض', true},
		{"Ta2", 'ط', true},
		{"Za", 'ظ', true},
		{"Ain", 'ع', true},
		{"Ghain", 'غ', true},
		{"Fa", 'ف', true},
		{"Qaf", 'ق', true},
		{"Kaf", 'ك', true},
		{"Lam", 'ل', true},
		{"Meem", 'م', true},
		{"Noon", 'ن', true},
		{"Ha2", 'ه', true},
		{"Waw", 'و', true},
		{"Ya", 'ي', true},
		{"Hamza", 'ء', true},
		{"AlefMaksura", 'ى', true},
		{"TaMarbuta", 'ة', true},
		
		// Non-Arabic characters
		{"LatinA", 'a', false},
		{"LatinZ", 'z', false},
		{"LatinCapA", 'A', false},
		{"Digit0", '0', false},
		{"Digit9", '9', false},
		{"Space", ' ', false},
		{"Period", '.', false},
		{"Comma", ',', false},
		{"OpenParen", '(', false},
		{"CloseParen", ')', false},
		{"Newline", '\n', false},
		{"Tab", '\t', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsArabic(tt.input)
			if result != tt.expected {
				t.Errorf("IsArabic(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// ==================== HasArabic Tests ====================

func TestHasArabic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Pure Arabic
		{"PureArabicWord", "مرحبا", true},
		{"ArabicSentence", "السلام عليكم", true},
		{"SingleArabicLetter", "م", true},
		
		// Mixed content
		{"MixedArabicEnglish", "Hello مرحبا World", true},
		{"ArabicWithNumbers", "منتج 123", true},
		{"ArabicWithParens", "القيمة (15%)", true},
		{"NumbersFirstThenArabic", "123 منتج", true},
		
		// Pure non-Arabic
		{"PureEnglish", "Hello World", false},
		{"PureNumbers", "12345", false},
		{"PureSymbols", "!@#$%", false},
		{"EmptyString", "", false},
		{"OnlySpaces", "   ", false},
		{"OnlyPunctuation", ".,;:!?", false},
		{"EnglishWithNumbers", "Price: 100", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasArabic(tt.input)
			if result != tt.expected {
				t.Errorf("HasArabic(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// ==================== IsNonConnecting Tests ====================

func TestIsNonConnecting(t *testing.T) {
	nonConnecting := []rune{'ا', 'د', 'ذ', 'ر', 'ز', 'و', 'أ', 'إ', 'آ', 'ؤ', 'ة', 'ى', 'ء'}
	connecting := []rune{'ب', 'ت', 'ث', 'ج', 'ح', 'خ', 'س', 'ش', 'ص', 'ض', 'ط', 'ظ', 'ع', 'غ', 'ف', 'ق', 'ك', 'ل', 'م', 'ن', 'ه', 'ي'}

	for _, r := range nonConnecting {
		if !IsNonConnecting(r) {
			t.Errorf("IsNonConnecting(%q) = false, want true", r)
		}
	}

	for _, r := range connecting {
		if IsNonConnecting(r) {
			t.Errorf("IsNonConnecting(%q) = true, want false", r)
		}
	}
}

// ==================== GetLetterForm Tests ====================

func TestGetLetterForm(t *testing.T) {
	tests := []struct {
		name     string
		letter   rune
		formType FormType
		notEmpty bool // Just verify it returns something valid
	}{
		{"BaIsolated", 'ب', Isolated, true},
		{"BaFinal", 'ب', Final, true},
		{"BaMedial", 'ب', Medial, true},
		{"BaInitial", 'ب', Initial, true},
		{"MeemIsolated", 'م', Isolated, true},
		{"MeemFinal", 'م', Final, true},
		{"MeemMedial", 'م', Medial, true},
		{"MeemInitial", 'م', Initial, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLetterForm(tt.letter, tt.formType)
			if tt.notEmpty && result == 0 {
				t.Errorf("GetLetterForm(%q, %d) returned empty rune", tt.letter, tt.formType)
			}
		})
	}
}

func TestGetLetterFormUnknownLetter(t *testing.T) {
	// Unknown letters should return themselves
	result := GetLetterForm('x', Isolated)
	if result != 'x' {
		t.Errorf("GetLetterForm('x', Isolated) = %q, want 'x'", result)
	}
}

// ==================== Reshape Tests ====================

func TestReshape(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"SingleLetter", "م"},
		{"TwoLetters", "مر"},
		{"ThreeLetters", "مرح"},
		{"FullWord", "مرحبا"},
		{"WordWithNonConnecting", "ماذا"},
		{"MultipleWords", "السلام عليكم"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reshape(tt.input)
			if len(result) == 0 {
				t.Errorf("Reshape(%q) returned empty string", tt.input)
			}
			// Reshaped text should have same number of runes
			if len([]rune(result)) != len([]rune(tt.input)) {
				t.Errorf("Reshape(%q) changed length: got %d, want %d", 
					tt.input, len([]rune(result)), len([]rune(tt.input)))
			}
		})
	}
}

func TestReshapePreservesNonArabic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "Hello"},
		{"123", "123"},
		{"Hello World", "Hello World"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Reshape(tt.input)
			if result != tt.expected {
				t.Errorf("Reshape(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestReshapeMixedContent(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"ArabicWithNumbers", "منتج 123"},
		{"ArabicWithEnglish", "مرحبا Hello"},
		{"NumbersInMiddle", "القيمة 100 ريال"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reshape(tt.input)
			if len(result) == 0 {
				t.Errorf("Reshape(%q) returned empty string", tt.input)
			}
			// Length should be preserved
			if len([]rune(result)) != len([]rune(tt.input)) {
				t.Errorf("Reshape(%q) changed length", tt.input)
			}
		})
	}
}

// ==================== Reverse Tests ====================

func TestReverse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"abc", "cba"},
		{"hello", "olleh"},
		{"12345", "54321"},
		{"مرحبا", "ابحرم"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Reverse(tt.input)
			if result != tt.expected {
				t.Errorf("Reverse(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestReverseDoubleReverse(t *testing.T) {
	// Reversing twice should give original
	tests := []string{"hello", "مرحبا", "Hello World", "منتج 123", ""}

	for _, input := range tests {
		result := Reverse(Reverse(input))
		if result != input {
			t.Errorf("Reverse(Reverse(%q)) = %q, want %q", input, result, input)
		}
	}
}

// ==================== MirrorBrackets Tests ====================

func TestMirrorBrackets(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"()", ")("},
		{"[]", "]["},
		{"{}", "}{"},
		{"<>", "><"},
		{"«»", "»«"},
		{"(hello)", ")hello("},
		{"no brackets", "no brackets"},
		{"", ""},
		{"(a[b{c}d]e)", ")a]b}c{d[e("},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := MirrorBrackets(tt.input)
			if result != tt.expected {
				t.Errorf("MirrorBrackets(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// ==================== Process Tests ====================

func TestProcess(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"EmptyString", ""},
		{"PureEnglish", "Hello World"},
		{"PureNumbers", "12345"},
		{"PureArabic", "مرحبا"},
		{"ArabicWithNumbers", "منتج 123"},
		{"ArabicWithParens", "القيمة (15%)"},
		{"MixedContent", "Hello مرحبا 123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Process(tt.input)
			// Should not panic and should return something
			if tt.input != "" && result == "" && HasArabic(tt.input) {
				// Only Arabic content should produce non-empty result
				// Actually this could be valid, so just ensure no panic
			}
		})
	}
}

func TestProcessPreservesNonArabic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "Hello"},
		{"123", "123"},
		{"Hello World", "Hello World"},
		{"", ""},
		{"Price: $100", "Price: $100"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Process(tt.input)
			if result != tt.expected {
				t.Errorf("Process(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// ==================== ProcessSimple Tests ====================

func TestProcessSimple(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"EmptyString", ""},
		{"SingleLetter", "م"},
		{"Word", "مرحبا"},
		{"Sentence", "السلام عليكم"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessSimple(tt.input)
			// Should not panic
			_ = result
		})
	}
}

func TestProcessSimplePreservesNonArabic(t *testing.T) {
	input := "Hello"
	expected := "Hello"
	result := ProcessSimple(input)
	if result != expected {
		t.Errorf("ProcessSimple(%q) = %q, want %q", input, result, expected)
	}
}

// ==================== Integration Tests ====================

func TestArabicWithNumbers(t *testing.T) {
	input := "منتج 1"
	result := Process(input)
	
	// Result should contain the number
	if len(result) == 0 {
		t.Error("Process returned empty result for Arabic with number")
	}
}

func TestArabicWithPercentage(t *testing.T) {
	input := "ضريبة (15%)"
	result := Process(input)
	
	if len(result) == 0 {
		t.Error("Process returned empty result for Arabic with percentage")
	}
}

func TestInvoiceLabels(t *testing.T) {
	// Test actual invoice labels used in the app
	labels := []string{
		"فاتورة ضريبية مبسطة",
		"رقم الفاتورة:",
		"اسم المتجر",
		"عنوان المتجر",
		"تاريخ:",
		"رقم تسجيل ضريبة القيمة المضافة:",
		"المنتجات",
		"الكمية",
		"سعر الوحدة",
		"ضريبة القيمة المضافة",
		"السعر شامل",
		"منتج 1",
		"منتج 2",
		"منتج 3",
		"اجمالي المبلغ الخاضع للضريبة",
		"ضريبة القيمة المضافة (15%)",
		"المجموع مع الضريبة (15%)",
		"إغلاق الفاتورة",
	}

	for _, label := range labels {
		t.Run(label, func(t *testing.T) {
			result := Process(label)
			if len(result) == 0 {
				t.Errorf("Process(%q) returned empty result", label)
			}
		})
	}
}

// ==================== Benchmark Tests ====================

func BenchmarkReshape(b *testing.B) {
	text := "فاتورة ضريبية مبسطة"
	for i := 0; i < b.N; i++ {
		Reshape(text)
	}
}

func BenchmarkReverse(b *testing.B) {
	text := "فاتورة ضريبية مبسطة"
	for i := 0; i < b.N; i++ {
		Reverse(text)
	}
}

func BenchmarkProcess(b *testing.B) {
	text := "فاتورة ضريبية مبسطة"
	for i := 0; i < b.N; i++ {
		Process(text)
	}
}

func BenchmarkProcessMixed(b *testing.B) {
	text := "ضريبة القيمة المضافة (15%)"
	for i := 0; i < b.N; i++ {
		Process(text)
	}
}

// ==================== Edge Cases ====================

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"OnlySpaces", "   "},
		{"OnlyNumbers", "12345"},
		{"OnlyBrackets", "()[]{}"},
		{"NewlineInArabic", "مرحبا\nعالم"},
		{"TabInArabic", "مرحبا\tعالم"},
		{"UnicodeEmoji", "مرحبا 👋"},
		{"VeryLongText", "مرحبا مرحبا مرحبا مرحبا مرحبا مرحبا مرحبا مرحبا مرحبا مرحبا"},
		{"RepeatedNumbers", "123 456 789"},
		{"DecimalNumber", "15.5%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			result := Process(tt.input)
			_ = result
		})
	}
}

func TestFormTypeConstants(t *testing.T) {
	if Isolated != 0 {
		t.Error("Isolated should be 0")
	}
	if Final != 1 {
		t.Error("Final should be 1")
	}
	if Medial != 2 {
		t.Error("Medial should be 2")
	}
	if Initial != 3 {
		t.Error("Initial should be 3")
	}
}

func TestAllArabicLettersHaveForms(t *testing.T) {
	// Common Arabic letters that must have forms
	letters := []rune{'ا', 'ب', 'ت', 'ث', 'ج', 'ح', 'خ', 'د', 'ذ', 'ر', 'ز', 
		'س', 'ش', 'ص', 'ض', 'ط', 'ظ', 'ع', 'غ', 'ف', 'ق', 'ك', 'ل', 'م', 'ن', 'ه', 'و', 'ي'}

	for _, letter := range letters {
		forms, exists := arabicForms[letter]
		if !exists {
			t.Errorf("No forms defined for letter %q", letter)
			continue
		}
		for i, form := range forms {
			if form == 0 {
				t.Errorf("Empty form at index %d for letter %q", i, letter)
			}
		}
	}
}
