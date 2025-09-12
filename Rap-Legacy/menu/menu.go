package menu

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Menu struct {
	menuSelected int
	menuOptions  []string
}

func NewMenu() *Menu {
	return &Menu{menuOptions: []string{"Start Game", "Settings", "Quit"}}
}

func (m *Menu) Update() string {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		m.menuSelected--
		if m.menuSelected < 0 {
			m.menuSelected = len(m.menuOptions) - 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		m.menuSelected++
		if m.menuSelected >= len(m.menuOptions) {
			m.menuSelected = 0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		switch m.menuSelected {
		case 0:
			return "play"
		case 1: // settings
		case 2:
			os.Exit(0)
		}
	}
	return ""
}

func (m *Menu) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})
	for i, option := range m.menuOptions {
		x, y := 400.0, 200.0+float64(i*60)
		text := option
		if i == m.menuSelected {
			text = "> " + option
		}
		ebitenutil.DebugPrintAt(screen, text, int(x), int(y))
	}
}
