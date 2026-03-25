package table

import (
	"github.com/charmbracelet/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	headerColour     = tcell.NewHexColor(0x7571F9)
	cellColour       = tcell.NewHexColor(0xFFFDF5)
	backgroundColour = tcell.NewHexColor(0x1A1B26)
)

type Table struct {
	Headers     []string
	Data        [][]string
	headerStyle tcell.Style
	dataStyle   tcell.Style
}

func NewTable(headers []string, data [][]string) *Table {
	headerStyle := tcell.Style{}.Bold(true).
		Background(backgroundColour)

	dataStyle := tcell.Style{}.
		Background(backgroundColour)

	return &Table{
		Headers:     headers,
		Data:        data,
		headerStyle: headerStyle,
		dataStyle:   dataStyle,
	}
}

func (t *Table) Print() error {
	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)

	for c, value := range t.Headers {
		table.SetCell(0, c,
			tview.NewTableCell(value).SetStyle(t.headerStyle).
				SetTextColor(headerColour).
				SetAlign(tview.AlignCenter))
	}

	for r, row := range t.Data {
		for c, value := range row {

			table.SetCell(r+1, c,
				tview.NewTableCell(value).SetStyle(t.dataStyle).
					SetTextColor(cellColour).
					SetAlign(tview.AlignCenter))
		}
	}

	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		log.Debugf("error occurred while rendering table: %v", err)
		return err
	}

	return nil
}
