package config

import (
	"github.com/jung-kurt/gofpdf"
)

func RegisterBookmanOldStyle(pdf *gofpdf.Fpdf) (err error) {
	pdf.AddUTF8Font("Bookman", "", "./font/bookman-old-style.ttf")
	pdf.AddUTF8Font("Bookman", "B", "./font/bookman-old-style-bold.ttf")
	pdf.AddUTF8Font("Bookman", "BI", "./font/bookman-old-style-bold-italic.ttf")
	pdf.AddUTF8Font("Bookman", "I", "./font/bookman-old-style-italic.ttf")

	return nil
}
