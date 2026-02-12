package report

import (
	"image"
	"image/color"
	"io"
)

// IM is a map[string]interface{} type short alias
type IM = map[string]interface{}

// SM is a map[string]string type short alias
type SM = map[string]string

// Generator the PDF generator interface
type Generator interface {
	Init(rpt *Report)
	GetPageSize() (width, height float64)
	PageNo() int
	AddPage()
	AddImage(image *Image, x, y float64, options IM)
	LoadImage(img image.Image, x, y, h, w float64) error
	AddFont(familyStr, styleStr, fileStr string, rd io.Reader)
	GetFontSize() (ptSize float64)
	SetFont(familyStr, styleStr string, size float64)
	SetFontSize(size float64)
	GetTextWidth(s string) float64
	SetDrawColor(r, g, b int)
	SetFillColor(r, g, b int)
	SetTextColor(r, g, b int)
	SetProperties(rpt *Report)
	Text(txtStr string, pageBreak float64)
	Rect(x, y, w, h float64, styleStr string)
	Line(x1, y1, x2, y2 float64)
	GetX() float64
	GetY() float64
	SetX(x float64)
	SetY(y float64)
	SetXY(x, y float64)
	SetText(x, y float64, value string) error
	Ln(h float64)
	Cell(options IM)
	MultiCell(options IM)
	Save2Pdf() ([]byte, error)
	Save2PdfFile(filename string) error
}

// PageItem - element interface wrapper
type PageItem struct {
	ItemType string
	Item     interface{}
}

// Row - Horizontal logical group. The last element width extends up to the right margin.
type Row struct {
	Height  float64    `xml:"height,attr" json:"height"`
	HGap    float64    `xml:"hgap,attr" json:"hgap"`
	Visible string     `xml:"visible,attr" json:"visible"`
	Columns []PageItem `xml:"columns,attr" json:"columns"`
}

// Cell - Row unit
type Cell struct {
	Name            string     `xml:"name,attr" json:"name"`
	Value           string     `xml:"value,attr" json:"value"`
	Width           string     `xml:"width,attr" json:"width"`
	Border          string     `xml:"border,attr" json:"border"`
	Align           string     `xml:"align,attr" json:"align"`
	Multiline       bool       `xml:"multiline,attr" json:"multiline"`
	FontStyle       string     `xml:"font-style,attr" json:"font-style"`
	FontSize        float64    `xml:"font-size,attr" json:"font-size"`
	TextColor       color.RGBA `xml:"color,attr" json:"color"`
	BorderColor     color.RGBA `xml:"border-color,attr" json:"border-color"`
	BackgroundColor color.RGBA `xml:"background-color,attr" json:"background-color"`
}

// Image - Row unit
type Image struct {
	Src       string  `xml:"src,attr" json:"src"`
	Data      []byte  `xml:"data,attr" json:"data"`
	MaxWidth  float64 `xml:"max-width,attr" json:"max-width"`
	MaxHeight float64 `xml:"max-height,attr" json:"max-height"`
	Height    float64 `xml:"height,attr" json:"height"`
	Width     float64 `xml:"width,attr" json:"width"`
}

// Barcode - Row unit
type Barcode struct {
	CodeType     string  `xml:"code-type,attr" json:"code-type"`
	Value        string  `xml:"value,attr" json:"value"`
	VisibleValue bool    `xml:"visible-value,attr" json:"visible-value"`
	Width        float64 `xml:"wide,attr" json:"wide"`
	Height       float64 `xml:"narrow,attr" json:"narrow"`
	Extend       bool    `xml:"extend,attr" json:"extend"`
}

// Separator - Row unit, A horizontal separator line.
type Separator struct {
	Gap float64 `xml:"gap,attr" json:"gap"`
}

// VGap - a vertical gap.
type VGap struct {
	Height    float64 `xml:"height,attr" json:"height"`
	PageBreak bool    `xml:"page-break,attr" json:"page-break"`
}

// HLine - a horizontal line.
type HLine struct {
	Width       string     `xml:"width,attr" json:"width"`
	Gap         float64    `xml:"gap,attr" json:"gap"`
	BorderColor color.RGBA `xml:"border-color,attr" json:"border-color"`
}

// HTML - a basic HTML elements rendering.
type HTML struct {
	Fieldname string `xml:"fieldname,attr" json:"fieldname"`
	Value     string `xml:",cdata" json:"html"`
}

// Datagrid - Create a table from a data list.
type Datagrid struct {
	Name             string     `xml:"name,attr" json:"name"`
	Databind         string     `xml:"databind,attr" json:"databind"`
	Width            string     `xml:"width,attr" json:"width"`
	Merge            bool       `xml:"merge,attr" json:"merge"`
	Border           string     `xml:"border,attr" json:"border"`
	FontSize         float64    `xml:"font-size,attr" json:"font-size"`
	TextColor        color.RGBA `xml:"color,attr" json:"color"`
	BorderColor      color.RGBA `xml:"border-color,attr" json:"border-color"`
	BackgroundColor  color.RGBA `xml:"background-color,attr" json:"background-color"`
	HeaderBackground color.RGBA `xml:"header-background,attr" json:"header-background"`
	FooterBackground color.RGBA `xml:"footer-background,attr" json:"footer-background"`
	Columns          []PageItem `xml:"columns" json:"columns"`
}

// Column - Datagrid unit
type Column struct {
	Fieldname   string `xml:"fieldname,attr" json:"fieldname"`
	Label       string `xml:"label,attr" json:"label"`
	Width       string `xml:"width,attr" json:"width"`
	Align       string `xml:"align,attr" json:"align"`
	HeaderAlign string `xml:"header-align,attr" json:"header-align"`
	FooterAlign string `xml:"footer-align,attr" json:"footer-align"`
	Footer      string `xml:"footer,attr" json:"footer"`
}

// Report is the principal structure for creating a single PDF document
type Report struct {
	Pdf                                                 Generator
	orientation, format, fontDir, xmlHeader, xmlDetails string
	header, details, footer                             []PageItem
	data                                                IM
	footerHeight, pageBreak                             float64
	Title                                               string     `xml:"title,attr" json:"title"`
	Author                                              string     `xml:"author,attr" json:"author"`
	Creator                                             string     `xml:"creator,attr" json:"creator"`
	Subject                                             string     `xml:"subject,attr" json:"subject"`
	Keywords                                            string     `xml:"keywords,attr" json:"keywords"`
	LeftMargin                                          float64    `xml:"left-margin,attr" json:"left-margin"`
	RightMargin                                         float64    `xml:"right-margin,attr" json:"right-margin"`
	TopMargin                                           float64    `xml:"top-margin,attr" json:"top-margin"`
	BottomMargin                                        float64    `xml:"bottom-margin,attr" json:"bottom-margin"`
	FontFamily                                          string     `xml:"font-family,attr" json:"font-family"`
	FontStyle                                           string     `xml:"font-style,attr" json:"font-style"`
	FontSize                                            float64    `xml:"font-size,attr" json:"font-size"`
	TextColor                                           color.RGBA `xml:"color,attr" json:"color"`
	BorderColor                                         color.RGBA `xml:"border-color,attr" json:"border-color"`
	BackgroundColor                                     color.RGBA `xml:"background-color,attr" json:"background-color"`
	ImagePath                                           string     `xml:"image-path,attr" json:"image-path"`
}
