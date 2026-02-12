/*
Go PDF report generation library

  - Fully declarative: can be easily modified and used for relative layout (no need to specify the x and y coordinates)
  - Powerful layout engine: row, datagrid, column, cell, image, separator, html, barcode, hline, vgap elements
  - Creating a PDF from a JSON template or GO language code

Quick start: example/example.go
*/
package report

import (
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"image/color"
	"os"
	"path"
	"strings"
)

//go:embed fonts
var Fonts embed.FS

func (rpt *Report) addToXML(section string, values []string) {
	if values[0] != "label" {
		switch section {
		case "header":
			rpt.xmlHeader += fmt.Sprintf("\n    <%s><![CDATA[%s]]></%s>", values[0], values[1], values[2])
		case "details":
			if len(values) == 1 {
				rpt.xmlDetails += fmt.Sprintf("\n    <%s>", values[0])
				return
			}
			rpt.xmlDetails += fmt.Sprintf("\n    <%s><![CDATA[%s]]></%s>", values[0], values[1], values[2])
		case "footer":
			rpt.xmlDetails += fmt.Sprintf("\n    <%s_footer"+"><![CDATA[%s]]></%s_footer"+">",
				values[0], values[1], values[2])
		}
	}
}

func (rpt *Report) createElement(section string, element interface{}) {
	switch v := element.(type) {
	case *Row:
		if v.Visible != "" {
			if _, found := rpt.data[v.Visible]; found {
				if srows, valid := rpt.data[v.Visible].([]SM); !valid || len(srows) == 0 {
					return
				}
			}
		}
		rpt.createRow(section, v, false)
	case *VGap:
		if v.PageBreak {
			rpt.addPage()
		}
		if rpt.checkPageBreak(v.Height) {
			rpt.addPage()
		}
		rpt.Pdf.Ln(v.Height)
	case *HLine:
		rpt.createLine(v, false)
	case *HTML:
		rpt.createHTML(v)
	case *Datagrid:
		rpt.createDatagrid(v, false)
	}
}

func (rpt *Report) setFont() bool {
	fontFile := map[string]func(family string) string{
		"REGULAR":    func(family string) string { return family + "-Regular.ttf" },
		"BOLD":       func(family string) string { return family + "-Bold.ttf" },
		"ITALIC":     func(family string) string { return family + "-Italic.ttf" },
		"BOLDITALIC": func(family string) string { return family + "-BoldItalic.ttf" },
	}

	checkCustom := func() bool {
		if _, err := os.Stat(path.Join(rpt.fontDir, fontFile["REGULAR"](rpt.FontFamily))); err != nil {
			return false
		}
		if _, err := os.Stat(path.Join(rpt.fontDir, fontFile["BOLD"](rpt.FontFamily))); err != nil {
			return false
		}
		if _, err := os.Stat(path.Join(rpt.fontDir, fontFile["ITALIC"](rpt.FontFamily))); err != nil {
			return false
		}
		if _, err := os.Stat(path.Join(rpt.fontDir, fontFile["BOLDITALIC"](rpt.FontFamily))); err != nil {
			return false
		}
		return true
	}

	custom := false
	if rpt.FontFamily != _fontFamily && rpt.fontDir != "" {
		custom = checkCustom()
	}

	if custom {
		rpt.Pdf.AddFont(rpt.FontFamily, "", path.Join(rpt.fontDir, fontFile["REGULAR"](rpt.FontFamily)), nil)
		rpt.Pdf.AddFont(rpt.FontFamily, "B", path.Join(rpt.fontDir, fontFile["BOLD"](rpt.FontFamily)), nil)
		rpt.Pdf.AddFont(rpt.FontFamily, "I", path.Join(rpt.fontDir, fontFile["ITALIC"](rpt.FontFamily)), nil)
		rpt.Pdf.AddFont(rpt.FontFamily, "BI", path.Join(rpt.fontDir, fontFile["BOLDITALIC"](rpt.FontFamily)), nil)
	} else {
		rpt.FontFamily = _fontFamily
		font, _ := Fonts.Open(path.Join("fonts", fontFile["REGULAR"](_fontFamily)))
		rpt.Pdf.AddFont(rpt.FontFamily, "", "", font)
		font, _ = Fonts.Open(path.Join("fonts", fontFile["BOLD"](_fontFamily)))
		rpt.Pdf.AddFont(rpt.FontFamily, "B", "", font)
		font, _ = Fonts.Open(path.Join("fonts", fontFile["ITALIC"](_fontFamily)))
		rpt.Pdf.AddFont(rpt.FontFamily, "I", "", font)
		font, _ = Fonts.Open(path.Join("fonts", fontFile["BOLDITALIC"](_fontFamily)))
		rpt.Pdf.AddFont(rpt.FontFamily, "BI", "", font)
	}
	return true
}

