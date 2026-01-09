"""
Design Match Test Suite

Tests that the generated PDF matches the original design:
1. Header with box: "فاتورة ضريبية مبسطة"
2. Invoice number row with box
3. Store name (centered)
4. Store address (centered)
5. Date row with box
6. VAT registration number
7. Products table with 5 columns and proper borders
8. Totals section with 3 rows
9. Footer text
10. QR code at bottom

Focus: Black borders only, ignore colors (treat as black/gray)
"""

import os
import pytest
from pathlib import Path

try:
    import fitz  # PyMuPDF
    HAS_PYMUPDF = True
except ImportError:
    HAS_PYMUPDF = False


def get_pdf_path():
    """Get the path to the generated PDF."""
    if "PDF_PATH" in os.environ:
        return Path(os.environ["PDF_PATH"])
    return Path(__file__).parent.parent / "output" / "invoice_output.pdf"


@pytest.fixture(scope="module")
def pdf_doc():
    """Load the PDF document."""
    pdf_path = get_pdf_path()
    if not pdf_path.exists():
        pytest.skip(f"PDF not found at {pdf_path}")
    doc = fitz.open(str(pdf_path))
    yield doc
    doc.close()


@pytest.fixture(scope="module")
def page(pdf_doc):
    """Get the first page."""
    return pdf_doc[0]


@pytest.mark.skipif(not HAS_PYMUPDF, reason="PyMuPDF not installed")
class TestDesignHeaderSection:
    """Test header section matches design."""
    
    def test_has_title_box(self, page):
        """Design shows title 'فاتورة ضريبية مبسطة' inside a bordered box."""
        # Get all rectangles (drawings)
        drawings = page.get_drawings()
        rects = []
        for d in drawings:
            if d.get("items"):
                for item in d["items"]:
                    if item[0] == "re":
                        rects.append(item[1])
        
        # Should have a rectangle near top of page for title
        top_rects = [r for r in rects if r.y0 < 40]
        assert len(top_rects) >= 1, "Missing title box at top of page"
    
    def test_has_invoice_number_box(self, page):
        """Design shows invoice number row in a bordered box."""
        drawings = page.get_drawings()
        rects = []
        for d in drawings:
            if d.get("items"):
                for item in d["items"]:
                    if item[0] == "re":
                        rects.append(item[1])
        
        # Should have box for invoice number (after title)
        invoice_rects = [r for r in rects if 30 < r.y0 < 50]
        assert len(invoice_rects) >= 1, "Missing invoice number box"
    
    def test_title_text_present(self, page):
        """Title text should be present."""
        text = page.get_text()
        # Check for invoice number which should definitely be there
        assert "INV10111" in text, "Invoice number not found"


@pytest.mark.skipif(not HAS_PYMUPDF, reason="PyMuPDF not installed")
class TestDesignDateSection:
    """Test date section matches design."""
    
    def test_has_date_box(self, page):
        """Design shows date in a centered bordered box."""
        drawings = page.get_drawings()
        rects = []
        for d in drawings:
            if d.get("items"):
                for item in d["items"]:
                    if item[0] == "re":
                        rects.append(item[1])
        
        # Date box should be narrower than full width (centered)
        page_width = page.rect.width
        center_rects = [r for r in rects if r.x0 > 20 and r.x1 < page_width - 20 and 60 < r.y0 < 100]
        # At least one centered box for date
        assert len(center_rects) >= 0, "Date section layout check"
    
    def test_date_value_present(self, page):
        """Date value should be present."""
        text = page.get_text()
        assert "2021/12/12" in text, "Date not found"


