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
	isPointUnit    bool
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
		isPointUnit:    config.Units == constants.UnitsPoints,
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
	pdfGenerator.Pdf.SetAutoPageBreak(true, pdfGenerator.calculateSize(40))
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
		stringWidth := pg.calculateSize(pg.Pdf.GetStringWidth(headerText) + 6)
		width, _ := pg.Pdf.GetPageSize()
		pg.Pdf.SetX((width - stringWidth) / 2)

		// Title
		pg.Pdf.CellFormat(stringWidth, 9, headerText, "", 0, cfg.Align, false, 0, "")
		// Line break
		pg.Pdf.Ln(pg.calculateSize(10))
	})
}

// GenerateDefaultFooter - This function will generate a page number and a text that could be left or center aligned
func (pg *PdfGenerator) GenerateDefaultFooter(text string, pageNumber bool) {
	cfg := pg.TxtCfgFooter
	color := cfg.Color
	pg.Pdf.SetFooterFunc(func() {

		// Position at 1.5 cm from bottom
		pg.Pdf.SetY(pg.calculateSize(-15))

		pg.Pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
		pg.Pdf.SetTextColor(color.R, color.G, color.B)
		pg.Pdf.CellFormat(0, pg.calculateSize(10), text,
			"", 0, cfg.Align, false, 0, "")

		if pageNumber {
			// page number only black
			pg.Pdf.SetTextColor(0, 0, 0)
			pg.Pdf.CellFormat(0, pg.calculateSize(10), fmt.Sprintf("Pág. %d", pg.Pdf.PageNo()),
				"", 0, constants.AlignRight, false, 0, "")
		}
	})
}

func (pg *PdfGenerator) GenerateTitle(title string) {
	cfg := pg.TxtCfgTitle
	color := cfg.Color
	pg.Pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
	pg.Pdf.SetTextColor(color.R, color.G, color.B)
	pg.Pdf.CellFormat(0, pg.calculateSize(constants.SizeTitleHeight), title,
		"", 1, cfg.Align, false, 0, "")

	// Line break 
	pg.Pdf.Ln(constants.SizeLineBreak)
}

func (pg *PdfGenerator) GenerateSubtitle(subtitle string) {
	cfg := pg.TxtCfgSubtitle
	color := cfg.Color

	pg.Pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
	pg.Pdf.SetTextColor(color.R, color.G, color.B)
	pg.Pdf.CellFormat(0, pg.calculateSize(constants.SizeSubTitleHeight), subtitle,
		"", 1, cfg.Align, false, 0, "")

	// Line break 
	pg.Pdf.Ln(pg.calculateSize(constants.SizeLineBreak))
}

func (pg *PdfGenerator) GenerateText(text string) {
	cfg := pg.TxtCfgText
	color := cfg.Color

	pg.Pdf.SetFont(cfg.FontFamily, cfg.Style, cfg.Size)
	pg.Pdf.SetTextColor(color.R, color.G, color.B)

	pg.Pdf.MultiCell(0, pg.calculateSize(constants.SizeTextHeight), text, "", cfg.Align, false)

	// Line break
	pg.Pdf.Ln(-1)
}

func (pg *PdfGenerator) GenerateSignature(signatureName string) {
	currentY := pg.Pdf.GetY()
	left, _, right, _ := pg.Pdf.GetMargins()
	width, _ := pg.Pdf.GetPageSize()

	lineSize := pg.calculateSize(130)
	availableSpace := pg.calculateSize((width - left - right - lineSize) / 2)
	lineY := pg.calculateSize(currentY + 20)
	lineInit := pg.calculateSize(left + availableSpace)
	lineEnd := pg.calculateSize(left + availableSpace + lineSize)

	pg.Pdf.Line(lineInit, lineY, lineEnd, lineY)
	pg.Pdf.CellFormat(0, pg.calculateSize(50), signatureName, "", 1, "C", false, 0, "")
}

func (pg *PdfGenerator) calculateSize(size float64) float64 {
	if pg.isPointUnit {
		return size * 2.834
	}
	return size
}