package report

// elementRenderer is an internal interface for rendering a single element.
// It allows new element types to be added without modifying createElement's switch.
// Current code paths continue to use createElement directly; this is for future extensibility.
//type elementRenderer interface {
//	Render(rpt *Report, element interface{}, section string, virtual bool) (height float64, err error)
//}

// rowRenderer renders a Row element.
type rowRenderer struct{}

func (rowRenderer) Render(rpt *Report, element interface{}, section string, virtual bool) (float64, error) {
	return rpt.createRow(section, element.(*Row), virtual), nil
}

// datagridRenderer renders a Datagrid element.
type datagridRenderer struct{}

func (datagridRenderer) Render(rpt *Report, element interface{}, section string, virtual bool) (float64, error) {
	if rpt.createDatagrid(element.(*Datagrid), virtual) {
		return 0, nil
	}
	return 0, nil
}

// lineRenderer renders an HLine element.
type lineRenderer struct{}

func (lineRenderer) Render(rpt *Report, element interface{}, section string, virtual bool) (float64, error) {
	rpt.createLine(element.(*HLine), virtual)
	return 0, nil
}

// htmlRenderer renders an HTML element.
type htmlRenderer struct{}

func (htmlRenderer) Render(rpt *Report, element interface{}, section string, virtual bool) (float64, error) {
	rpt.createHTML(element.(*HTML))
	return 0, nil
}