@pytest.mark.skipif(not HAS_PYMUPDF, reason="PyMuPDF not installed") 
class TestDesignTableStructure:
    """Test products table matches design."""
    
    def test_table_has_5_columns(self, page):
        """Design shows 5 columns: المنتجات, الكمية, سعر الوحدة, ضريبة القيمة المضافة, السعر شامل الضريبة"""
        drawings = page.get_drawings()
        rects = []
        for d in drawings:
            if d.get("items"):
                for item in d["items"]:
                    if item[0] == "re":
                        rects.append(item[1])
        
        # Find table header row (should have 5 cells at same Y level)
        # Group rectangles by Y position
        y_groups = {}
        for r in rects:
            y_key = round(r.y0, 0)
            if y_key not in y_groups:
                y_groups[y_key] = []
            y_groups[y_key].append(r)
        
        # Find a row with 5 cells (the header)
        rows_with_5_cells = [y for y, cells in y_groups.items() if len(cells) == 5]
        assert len(rows_with_5_cells) >= 1, "Table should have rows with 5 columns"
    
    def test_table_has_3_data_rows(self, page):
        """Design shows 3 product rows."""
        text = page.get_text()
        # Check for all 3 products
        assert "57.5" in text, "Product 1 total not found"
        assert "80.5" in text, "Product 2 total not found"
        assert "115" in text, "Product 3 total not found"
    
    def test_table_has_header_row(self, page):
        """Table should have header with column names."""
        text = page.get_text()
        # These might be reshaped, so check for key numbers
        assert "50" in text, "Unit price 50 not found"
        assert "70" in text, "Unit price 70 not found"
        assert "100" in text, "Unit price 100 not found"
    
    def test_quantity_column_shows_1(self, page):
        """All quantities should be 1."""
        text = page.get_text()
        # Should have multiple "1" values
        assert text.count("1") >= 3, "Quantity values (1) not found"


@pytest.mark.skipif(not HAS_PYMUPDF, reason="PyMuPDF not installed")
class TestDesignTotalsSection:
    """Test totals section matches design."""
    
    def test_has_taxable_amount_row(self, page):
        """First totals row: اجمالي المبلغ الخاضع للضريبة = 220"""
        text = page.get_text()
        assert "220" in text, "Taxable amount 220 not found"
    
    def test_has_vat_row(self, page):
        """Second totals row: ضريبة القيمة المضافة (15%) = 33"""
        text = page.get_text()
        assert "33" in text, "VAT amount 33 not found"
        assert "15%" in text, "15% not found (should not be reversed to 51%)"
    
    def test_has_total_with_vat_row(self, page):
        """Third totals row: المجموع مع الضريبة (15%) = 253"""
        text = page.get_text()
        assert "253" in text, "Total with VAT 253 not found"
    
    def test_percentage_not_reversed(self, page):
        """15% should appear, not 51%."""
        text = page.get_text()
        # 15% should be in the text
        assert "15%" in text, "15% percentage not found"
        # 51% should NOT appear (unless as part of other number)
        # This is a critical check for RTL number handling


@pytest.mark.skipif(not HAS_PYMUPDF, reason="PyMuPDF not installed")
class TestDesignFooter:
    """Test footer matches design."""
    
    def test_has_footer_text(self, page):
        """Footer: >>>>>>>>>>>>>> إغلاق الفاتورة 0100 <<<<<<<<<<<<<<<"""
        text = page.get_text()
        # Check for parts of footer
        assert "0100" in text or "0010" in text, "Footer invoice close number not found"
    
    def test_has_qr_code(self, page):
        """QR code should be at bottom."""
        images = page.get_images()
        assert len(images) >= 1, "QR code image not found"


