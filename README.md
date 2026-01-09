# Arabic Tax Invoice PDF Generator (فاتورة ضريبية مبسطة)

Pixel-perfect Arabic tax invoice generator matching Saudi ZATCA requirements.

## Features

- RTL Arabic text support
- 15% VAT calculation
- QR code generation for e-invoicing
- Thermal printer format (80mm width)
- Docker support for portability

## Quick Start

### 1. Download Fonts

```powershell
# Windows
.\scripts\download-fonts.ps1

# Linux/Mac
chmod +x scripts/download-fonts.sh
./scripts/download-fonts.sh
```

### 2. Run Locally

```bash
go mod tidy
go run main.go
```

### 3. Run with Docker

```bash
docker-compose up bill-generator
```

### 4. Run Tests

```bash
# Local
go test -v ./...

# Docker
docker-compose up test
```

## Output

Generated PDF: `invoice_output.pdf`

## Structure

```
bill/
├── main.go              # Main application
├── main_test.go         # Test suite
├── go.mod               # Go modules
├── Dockerfile           # Production container
├── Dockerfile.test      # Test container
├── docker-compose.yml   # Container orchestration
├── fonts/               # Arabic fonts (Amiri)
│   ├── Amiri-Regular.ttf
│   └── Amiri-Bold.ttf
├── scripts/
│   ├── download-fonts.sh
│   └── download-fonts.ps1
└── output/              # Generated PDFs
```

## Customization

Edit `generateSampleInvoice()` in `main.go` to modify:
- Store name and address
- Products and prices
- VAT registration number
- Invoice number
