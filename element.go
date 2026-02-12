package report

import (
	"errors"
	"strings"
)

// ElementPropertySetter sets a property on an element. Used by the element registry.
type ElementPropertySetter func(item interface{}, value interface{}) error

var generators = make(map[string]Generator)

// elementSetters maps element type -> property name -> setter.
var elementSetters map[string]map[string]ElementPropertySetter

func init() {
	elementSetters = map[string]map[string]ElementPropertySetter{
		"row": {
			"Height":  func(item, value interface{}) error { item.(*Row).Height = ToFloat(value, 0); return nil },
			"HGap":    func(item, value interface{}) error { item.(*Row).HGap = ToFloat(value, 0); return nil },
			"Visible": func(item, value interface{}) error { item.(*Row).Visible = ToString(value, ""); return nil },
		},
		"cell": {
			"Name":            func(item, value interface{}) error { item.(*Cell).Name = ToString(value, ""); return nil },
			"Value":           func(item, value interface{}) error { item.(*Cell).Value = ToString(value, ""); return nil },
			"Width":           func(item, value interface{}) error { item.(*Cell).Width = ToString(value, ""); return nil },
			"Border":          func(item, value interface{}) error { item.(*Cell).Border = ToString(value, ""); return nil },
			"Align":           func(item, value interface{}) error { item.(*Cell).Align = ToString(value, "L"); return nil },
			"Multiline":       func(item, value interface{}) error { item.(*Cell).Multiline = ToBoolean(value, false); return nil },
			"FontStyle":       func(item, value interface{}) error { item.(*Cell).FontStyle = ToString(value, ""); return nil },
			"FontSize":        func(item, value interface{}) error { item.(*Cell).FontSize = ToFloat(value, item.(*Cell).FontSize); return nil },
			"TextColor":       func(item, value interface{}) error { item.(*Cell).TextColor = ToRGBA(value, item.(*Cell).TextColor); return nil },
			"BorderColor":     func(item, value interface{}) error { item.(*Cell).BorderColor = ToRGBA(value, item.(*Cell).BorderColor); return nil },
			"BackgroundColor": func(item, value interface{}) error { item.(*Cell).BackgroundColor = ToRGBA(value, item.(*Cell).BackgroundColor); return nil },
		},
		"image": {
			"Src":    func(item, value interface{}) error { item.(*Image).Src = ToString(value, ""); return nil },
			"Height": func(item, value interface{}) error { item.(*Image).Height = ToFloat(value, 0); return nil },
		},
		"barcode": {
			"CodeType":     func(item, value interface{}) error { item.(*Barcode).CodeType = ToString(value, ""); return nil },
			"Value":        func(item, value interface{}) error { item.(*Barcode).Value = ToString(value, ""); return nil },
			"VisibleValue": func(item, value interface{}) error { item.(*Barcode).VisibleValue = ToBoolean(value, false); return nil },
			"Width":        func(item, value interface{}) error { item.(*Barcode).Width = ToFloat(value, 0); return nil },
			"Height":       func(item, value interface{}) error { item.(*Barcode).Height = ToFloat(value, 0); return nil },
			"Extend":       func(item, value interface{}) error { item.(*Barcode).Extend = ToBoolean(value, false); return nil },
		},
		"separator": {
			"Gap": func(item, value interface{}) error { item.(*Separator).Gap = ToFloat(value, 0); return nil },
		},
		"vgap": {
			"Height":    func(item, value interface{}) error { item.(*VGap).Height = ToFloat(value, 0); return nil },
			"PageBreak": func(item, value interface{}) error { item.(*VGap).PageBreak = ToBoolean(value, false); return nil },
			"Visible":   func(item, value interface{}) error { return nil }, // Deprecated
		},
		"hline": {
			"Width":       func(item, value interface{}) error { item.(*HLine).Width = ToString(value, ""); return nil },
			"Gap":         func(item, value interface{}) error { item.(*HLine).Gap = ToFloat(value, 0); return nil },
			"BorderColor": func(item, value interface{}) error { item.(*HLine).BorderColor = ToRGBA(value, item.(*HLine).BorderColor); return nil },
			"Visible":     func(item, value interface{}) error { return nil }, // Deprecated
		},
		"html": {
			"Fieldname": func(item, value interface{}) error { item.(*HTML).Fieldname = ToString(value, ""); return nil },
			"Value":     func(item, value interface{}) error { item.(*HTML).Value = ToString(value, ""); return nil },
		},
		"datagrid": {
			"Name":             func(item, value interface{}) error { item.(*Datagrid).Name = ToString(value, ""); return nil },
			"Databind":         func(item, value interface{}) error { item.(*Datagrid).Databind = ToString(value, ""); return nil },
			"Width":            func(item, value interface{}) error { item.(*Datagrid).Width = ToString(value, ""); return nil },
			"Merge":            func(item, value interface{}) error { item.(*Datagrid).Merge = ToBoolean(value, false); return nil },
			"Border":           func(item, value interface{}) error { item.(*Datagrid).Border = ToString(value, "1"); return nil },
			"FontSize":         func(item, value interface{}) error { item.(*Datagrid).FontSize = ToFloat(value, item.(*Datagrid).FontSize); return nil },
			"TextColor":       func(item, value interface{}) error { item.(*Datagrid).TextColor = ToRGBA(value, item.(*Datagrid).TextColor); return nil },
			"BorderColor":     func(item, value interface{}) error { item.(*Datagrid).BorderColor = ToRGBA(value, item.(*Datagrid).BorderColor); return nil },
			"BackgroundColor": func(item, value interface{}) error { item.(*Datagrid).BackgroundColor = ToRGBA(value, item.(*Datagrid).BackgroundColor); return nil },
			"HeaderBackground": func(item, value interface{}) error { item.(*Datagrid).HeaderBackground = ToRGBA(value, item.(*Datagrid).HeaderBackground); return nil },
			"FooterBackground": func(item, value interface{}) error { item.(*Datagrid).FooterBackground = ToRGBA(value, item.(*Datagrid).FooterBackground); return nil },
		},
		"column": {
			"Fieldname":   func(item, value interface{}) error { item.(*Column).Fieldname = ToString(value, ""); return nil },
			"Label":       func(item, value interface{}) error { item.(*Column).Label = ToString(value, ""); return nil },
			"Width":       func(item, value interface{}) error { item.(*Column).Width = ToString(value, ""); return nil },
			"Align":       func(item, value interface{}) error { item.(*Column).Align = ToString(value, "L"); return nil },
			"HeaderAlign": func(item, value interface{}) error { item.(*Column).HeaderAlign = ToString(value, "L"); return nil },
			"FooterAlign": func(item, value interface{}) error { item.(*Column).FooterAlign = ToString(value, "L"); return nil },
			"Footer":      func(item, value interface{}) error { item.(*Column).Footer = ToString(value, ""); return nil },
		},
	}
}

