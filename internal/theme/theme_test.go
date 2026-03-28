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
	assert.Equal(t, tview.Styles.TitleColor, tcell.GetColor(AccentColor))
}

func TestXSH(t *testing.T) {
	style := XSH(false)

	assert.Equal(t, style.Group.Title.GetForeground(), lipgloss.Color(AccentColor))
}
