package report

func (rpt *Report) createHTML(v *HTML) {
	lineHt := rpt.Pdf.GetFontSize()
	htmlStr := v.Value
	fieldname := ToString(v.Fieldname, "head")
	options := IM{
		"fontFamily":      rpt.FontFamily,
		"fontStyle":       rpt.FontStyle,
		"fontSize":        rpt.FontSize,
		"textColor":       rpt.TextColor,
		"borderColor":     rpt.BorderColor,
		"backgroundColor": rpt.BackgroundColor}
	htmlStr = rpt.setHTMLValue(htmlStr, fieldname)
	rpt.setPageStyle(options)
	rpt.writeHTML(lineHt, htmlStr)
	rpt.Pdf.SetXY(rpt.LeftMargin, rpt.Pdf.GetY()+lineHt+_padding)
}
