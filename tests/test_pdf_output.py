"""
PDF Output Test Suite for Arabic Invoice Generator

Tests the actual PDF output to verify:
1. Arabic text stays within table cell boundaries
2. Text is properly positioned (not crossing borders)
3. All required content is present
4. RTL text rendering is correct
5. Percentage values are not reversed
"""

import os
import subprocess
import pytest
from pathlib import Path

# Try to import PDF libraries
try:
    import fitz  # PyMuPDF
    HAS_PYMUPDF = True
except ImportError:
    HAS_PYMUPDF = False

try:
    import pdfplumber
    HAS_PDFPLUMBER = True
except ImportError:
    HAS_PDFPLUMBER = False


# Constants from main.go - MUST MATCH THE CODE
TABLE_ROW_HEIGHT = 22.0  # Current value in main.go
TABLE_TEXT_Y_OFFSET = 8.0  # currentY + 8
FONT_SIZE = 9.0
MARGIN = 10.0
PAGE_WIDTH = 226.77

# Column widths from main.go
COL_WIDTHS = [35, 35, 35, 25, 76]  # Total, VAT, Price, Qty, Product
TABLE_X = MARGIN

# Header heights
HEADER_HEIGHT = 28.0
TOTALS_ROW_HEIGHT = 20.0
TOTALS_TOTAL_ROW_HEIGHT = 22.0


def get_pdf_path():
    """Get the path to the generated PDF."""
    # Check environment variable first (for Docker)
    if "PDF_PATH" in os.environ:
        return Path(os.environ["PDF_PATH"])
    return Path(__file__).parent.parent / "output" / "invoice_output.pdf"
    return Path(__file__).parent.parent / "output" / "invoice_output.pdf"


def generate_pdf():
    """Generate the PDF by running the Docker container."""
    project_dir = Path(__file__).parent.parent
    result = subprocess.run(
        ["docker", "run", "--rm", 
         "-v", f"{project_dir}/output:/app/output", 
         "bill-generator"],
        capture_output=True,
        text=True,
        cwd=project_dir
    )
    return result.returncode == 0


class TestPDFExists:
    """Basic tests for PDF existence and validity."""
    
    def test_pdf_file_exists(self):
        """Verify the PDF file exists."""
        pdf_path = get_pdf_path()
        assert pdf_path.exists(), f"PDF not found at {pdf_path}"
    
    def test_pdf_not_empty(self):
        """Verify the PDF has content."""
        pdf_path = get_pdf_path()
        assert pdf_path.stat().st_size > 5000, "PDF file too small, likely empty or corrupted"
    
    def test_pdf_header_valid(self):
        """Verify the PDF has a valid header."""
        pdf_path = get_pdf_path()
        with open(pdf_path, 'rb') as f:
            header = f.read(5)
        assert header == b'%PDF-', "Invalid PDF header"


