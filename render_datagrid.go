package report

import (
	"strconv"
	"strings"
)

func (rpt *Report) createDatagrid(gridElement *Datagrid, virtual bool) bool {
	if len(gridElement.Columns) == 0 {
		return false
	}
	rows, valid := rpt.data[gridElement.Databind].([]SM)
	if !valid || len(rows) == 0 {
		return false
	}
	gridOptions, headerOptions, footerOptions := rpt.buildDatagridOptions(gridElement, virtual)
	footers := rpt.buildDatagridColumns(gridElement, gridOptions, headerOptions)
	if footers == nil {
		return false
	}
	if !headerOptions["merge"].(bool) {
		rpt.createGridHeader(headerOptions)
	}
	rpt.renderDatagridRows(rows, gridOptions, headerOptions)
	rpt.renderDatagridFooters(gridOptions, headerOptions, footerOptions, footers)
	return true
}

func (rpt *Report) buildDatagridOptions(gridElement *Datagrid, virtual bool) (IM, IM, IM) {
	pageWidth, _ := rpt.Pdf.GetPageSize()
	nwidth := pageWidth - rpt.RightMargin - rpt.LeftMargin
	gridWidth := rpt.computeGridWidth(gridElement, nwidth)

	gridOptions := IM{
		"xname":           ToString(gridElement.Name, "items"),
		"border":          ToString(gridElement.Border, "1"),
		"fontFamily":      rpt.FontFamily,
		"fontStyle":       rpt.FontStyle,
		"fontSize":        gridElement.FontSize,
		"textColor":       gridElement.TextColor,
		"borderColor":     gridElement.BorderColor,
		"backgroundColor": gridElement.BackgroundColor,
		"virtual":         virtual,
	}
	headerOptions := IM{
		"fontSize": gridElement.FontSize, "textColor": gridElement.TextColor,
		"borderColor": gridElement.BorderColor, "border": ToString(gridElement.Border, "1"),
		"backgroundColor": ToRGBA(gridElement.HeaderBackground, gridElement.BackgroundColor),
		"merge":           ToBoolean(gridElement.Merge, false),
		"fontFamily":      rpt.FontFamily, "fontStyle": "B", "text": "", "height": float64(0),
		"columns":      make([]IM, 0),
		"columnsWidth": float64(0), "gridWidth": gridWidth, "multiline": false, "virtual": virtual,
		"extend": gridWidth == nwidth,
	}
	gridOptions["extend"] = headerOptions["extend"]
	footerOptions := IM{
		"fontSize": gridElement.FontSize, "textColor": gridElement.TextColor,
		"borderColor": gridElement.BorderColor, "border": ToString(gridElement.Border, "1"),
		"fontFamily": rpt.FontFamily, "fontStyle": "B", "text": "", "height": float64(0),
		"backgroundColor": ToRGBA(gridElement.FooterBackground, gridElement.BackgroundColor),
		"extend":          headerOptions["extend"], "multiline": false, "virtual": virtual,
	}
	return gridOptions, headerOptions, footerOptions
}

func (rpt *Report) computeGridWidth(gridElement *Datagrid, nwidth float64) float64 {
	gridElement.Width = ToString(gridElement.Width, "100%")
	if strings.HasSuffix(gridElement.Width, "%") {
		return nwidth * ToFloat(strings.Replace(gridElement.Width, "%", "", -1), 0) / 100
	}
	w := ToFloat(gridElement.Width, 0)
	if w > nwidth {
		return nwidth
	}
	return w
}

