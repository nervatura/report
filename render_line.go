package report

import "strings"

func (rpt *Report) createLine(v *HLine, virtual bool) {
	width := float64(0)
	pageWidth, _ := rpt.Pdf.GetPageSize()
	if strings.HasSuffix(v.Width, "%") {
		width = ToFloat(strings.Replace(v.Width, "%", "", -1), width) / 100
		width = (pageWidth - rpt.LeftMargin - rpt.RightMargin) * width
	} else {
		width = ToFloat(v.Width, width)
	}
	if width == 0 {
		width = pageWidth - rpt.LeftMargin - rpt.RightMargin
	}
	options := IM{"borderColor": v.BorderColor}
	rpt.setPageStyle(options)
	if !virtual {
		rpt.Pdf.Line(rpt.Pdf.GetX(), rpt.Pdf.GetY(), rpt.Pdf.GetX()+width, rpt.Pdf.GetY())
		if v.Gap > 0 {
			rpt.Pdf.Line(rpt.Pdf.GetX(), rpt.Pdf.GetY()+v.Gap, rpt.Pdf.GetX()+width, rpt.Pdf.GetY()+v.Gap)
		}
	}
}