/*
New returns a pointer to a new Report instance. Options:
  - orientation - Optional. Default value:"P" Values: "P","portrait","L","landscape".
  - format - Optional. Defaut value: "A4" Values: "A3","A4","A5","letter","legal".
  - fontFamily - Optional Default: Cabin
  - fontDir - Optional Default: ""

Example:

	rpt := report.New("P", "A4")
*/
func New(options ...string) (rpt *Report) {
	rpt = new(Report)
	if len(options) > 0 {
		rpt.orientation = rpt.parseValue("Orientation", options[0]).(string)
	} else {
		rpt.orientation = rpt.parseValue("Orientation", _orientation).(string)
	}
	if len(options) > 1 {
		rpt.format = rpt.parseValue("Format", options[1]).(string)
	} else {
		rpt.format = rpt.parseValue("Format", _format).(string)
	}
	rpt.FontFamily = _fontFamily
	if len(options) > 2 {
		rpt.FontFamily = options[2]
	}
	if len(options) > 3 {
		rpt.fontDir = options[3]
	}

	rpt.Title = _title
	rpt.LeftMargin = _margin
	rpt.TopMargin = _margin
	rpt.RightMargin = _margin
	rpt.BottomMargin = _margin
	rpt.FontStyle = _fontStyle
	rpt.FontSize = _fontSize
	rpt.TextColor = color.RGBA{_textColor, _textColor, _textColor, 0}
	rpt.BorderColor = color.RGBA{_borderColor, _borderColor, _borderColor, 0}
	rpt.BackgroundColor = color.RGBA{_backgroundColor, _backgroundColor, _backgroundColor, 0}
	rpt.header = make([]PageItem, 0)
	rpt.details = make([]PageItem, 0)
	rpt.footer = make([]PageItem, 0)
	rpt.data = make(IM)

	rpt.Pdf = generators[_generator]
	rpt.Pdf.Init(rpt)
	rpt.setFont()

	return
}

// CreateReport - the report template processing, databind replacement.
func (rpt *Report) CreateReport() bool {
	rpt.Pdf.SetProperties(rpt)
	rpt.setPageStyle(make(IM))
	rpt.footerHeight = rpt.getFooterHeight()
	rpt.addPage()
	for index := 0; index < len(rpt.details); index++ {
		rpt.createElement("details", rpt.details[index].Item)
	}
	return true
}

