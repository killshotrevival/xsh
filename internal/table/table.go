package table

import (
	"github.com/charmbracelet/log"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	headerColour = tcell.NewHexColor(0x6562db)
	cellColour   = tcell.NewHexColor(0xFFFDF5)
)

type Table struct {
	Headers []string
	Data    [][]string
}

func NewTable(headers []string, data [][]string) *Table {
	return &Table{
		Headers: headers,
		Data:    data,
	}
}

func (t *Table) Print() error {
	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)

	for c, value := range t.Headers {
		table.SetCell(0, c,
			tview.NewTableCell(value).
				SetTextColor(headerColour).
				SetAlign(tview.AlignCenter))
	}

	for r, row := range t.Data {
		for c, value := range row {

			table.SetCell(r+1, c,
				tview.NewTableCell(value).
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
