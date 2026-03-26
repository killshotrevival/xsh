package theme

import (
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	BgColor        = "#1e1e2e" // Deep Slate
	FgColor        = "#cdd6f4" // Off-white text
	AccentColor    = "#89b4fa" // Cyber Blue
	SuccessColor   = "#a6e3a1" // Matrix Green
	ErrorColor     = "#f38ba8" // Soft Red
	MutedColor     = "#585b70" // Grey for borders/descriptions
	SelectionColor = "#f5c2e7" // Pink/Purple for highlights
)

func ApplyTviewTheme() {
	// Apply global styles for borders and backgrounds
	tview.Styles.PrimitiveBackgroundColor = tcell.GetColor(BgColor)
	tview.Styles.ContrastBackgroundColor = tcell.GetColor("#313244") // Slightly lighter for contrast
	tview.Styles.BorderColor = tcell.GetColor(MutedColor)
	tview.Styles.TitleColor = tcell.GetColor(AccentColor)
	tview.Styles.PrimaryTextColor = tcell.GetColor(FgColor)
	tview.Styles.SecondaryTextColor = tcell.GetColor(MutedColor)
}

func XSH(isDark bool) *huh.Styles {
	t := huh.ThemeBase(isDark)

	// 1. General Form & Group styling
	t.Form.Base = lipgloss.NewStyle().Padding(1, 2)
	t.Group.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(AccentColor)).Bold(true).MarginBottom(1)
	t.Group.Description = lipgloss.NewStyle().Foreground(lipgloss.Color(MutedColor))

	// 2. Focused Field Styling (When user is typing/selecting)
	t.Focused.Base = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderLeftForeground(lipgloss.Color(AccentColor)).PaddingLeft(1)
	t.Focused.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(AccentColor)).Bold(true)
	t.Focused.Description = lipgloss.NewStyle().Foreground(lipgloss.Color(MutedColor))

	// Selection Indicators
	t.Focused.SelectSelector = lipgloss.NewStyle().Foreground(lipgloss.Color(SelectionColor)).SetString(" ")
	t.Focused.Option = lipgloss.NewStyle().Foreground(lipgloss.Color(FgColor))
	t.Focused.MultiSelectSelector = lipgloss.NewStyle().Foreground(lipgloss.Color(SelectionColor)).SetString("󰄬 ")
	t.Focused.SelectedOption = lipgloss.NewStyle().Foreground(lipgloss.Color(SuccessColor))

	// 3. Error Styling
	t.Focused.ErrorIndicator = lipgloss.NewStyle().Foreground(lipgloss.Color(ErrorColor)).SetString(" ✘")
	t.Focused.ErrorMessage = lipgloss.NewStyle().Foreground(lipgloss.Color(ErrorColor))

	// 4. Blurred State (When field is not active)
	t.Blurred.Base = lipgloss.NewStyle().PaddingLeft(2) // No border
	t.Blurred.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(MutedColor))
	t.Blurred.Description = lipgloss.NewStyle().Foreground(lipgloss.Color("#45475a"))

	return t
}
