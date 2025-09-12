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
	// coordonnées et dimensions des boutons
	buttonX, buttonY, buttonW, buttonH int
}

func NewMenu() *Menu {
	return &Menu{
		menuOptions: []string{"New Game", "Options", "Quit"},
		buttonX:     700, // position X du premier bouton
		buttonY:     400, // position Y du premier bouton
		buttonW:     500, // largeur bouton
		buttonH:     80,  // hauteur bouton
	}
}

func (m *Menu) Update() string {
	// --- Navigation clavier (optionnel si tu veux garder) ---
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
		case 1:
			return "settings"
		case 2:
			os.Exit(0)
		}
	}

	// --- Détection clic souris ---
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		for i := range m.menuOptions {
			bx := m.buttonX
			by := m.buttonY + i*(m.buttonH+20)
			bw := m.buttonW
			bh := m.buttonH

			if x >= bx && x <= bx+bw && y >= by && y <= by+bh {
				switch i {
				case 0:
					return "play"
				case 1:
					return "settings"
				case 2:
					os.Exit(0)
				}
			}
		}
	}

	return ""
}

func (m *Menu) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})
	for i, option := range m.menuOptions {
		x := m.buttonX
		y := m.buttonY + i*(m.buttonH+20)

		text := option
		if i == m.menuSelected {
			text = "> " + option
		}

		// texte affiché
		ebitenutil.DebugPrintAt(screen, text, x+20, y+20)
	}
}
