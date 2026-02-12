package report

import (
	"regexp"
	"strconv"
	"strings"
)

func (rpt *Report) setHTMLValue(value, fieldname string) string {
	r, _ := regexp.Compile(_regValue)
	if r.MatchString(value) {
		valueKey := r.FindString(value)
		valueF := rpt.setValue(valueKey)
		value = strings.Replace(value, valueKey, valueF, -1)
		rpt.addToXML("details", []string{fieldname, valueF, fieldname})
		if r.MatchString(value) {
			value = rpt.setHTMLValue(value, fieldname)
		}
		return value
	}
	return value
}

func (rpt *Report) setValue(value string) string {
	var getValue = func(valueGet string) string {
		if matched, _ := regexp.MatchString("{{page}}", valueGet); matched {
			valueGet = strings.ReplaceAll(valueGet, "{{page}}", strconv.Itoa(rpt.Pdf.PageNo()))
		}
		dbv := strings.Split(valueGet, ".")
		storeData, isData := rpt.data[dbv[0]]
		if isData {
			if data, valid := storeData.([]SM); valid {
				if len(dbv) > 2 {
					row := ToInteger(dbv[1], 0)
					if len(data) > int(row) {
						rowValue, isData := data[row][dbv[2]]
						if isData {
							return rowValue
						}
					}
					return ""
				}
				return valueGet
			}
			if data, valid := storeData.(SM); valid {
				if len(dbv) > 1 {
					dictValue, isData := data[dbv[1]]
					if isData {
						return dictValue
					}
					return ""
				}
				return valueGet
			}
			if data, valid := storeData.(string); valid {
				return data
			}
		}
		return valueGet
	}
	if matched, _ := regexp.MatchString(_regValue, value); matched {
		valueSet := value[strings.Index(value, "={{")+3 : strings.Index(value, "}}")]
		value = strings.Replace(value, "={{"+valueSet+"}}", getValue(valueSet), strings.Index(value, "}}")+2)
		if matched, _ := regexp.MatchString(_regValue, value); matched {
			return rpt.setValue(value)
		}
		return value
	}
	return getValue(value)
}

func (rpt *Report) parseValue(vname string, value interface{}) interface{} {
	parseStringMap := func(value interface{}, defValue string) interface{} {
		svalue := ToString(value, defValue)
		valid := SM{
			"R": "R", "C": "C", "L": "L", "J": "J", "left": "L", "center": "C", "right": "R", "justify": "L",
			"B": "B", "I": "I", "BI": "BI", "IB": "IB", "bold": "B", "italic": "I", "bolditalic": "BI", "normal": "",
			"p": "p", "l": "l", "portrait": "p", "landscape": "l",
			"a3": "a3", "a4": "a4", "a5": "a5", "letter": "letter", "legal": "legal",
		}
		if _, found := valid[svalue]; found {
			return valid[svalue]
		}
		return defValue
	}
	checkValue := map[string]func(value interface{}) interface{}{
		"Format": func(value interface{}) interface{} {
			return parseStringMap(value, _format)
		},
		"Orientation": func(value interface{}) interface{} {
			return parseStringMap(value, _orientation)
		},
		"FontSize": func(value interface{}) interface{} {
			return ToFloat(value, rpt.FontSize)
		},
		"Height": func(value interface{}) interface{} {
			return ToFloat(value, 0) * _mmPt
		},
		"Gap": func(value interface{}) interface{} {
			return ToFloat(value, 0) * _mmPt
		},
		"HGap": func(value interface{}) interface{} {
			return ToFloat(value, 0) * _mmPt
		},
		"Width": func(value interface{}) interface{} {
			switch v := value.(type) {
			case string:
				if strings.HasSuffix(v, "%") {
					svalue := strings.Replace(v, "%", "", -1)
					ivalue := ToInteger(svalue, 0)
					if ivalue > 100 {
						ivalue = 100
					}
					return ToString(ivalue, "0") + "%"
				}
				return ToString(ToFloat(value, 0)*_mmPt, "0")
			default:
				return ToString(ToFloat(value, 0)*_mmPt, "0")
			}
		},
		"Merge": func(value interface{}) interface{} {
			return ToBoolean(value, false)
		},
		"VisibleValue": func(value interface{}) interface{} {
			return ToBoolean(value, false)
		},
		"Multiline": func(value interface{}) interface{} {
			return ToBoolean(value, false)
		},
		"TextColor": func(value interface{}) interface{} {
			return ToRGBA(value, rpt.TextColor)
		},
		"Border": func(value interface{}) interface{} {
			borderStr := ToString(value, "")
			if borderStr == "1" || strings.Contains(borderStr, "T") || strings.Contains(borderStr, "L") ||
				strings.Contains(borderStr, "R") || strings.Contains(borderStr, "B") {
				return value
			}
			return ""
		},
		"BorderColor": func(value interface{}) interface{} {
			return ToRGBA(value, rpt.BorderColor)
		},
		"BackgroundColor": func(value interface{}) interface{} {
			return ToRGBA(value, rpt.BackgroundColor)
		},
		"FooterBackground": func(value interface{}) interface{} {
			return ToRGBA(value, rpt.BackgroundColor)
		},
		"HeaderBackground": func(value interface{}) interface{} {
			return ToRGBA(value, rpt.BackgroundColor)
		},
		"FontStyle": func(value interface{}) interface{} {
			return parseStringMap(value, _fontStyle)
		},
		"Align": func(value interface{}) interface{} {
			return parseStringMap(value, _align)
		},
		"HeaderAlign": func(value interface{}) interface{} {
			return parseStringMap(value, _align)
		},
		"FooterAlign": func(value interface{}) interface{} {
			return parseStringMap(value, _align)
		},
	}

	if fn, found := checkValue[vname]; found {
		return fn(value)
	}
	return value
}
