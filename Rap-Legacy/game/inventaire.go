package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// -----------------
// Inventaire struct
// -----------------
type Inventaire struct {
	Open     bool
	Bg       *ebiten.Image
	Items    []string
	selected int
}

// -----------------
// NewInventaire
// -----------------
func NewInventaire() *Inventaire {
	return &Inventaire{
		Open:  false,
		Bg:    LoadImage("assets/inventaire.png"),
		Items: []string{"Téléphone", "JNR - 9k", "Cristalline - mystérieuse"},
	}
}

// -----------------
// Update
// -----------------
func (inv *Inventaire) Update() {
	// Rien de spécial pour l’instant (navigation possible plus tard)
}

// -----------------
// Draw
// -----------------
func (inv *Inventaire) Draw(screen *ebiten.Image) {
	if !inv.Open {
		return
	}

	// Dessiner le fond
	if inv.Bg != nil {
		opts := &ebiten.DrawImageOptions{}
		screen.DrawImage(inv.Bg, opts)
	} else {
		// Si l'image est manquante, écran noir par défaut
		screen.Fill(color.RGBA{0, 0, 0, 200})
	}

	// Affichage de la liste des objets au centre
	startX, startY := 800, 400
	for i, item := range inv.Items {
		ebitenutil.DebugPrintAt(screen, item, startX, startY+(i*30))
	}
}

// -----------------
// DrawNote (toujours affichée en haut à gauche)
// -----------------
func (inv *Inventaire) DrawNote(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Appuie sur [TAB] pour ouvrir lA FAUSSSE sacoche Gucci", 20, 20)
}
