package go_pdf_generator

import (
	"errors"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/otaviobaldan/go-pdf-generator/config"
	"github.com/otaviobaldan/go-pdf-generator/constants"
)

type PdfGenerator struct {
	pdf            *gofpdf.Fpdf
	TxtCfgFooter   *config.TextConfig
	TxtCfgTitle    *config.TextConfig
	TxtCfgSubtitle *config.TextConfig
	TxtCfgText     *config.TextConfig
}

func NewPdfGenerator(
	config *config.PdfConfig,
	TxtCfgFooter *config.TextConfig,
	TxtCfgTitle *config.TextConfig,
	TxtCfgSubtitle *config.TextConfig,
	TxtCfgText *config.TextConfig,
) (pdfGenerator *PdfGenerator, err error) {
	pdfGenerator = new(PdfGenerator)
	pdfGenerator.pdf = gofpdf.New(config.Orientation, config.Units, config.PaperSize, "")

	if config.RegisterFonts {
		err = config.RegisterExternalFonts(pdfGenerator.pdf)
		if err != nil {
			return nil, errors.New("error loading fonts")
		}
	}

	margins := config.Margins
	pdfGenerator.pdf.SetMargins(margins.Left, margins.Top, margins.Right)

	return pdfGenerator, nil
}

// GenerateDefaultHeader - This function will generate a default header, for now without image
func (pg *PdfGenerator) GenerateDefaultHeader(headerText string) {
	cfg := pg.TxtCfgFooter
	color := cfg.Color
	pg.pdf.SetHeaderFunc(func() {
		pg.pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
		pg.pdf.SetTextColor(color.R, color.G, color.B)

		// Calculate width of title and position
		wd := pg.pdf.GetStringWidth(headerText) + 6
		pg.pdf.SetX((210 - wd) / 2)

		// Title
		pg.pdf.CellFormat(wd, 9, headerText, "", 0, cfg.Align, false, 0, "")
		// Line break
		pg.pdf.Ln(10)
	})
}

// GenerateDefaultFooter - This function will generate a page number and a text that could be left or center aligned
func (pg *PdfGenerator) GenerateDefaultFooter(text string, pageNumber bool) {
	cfg := pg.TxtCfgFooter
	color := cfg.Color
	pg.pdf.SetFooterFunc(func() {
		// Position at 1.5 cm from bottom
		pg.pdf.SetY(constants.SizeLessOneHalfCmInPoints)

		pg.pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
		pg.pdf.SetTextColor(color.R, color.G, color.B)
		pg.pdf.CellFormat(0, 28.34, text,
			"", 0, cfg.Align, false, 0, "")

		if pageNumber {
			// page number only black
			pg.pdf.SetTextColor(0, 0, 0)
			pg.pdf.CellFormat(0, 28.34, fmt.Sprintf("Pág. %d", pg.pdf.PageNo()),
				"", 0, constants.AlignRight, false, 0, "")
		}
	})
}

func (pg *PdfGenerator) GenerateTitle(title string) {
	cfg := pg.TxtCfgTitle
	color := cfg.Color
	pg.pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
	pg.pdf.SetTextColor(color.R, color.G, color.B)
	pg.pdf.CellFormat(0, constants.SizeTitleHeight, title,
		"", 1, cfg.Align, false, 0, "")

	// Line break 
	pg.pdf.Ln(constants.SizeLineBreak)
}

func (pg *PdfGenerator) GenerateSubtitle(subtitle string) {
	cfg := pg.TxtCfgSubtitle
	color := cfg.Color

	pg.pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
	pg.pdf.SetTextColor(color.R, color.G, color.B)
	pg.pdf.CellFormat(0, constants.SizeSubTitleHeight, subtitle,
		"", 1, cfg.Align, false, 0, "")

	// Line break 
	pg.pdf.Ln(constants.SizeLineBreak)
}

func (pg *PdfGenerator) GenerateText(text string) {
	cfg := pg.TxtCfgText
	color := cfg.Color

	pg.pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
	pg.pdf.SetTextColor(color.R, color.G, color.B)

	pg.pdf.MultiCell(0, constants.SizeTextHeight, text, "", cfg.Align, false)

	// Line break
	pg.pdf.Ln(-1)
}