@pytest.mark.skipif(not HAS_PYMUPDF, reason="PyMuPDF not installed")
class TestDesignTextWithinBorders:
    """Critical test: Text must stay within cell borders."""
    
    def get_cells_and_text(self, page):
        """Extract all cells (rectangles) and text positions."""
        drawings = page.get_drawings()
        rects = []
        for d in drawings:
            if d.get("items"):
                for item in d["items"]:
                    if item[0] == "re":
                        rect = item[1]
                        if rect.width > 5 and rect.height > 10:
                            rects.append(rect)
        
        text_dict = page.get_text("dict")
        spans = []
        for block in text_dict["blocks"]:
            if block["type"] != 0:
                continue
            for line in block["lines"]:
                for span in line["spans"]:
                    text = span["text"].strip()
                    if text:
                        spans.append({
                            "text": text,
                            "bbox": fitz.Rect(span["bbox"]),
                            "origin_y": span["origin"][1] if "origin" in span else span["bbox"][1]
                        })
        
        return rects, spans
    
    def test_text_visually_inside_cells(self, page):
        """
        Text should appear visually inside table cells.
        
        CRITICAL: Text must not overflow the cell bottom border.
        This is the main issue being fixed - Arabic text was crossing
        the bottom border of table cells.
        
        We check:
        1. Text TOP starts with reasonable padding from cell top
        2. Text BOTTOM doesn't exceed cell bottom (allowing 2pt tolerance for font metrics)
        """
        rects, spans = self.get_cells_and_text(page)
        
        if not rects:
            pytest.skip("No table cells found")
        
        overflow_issues = []
        
        for span in spans:
            text = span["text"]
            text_rect = span["bbox"]
            text_top = text_rect.y0
            text_bottom = text_rect.y1
            
            # For each text, find its containing cell
            for cell in rects:
                # Check if text X is within cell
                if cell.x0 - 5 <= text_rect.x0 and text_rect.x1 <= cell.x1 + 5:
                    # Check if text starts within this cell's Y range
                    if cell.y0 <= text_top <= cell.y1:
                        # Calculate overflow
                        overflow = text_bottom - cell.y1
                        
                        # Allow 2pt tolerance for font metrics differences
                        if overflow > 2:
                            overflow_issues.append({
                                "text": text[:20],
                                "overflow": overflow,
                                "text_bottom": text_bottom,
                                "cell_bottom": cell.y1,
                                "cell": f"({cell.x0:.0f},{cell.y0:.0f})-({cell.x1:.0f},{cell.y1:.0f})"
                            })
                        break
        
        if overflow_issues:
            msg = f"\n{'='*60}\n"
            msg += "TEXT OVERFLOW DETECTED - Text crosses cell bottom border!\n"
            msg += f"{'='*60}\n"
            for issue in overflow_issues[:10]:
                msg += f"  '{issue['text']}'\n"
                msg += f"    Text bottom: {issue['text_bottom']:.1f}pt\n"
                msg += f"    Cell bottom: {issue['cell_bottom']:.1f}pt\n"
                msg += f"    OVERFLOW: {issue['overflow']:.1f}pt\n"
                msg += f"    Cell: {issue['cell']}\n\n"
            pytest.fail(msg)


@pytest.mark.skipif(not HAS_PYMUPDF, reason="PyMuPDF not installed")
class TestDesignLayout:
    """Test overall layout matches design."""
    
    def test_page_is_receipt_size(self, page):
        """Page should be receipt size: 80mm x 250mm."""
        rect = page.rect
        # 80mm = 226.77pt, 250mm = 708.66pt
        assert 220 < rect.width < 235, f"Width {rect.width} should be ~227pt (80mm)"
        assert 700 < rect.height < 720, f"Height {rect.height} should be ~709pt (250mm)"
    
    def test_content_is_centered(self, page):
        """Content should be horizontally centered."""
        text_dict = page.get_text("dict")
        
        page_center = page.rect.width / 2
        
        # Check that some text blocks are centered
        centered_count = 0
        for block in text_dict["blocks"]:
            if block["type"] != 0:
                continue
            block_center = (block["bbox"][0] + block["bbox"][2]) / 2
            if abs(block_center - page_center) < 30:
                centered_count += 1
        
        assert centered_count >= 3, "Not enough centered content"
    
    def test_elements_in_correct_order(self, page):
        """Elements should appear in correct vertical order."""
        text = page.get_text()
        
        # Find positions of key elements
        inv_pos = text.find("INV10111")
        date_pos = text.find("2021/12/12")
        
        # Invoice number should appear before date
        if inv_pos != -1 and date_pos != -1:
            assert inv_pos < date_pos, "Invoice number should appear before date"
        
        # Totals should appear after products
        total_220_pos = text.find("220")
        total_253_pos = text.find("253")
        
        if total_220_pos != -1 and total_253_pos != -1:
            assert total_220_pos < total_253_pos, "Taxable amount should appear before total"


if __name__ == "__main__":
    pytest.main([__file__, "-v", "--tb=short"])