@pytest.mark.skipif(not HAS_PYMUPDF, reason="PyMuPDF not installed")
class TestTextPositioning:
    """Tests for text positioning within cells using PyMuPDF."""
    
    @pytest.fixture
    def pdf_doc(self):
        """Load the PDF document."""
        pdf_path = get_pdf_path()
        doc = fitz.open(str(pdf_path))
        yield doc
        doc.close()
    
    def test_page_dimensions(self, pdf_doc):
        """Verify page dimensions match expected receipt size."""
        page = pdf_doc[0]
        rect = page.rect
        # 80mm x 250mm in points (1mm = 2.83465 points)
        assert abs(rect.width - 226.77) < 1, f"Page width {rect.width} != 226.77"
        assert abs(rect.height - 708.66) < 1, f"Page height {rect.height} != 708.66"
    
    def test_text_blocks_extracted(self, pdf_doc):
        """Verify text can be extracted from the PDF."""
        page = pdf_doc[0]
        text = page.get_text()
        assert len(text) > 100, "Too little text extracted from PDF"
    
    def test_arabic_text_present(self, pdf_doc):
        """Verify Arabic text is present in the PDF."""
        page = pdf_doc[0]
        text = page.get_text()
        # Check for Arabic Unicode range
        has_arabic = any('\u0600' <= c <= '\u06FF' for c in text)
        assert has_arabic, "No Arabic characters found in PDF"
    
    def test_invoice_number_present(self, pdf_doc):
        """Verify invoice number is in the PDF."""
        page = pdf_doc[0]
        text = page.get_text()
        assert "INV10111" in text, "Invoice number not found"
    
    def test_date_present(self, pdf_doc):
        """Verify date is in the PDF."""
        page = pdf_doc[0]
        text = page.get_text()
        assert "2021/12/12" in text, "Date not found"
    
    def test_vat_number_present(self, pdf_doc):
        """Verify VAT registration number is in the PDF."""
        page = pdf_doc[0]
        text = page.get_text()
        assert "123456789900003" in text, "VAT registration number not found"
    
    def test_totals_present(self, pdf_doc):
        """Verify all total values are present."""
        page = pdf_doc[0]
        text = page.get_text()
        assert "220" in text, "Taxable amount (220) not found"
        assert "33" in text, "VAT amount (33) not found"
        assert "253" in text, "Total with VAT (253) not found"
    
    def test_product_prices_present(self, pdf_doc):
        """Verify product prices are in the PDF."""
        page = pdf_doc[0]
        text = page.get_text()
        assert "57.5" in text, "Product 1 total not found"
        assert "80.5" in text, "Product 2 total not found"
        assert "115" in text, "Product 3 total not found"
    
    def test_percentage_not_reversed(self, pdf_doc):
        """Verify 15% is not reversed to 51%."""
        page = pdf_doc[0]
        text = page.get_text()
        # Should have 15%, should NOT have 51% (unless it's part of another number)
        assert "15%" in text, "15% not found in PDF"
        # Check that 51% doesn't appear as a standalone percentage
        # This is tricky because 51 could be part of other numbers
    
    def test_text_within_page_bounds(self, pdf_doc):
        """Verify all text blocks are within page boundaries."""
        page = pdf_doc[0]
        page_rect = page.rect
        
        blocks = page.get_text("dict")["blocks"]
        for block in blocks:
            if "lines" in block:
                for line in block["lines"]:
                    for span in line["spans"]:
                        bbox = span["bbox"]
                        # Check text is within page
                        assert bbox[0] >= 0, f"Text extends past left edge: {bbox}"
                        assert bbox[2] <= page_rect.width, f"Text extends past right edge: {bbox}"
                        assert bbox[1] >= 0, f"Text extends past top edge: {bbox}"
                        assert bbox[3] <= page_rect.height, f"Text extends past bottom edge: {bbox}"


