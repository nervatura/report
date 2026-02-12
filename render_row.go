package report

func (rpt *Report) createRow(section string, rowElement *Row, virtual bool) float64 {
	maxHeight := rowElement.Height
	for index := 0; index < len(rowElement.Columns); index++ {
		startY := rpt.Pdf.GetY()
		if rpt.Pdf.GetX() != rpt.LeftMargin {
			rpt.Pdf.SetX(rpt.Pdf.GetX() + rowElement.HGap)
		}
		startX := rpt.Pdf.GetX()
		ln := len(rowElement.Columns)-1 == index
		element := rowElement.Columns[index].Item
		maxHeight = rpt.createRowElement(section, element, maxHeight, startX, startY, ln, index, len(rowElement.Columns), virtual)
	}
	return maxHeight
}

func (rpt *Report) createRowElement(section string, element interface{}, maxHeight, startX, startY float64, ln bool, index, colCount int, virtual bool) float64 {
	switch v := element.(type) {
	case *Cell:
		return rpt.createRowCell(section, v, maxHeight, ln, virtual)
	case *Image:
		return rpt.createRowImage(v, maxHeight, startX, startY, index, colCount, virtual)
	case *Barcode:
		return rpt.createRowBarcode(v, maxHeight, startX, startY, ln, index, colCount, virtual)
	case *Separator:
		return rpt.createRowSeparator(v, maxHeight, index, colCount, virtual)
	}
	return maxHeight
}

func (rpt *Report) createRowCell(section string, v *Cell, maxHeight float64, ln, virtual bool) float64 {
	options := IM{
		"height": maxHeight, "fontFamily": rpt.FontFamily, "fontStyle": v.FontStyle,
		"fontSize": v.FontSize, "textColor": v.TextColor, "borderColor": v.BorderColor,
		"backgroundColor": v.BackgroundColor, "text": rpt.setValue(v.Value),
		"width": v.Width, "border": v.Border, "align": v.Align,
		"multiline": false, "extend": true, "virtual": virtual, "ln": ln,
	}
	if section == "details" {
		options["multiline"] = v.Multiline
		if v.Multiline {
			options["height"] = 0
		}
	}
	cellHeight := rpt.createCell(options)
	if cellHeight > maxHeight || maxHeight == 0 {
		maxHeight = cellHeight
	}
	rpt.addToXML(section, []string{ToString(v.Name, "head"), options["text"].(string), ToString(v.Name, "head")})
	return maxHeight
}

func (rpt *Report) createRowImage(v *Image, maxHeight, startX, startY float64, index, colCount int, virtual bool) float64 {
	if v.Src == "" {
		return maxHeight
	}
	height, width := rpt.createImage(v, maxHeight, virtual)
	if height > maxHeight {
		maxHeight = height
	}
	if colCount-1 == index {
		rpt.Pdf.SetXY(rpt.LeftMargin, startY+maxHeight)
	} else {
		rpt.Pdf.SetXY(startX+width, startY)
	}
	return maxHeight
}

func (rpt *Report) createRowBarcode(v *Barcode, maxHeight, startX, startY float64, ln bool, index, colCount int, virtual bool) float64 {
	height, width := rpt.createBarcode(v, virtual, ln)
	if height > maxHeight || maxHeight == 0 {
		maxHeight = height
	}
	if colCount-1 == index {
		rpt.Pdf.SetXY(rpt.LeftMargin, startY+maxHeight)
	} else {
		rpt.Pdf.SetXY(startX+width+_padding, startY)
	}
	return maxHeight
}

func (rpt *Report) createRowSeparator(v *Separator, maxHeight float64, index, colCount int, virtual bool) float64 {
	if !virtual {
		rpt.Pdf.Line(rpt.Pdf.GetX()+v.Gap, rpt.Pdf.GetY(), rpt.Pdf.GetX()+v.Gap, rpt.Pdf.GetY()+maxHeight)
	}
	if colCount-1 == index {
		rpt.Pdf.SetX(rpt.Pdf.GetX() + v.Gap)
	}
	if v.Gap > maxHeight || maxHeight == 0 {
		return v.Gap
	}
	return maxHeight
}
