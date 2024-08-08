package main

import (
	"log"
	"os"
	"path"

	rp "github.com/nervatura/report"
)

func main() {
	var err error

	// Creating a PDF from a JSON template
	json, _ := os.ReadFile(path.Join("sample.json"))
	rpt := rp.New("L")
	if err = rpt.LoadJSONDefinition(string(json)); err != nil {
		log.Fatal(err)
		return
	}
	rpt.CreateReport()

	var pdf []byte
	if pdf, err = rpt.Save2Pdf(); err == nil {
		err = os.WriteFile("out/json.pdf", pdf, 0644)
	}
	if err != nil {
		log.Fatal(err)
		return
	}

	var base64Str string
	if base64Str, err = rpt.Save2DataURLString("Report.pdf"); err == nil {
		err = os.WriteFile("out/base64.txt", []byte(base64Str), 0644)
	}
	if err != nil {
		log.Fatal(err)
		return
	}

	// PDF creation from GO language code
	rpt = createGoReport()
	if err := rpt.Save2PdfFile("out/go.pdf"); err != nil {
		log.Fatal(err)
		return
	}

	dataXml := rpt.Save2Xml()
	if err := os.WriteFile("out/data.xml", []byte(dataXml), 0644); err != nil {
		log.Fatal(err)
		return
	}
}
