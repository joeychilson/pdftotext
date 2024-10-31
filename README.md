# pdftotextgo

A Go library for converting PDF files to text using the `pdftotext` utility.

## Prerequisites

- `pdftotext` utility installed on your system (usually part of the `poppler-utils` package)

### Installing pdftotext

**Ubuntu/Debian:**
```bash
sudo apt-get install poppler-utils
```

**macOS:**
```bash
brew install poppler
```

## Installation

```bash
go get github.com/joeychilson/pdftotextgo
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/joeychilson/pdftotextgo"
)

func main() {
    options := pdftotextgo.Options{
        Layout:   true,
        Encoding: "UTF-8",
    }

    converter, err := pdftotextgo.New(options)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    text, err := converter.Convert(ctx, "input.pdf")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(text)
}
```

## Converting to File

```go
err = converter.ConvertToFile(ctx, "input.pdf", "output.txt")
if err != nil {
    log.Fatal(err)
}
```

## Available Options

```go
type Options struct {
	// FirstPage is the first page to convert
	FirstPage int
	// LastPage is the last page to convert
	LastPage int
	// Resolution is the resolution in DPI (default 72)
	Resolution int
	// CropX is the X-coordinate of crop area
	CropX int
	// CropY is the Y-coordinate of crop area
	CropY int
	// CropWidth is the width of crop area
	CropWidth int
	// CropHeight is the height of crop area
	CropHeight int
	// Layout maintains the original layout
	Layout bool
	// FixedPitch keeps the text in a fixed-pitch font
	FixedPitch float64
	// Raw keeps text in content stream order
	Raw bool
	// NoDiagonal discards diagonal text
	NoDiagonal bool
	// HTMLMeta generates HTML with meta information
	HTMLMeta bool
	// BBox generates XHTML with word bounding boxes
	BBox bool
	// BBoxLayout generates XHTML with block/line/word bounding boxes
	BBoxLayout bool
	// TSV generates TSV with bounding box information
	TSV bool
	// CropBox uses crop box instead of media box
	CropBox bool
	// ColSpacing is the column spacing (default 0.7)
	ColSpacing float64
	// Encoding is the text output encoding (default UTF-8)
	Encoding string
	// EOL is the end-of-line convention (default Unix)
	EOL EOLType
	// NoPageBreaks don't insert page breaks
	NoPageBreaks bool
	// OwnerPassword is the PDF owner password
	OwnerPassword string
	// UserPassword is the PDF user password
	UserPassword string
	// Quiet suppresses messages and errors
	Quiet bool
}
```

## Error Handling

The library provides specific error types for common failure cases:

```go
var (
    ErrPDFOpen        = errors.New("error opening PDF file")
    ErrOutputFile     = errors.New("error opening output file")
    ErrPermissions    = errors.New("error related to PDF permissions")
    ErrInvalidPage    = errors.New("invalid page number")
    ErrInvalidRange   = errors.New("invalid page range")
    ErrCommandFailed  = errors.New("pdftotext command failed")
    ErrBinaryNotFound = errors.New("pdftotext binary not found")
)
```