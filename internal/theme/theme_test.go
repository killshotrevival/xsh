package theme

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestApplyTheme(t *testing.T) {
	ApplyTviewTheme()
	assert.Equal(t, tcell.GetColor(BgColor), tview.Styles.TitleColor)
}

func TestXSH(t *testing.T) {
	style := XSH(false)

	assert.Equal(t, lipgloss.Color(AccentColor), style.Group.Title.GetForeground())
}