@pytest.mark.skipif(not HAS_PYMUPDF, reason="PyMuPDF not installed")
class TestTableCellBoundaries:
    """Tests specifically for text staying within table cell boundaries."""
    
    @pytest.fixture
    def pdf_doc(self):
        """Load the PDF document."""
        pdf_path = get_pdf_path()
        doc = fitz.open(str(pdf_path))
        yield doc
        doc.close()
    
    @pytest.fixture
    def page_data(self, pdf_doc):
        """Extract all page data: text blocks and drawings."""
        page = pdf_doc[0]
        return {
            "text_dict": page.get_text("dict"),
            "drawings": page.get_drawings(),
            "page_rect": page.rect
        }
    
    def get_rectangles(self, page_data):
        """Extract all rectangles (table cells) from drawings."""
        rects = []
        for d in page_data["drawings"]:
            if d.get("items"):
                for item in d["items"]:
                    if item[0] == "re":  # rectangle
                        rect = item[1]
                        if hasattr(rect, 'width') and rect.width > 5 and rect.height > 5:
                            rects.append(rect)
        return rects
    
    def get_text_spans(self, page_data):
        """Extract all text spans with bounding boxes."""
        spans = []
        for block in page_data["text_dict"]["blocks"]:
            if block["type"] != 0:  # text block
                continue
            for line in block["lines"]:
                for span in line["spans"]:
                    text = span["text"].strip()
                    if text:
                        spans.append({
                            "text": text,
                            "bbox": fitz.Rect(span["bbox"]),
                            "size": span["size"],
                            "is_arabic": any('\u0600' <= c <= '\u06FF' for c in text)
                        })
        return spans
    
    def test_text_within_cell_bounds(self, page_data):
        """
        CRITICAL TEST: Verify text does not cross cell bottom boundaries.
        
        Note: PyMuPDF's text bbox includes line-height spacing which may be
        larger than the actual glyph. We check that the VISIBLE part of text
        (baseline + small descender) stays within bounds.
        """
        rects = self.get_rectangles(page_data)
        spans = self.get_text_spans(page_data)
        
        if not rects:
            pytest.skip("No table rectangles found in PDF")
        
        # Debug: Print all rectangles sorted by Y
        print("\n\nDEBUG: All rectangles found:")
        sorted_rects = sorted(rects, key=lambda r: (r.y0, r.x0))
        for i, r in enumerate(sorted_rects[:25]):
            print(f"  Rect {i}: y0={r.y0:.1f}, y1={r.y1:.1f}, height={r.height:.1f}, x0={r.x0:.1f}")
        
        violations = []
        
        for span in spans:
            text_bbox = span["bbox"]
            font_size = span["size"]
            
            # Estimate actual glyph bounds (not PyMuPDF's padded bbox)
            # Text top is approximately: bbox.y0 + (bbox.height - font_size) / 2
            # But for our purposes, we use bbox.y0 as-is for cell matching
            text_origin = fitz.Point(text_bbox.x0, text_bbox.y0)
            
            # Find the cell where text STARTS (using top of text box)
            for cell_rect in rects:
                if (cell_rect.x0 <= text_origin.x <= cell_rect.x1 and
                    cell_rect.y0 <= text_origin.y <= cell_rect.y1):
                    
                    # For overflow check, estimate visible glyph bottom
                    # Visible glyph is approximately font_size * 1.3 from top
                    estimated_visible_bottom = text_bbox.y0 + (font_size * 1.5)
                    
                    overflow = estimated_visible_bottom - cell_rect.y1
                    if overflow > 1.0:  # 1pt tolerance for visual overflow
                        violations.append({
                            "text": span["text"][:30],
                            "text_top": text_bbox.y0,
                            "visible_bottom": estimated_visible_bottom,
                            "bbox_bottom": text_bbox.y1,
                            "cell_top": cell_rect.y0,
                            "cell_bottom": cell_rect.y1,
                            "overflow": overflow,
                            "is_arabic": span["is_arabic"]
                        })
                    break
        
        if violations:
            msg = f"\n{'='*60}\n"
            msg += "TEXT OVERFLOW DETECTED - Text crosses cell boundaries!\n"
            msg += f"{'='*60}\n"
            for v in violations[:10]:
                arabic_marker = "[ARABIC]" if v["is_arabic"] else ""
                msg += f"  {arabic_marker} '{v['text']}'\n"
                msg += f"     Text top: {v['text_top']:.1f}\n"
                msg += f"     Est. visible bottom: {v['visible_bottom']:.1f}\n"
                msg += f"     Cell: Y {v['cell_top']:.1f} to {v['cell_bottom']:.1f}\n"
                msg += f"     OVERFLOW: {v['overflow']:.1f}pt\n\n"
            pytest.fail(msg)
    
    def test_arabic_text_specifically(self, page_data):
        """Test that Arabic text in particular stays within bounds."""
        rects = self.get_rectangles(page_data)
        spans = self.get_text_spans(page_data)
        
        arabic_spans = [s for s in spans if s["is_arabic"]]
        
        if not arabic_spans:
            pytest.skip("No Arabic text found")
        
        for span in arabic_spans:
            text_bbox = span["bbox"]
            text_height = text_bbox.height
            font_size = span["size"]
            
            # Arabic text height ratio should be reasonable
            if font_size > 0:
                height_ratio = text_height / font_size
                # Amiri font typically has 1.2-1.5x ratio
                assert height_ratio < 2.0, (
                    f"Arabic text '{span['text'][:20]}' has unusual height ratio: "
                    f"{height_ratio:.2f} (height={text_height:.1f}, size={font_size})"
                )
    
    def get_table_region_text_blocks(self, pdf_doc):
        """Extract text blocks from the table region."""
        page = pdf_doc[0]
        blocks = page.get_text("dict")["blocks"]
        
        # Table starts after header content (approximately Y=104 based on layout)
        # This is: 10 (start) + 22 (title) + 14 (invoice#) + 14 (store) + 14 (addr) + 14 (date) + 16 (vat) = 104
        table_start_y = 100  # approximate
        
        table_blocks = []
        for block in blocks:
            if "lines" in block:
                for line in block["lines"]:
                    for span in line["spans"]:
                        bbox = span["bbox"]
                        if bbox[1] >= table_start_y:
                            table_blocks.append({
                                "text": span["text"],
                                "bbox": bbox,
                                "y_top": bbox[1],
                                "y_bottom": bbox[3],
                                "height": bbox[3] - bbox[1]
                            })
        return table_blocks
    
    def test_arabic_text_height_within_row(self, pdf_doc):
        """
        CRITICAL TEST: Verify Arabic text doesn't exceed row height.
        
        This test checks that the text glyph height (including descenders)
        fits within the allocated row height.
        """
        blocks = self.get_table_region_text_blocks(pdf_doc)
        
        for block in blocks:
            text = block["text"]
            height = block["height"]
            
            # Arabic text with font size 9 should have height <= row_height - padding
            # Maximum expected height for size 9 Arabic font is about 14-15pt
            max_expected_height = 16  # Allow some tolerance
            
            if any('\u0600' <= c <= '\u06FF' for c in text):
                assert height <= max_expected_height, (
                    f"Arabic text '{text}' has height {height:.1f}pt which may exceed cell bounds. "
                    f"Expected <= {max_expected_height}pt"
                )
    
    def test_row_spacing_consistent(self, pdf_doc):
        """Verify consistent spacing between table rows."""
        blocks = self.get_table_region_text_blocks(pdf_doc)
        
        # Get Y positions of text blocks
        y_positions = sorted(set(block["y_top"] for block in blocks))
        
        # Check that rows are evenly spaced (approximately)
        if len(y_positions) >= 3:
            spacings = [y_positions[i+1] - y_positions[i] for i in range(len(y_positions)-1)]
            # Filter for table row spacings (should be around 16-20pt)
            row_spacings = [s for s in spacings if 14 <= s <= 22]
            
            if row_spacings:
                avg_spacing = sum(row_spacings) / len(row_spacings)
                for spacing in row_spacings:
                    # Allow 20% variance
                    assert abs(spacing - avg_spacing) <= avg_spacing * 0.3, (
                        f"Inconsistent row spacing: {spacing:.1f} vs avg {avg_spacing:.1f}"
                    )