func (rpt *Report) buildDatagridColumns(gridElement *Datagrid, gridOptions, headerOptions IM) []IM {
	footers := make([]IM, 0)
	footerWidth := float64(0)
	xCol := rpt.LeftMargin
	lnWidth := headerOptions["gridWidth"].(float64)
	columns := headerOptions["columns"].([]IM)

	for index := 0; index < len(gridElement.Columns); index++ {
		column := gridElement.Columns[index].Item.(*Column)
		if headerOptions["columnsWidth"].(float64) >= headerOptions["gridWidth"].(float64) {
			return nil
		}
		columnOptions := rpt.buildColumnOptions(column, gridOptions, headerOptions, index, lnWidth, len(gridElement.Columns), &xCol)
		headerOptions["columnsWidth"] = headerOptions["columnsWidth"].(float64) + columnOptions["columnWidth"].(float64)
		lnWidth = lnWidth - columnOptions["columnWidth"].(float64)
		if !headerOptions["merge"].(bool) {
			cheight := rpt.getCellHeight(columnOptions["label"].(string), columnOptions["columnWidth"].(float64), columnOptions)
			if cheight > headerOptions["height"].(float64) {
				headerOptions["height"] = cheight
			}
		}
		footers, footerWidth = rpt.appendDatagridFooter(footers, column, columnOptions, footerWidth)
		columnOptions["ln"] = 0
		if len(gridElement.Columns)-1 == index {
			columnOptions["ln"] = 1
		}
		columns = append(columns, columnOptions)
	}
	headerOptions["columns"] = columns
	return footers
}

func (rpt *Report) buildColumnOptions(column *Column, gridOptions, headerOptions IM, index int, lnWidth float64, colCount int, xCol *float64) IM {
	label := rpt.setValue(column.Label)
	columnOptions := IM{
		"fontFamily": gridOptions["fontFamily"], "fontStyle": gridOptions["fontStyle"],
		"fontSize": gridOptions["fontSize"], "textColor": gridOptions["textColor"],
		"borderColor": gridOptions["borderColor"], "border": gridOptions["border"],
		"fieldname": column.Fieldname, "multiline": true,
		"label":       label,
		"columnWidth": rpt.computeColumnWidth(column, headerOptions, index, colCount, lnWidth, label),
		"headerAlign": ToString(column.HeaderAlign, "L"),
		"align":       ToString(column.Align, "L"),
	}
	if !headerOptions["merge"].(bool) {
		columnOptions["xCol"] = *xCol
		*xCol += columnOptions["columnWidth"].(float64)
		if headerOptions["columnsWidth"].(float64)+columnOptions["columnWidth"].(float64) >= headerOptions["gridWidth"].(float64) {
			columnOptions["columnWidth"] = headerOptions["gridWidth"].(float64) - headerOptions["columnsWidth"].(float64)
		}
	}
	return columnOptions
}

func (rpt *Report) computeColumnWidth(column *Column, headerOptions IM, index, colCount int, lnWidth float64, label string) float64 {
	columnWidth := ToString(column.Width, "")
	if columnWidth != "" {
		if strings.HasSuffix(columnWidth, "%") {
			return headerOptions["gridWidth"].(float64) * ToFloat(strings.Replace(columnWidth, "%", "", -1), 0) / 100
		}
		return ToFloat(columnWidth, 0)
	}
	if colCount-1 == index {
		return lnWidth
	}
	return rpt.Pdf.GetTextWidth(label) + _padding
}

func (rpt *Report) appendDatagridFooter(footers []IM, column *Column, columnOptions IM, footerWidth float64) ([]IM, float64) {
	footerValue := rpt.setValue(ToString(column.Footer, ""))
	footerAlign := ToString(column.FooterAlign, "L")
	if footerValue != "" {
		cw := columnOptions["columnWidth"].(float64)
		if len(footers) == 0 {
			return append(footers, IM{"text": footerValue, "align": footerAlign, "columnWidth": footerWidth + cw}), 0
		}
		footers[len(footers)-1]["columnWidth"] = footers[len(footers)-1]["columnWidth"].(float64) + footerWidth
		return append(footers, IM{"text": footerValue, "align": footerAlign, "columnWidth": cw}), 0
	}
	return footers, footerWidth + columnOptions["columnWidth"].(float64)
}