func (rpt *Report) getJSONElements(edata interface{}) (el PageItem, err error) {
	for eName, eValue := range edata.(IM) {
		if el, err = rpt.getPageItem(eName); err != nil {
			return el, err
		}
		for ekey, eValueData := range eValue.(IM) {
			if ekey == "columns" {
				for colIndex := 0; colIndex < len(eValueData.([]interface{})); colIndex++ {
					coldata := eValueData.([]interface{})[colIndex]
					for cName, cValue := range coldata.(IM) {
						switch cName {
						case "cell", "image", "barcode", "separator", "column":
							el2, _ := rpt.getPageItem(cName)
							for ckey, cValueData := range cValue.(IM) {
								if err := el2.setPageItem(ckey, rpt.parseValue(propMap[strings.ToLower(ckey)], cValueData)); err != nil {
									return el, err
								}
							}
							if eName == "row" {
								el.Item.(*Row).Columns = append(el.Item.(*Row).Columns, el2)
							} else {
								el.Item.(*Datagrid).Columns = append(el.Item.(*Datagrid).Columns, el2)
							}
						default:
							return el, errors.New(invalidErr("Columns", cName))
						}
					}
				}
			} else {
				if err := el.setPageItem(ekey, rpt.parseValue(propMap[strings.ToLower(ekey)], eValueData)); err != nil {
					return el, err
				}
			}
		}
	}
	return el, nil
}

// LoadJSONDefinition load to the report an JSON definition.
func (rpt *Report) LoadJSONDefinition(jsonString string) error {
	if jsonString == "" {
		return errors.New("missing JSON")
	}
	var jsonData IM
	if err := ConvertFromByte([]byte(jsonString), &jsonData); err != nil {
		return err
	}
	if err := rpt.loadJSONReportValues(jsonData); err != nil {
		return err
	}
	if err := rpt.loadJSONSection(jsonData, "header", &rpt.header); err != nil {
		return err
	}
	if err := rpt.loadJSONSection(jsonData, "details", &rpt.details); err != nil {
		return err
	}
	if err := rpt.loadJSONSection(jsonData, "footer", &rpt.footer); err != nil {
		return err
	}
	return rpt.loadJSONData(jsonData)
}

func (rpt *Report) loadJSONReportValues(jsonData IM) error {
	report, found := jsonData["report"]
	if !found {
		return nil
	}
	for valueKey, valueData := range report.(IM) {
		if err := rpt.SetReportValue(valueKey, rpt.parseValue(propMap[strings.ToLower(valueKey)], valueData)); err != nil {
			return err
		}
	}
	return nil
}

func (rpt *Report) loadJSONSection(jsonData IM, section string, items *[]PageItem) error {
	sectionData, found := jsonData[section]
	if !found {
		return nil
	}
	for index := 0; index < len(sectionData.([]interface{})); index++ {
		el, err := rpt.getJSONElements(sectionData.([]interface{})[index])
		if err != nil {
			return err
		}
		*items = append(*items, el)
	}
	return nil
}

func (rpt *Report) loadJSONData(jsonData IM) error {
	data, found := jsonData["data"]
	if !found {
		return nil
	}
	for dKey, dValue := range data.(IM) {
		if err := rpt.loadJSONDataValue(dKey, dValue); err != nil {
			return err
		}
	}
	return nil
}

func (rpt *Report) loadJSONDataValue(dKey string, dValue interface{}) error {
	switch dValue.(type) {
	case []interface{}:
		rpt.data[dKey] = make([]SM, 0)
		for index := 0; index < len(dValue.([]interface{})); index++ {
			jRow := dValue.([]interface{})[index]
			dRow := SM{}
			for key, value := range jRow.(IM) {
				dRow[key] = ToString(value, "")
			}
			rpt.data[dKey] = append(rpt.data[dKey].([]SM), dRow)
		}
	case IM:
		rpt.data[dKey] = SM{}
		for key, value := range dValue.(IM) {
			rpt.data[dKey].(SM)[key] = ToString(value, "")
		}
	case string, SM, []SM:
		rpt.data[dKey] = dValue
	default:
		return errors.New("valid data types: string, map[string][string], []map[string][string] ")
	}
	return nil
}