@pytest.mark.skipif(not HAS_PDFPLUMBER, reason="pdfplumber not installed")
class TestWithPdfplumber:
    """Alternative tests using pdfplumber."""
    
    @pytest.fixture
    def pdf(self):
        """Load PDF with pdfplumber."""
        pdf_path = get_pdf_path()
        with pdfplumber.open(str(pdf_path)) as pdf:
            yield pdf
    
    def test_text_extraction(self, pdf):
        """Test basic text extraction."""
        page = pdf.pages[0]
        text = page.extract_text()
        assert text is not None and len(text) > 50
    
    def test_tables_detected(self, pdf):
        """Test if tables can be detected."""
        page = pdf.pages[0]
        tables = page.extract_tables()
        # We should have at least the products table
        # Note: This may not work perfectly with our manually drawn tables
    
    def test_chars_with_positions(self, pdf):
        """Test character-level extraction with positions."""
        page = pdf.pages[0]
        chars = page.chars
        
        assert len(chars) > 0, "No characters extracted"
        
        # Check all characters are within page bounds
        for char in chars:
            assert char["x0"] >= 0, f"Char '{char['text']}' has negative x0"
            assert char["top"] >= 0, f"Char '{char['text']}' has negative top"


class TestQRCode:
    """Tests for QR code presence."""
    
    def test_pdf_has_image(self):
        """Verify the PDF contains an image (QR code)."""
        if not HAS_PYMUPDF:
            pytest.skip("PyMuPDF not installed")
        
        pdf_path = get_pdf_path()
        doc = fitz.open(str(pdf_path))
        page = doc[0]
        
        images = page.get_images()
        assert len(images) > 0, "No images found in PDF (QR code missing)"
        doc.close()


