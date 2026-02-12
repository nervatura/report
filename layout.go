package report

import "image/color"

func (rpt *Report) createHeaderAndFooter() {
	createSection := func(section string, elements []PageItem) {
		for index := 0; index < len(elements); index++ {
			switch elements[index].Item.(type) {
			case *Row, *VGap, *HLine:
				rpt.createElement(section, elements[index].Item)
			}
		}
	}
	createSection("header", rpt.header)
	cx := rpt.Pdf.GetX()
	cy := rpt.Pdf.GetY()
	_, pageHeight := rpt.Pdf.GetPageSize()
	rpt.Pdf.SetY(pageHeight - rpt.BottomMargin - rpt.footerHeight)
	createSection("footer", rpt.footer)
	rpt.Pdf.SetXY(cx, cy)
}

func (rpt *Report) checkPageBreak(nextHeight float64) bool {
	cy := rpt.Pdf.GetY()
	dLine := rpt.pageBreak
	if cy < rpt.pageBreak-rpt.footerHeight {
		dLine -= rpt.footerHeight
	}
	if cy+nextHeight > dLine {
		return true
	}
	return false
}

func (rpt *Report) getFooterHeight() (fHeight float64) {
	for index := 0; index < len(rpt.footer); index++ {
		switch v := rpt.footer[index].Item.(type) {
		case *Row:
			fHeight += rpt.createRow("footer", v, true)
		case *VGap:
			fHeight += v.Height
		case *HLine:
			fHeight += (1 + v.Gap)
		}
	}
	_, pageHeight := rpt.Pdf.GetPageSize()
	rpt.pageBreak = pageHeight - rpt.BottomMargin
	return fHeight
}

func (rpt *Report) getCellHeight(text string, width float64, options IM) float64 {
	if text == "" {
		text = "X"
	}
	rpt.Pdf.SetFont(rpt.FontFamily, options["fontStyle"].(string), options["fontSize"].(float64))
	lines := rpt.wrapTextLines(text, width-_padding)
	lineHt := rpt.Pdf.GetFontSize()
	return float64(len(lines)) * (lineHt + _padding)
}

func (rpt *Report) createGridHeader(headerOptions IM) {
	headerOptions["height"] = float64(0)
	for colIndex := 0; colIndex < len(headerOptions["columns"].([]IM)); colIndex++ {
		column := headerOptions["columns"].([]IM)[colIndex]
		if !headerOptions["merge"].(bool) {
			headerOptions["text"] = column["label"]
			headerOptions["ln"] = column["ln"]
			if column["columnWidth"] == 0 {
				column["columnWidth"] = (headerOptions["gridWidth"].(float64) - headerOptions["columnsWidth"].(float64)) / float64(len(headerOptions["columns"].([]IM)))
			}
			headerOptions["columnWidth"] = column["columnWidth"]
		} else {
			headerOptions["text"] = headerOptions["text"].(string) + " " + column["label"].(string)
		}
		headerOptions["align"] = column["headerAlign"]
		rpt.createCell(headerOptions)
	}
}

func (rpt *Report) setPageStyle(options IM) {
	fontStyle := ToString(options["fontStyle"], "")
	fontSize := ToFloat(options["fontSize"], rpt.FontSize)
	rpt.Pdf.SetFont(rpt.FontFamily, fontStyle, fontSize)

	if textColor, textKey := options["textColor"]; textKey {
		rpt.Pdf.SetTextColor(int(textColor.(color.RGBA).R), int(textColor.(color.RGBA).G), int(textColor.(color.RGBA).B))
	}
	if borderColor, borderKey := options["borderColor"]; borderKey {
		rpt.Pdf.SetDrawColor(int(borderColor.(color.RGBA).R), int(borderColor.(color.RGBA).G), int(borderColor.(color.RGBA).B))
	}
	if backgroundColor, backgroundKey := options["backgroundColor"]; backgroundKey {
		rpt.Pdf.SetFillColor(int(backgroundColor.(color.RGBA).R), int(backgroundColor.(color.RGBA).G), int(backgroundColor.(color.RGBA).B))
	}
}

func (rpt *Report) addPage() {
	rpt.Pdf.AddPage()
	rpt.Pdf.SetXY(rpt.LeftMargin, rpt.TopMargin)
	rpt.createHeaderAndFooter()
}

func (rpt *Report) onPage() {
	rpt.addPage()
}
