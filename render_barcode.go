package report

import (
	"bytes"
	"image/jpeg"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/boombuler/barcode/code39"
	"github.com/boombuler/barcode/ean"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode/twooffive"
)

func (rpt *Report) createBarcode(v *Barcode, virtual, ln bool) (float64, float64) {
	pageWidth, _ := rpt.Pdf.GetPageSize()
	rpt.Pdf.SetTextColor(int(rpt.TextColor.R), int(rpt.TextColor.G), int(rpt.TextColor.B))
	width := v.Width
	strWidth := rpt.Pdf.GetTextWidth(v.Value)
	if width == 0 {
		width = strWidth + 1.5*_padding
	}
	height := v.Height
	if height == 0 {
		height = 10 * _mmPt
	}
	startX := rpt.Pdf.GetX()
	startY := rpt.Pdf.GetY()
	lineHt := rpt.Pdf.GetFontSize()
	if ln {
		if v.Extend {
			width = pageWidth - startX - rpt.RightMargin - _padding
		}
	}
	if rpt.checkPageBreak(height) && !virtual {
		rpt.addPage()
	}
	var bcode barcode.Barcode
	switch v.CodeType {
	case "CODE_39", "code39":
		bcode, _ = code39.Encode(v.Value, true, true)

	case "ITF", "i2of5":
		bcode, _ = twooffive.Encode(v.Value, true)

	case "CODE_128", "code128":
		bcode, _ = code128.Encode(v.Value)

	case "EAN", "ean":
		bcode, _ = ean.Encode(v.Value)

	case "QR", "qr":
		width = height
		bcode, _ = qr.Encode(v.Value, qr.H, qr.Unicode)

	}

	if bcode != nil {
		buf := new(bytes.Buffer)
		err := jpeg.Encode(buf, bcode, &jpeg.Options{Quality: 100})
		if err == nil {
			rpt.Pdf.AddImage(&Image{Data: buf.Bytes(), Width: width, Height: height}, startX, startY, IM{})
		}
		if v.VisibleValue {
			rpt.Pdf.SetXY(startX+(width-strWidth)/2, startY+height+1.5*_padding)
			rpt.Pdf.Text(v.Value, rpt.pageBreak-rpt.footerHeight)
			height += lineHt + 1.5*_padding
		}
	}
	return height, width
}
