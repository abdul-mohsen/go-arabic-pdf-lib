# Bill Generator Library

A Go library for generating PDF invoices with full Arabic RTL support, English LTR support, VAT calculations, and discount handling.

![Go Version](https://img.shields.io/badge/Go-1.21-blue)
![Docker](https://img.shields.io/badge/docker-ready-blue)

## Features

- 🌍 **Bilingual Support**: Arabic (RTL) and English (LTR)
- 💰 **VAT Calculations**: Automatic VAT computation
- 🏷️ **Discount Support**: Percentage and fixed amount discounts
- 📄 **PDF Generation**: Clean, print-ready thermal receipt style invoices
- 📱 **QR Code**: Embedded QR codes for e-invoicing compliance
- 🎨 **Customizable Labels**: Full i18n support for all text labels
- 🔧 **Component-Based Architecture**: Extensible design with reusable components

## Installation

```bash
go get github.com/abdul-mohsen/github.com/abdul-mohsen/go-arabic-pdf-lib
```

## Quick Start

### Option 1: Generate from JSON file

```go
package main

import (
    "log"
    "github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/invoice"
)

func main() {
    // Set font path (required for Arabic support)
    invoice.SetDefaultFontPath("/path/to/fonts")
    
    // Generate invoice from JSON file
    err := invoice.GenerateFromFile("invoice_data.json", "output/invoice.pdf")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Option 2: Use the Builder pattern

```go
package main

import (
    "log"
    "github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/invoice"
)

func main() {
    gen := invoice.NewGenerator(
        invoice.WithFontPath("/fonts"),
    )
    
    // Build invoice programmatically
    inv := invoice.NewBuilder().
        WithTitle("Tax Invoice").
        WithInvoiceNumber("INV-001").
        WithStoreName("My Store").
        WithStoreAddress("123 Main Street").
        WithDate("2024/01/15").
        WithVATRegistration("123456789").
        WithVATPercentage(15.0).
        WithQRCode("base64-encoded-qr-data").
        WithEnglish(). // Use English (LTR), omit for Arabic (RTL)
        WithLabels(invoice.DefaultEnglishLabels()).
        AddProduct("Product A", 2, 50.00).
        AddProductWithDiscount("Product B", 1, 100.00, 10, 0). // 10% discount
        AddProductWithDiscount("Product C", 1, 80.00, 0, 5.00). // $5 fixed discount
        Build()
    
    err := gen.Generate(inv, "output/invoice.pdf")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Option 3: Generate PDF as bytes (for web services)

```go
package main

import (
    "net/http"
    "github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/invoice"
)

func handler(w http.ResponseWriter, r *http.Request) {
    gen := invoice.NewGenerator(invoice.WithFontPath("/fonts"))
    
    inv := invoice.NewBuilder().
        WithTitle("Tax Invoice").
        WithInvoiceNumber("INV-001").
        // ... configure invoice
        Build()
    
    pdfBytes, err := gen.GenerateBytes(inv)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    
    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "attachment; filename=invoice.pdf")
    w.Write(pdfBytes)
}
```

## Run with Docker

```powershell
# Windows
.\run-docker.ps1

# Or using Docker directly
docker build -t bill-generator .
docker run --rm -v "${PWD}/output:/app/output" bill-generator
```

## JSON Schema

```json
{
  "config": {
    "vatPercentage": 15,
    "currencySymbol": "SAR",
    "dateFormat": "2006/01/02",
    "english": true
  },
  "invoice": {
    "title": "Tax Invoice",
    "invoiceNumber": "INV-001",
    "storeName": "Store Name",
    "storeAddress": "Store Address",
    "date": "2024/01/15",
    "vatRegistrationNo": "123456789",
    "qrCodeData": "base64-encoded-data"
  },
  "products": [
    {
      "name": "Product Name",
      "quantity": 2.0,
      "unitPrice": 50.00,
      "discountPercent": 10,
      "discountAmount": 0
    }
  ],
  "labels": {
    "invoiceNumber": "Invoice Number:",
    "date": "Date:",
    "vatRegistration": "VAT Registration:",
    "totalTaxable": "Total Taxable",
    "totalWithVat": "Total with VAT",
    "productColumn": "Product",
    "quantityColumn": "Qty",
    "unitPriceColumn": "Unit Price",
    "discountColumn": "Discount",
    "vatAmountColumn": "VAT",
    "totalColumn": "Total",
    "totalDiscount": "Total Discount",
    "footer": "Thank you for your business"
  }
}
```

## Discount Calculations

Discounts are applied before VAT calculation:

1. **Gross Amount** = Quantity × Unit Price
2. **Discount** = (Gross × Discount%) + Fixed Discount Amount
3. **Net Amount** = Gross - Discount
4. **VAT** = Net Amount × VAT Rate
5. **Total** = Net Amount + VAT

Example:
- Product: 100 SAR, 10% discount, 15% VAT
- Gross: 100.00
- Discount: 10.00 (10%)
- Net: 90.00
- VAT: 13.50 (15% of 90)
- Total: 103.50

## Language Support

### Arabic (Default - RTL)

```go
// Arabic is the default, no flag needed
inv := invoice.NewBuilder().
    WithTitle("فاتورة ضريبية مبسطة").
    WithLabels(invoice.DefaultArabicLabels()).
    // ...
    Build()
```

### English (LTR)

```go
// Set english flag
inv := invoice.NewBuilder().
    WithEnglish().
    WithLabels(invoice.DefaultEnglishLabels()).
    // ...
    Build()
```

Or in JSON:
```json
{
  "config": {
    "english": true
  }
}
```

## Component Architecture

The library uses a component-based architecture for extensibility:

```go
import "github.com/abdul-mohsen/go-arabic-pdf-lib/pkg/component"

// Components implement the Component interface
type Component interface {
    Draw(pdf *gopdf.GoPdf) float64
}

// Available components:
// - TextBlock: Simple text rendering
// - LabelValuePair: Label-value pairs
// - Header: Centered header text
// - WrappedText: Multi-line text with wrapping
// - Table: Data tables with columns
// - TotalsTable: Summary/totals tables
// - QRCode: QR code rendering

// Use functional options for configuration
text := component.NewTextBlock("Hello World",
    component.WithPosition(10, 10),
    component.WithFontSize(12),
    component.WithBold(true),
    component.WithAlignment(component.AlignCenter),
)
```

## Project Structure

```
go-arabic-pdf-lib/
├── pkg/
│   ├── invoice/      # Main library API
│   ├── component/    # Reusable PDF components
│   ├── models/       # Data structures
│   ├── loader/       # JSON loading and parsing
│   ├── pdf/          # PDF generation
│   └── textutil/     # Text processing utilities
├── arabictext/       # Arabic text reshaping
├── cmd/
│   └── generator/    # CLI tool
├── invoice_data.json     # Arabic sample
├── invoice_data_en.json  # English sample
└── output/               # Generated PDFs
```

## API Reference

### Generator

```go
// Create a new generator
gen := invoice.NewGenerator(
    invoice.WithFontPath("/fonts"),
)

// Generate to file
err := gen.GenerateFromFile("input.json", "output.pdf")
err := gen.GenerateFromJSON(jsonData, "output.pdf")
err := gen.Generate(inv, "output.pdf")

// Generate to bytes
pdfBytes, err := gen.GenerateBytes(inv)
```

### Builder

```go
builder := invoice.NewBuilder()

// Configuration
builder.WithTitle(title string)
builder.WithInvoiceNumber(number string)
builder.WithStoreName(name string)
builder.WithStoreAddress(address string)
builder.WithDate(date string)
builder.WithVATRegistration(vatNo string)
builder.WithQRCode(data string)
builder.WithVATPercentage(percentage float64)
builder.WithEnglish()
builder.WithArabic()
builder.WithLabels(labels Labels)

// Products
builder.AddProduct(name string, quantity, unitPrice float64)
builder.AddProductWithDiscount(name string, quantity, unitPrice, discountPercent, discountAmount float64)

// Build
inv := builder.Build()
```

### Models

```go
// Invoice - complete calculated invoice
type Invoice struct {
    Title             string
    InvoiceNumber     string
    StoreName         string
    StoreAddress      string
    Date              string
    VATRegistrationNo string
    Products          []Product
    TotalGross        float64
    TotalDiscount     float64
    TotalTaxableAmt   float64
    TotalVAT          float64
    TotalWithVAT      float64
    QRCodeData        string
    VATPercentage     float64
    Labels            Labels
    Language          string
    IsRTL             bool
}

// Product - calculated product line
type Product struct {
    Name            string
    Quantity        float64
    UnitPrice       float64
    DiscountPercent float64
    DiscountAmount  float64
    GrossAmount     float64
    NetAmount       float64
    TaxableAmt      float64
    VATAmount       float64
    TotalWithVAT    float64
}
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/signintech/gopdf` | PDF generation with Unicode support |
| `github.com/skip2/go-qrcode` | QR code generation |
| `github.com/stretchr/testify` | Testing assertions |

## License

MIT License

