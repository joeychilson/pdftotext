package pdftotextgo

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const expectedContent = `This is a test PDF document.
If you can read this, you have Adobe Acrobat Reader installed on your computer.`

func TestConverter_Convert(t *testing.T) {
	testPDFPath := filepath.Join("testdata", "test.pdf")

	tests := []struct {
		name          string
		options       *Options
		inputPath     string
		expectedError error
		expectedText  string
		checkContains bool
	}{
		{
			name:          "Non-existent file",
			options:       nil,
			inputPath:     "nonexistent.pdf",
			expectedError: ErrPDFOpen,
		},
		{
			name: "Basic conversion",
			options: &Options{
				Layout:   true,
				Encoding: "UTF-8",
			},
			inputPath:     testPDFPath,
			expectedText:  expectedContent,
			checkContains: true,
		},
		{
			name: "With specific pages",
			options: &Options{
				FirstPage: 1,
				LastPage:  1,
				Layout:    true,
				Encoding:  "UTF-8",
			},
			inputPath:     testPDFPath,
			expectedText:  expectedContent,
			checkContains: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter, err := New()
			if err != nil {
				t.Fatalf("failed to create converter: %v", err)
			}

			ctx := context.Background()
			text, err := converter.Convert(ctx, tt.inputPath, tt.options)

			if tt.expectedError != nil {
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.checkContains {
				normalizedText := strings.ReplaceAll(strings.TrimSpace(text), "\r\n", "\n")
				normalizedExpected := strings.ReplaceAll(strings.TrimSpace(tt.expectedText), "\r\n", "\n")

				if !strings.Contains(normalizedText, normalizedExpected) {
					t.Errorf("expected text to contain:\n%s\n\ngot:\n%s", normalizedExpected, normalizedText)
				}
			}
		})
	}
}

func TestConverter_ConvertToFile(t *testing.T) {
	testPDFPath := filepath.Join("testdata", "test.pdf")
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		options       *Options
		inputPath     string
		outputPath    string
		expectedError error
		checkContent  bool
	}{
		{
			name:          "Non-existent input file",
			options:       nil,
			inputPath:     "nonexistent.pdf",
			outputPath:    filepath.Join(tmpDir, "output1.txt"),
			expectedError: ErrPDFOpen,
		},
		{
			name: "Valid conversion",
			options: &Options{
				Layout:   true,
				Encoding: "UTF-8",
			},
			inputPath:    testPDFPath,
			outputPath:   filepath.Join(tmpDir, "output2.txt"),
			checkContent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter, err := New()
			if err != nil {
				t.Fatalf("failed to create converter: %v", err)
			}

			ctx := context.Background()
			err = converter.ConvertToFile(ctx, tt.inputPath, tt.outputPath, tt.options)

			if tt.expectedError != nil {
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.checkContent {
				content, err := os.ReadFile(tt.outputPath)
				if err != nil {
					t.Fatalf("failed to read output file: %v", err)
				}

				normalizedContent := strings.ReplaceAll(strings.TrimSpace(string(content)), "\r\n", "\n")
				normalizedExpected := strings.ReplaceAll(strings.TrimSpace(expectedContent), "\r\n", "\n")

				if !strings.Contains(normalizedContent, normalizedExpected) {
					t.Errorf("expected output file to contain:\n%s\n\ngot:\n%s", normalizedExpected, normalizedContent)
				}
			}
		})
	}
}

func TestConverter_BuildArgs(t *testing.T) {
	tests := []struct {
		name         string
		options      *Options
		inputPath    string
		outputPath   string
		expectedArgs []string
	}{
		{
			name: "All options",
			options: &Options{
				FirstPage:     1,
				LastPage:      10,
				Resolution:    300,
				CropX:         50,
				CropY:         50,
				CropWidth:     500,
				CropHeight:    700,
				Layout:        true,
				FixedPitch:    12.0,
				Raw:           true,
				NoDiagonal:    true,
				HTMLMeta:      true,
				BBox:          true,
				BBoxLayout:    false,
				TSV:           true,
				CropBox:       true,
				ColSpacing:    0.7,
				Encoding:      "UTF-8",
				EOL:           EOLUnix,
				NoPageBreaks:  true,
				OwnerPassword: "owner123",
				UserPassword:  "user123",
				Quiet:         true,
			},
			inputPath:  "input.pdf",
			outputPath: "output.txt",
			expectedArgs: []string{
				"-f", "1",
				"-l", "10",
				"-r", "300",
				"-x", "50",
				"-y", "50",
				"-W", "500",
				"-H", "700",
				"-layout",
				"-fixed", "12",
				"-raw",
				"-nodiag",
				"-htmlmeta",
				"-bbox",
				"-tsv",
				"-cropbox",
				"-colspacing", "0.7",
				"-enc", "UTF-8",
				"-eol", "unix",
				"-nopgbrk",
				"-opw", "owner123",
				"-upw", "user123",
				"-q",
				"input.pdf",
				"output.txt",
			},
		},
		{
			name:         "Minimal options",
			options:      nil,
			inputPath:    "input.pdf",
			outputPath:   "output.txt",
			expectedArgs: []string{"input.pdf", "output.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter, err := New()
			if err != nil {
				t.Fatalf("failed to create converter: %v", err)
			}

			args := converter.buildArgs(tt.options, tt.inputPath, tt.outputPath)

			if len(args) != len(tt.expectedArgs) {
				t.Errorf("expected %d args, got %d", len(tt.expectedArgs), len(args))
				t.Errorf("expected: %v", tt.expectedArgs)
				t.Errorf("got: %v", args)
				return
			}

			for i := range args {
				if args[i] != tt.expectedArgs[i] {
					t.Errorf("arg %d: expected %q, got %q", i, tt.expectedArgs[i], args[i])
				}
			}
		})
	}
}
