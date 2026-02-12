package report

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (rpt *Report) encodeImage(data string, v *Image) {
	rawImage := string(data)[strings.Index(string(data), ",")+1:]
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(rawImage))
	m, iFormat, err := image.Decode(reader)
	if err == nil {
		v.MaxWidth = float64(m.Bounds().Max.X)
		v.MaxHeight = float64(m.Bounds().Max.Y)
		buf := new(bytes.Buffer)
		switch iFormat {
		case "jpeg":
			err = jpeg.Encode(buf, m, &jpeg.Options{Quality: 85})
		case "png":
			err = png.Encode(buf, m)
		}
		if err == nil && len(buf.Bytes()) > 0 {
			v.Data = buf.Bytes()
		}
	}
}

func (rpt *Report) setImageSize(v *Image) {
	src := v.Src
	if rpt.ImagePath != "" {
		src = path.Join(rpt.ImagePath, v.Src)
	}
	reader, err := os.Open(filepath.Clean(src))
	if err == nil {
		defer reader.Close()
		m, _, err := image.Decode(reader)
		if err == nil {
			v.MaxWidth = float64(m.Bounds().Max.X)
			v.MaxHeight = float64(m.Bounds().Max.Y)
		}
	} else {
		v.Src = ""
	}
}

func (rpt *Report) drawImage(v *Image) {
	if v.Width == 0 {
		v.Width = v.Height * (v.MaxWidth / v.MaxHeight)
	}
	if rpt.checkPageBreak(v.Height) {
		rpt.addPage()
	}
	rpt.Pdf.AddImage(v, rpt.Pdf.GetX(), rpt.Pdf.GetY(), IM{"ImagePath": rpt.ImagePath})
}

func (rpt *Report) createImage(v *Image, rowHeight float64, virtual bool) (float64, float64) {
	data := rpt.setValue(v.Src)
	if strings.HasPrefix(data, "data:image") && v.Data == nil {
		rpt.encodeImage(data, v)
	}
	if v.Data == nil && v.Width == 0 {
		rpt.setImageSize(v)
	}
	if v.Data != nil || v.Src != "" {
		height := float64(-1)
		if height <= 0 && v.Height > 0 {
			height = v.Height
			if height > rowHeight {
				rowHeight = height
			}
		}
		if height <= 0 && rowHeight > 0 {
			height = rowHeight - _padding/3
		}
		v.Height = height
		if !virtual {
			rpt.drawImage(v)
		}
		if height > rowHeight || rowHeight == 0 {
			rowHeight = height
		}
	}
	return rowHeight, v.Width
}
