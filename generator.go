package go_pdf_generator

import (
	"errors"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/otaviobaldan/go-pdf-generator/config"
	"github.com/otaviobaldan/go-pdf-generator/constants"
)

type PdfGenerator struct {
	Pdf            *gofpdf.Fpdf
	TxtCfgHeader   *config.TextConfig
	TxtCfgFooter   *config.TextConfig
	TxtCfgTitle    *config.TextConfig
	TxtCfgSubtitle *config.TextConfig
	TxtCfgText     *config.TextConfig
}

func NewPdfGenerator(
	config *config.PdfConfig,
	TxtCfgHeader *config.TextConfig,
	TxtCfgFooter *config.TextConfig,
	TxtCfgTitle *config.TextConfig,
	TxtCfgSubtitle *config.TextConfig,
	TxtCfgText *config.TextConfig,
) (pdfGenerator *PdfGenerator, err error) {
	pdfGenerator = &PdfGenerator{
		Pdf:            gofpdf.New(config.Orientation, config.Units, config.PaperSize, ""),
		TxtCfgHeader:   TxtCfgHeader,
		TxtCfgTitle:    TxtCfgTitle,
		TxtCfgSubtitle: TxtCfgSubtitle,
		TxtCfgText:     TxtCfgText,
		TxtCfgFooter:   TxtCfgFooter,
	}

	if config.RegisterFonts {
		err = config.RegisterExternalFonts(pdfGenerator.Pdf)
		if err != nil {
			return nil, errors.New("error loading fonts")
		}
	}

	margins := config.Margins
	pdfGenerator.Pdf.SetMargins(margins.Left, margins.Top, margins.Right)
	pdfGenerator.Pdf.SetAutoPageBreak(true, 40)
	pdfGenerator.Pdf.AddPage()

	return pdfGenerator, nil
}

// GenerateDefaultHeader - This function will generate a default header, for now without image
func (pg *PdfGenerator) GenerateDefaultHeader(headerText string) {
	cfg := pg.TxtCfgHeader
	color := cfg.Color
	pg.Pdf.SetHeaderFunc(func() {
		pg.Pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
		pg.Pdf.SetTextColor(color.R, color.G, color.B)

		// Calculate width of title and position
		wd := pg.Pdf.GetStringWidth(headerText) + 6
		pg.Pdf.SetX((210 - wd) / 2)

		// Title
		pg.Pdf.CellFormat(wd, 9, headerText, "", 0, cfg.Align, false, 0, "")
		// Line break
		pg.Pdf.Ln(10)
	})
}

// GenerateDefaultFooter - This function will generate a page number and a text that could be left or center aligned
func (pg *PdfGenerator) GenerateDefaultFooter(text string, pageNumber bool) {
	cfg := pg.TxtCfgFooter
	color := cfg.Color
	pg.Pdf.SetFooterFunc(func() {
		// Position at 1.5 cm from bottom
		pg.Pdf.SetY(constants.SizeLessOneHalfCmInPoints)

		pg.Pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
		pg.Pdf.SetTextColor(color.R, color.G, color.B)
		pg.Pdf.CellFormat(0, 28.34, text,
			"", 0, cfg.Align, false, 0, "")

		if pageNumber {
			// page number only black
			pg.Pdf.SetTextColor(0, 0, 0)
			pg.Pdf.CellFormat(0, 28.34, fmt.Sprintf("Pág. %d", pg.Pdf.PageNo()),
				"", 0, constants.AlignRight, false, 0, "")
		}
	})
}

func (pg *PdfGenerator) GenerateTitle(title string) {
	cfg := pg.TxtCfgTitle
	color := cfg.Color
	pg.Pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
	pg.Pdf.SetTextColor(color.R, color.G, color.B)
	pg.Pdf.CellFormat(0, constants.SizeTitleHeight, title,
		"", 1, cfg.Align, false, 0, "")

	// Line break 
	pg.Pdf.Ln(constants.SizeLineBreak)
}

func (pg *PdfGenerator) GenerateSubtitle(subtitle string) {
	cfg := pg.TxtCfgSubtitle
	color := cfg.Color

	pg.Pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
	pg.Pdf.SetTextColor(color.R, color.G, color.B)
	pg.Pdf.CellFormat(0, constants.SizeSubTitleHeight, subtitle,
		"", 1, cfg.Align, false, 0, "")

	// Line break 
	pg.Pdf.Ln(constants.SizeLineBreak)
}

func (pg *PdfGenerator) GenerateText(text string) {
	cfg := pg.TxtCfgText
	color := cfg.Color

	pg.Pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
	pg.Pdf.SetTextColor(color.R, color.G, color.B)

	pg.Pdf.MultiCell(0, constants.SizeTextHeight, text, "", cfg.Align, false)

	// Line break
	pg.Pdf.Ln(-1)
}

func (pg *PdfGenerator) GenerateSignature(signatureName string) {
	currentY := pg.Pdf.GetY()
	left, _, right, _ := pg.Pdf.GetMargins()
	width, _ := pg.Pdf.GetPageSize()

	lineSize := float64(130)
	availableSpace := (width - left - right - lineSize) / 2
	lineY := currentY + 20

	pg.Pdf.Line(left+availableSpace, lineY, left+availableSpace+lineSize, lineY)
	pg.Pdf.CellFormat(0, 50, signatureName, "", 1, "C", false, 0, "")
}