func registerGenerators(name string, gen Generator) {
	generators[name] = gen
}

func (pi *PageItem) setPageItem(fieldname string, value interface{}) error {
	propName := propMap[strings.ToLower(fieldname)]
	if setters, ok := elementSetters[pi.ItemType]; ok {
		if setter, found := setters[propName]; found {
			return setter(pi.Item, value)
		}
	}
	return errors.New(invalidErr(pi.ItemType, fieldname))
}

func (rpt *Report) getPageItem(etype string) (PageItem, error) {
	switch etype {
	case "barcode":
		return PageItem{ItemType: etype, Item: &Barcode{CodeType: "CODE_39"}}, nil
	case "cell":
		return PageItem{
			ItemType: etype,
			Item: &Cell{
				Border:          "",
				Align:           "L",
				FontStyle:       rpt.FontStyle,
				FontSize:        rpt.FontSize,
				TextColor:       rpt.TextColor,
				BorderColor:     rpt.BorderColor,
				BackgroundColor: rpt.BackgroundColor,
			},
		}, nil
	case "column":
		return PageItem{
			ItemType: etype,
			Item:     &Column{Align: "L", HeaderAlign: "L", FooterAlign: "L"},
		}, nil
	case "datagrid":
		return PageItem{
			ItemType: etype,
			Item: &Datagrid{
				Border:           "",
				FontSize:         rpt.FontSize,
				TextColor:        rpt.TextColor,
				BorderColor:      rpt.BorderColor,
				BackgroundColor:  rpt.BackgroundColor,
				HeaderBackground: rpt.BackgroundColor,
				FooterBackground: rpt.BackgroundColor,
				Columns:          make([]PageItem, 0),
			},
		}, nil
	case "hline":
		return PageItem{ItemType: etype, Item: &HLine{BorderColor: rpt.BorderColor}}, nil
	case "html":
		return PageItem{ItemType: etype, Item: &HTML{}}, nil
	case "image":
		return PageItem{ItemType: etype, Item: &Image{}}, nil
	case "row":
		return PageItem{ItemType: etype, Item: &Row{Columns: make([]PageItem, 0)}}, nil
	case "separator":
		return PageItem{ItemType: etype, Item: &Separator{}}, nil
	case "vgap":
		return PageItem{ItemType: etype, Item: &VGap{}}, nil
	}
	return PageItem{}, errors.New(invalidErr("Element", etype))
}