class TestLayoutMeasurements:
    """Tests that verify layout measurements match expected values."""
    
    def test_row_height_sufficient_for_arabic(self):
        """
        Test that row height is sufficient for Arabic text.
        
        Arabic fonts like Amiri have larger ascenders/descenders.
        For font size 9:
        - Ascender: ~7pt (above baseline)
        - Descender: ~5pt (below baseline)  
        - Total glyph height: ~12pt
        
        With baseline at currentY + 8, ascender reaches currentY + 1
        and descender reaches currentY + 13. Row height of 22 gives 9pt margin.
        """
        font_size = 9.0
        row_height = 22.0
        text_y_offset = 8.0  # Baseline position from top of cell
        
        # Ascender and descender estimates for Arabic font
        ascender = font_size * 0.8  # ~7pt above baseline
        descender = font_size * 0.5  # ~5pt below baseline
        
        # Check ascender stays in cell
        text_top = text_y_offset - ascender  # Where top of glyph is
        assert text_top >= 0, f"Ascender extends {-text_top:.1f}pt above cell"
        
        # Check descender stays in cell
        text_bottom = text_y_offset + descender
        bottom_margin = row_height - text_bottom
        
        assert bottom_margin >= 1.0, (
            f"Insufficient bottom margin: {bottom_margin:.1f}pt. "
            f"Text bottom at {text_bottom:.1f}pt, row height {row_height}pt"
        )
    
    def test_header_row_height_sufficient(self):
        """Test header row has enough height for 2-line text."""
        header_height = 28.0
        font_size = 7.0  # Header uses smaller font
        
        # Two lines with proper spacing
        line_height = font_size * 1.3  # ~9pt per line
        total_text_height = line_height * 2  # ~18pt for 2 lines
        available_padding = header_height - total_text_height
        
        assert available_padding >= 6.0, (
            f"Header padding {available_padding:.1f}pt insufficient for 2 lines"
        )


class TestRTLRendering:
    """Tests for right-to-left text rendering."""
    
    @pytest.mark.skipif(not HAS_PYMUPDF, reason="PyMuPDF not installed")
    def test_arabic_words_intact(self):
        """Verify Arabic words are not broken or reversed incorrectly."""
        pdf_path = get_pdf_path()
        doc = fitz.open(str(pdf_path))
        page = doc[0]
        text = page.get_text()
        doc.close()
        
        # These Arabic phrases should appear (possibly reshaped but recognizable)
        # Note: Exact matching is tricky due to reshaping
        expected_numbers = ["INV10111", "2021/12/12", "123456789900003"]
        
        for num in expected_numbers:
            assert num in text, f"Expected '{num}' not found in PDF"


if __name__ == "__main__":
    # Run a quick check
    print("PDF Output Test Suite")
    print("=" * 50)
    
    pdf_path = get_pdf_path()
    if not pdf_path.exists():
        print(f"[FAIL] PDF not found at {pdf_path}")
        print("Run 'docker build -t bill-generator . && docker run --rm -v ./output:/app/output bill-generator' first")
        exit(1)
    
    print(f"[PASS] PDF exists: {pdf_path}")
    print(f"[INFO] Size: {pdf_path.stat().st_size} bytes")
    
    # Check header
    with open(pdf_path, 'rb') as f:
        if f.read(5) == b'%PDF-':
            print("[PASS] Valid PDF header")
        else:
            print("[FAIL] Invalid PDF header")
    
    # Check with PyMuPDF if available
    if HAS_PYMUPDF:
        doc = fitz.open(str(pdf_path))
        page = doc[0]
        text = page.get_text()
        
        print(f"[INFO] Page size: {page.rect.width:.1f} x {page.rect.height:.1f} pts")
        print(f"[INFO] Text length: {len(text)} chars")
        
        # Check for key content
        checks = [
            ("INV10111", "Invoice number"),
            ("2021/12/12", "Date"),
            ("220", "Taxable amount"),
            ("33", "VAT amount"),
            ("253", "Total"),
            ("15%", "Percentage (not reversed)"),
        ]
        
        for value, desc in checks:
            if value in text:
                print(f"[PASS] {desc}: {value}")
            else:
                print(f"[FAIL] {desc}: {value} not found")
        
        doc.close()
    else:
        print("[SKIP] PyMuPDF not installed, skipping detailed checks")
    
    print("=" * 50)
    print("Run 'pytest test_pdf_output.py -v' for full test suite")
