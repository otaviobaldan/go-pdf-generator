package config

import (
	"github.com/jung-kurt/gofpdf"
)

type PdfConfig struct {
	Orientation string
	Units string
	PaperSize string
	Margins Margins
	RegisterFonts bool
}

type Margins struct {
	Left float64
	Top float64
	Right float64
}

type TextConfig struct {
	FontFamily string
	Align string
	Style string
	Size float64
	Color Color
}

type Color struct {
	R,G,B int
}

func NewPdfConfig(
	orientation, units, paperSize string,
	leftMargin, rightMargin, topMargin float64,
	registerFonts bool,
) (pdfConfig *PdfConfig) {
	return &PdfConfig{
		Orientation:   orientation,
		Units:         units,
		PaperSize:     paperSize,
		Margins:       Margins{
			Left:  leftMargin,
			Top:  topMargin,
			Right: rightMargin,
		},
		RegisterFonts: registerFonts,
	}
}

func (cfg *PdfConfig) RegisterExternalFonts(pdf *gofpdf.Fpdf) (err error) {
	return RegisterBookmanOldStyle(pdf)
}