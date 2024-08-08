package main

import (
	"log"

	rp "github.com/nervatura/report"
)

func createGoReport() (rpt *rp.Report) {
	var appendElement = func(parent interface{}, ename string, values rp.IM) *[]rp.PageItem {
		eParent, err := rpt.AppendElement(parent, ename, values)
		if err != nil {
			log.Fatal(err)
		}
		return eParent
	}
	var setData = func(key string, value interface{}) {
		if _, err := rpt.SetData(key, value); err != nil {
			log.Fatal(err)
		}
	}

	//rpt = New("p", "A4", "NotoSans", "../../data/fonts")
	rpt = rp.New("p", "A4")

	//default values
	rpt.SetReportValue("title", "Go Report")
	rpt.SetReportValue("fontSize", float64(9))

	//header
	rowData := appendElement("header", "row", rp.IM{"height": 10})
	appendElement(rowData, "image", rp.IM{"src": "logo"})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "value": "labels.title",
		"font-style": "bolditalic", "font-size": 26, "color": "#D8DBDA"})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "value": "Go Sample", "font-style": "bold", "align": "right"})
	appendElement("header", "vgap", rp.IM{"height": 2})
	appendElement("header", "hline", rp.IM{"border-color": 218})
	appendElement("header", "vgap", rp.IM{"height": 2})

	//details
	appendElement("details", "vgap", rp.IM{"height": 2})
	rowData = appendElement("details", "row", rp.IM{})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "width": "50%", "font-style": "bold", "value": "labels.left_text", "border": "LT",
		"border-color": 218, "background-color": 245})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "font-style": "bold", "value": "labels.left_text", "border": "LTR",
		"border-color": 218, "background-color": 245})

	rowData = appendElement("details", "row", rp.IM{})
	appendElement(rowData, "cell", rp.IM{
		"name": "short_text", "width": "50%", "value": "head.short_text", "border": "L", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "short_text", "value": "head.short_text", "border": "LR", "border-color": 218})
	rowData = appendElement("details", "row", rp.IM{})
	appendElement(rowData, "cell", rp.IM{
		"name": "short_text", "width": "50%", "value": "head.short_text", "border": "LB", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "short_text", "value": "head.short_text", "border": "LBR", "border-color": 218})

	rowData = appendElement("details", "row", rp.IM{})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "width": "40", "font-style": "bold", "value": "labels.left_text", "border": "LB", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "align": "center", "width": "30", "font-style": "bold", "value": "labels.center_text",
		"border": "LB", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "align": "right", "width": "40", "font-style": "bold", "value": "labels.right_text",
		"border": "LB", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "font-style": "bold", "value": "labels.left_text", "border": "LBR", "border-color": 218})

	rowData = appendElement("details", "row", rp.IM{})
	appendElement(rowData, "cell", rp.IM{
		"name": "short_text", "width": "40", "value": "head.short_text", "border": "LB", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "date", "align": "center", "width": "30", "value": "head.date", "border": "LB", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "amount", "align": "right", "width": "40", "value": "head.number", "border": "LB", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "short_text", "value": "head.short_text", "border": "LBR", "border-color": 218})

	rowData = appendElement("details", "row", rp.IM{})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "font-style": "bold", "value": "labels.left_text", "border": "LB", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "short_text", "width": "50", "value": "head.short_text", "border": "LB", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "font-style": "bold", "value": "labels.left_text", "border": "LB", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "short_text", "value": "head.short_text", "border": "LBR", "border-color": 218})

	rowData = appendElement("details", "row", rp.IM{})
	appendElement(rowData, "cell", rp.IM{
		"name": "long_text", "multiline": true, "value": "head.long_text",
		"border": "LBR", "border-color": 218})

	appendElement("details", "vgap", rp.IM{"height": 2})
	rowData = appendElement("details", "row", rp.IM{"hgap": 2})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "value": "labels.left_text", "font-style": "bold", "border": "1", "border-color": 218,
		"background-color": 245})
	appendElement(rowData, "cell", rp.IM{
		"name": "short_text", "value": "head.short_text", "border": "1", "border-color": 218})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "value": "labels.left_text", "font-style": "bold", "border": "1", "border-color": 218, "background-color": 245})
	appendElement(rowData, "cell", rp.IM{
		"name": "short_text", "value": "head.short_text", "border": "1", "border-color": 218})

	appendElement("details", "vgap", rp.IM{"height": 2})
	rowData = appendElement("details", "row", rp.IM{"hgap": 2})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "value": "labels.long_text", "font-style": "bold", "border": "1", "border-color": 218, "background-color": 245})
	appendElement(rowData, "cell", rp.IM{
		"name": "long_text", "multiline": true, "value": "head.long_text", "border": "1", "border-color": 218})

	appendElement("details", "vgap", rp.IM{"height": 2})
	appendElement("details", "hline", rp.IM{"border-color": 218})
	appendElement("details", "vgap", rp.IM{"height": 2})

	rowData = appendElement("details", "row", rp.IM{"hgap": 5})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "value": "Barcode (code 128)", "font-style": "bold", "font-size": 9, "width": "40",
		"border": "1", "border-color": 245, "background-color": 245})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "value": "Barcode (Interleaved 2of5)", "font-style": "bold", "font-size": 9, "width": "40",
		"border": "1", "border-color": 245, "background-color": 245})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "value": "Barcode (EAN)", "font-style": "bold", "font-size": 9, "width": "40",
		"border": "1", "border-color": 245, "background-color": 245})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "value": "Barcode (Code 39)", "font-style": "bold", "font-size": 9, "width": "40",
		"border": "1", "border-color": 245, "background-color": 245})

	rowData = appendElement("details", "row", rp.IM{"hgap": 5})
	appendElement(rowData, "barcode",
		rp.IM{"code-type": "CODE_128", "value": "1234567890ABCDEF", "width": 40, "visible-value": 1})
	appendElement(rowData, "barcode",
		rp.IM{"code-type": "ITF", "value": "1234567890", "width": 40, "visible-value": 1})
	appendElement(rowData, "barcode",
		rp.IM{"code-type": "EAN", "value": "96385074", "width": 40, "visible-value": 1})
	appendElement(rowData, "barcode",
		rp.IM{"code-type": "CODE_39", "value": "1234567890ABCDEF", "width": 40, "extend": true, "visible-value": 1})

	appendElement("details", "vgap", rp.IM{"height": 3})

	rowData = appendElement("details", "row", rp.IM{"hgap": 5})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "value": "QR code: Hello Go Report!", "font-size": 9,
		"border": "1", "border-color": 245, "background-color": 245})
	appendElement(rowData, "barcode",
		rp.IM{"code-type": "QR", "value": "Hello Go Report!", "height": 10})
	appendElement(rowData, "cell", rp.IM{})
	appendElement("details", "vgap", rp.IM{"height": 0, "page-break": true})

	rowData = appendElement("details", "row", rp.IM{})
	appendElement(rowData, "cell", rp.IM{
		"name": "label", "value": "Datagrid Sample", "align": "center", "font-style": "bold",
		"border": "1", "border-color": 245, "background-color": 245})
	appendElement("details", "vgap", rp.IM{"height": 2})

	var gridData = appendElement("details", "datagrid", rp.IM{
		"name": "items", "databind": "items", "border": "1", "border-color": 218, "header-background": 245, "footer-background": 245})
	appendElement(gridData, "column", rp.IM{
		"width": "8%", "fieldname": "counter", "align": "right", "label": "labels.counter", "footer": "labels.total"})
	appendElement(gridData, "column", rp.IM{
		"width": "20%", "fieldname": "date", "align": "center", "label": "labels.center_text"})
	appendElement(gridData, "column", rp.IM{
		"width": "15%", "fieldname": "number", "align": "right", "label": "labels.right_text",
		"footer": "items_footer.items_total", "footer-align": "right"})
	appendElement(gridData, "column", rp.IM{
		"fieldname": "text", "label": "labels.left_text"})

	appendElement("details", "vgap", rp.IM{"height": 5})
	appendElement("details", "html", rp.IM{"fieldname": "html_text",
		"html": "<i>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</i> ={{html_text}} <p>Nulla a <b><i>pretium</i></b> nunc, in <u>cursus</u> quam.</p>"})

	//footer
	appendElement("footer", "vgap", rp.IM{"height": 2})
	appendElement("footer", "hline", rp.IM{"border-color": 218})
	rowData = appendElement("footer", "row", rp.IM{"height": 10})
	appendElement(rowData, "cell", rp.IM{"value": "Nervatura Report Template", "font-style": "bolditalic"})
	appendElement(rowData, "cell", rp.IM{"value": "{{page}}", "align": "right", "font-style": "bold"})

	//data
	setData("labels", rp.SM{"title": "REPORT TEMPLATE", "left_text": "Short text", "center_text": "Centered text",
		"right_text": "Right text", "long_text": "Long text", "counter": "No.", "total": "Total"})
	setData("head", rp.SM{"short_text": "Lorem éáőűúóüö dolor", "number": "123 456", "date": "2015.01.01",
		"long_text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque eu mattis diam, sed dapibus justo. In eget augue nisi. Cras eget odio vel mi vulputate interdum. Curabitur consequat sapien at lacus tincidunt, at sagittis felis lobortis. Aenean porta maximus quam eu porta. Fusce sed leo ut justo commodo facilisis. Vivamus vitae tempor erat, at ultrices enim. Nulla a pretium nunc, in cursus quam."})
	setData("html_text", "<p><b>Pellentesque eu mattis diam, sed dapibus justo. In eget augue nisi. Cras eget odio vel mi vulputate interdum. Curabitur consequat sapien at lacus tincidunt, at sagittis felis lobortis. Aenean porta maximus quam eu porta. Fusce sed leo ut justo commodo facilisis. Vivamus vitae tempor erat, at ultrices enim.</b></p>")
	setData("items_footer", rp.SM{"items_total": "3 703 680"})
	setData("logo", "data:image/jpg;base64,/9j/4AAQSkZJRgABAQIA7ADsAAD/2wBDAAoHBwgHBgoICAgLCgoLDhgQDg0NDh0VFhEYIx8lJCIfIiEmKzcvJik0KSEiMEExNDk7Pj4+JS5ESUM8SDc9Pjv/2wBDAQoLCw4NDhwQEBw7KCIoOzs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozv/wgARCABAAEADAREAAhEBAxEB/8QAGwABAAIDAQEAAAAAAAAAAAAAAAQGAQMFAgf/xAAYAQEBAQEBAAAAAAAAAAAAAAAAAQIDBP/aAAwDAQACEAMQAAABuYAAAAAAPCQZno3eCBMc+YjyRJnmzndenq6F3VMefg55Abl2W37p68Hznl4sHolXcWYs2+9h12wUPn5MGCLM9C7u/T1ejmzEaZr2eOiTr66WfXfbaAAAAP/EACMQAAICAgEDBQEAAAAAAAAAAAIDAQQABRATITAREhQVIjH/2gAIAQEAAQUC8xEIRXd8hvDLtdWFtl5O2LC2VksMzZNJfSq5sLkkfMh7QrK61jJ7D/eIiSn8oyZkp1aPQeG1yW/prHCd2ypVKyyIgYy1Fks+tCVMpvVMKZOV9YZYtYqDxf/EACERAAEDAwUBAQAAAAAAAAAAAAEAAhEQE0EDEiAwUSEy/9oACAEDAQE/Ae8Gal4CuhXVcciZTBApqPwOEQmiTy/NNMZqW/YUDK3eUY3dV27CtiEWEKCm6fqAjr//xAAbEQACAgMBAAAAAAAAAAAAAAABEQAQEiAwQP/aAAgBAgEBPwHxOZTKOhROo4C1oBZiipdP/8QAKRAAAQMCBAUEAwAAAAAAAAAAAQACEQMhEDFBYRIiMHGRICNRoTJCYv/aAAgBAQAGPwLrS4gDdPe38G2G+MF8n4C5abj3VqQ8qxDewUvcXd0wam5wNFhhoz39EuzOQTGb3wJxgCStHVPpqkmSjWOthi6l4k6L3Ks7MuuGm3gb9nD+BmUALAYNFB0A5ohziah/cqDTJ3F1am7wprco+NVwMEAdP//EACUQAQACAAUEAQUAAAAAAAAAAAEAERAhMUFRMGFxkdEgobHw8f/aAAgBAQABPyHrOTrdVBvNQ/Ld/GLlPyMF90BNn/N4bX69vLQTlXL0KD3YJyybN30VyV5muWcEtvDfB9gIqlW1wMuTQJWt2evyMcuTNWMdz9XDWU0Ms0APK2bd7O3vSUMzZNN+RwoGZ/Chl0KA2wrmyj3JrMv3JVUcNItSXgcXG7fV8QIHbHT/AP/aAAwDAQACAAMAAAAQAAAAAAANEty9DknkCjioltiArqvAAAAA/8QAIBEBAAICAQQDAAAAAAAAAAAAAQARECExMEFRYSCRsf/aAAgBAwEBPxDrKG2X1OM8kx7BL+ItE5MoTCL8Bqt5laZW94BWia9n8iq2yst3y4xK+X1HShRhV6gVowFRNiF35iXE9cV3ACjp/wD/xAAcEQACAgIDAAAAAAAAAAAAAAABEQAQIEEwMVH/2gAIAQIBAT8Q5wXZARKPCSYCFahgoDOXVDuyhUQHcflMttRVCQiMHuALj//EACUQAQABAwMDBQEBAAAAAAAAAAERACExQWGBEHGhMFGRscEg0f/aAAgBAQABPxD1svViA+auiElpLweHV2B5GZspY5aAZ7T8WaeMjpN9AqIC2f8AVbp1cO04rVUIIZuv2IOOi6wWIdY7GN/4ljjsNngDoZc4iXhJN2N/AjmgAgsU4GUT4pyyFVyvR9TwOVe1QZRsCHb6bWDekkDISrUKcbppbvLBx0QCJI2aOgLdIIkBj9ofKRqE5wOFqEJMI/nTYg26G0Kjt+zd8Zo+oglgYOjrVq4PZZdM4vTGblwoNInHvN6ZaDOXkxzFaDikH6rUecg/h5O1BLC32Lq7+n//2Q==")
	items := make([]rp.SM, 0)
	items = append(items, rp.SM{"text": "Lorem ipsum dolor", "number": "123 456", "date": "Lorem ipsum dolorjkjkjl jhkjh"})
	for index := 0; index < 20; index++ {
		items = append(items, rp.SM{"text": "Lorem ipsum dolor", "number": "123 456", "date": "2015.01.01"})
	}
	items = append(items, rp.SM{"text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque eu mattis diam, sed dapibus justo. In eget augue nisi. Cras eget odio vel mi vulputate interdum. Curabitur consequat sapien at lacus tincidunt, at sagittis felis lobortis. Aenean porta maximus quam eu porta. Fusce sed leo ut justo commodo facilisis. Vivamus vitae tempor erat, at ultrices enim. Nulla a pretium nunc, in cursus quam.", "number": "123 456", "date": "2015.01.01"})
	for index := 0; index < 20; index++ {
		items = append(items, rp.SM{"text": "Lorem ipsum dolor", "number": "123 456", "date": "2015.01.01"})
	}
	setData("items", items)

	rpt.CreateReport()
	return rpt
}
