package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Inventaire struct {
	Ouvert       bool
	Items        []string
	tabPrecedent bool
	escPrecedent bool
}

// Update vérifie les touches Tab / Échap
func (inv *Inventaire) Update() {
	// Toggle avec Tab
	if ebiten.IsKeyPressed(ebiten.KeyTab) && !inv.tabPrecedent {
		inv.Ouvert = !inv.Ouvert
	}
	inv.tabPrecedent = ebiten.IsKeyPressed(ebiten.KeyTab)

	// Fermer avec Échap
	if inv.Ouvert && ebiten.IsKeyPressed(ebiten.KeyEscape) && !inv.escPrecedent {
		inv.Ouvert = false
	}
	inv.escPrecedent = ebiten.IsKeyPressed(ebiten.KeyEscape)
}

// Draw affiche l'inventaire
func (inv *Inventaire) Draw(screen *ebiten.Image) {
	if inv.Ouvert {
		// Fond noir semi-transparent plein écran
		w, h := screen.Size()
		ebitenutil.DrawRect(screen, 0, 0, float64(w), float64(h), color.RGBA{0, 0, 0, 200})

		// Affichage des items
		ebitenutil.DebugPrintAt(screen, "=== INVENTAIRE === (Échap pour fermer)", 50, 50)
		for i, item := range inv.Items {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d. %s", i+1, item), 80, 100+30*i)
		}
	} else {
		// Petit message aide
		ebitenutil.DebugPrint(screen, "Appuie sur TAB pour ouvrir l'inventaire")
	}
}

func (inv *Inventaire) BindToGame(g *Game) {
	if g != nil {
		g.Inventaire = inv
	}
}
