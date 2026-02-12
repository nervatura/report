package report

import "strings"

func (rpt *Report) createCell(options IM) float64 {
	rpt.setPageStyle(options)
	virtual := ToBoolean(options["virtual"], false)
	text := ToString(options["text"], "")
	multiline := ToBoolean(options["multiline"], false)
	ln := ToBoolean(options["ln"], false)
	padding := rpt.getCellPadding(options)
	border := ToString(options["border"], "")
	align := ToString(options["align"], "L")
	fill := ToRGBA(options["backgroundColor"], rpt.BackgroundColor) != rpt.BackgroundColor

	pageWidth, _ := rpt.Pdf.GetPageSize()
	lineHt := rpt.Pdf.GetFontSize()
	startY := rpt.Pdf.GetY()
	startX := rpt.Pdf.GetX()
	rpt.setCellX(options, startX)

	width := rpt.getCellWidth(options, text, padding, pageWidth, startX)
	height := ToFloat(options["height"], 0)

	if multiline {
		return rpt.createCellMultiline(options, text, width, height, padding, lineHt, border, align, fill, ln, startY, virtual)
	}
	return rpt.createCellSingle(options, width, height, padding, lineHt, border, align, fill, ln, startY, virtual)
}

func (rpt *Report) getCellPadding(options IM) float64 {
	fill := ToRGBA(options["backgroundColor"], rpt.BackgroundColor) != rpt.BackgroundColor
	border := ToString(options["border"], "")
	multiline := ToBoolean(options["multiline"], false)
	if !fill && border == "" && !multiline {
		return 0
	}
	return _padding
}

func (rpt *Report) setCellX(options IM, startX float64) {
	xCol := ToFloat(options["xCol"], startX)
	if xCol != startX {
		rpt.Pdf.SetX(xCol)
	}
}

func (rpt *Report) getCellWidth(options IM, text string, padding, pageWidth, startX float64) float64 {
	ln := ToBoolean(options["ln"], false)
	xCol := ToFloat(options["xCol"], startX)
	width := ToFloat(options["columnWidth"], 0)
	if ln && ToBoolean(options["extend"], false) {
		return pageWidth - rpt.RightMargin - xCol
	}
	if width > 0 {
		return width
	}
	widthStr := ToString(options["width"], "")
	if strings.HasSuffix(widthStr, "%") {
		width = ToFloat(strings.Replace(widthStr, "%", "", -1), width) / 100
		return (pageWidth - rpt.LeftMargin - rpt.RightMargin) * width
	}
	width = ToFloat(widthStr, width)
	if width > 0 {
		width += padding
	}
	if width == 0 {
		width = rpt.Pdf.GetTextWidth(text) + padding
	}
	if startX+padding+width > pageWidth-rpt.RightMargin {
		return 0
	}
	return width
}

func (rpt *Report) createCellMultiline(options IM, text string, width, height, padding, lineHt float64, border, align string, fill, ln bool, startY float64, virtual bool) float64 {
	if height == 0 {
		height = rpt.getCellHeight(text, width, options)
	}
	if !virtual {
		rpt.Pdf.MultiCell(IM{
			"w": width, "h": height, "lineH": lineHt + padding, "padding": padding,
			"txtStr": text, "borderStr": border, "alignStr": align, "fill": fill,
		})
	}
	if ln {
		rpt.Pdf.SetY(startY + height)
	} else {
		rpt.Pdf.SetY(startY)
	}
	return height
}

func (rpt *Report) createCellSingle(options IM, width, height, padding, lineHt float64, border, align string, fill, ln bool, startY float64, virtual bool) float64 {
	if lineHt+padding > height {
		height = lineHt + padding
	}
	if !virtual {
		rpt.Pdf.Cell(IM{
			"w": width, "h": height, "padding": padding, "txtStr": ToString(options["text"], ""), "borderStr": border,
			"alignStr": align, "fill": fill, "ln": ln,
		})
	}
	if rpt.Pdf.GetY()-startY > height {
		return rpt.Pdf.GetY() - startY
	}
	return height
}