/*
AppendElement - Append an element in the template.
  - parent - Optional. The parent elemnt. Values: "header","details","footer" or result value (row, datagrid) Default value: "details"
  - ename - Optional. An Element type: "row", "datagrid", "vgap", "hline", "html", "column", "cell", "image", "separator", "barcode". Default value: "row"
  - values - Optional. Element attributes

Example:

	row_data := rpt.AppendElement("header", "row", map[string]interface{}{"height": 10})
	rpt.AppendElement(row_data, "image", map[string]interface{}{"src": "test/logo.jpg"})
*/
func (rpt *Report) AppendElement(options ...interface{}) (*[]PageItem, error) {
	el, _ := rpt.getPageItem("row")
	parent, err := rpt.resolveAppendParent(options, &el)
	if err != nil {
		return nil, err
	}
	if len(options) > 2 {
		if err := rpt.applyAppendElementValues(&el, options[2]); err != nil {
			return nil, err
		}
	}
	*parent = append(*parent, el)
	switch el.ItemType {
	case "row":
		return &el.Item.(*Row).Columns, nil
	case "datagrid":
		return &el.Item.(*Datagrid).Columns, nil
	}
	return parent, nil
}

func (rpt *Report) resolveAppendParent(options []interface{}, el *PageItem) (*[]PageItem, error) {
	parent := &rpt.details
	validHeaders := []string{"row", "vgap", "hline"}
	validDetails := []string{"row", "vgap", "hline", "html", "datagrid"}
	validFooter := []string{"row", "vgap", "hline"}
	validColumns := []string{"cell", "image", "barcode", "separator", "column"}

	if len(options) == 0 {
		return parent, nil
	}
	switch opt := options[0].(type) {
	case string:
		return rpt.resolveAppendParentString(opt, options, el, validHeaders, validDetails, validFooter)
	case *[]PageItem:
		parent = opt
		if len(options) > 1 {
			ename := ToString(options[1], "")
			if !Contains(validColumns, ename) {
				return nil, errors.New(invalidErr("columns", ename))
			}
			*el, _ = rpt.getPageItem(ename)
		}
		return parent, nil
	default:
		return nil, errors.New("valid parent values: 'header','details','footer' (string) or Columns of Row and Datagrid (*[]PageItem)")
	}
}

var sectionErrMap = map[string]string{"header": "Header", "details": "Details", "footer": "Footer"}

func (rpt *Report) resolveAppendParentString(section string, options []interface{}, el *PageItem, validH, validD, validF []string) (*[]PageItem, error) {
	var parent *[]PageItem
	var valid []string
	switch section {
	case "header":
		parent, valid = &rpt.header, validH
	case "details":
		parent, valid = &rpt.details, validD
	case "footer":
		parent, valid = &rpt.footer, validF
	case "body":
		parent, valid = &rpt.details, validD // body alias for details (not in sectionErrMap)
	default:
		return &rpt.details, nil
	}
	if len(options) > 1 {
		ename := ToString(options[1], "")
		if !Contains(valid, ename) {
			errSection := sectionErrMap[section]
			if errSection == "" {
				errSection = section
			}
			return nil, errors.New(invalidErr(errSection, ename))
		}
		*el, _ = rpt.getPageItem(ename)
	}
	return parent, nil
}

func (rpt *Report) applyAppendElementValues(el *PageItem, values interface{}) error {
	im, ok := values.(IM)
	if !ok {
		return errors.New("valid values type: map[string]interface{}")
	}
	for key, value := range im {
		if err := el.setPageItem(key, rpt.parseValue(propMap[strings.ToLower(key)], value)); err != nil {
			return err
		}
	}
	return nil
}