func (rpt *Report) setDatagridRowCellText(column IM, row SM, rowIndex int) {
	if column["fieldname"] == "counter" {
		column["text"] = strconv.Itoa(rowIndex + 1)
		return
	}
	if value, found := row[column["fieldname"].(string)]; found {
		column["text"] = value
	} else {
		column["text"] = ""
	}
}

func (rpt *Report) renderDatagridRows(rows []SM, gridOptions, headerOptions IM) {
	columns := headerOptions["columns"].([]IM)
	merge := headerOptions["merge"].(bool)

	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]
		rpt.addToXML("details", []string{gridOptions["xname"].(string)})
		gridOptions["height"] = float64(0)
		gridOptions["text"] = ""
		for colIndex := 0; colIndex < len(columns); colIndex++ {
			column := columns[colIndex]
			rpt.setDatagridRowCellText(column, row, rowIndex)
			if !merge {
				cheight := rpt.getCellHeight(column["text"].(string), column["columnWidth"].(float64), gridOptions)
				if cheight > gridOptions["height"].(float64) {
					gridOptions["height"] = cheight
				}
			} else {
				gridOptions["text"] = gridOptions["text"].(string) + " " + column["text"].(string)
				rpt.addToXML("details", []string{column["fieldname"].(string), column["text"].(string), column["fieldname"].(string)})
			}
		}
		if rpt.checkPageBreak(gridOptions["height"].(float64)) {
			rpt.addPage()
			if !merge {
				rpt.createGridHeader(headerOptions)
			}
		}
		rpt.renderDatagridRowCells(gridOptions, headerOptions, row)
		rpt.addToXML("details", []string{gridOptions["xname"].(string)})
	}
}

func (rpt *Report) renderDatagridRowCells(gridOptions, headerOptions IM, _ SM) {
	merge := headerOptions["merge"].(bool)
	columns := headerOptions["columns"].([]IM)
	if !merge {
		for colIndex := 0; colIndex < len(columns); colIndex++ {
			column := columns[colIndex]
			gridOptions["text"] = column["text"]
			gridOptions["columnWidth"] = column["columnWidth"]
			gridOptions["xCol"] = column["xCol"]
			gridOptions["align"] = column["align"]
			gridOptions["ln"] = column["ln"]
			gridOptions["multiline"] = column["multiline"]
			rpt.createCell(gridOptions)
			rpt.addToXML("details", []string{column["fieldname"].(string), column["text"].(string), column["fieldname"].(string)})
		}
		return
	}
	gridOptions["text"] = strings.Trim(gridOptions["text"].(string), " ")
	gridOptions["columnWidth"] = float64(0)
	gridOptions["ln"] = 1
	gridOptions["multiline"] = true
	gridOptions["height"] = float64(0)
	rpt.createCell(gridOptions)
}

func (rpt *Report) renderDatagridFooters(gridOptions, headerOptions, footerOptions IM, footers []IM) {
	if headerOptions["merge"].(bool) {
		return
	}
	for colIndex := 0; colIndex < len(footers); colIndex++ {
		column := footers[colIndex]
		cheight := rpt.getCellHeight(column["text"].(string), column["columnWidth"].(float64), footerOptions)
		if cheight > footerOptions["height"].(float64) {
			footerOptions["height"] = cheight
		}
	}
	for colIndex := 0; colIndex < len(footers); colIndex++ {
		column := footers[colIndex]
		footerOptions["text"] = column["text"]
		footerOptions["columnWidth"] = column["columnWidth"]
		footerOptions["align"] = column["align"]
		footerOptions["ln"] = 0
		if len(footers)-1 == colIndex {
			footerOptions["ln"] = 1
		}
		rpt.createCell(footerOptions)
		rpt.addToXML("footer", []string{gridOptions["xname"].(string), column["text"].(string), gridOptions["xname"].(string)})
	}
}
