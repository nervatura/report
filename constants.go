package report

import "fmt"

const (
	_generator       = "gopdf"
	_title           = "Nervatura Report"
	_margin          = float64(36.85)
	_fontFamily      = "Cabin"
	_fontStyle       = ""
	_fontSize        = float64(9)
	_textColor       = uint8(0)
	_borderColor     = uint8(0)
	_backgroundColor = uint8(255)
	_padding         = float64(6.4)
	_format          = "a4"
	_orientation     = "p"
	//_unit            = "pt"
	_mmPt  = 2.83465
	_align = "L"
	//_fontDir         = ""
	_regValue = "={{(\\S*?)[^}}]*}}.*?|={{.*? /}}"
)

var propMap SM = SM{
	"height": "Height", "hgap": "HGap", "visible": "Visible", "name": "Name", "value": "Value", "width": "Width",
	"border": "Border", "align": "Align", "multiline": "Multiline", "font-style": "FontStyle", "fontstyle": "FontStyle",
	"font-size": "FontSize", "fontsize": "FontSize", "color": "TextColor", "textcolor": "TextColor",
	"border-color": "BorderColor", "bordercolor": "BorderColor",
	"background-color": "BackgroundColor", "backgroundcolor": "BackgroundColor",
	"src": "Src", "code-type": "CodeType", "codetype": "CodeType",
	"visible-value": "VisibleValue", "visiblevalue": "VisibleValue", "wide": "Width", "narrow": "Height",
	"extend": "Extend", "gap": "Gap", "fieldname": "Fieldname", "html": "Value", "databind": "Databind",
	"merge": "Merge", "header-background": "HeaderBackground", "headerbackground": "HeaderBackground",
	"footer-background": "FooterBackground", "footerbackground": "FooterBackground", "label": "Label",
	"header-align": "HeaderAlign", "headeralign": "HeaderAlign",
	"footer-align": "FooterAlign", "footeralign": "FooterAlign", "footer": "Footer",
	"title": "Title", "author": "Author", "creator": "Creator",
	"subject": "Subject", "keywords": "Keywords", "leftmargin": "LeftMargin", "left-margin": "LeftMargin",
	"topmargin": "TopMargin", "top-margin": "TopMargin", "rightmargin": "RightMargin", "right-margin": "RightMargin",
	"bottommargin": "BottomMargin", "bottom-margin": "BottomMargin",
	"imagepath": "ImagePath", "image-path": "ImagePath",
	"fontfamily": "FontFamily", "font-family": "FontFamily",
	"page-break": "PageBreak",
}

func invalidErr(etype, evalue string) string {
	return fmt.Sprintf("invalid %s element: %s", etype, evalue)
}
