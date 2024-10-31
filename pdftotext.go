package pdftotextgo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var (
	// ErrPDFOpen is returned when there is an error opening the PDF file
	ErrPDFOpen = errors.New("error opening PDF file")
	// ErrOutputFile is returned when there is an error opening the output file
	ErrOutputFile = errors.New("error opening output file")
	// ErrPermissions is returned when there is an error related to PDF permissions
	ErrPermissions = errors.New("error related to PDF permissions")
	// ErrInvalidPage is returned when the page number is invalid
	ErrInvalidPage = errors.New("invalid page number")
	// ErrInvalidRange is returned when the page range is invalid
	ErrInvalidRange = errors.New("invalid page range")
	// ErrCommandFailed is returned when the pdftotext command fails
	ErrCommandFailed = errors.New("pdftotext command failed")
	// ErrBinaryNotFound is returned when the pdftotext binary is not found
	ErrBinaryNotFound = errors.New("pdftotext binary not found")
)

// EOLType represents the end-of-line convention
type EOLType string

const (
	// EOLUnix represents the Unix end-of-line convention
	EOLUnix EOLType = "unix"
	// EOLDos represents the DOS end-of-line convention
	EOLDos EOLType = "dos"
	// EOLMac represents the Mac end-of-line convention
	EOLMac EOLType = "mac"
)

// Options represents the configuration options for the PDF conversion
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

// Converter represents a PDF to text converter
type Converter struct{ binaryPath string }

// New creates a new Converter instance
func New() (*Converter, error) {
	binaryPath, err := exec.LookPath("pdftotext")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrBinaryNotFound, err)
	}
	return &Converter{binaryPath: binaryPath}, nil
}

// Convert converts a PDF file to text and returns the result
func (c *Converter) Convert(ctx context.Context, inputPath string, opts *Options) (string, error) {
	var stdout, stderr bytes.Buffer

	args := c.buildArgs(opts, inputPath, "-")
	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", c.handleError(err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

// ConvertToFile converts a PDF file to text and saves it to the specified output file
func (c *Converter) ConvertToFile(ctx context.Context, inputPath, outputPath string, opts *Options) error {
	var stderr bytes.Buffer

	args := c.buildArgs(opts, inputPath, outputPath)
	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return c.handleError(err, stderr.String())
	}
	return nil
}

func (c *Converter) handleError(err error, stderr string) error {
	if exitErr, ok := err.(*exec.ExitError); ok {
		switch exitErr.ExitCode() {
		case 1:
			return fmt.Errorf("%w: %s", ErrPDFOpen, stderr)
		case 2:
			return fmt.Errorf("%w: %s", ErrOutputFile, stderr)
		case 3:
			return fmt.Errorf("%w: %s", ErrPermissions, stderr)
		default:
			return fmt.Errorf("%w: %s", ErrCommandFailed, stderr)
		}
	}
	return fmt.Errorf("failed to run pdftotext: %w", err)
}

func (c *Converter) buildArgs(opts *Options, inputPath, outputPath string) []string {
	if opts == nil {
		opts = &Options{}
	}

	var args []string

	appendFlag := func(flag string, value any) {
		switch v := value.(type) {
		case int:
			if v > 0 {
				args = append(args, flag, strconv.Itoa(v))
			}
		case float64:
			if v > 0 {
				args = append(args, flag, strconv.FormatFloat(v, 'f', -1, 64))
			}
		case string:
			if v != "" {
				args = append(args, flag, v)
			}
		case bool:
			if v {
				args = append(args, flag)
			}
		}
	}

	appendFlag("-f", opts.FirstPage)
	appendFlag("-l", opts.LastPage)
	appendFlag("-r", opts.Resolution)
	appendFlag("-x", opts.CropX)
	appendFlag("-y", opts.CropY)
	appendFlag("-W", opts.CropWidth)
	appendFlag("-H", opts.CropHeight)
	appendFlag("-layout", opts.Layout)
	appendFlag("-fixed", opts.FixedPitch)
	appendFlag("-raw", opts.Raw)
	appendFlag("-nodiag", opts.NoDiagonal)
	appendFlag("-htmlmeta", opts.HTMLMeta)
	appendFlag("-bbox", opts.BBox)
	appendFlag("-bbox-layout", opts.BBoxLayout)
	appendFlag("-tsv", opts.TSV)
	appendFlag("-cropbox", opts.CropBox)
	appendFlag("-colspacing", opts.ColSpacing)
	appendFlag("-enc", opts.Encoding)
	appendFlag("-eol", string(opts.EOL))
	appendFlag("-nopgbrk", opts.NoPageBreaks)
	appendFlag("-opw", opts.OwnerPassword)
	appendFlag("-upw", opts.UserPassword)
	appendFlag("-q", opts.Quiet)

	args = append(args, inputPath)
	if outputPath != "" {
		args = append(args, outputPath)
	}
	return args
}