// SetReportValue - You can set the Report properties safely and type independent.
func (rpt *Report) SetReportValue(fieldname string, value interface{}) error {
	vmap := map[string]func(value interface{}){
		"Title": func(value interface{}) {
			rpt.Title = ToString(value, rpt.Title)
		},
		"Author": func(value interface{}) {
			rpt.Author = ToString(value, rpt.Author)
		},
		"Creator": func(value interface{}) {
			rpt.Creator = ToString(value, rpt.Creator)
		},
		"Subject": func(value interface{}) {
			rpt.Subject = ToString(value, rpt.Subject)
		},
		"Keywords": func(value interface{}) {
			rpt.Keywords = ToString(value, rpt.Keywords)
		},
		"LeftMargin": func(value interface{}) {
			rpt.LeftMargin = ToFloat(value, rpt.LeftMargin) * _mmPt
		},
		"TopMargin": func(value interface{}) {
			rpt.TopMargin = ToFloat(value, rpt.TopMargin) * _mmPt
		},
		"RightMargin": func(value interface{}) {
			rpt.RightMargin = ToFloat(value, rpt.RightMargin) * _mmPt
		},
		"BottomMargin": func(value interface{}) {
			rpt.BottomMargin = ToFloat(value, rpt.BottomMargin) * _mmPt
		},
		"FontStyle": func(value interface{}) {
			rpt.FontStyle = rpt.parseValue("FontStyle", value).(string)
		},
		"FontSize": func(value interface{}) {
			rpt.FontSize = ToFloat(value, rpt.FontSize)
		},
		"TextColor": func(value interface{}) {
			rpt.TextColor = ToRGBA(value, rpt.TextColor)
		},
		"BorderColor": func(value interface{}) {
			rpt.BorderColor = ToRGBA(value, rpt.BorderColor)
		},
		"BackgroundColor": func(value interface{}) {
			rpt.BackgroundColor = ToRGBA(value, rpt.BackgroundColor)
		},
		"ImagePath": func(value interface{}) {
			rpt.ImagePath = ToString(value, rpt.ImagePath)
		},
	}

	if fn, found := vmap[propMap[strings.ToLower(fieldname)]]; found {
		fn(value)
		return nil
	}
	return fmt.Errorf("missing report fieldname: %s", fieldname)
}

/*
SetData - Set the template data. Parameters:
  - key - string
  - value - interface{} Valid interface type: string or dictonary (map[string]string) or record list ([]map[string]string)

Example:

	rpt.SetData("items_footer", map[string]string{"items_total": "3 703 680"})
*/
func (rpt *Report) SetData(key string, value interface{}) (bool, error) {
	if _, found := rpt.data[key]; found {
		switch rpt.data[key].(type) {
		case SM:
			switch value.(type) {
			case SM:
				for valueKey, valueData := range value.(SM) {
					rpt.data[key].(SM)[valueKey] = valueData
				}
				return true, nil
			}
		}
	}
	switch value.(type) {
	case string, SM, []SM:
		rpt.data[key] = value
	default:
		return false, errors.New("valid value types: string, map[string][string], []map[string]string")
	}
	return true, nil
}

// Save2DataURLString creates a base64 data URI scheme.
func (rpt *Report) Save2DataURLString(filename string) (string, error) {
	pdf, err := rpt.Save2Pdf()
	if err != nil {
		return "", err
	}
	pdfStr := base64.URLEncoding.EncodeToString([]byte(pdf))
	if filename != "" {
		filename = "filename=" + filename + ";"
	}
	return "data:application/pdf;" + filename + "base64," + pdfStr, nil
}

// Save2Pdf creates a PDF output.
func (rpt *Report) Save2Pdf() ([]byte, error) {
	return rpt.Pdf.Save2Pdf()
}

// Save2PdfFile creates or truncates the file specified by fileStr and
// writes the PDF document to it.
func (rpt *Report) Save2PdfFile(filename string) error {
	return rpt.Pdf.Save2PdfFile(filename)
}

// Save2Xml creates an XML output. Only the values of cells and datagrid
// rows from header and details. The node name of the cell name (except when
// name="label"), or datagrid name/column fieldname.
func (rpt *Report) Save2Xml() string {
	return fmt.Sprintf("<data>%s%s\n</data>", rpt.xmlHeader, rpt.xmlDetails)
}
