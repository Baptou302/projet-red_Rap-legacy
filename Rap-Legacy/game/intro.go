package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// -----------------
// Intro (style Pokémon)
// -----------------

// updateIntroScene gère la logique de l’intro
func (g *Game) updateIntroScene() {
	if IsKeyJustPressed(ebiten.KeyEnter) {
		g.introStep++
		if g.introStep >= len(g.introTexts) {
			g.state = StateSaveSelect // ou StateCreateSave directement
		}
	}

}

// drawIntroScene affiche l’intro
func (g *Game) drawIntroScene(screen *ebiten.Image) {
	// fond noir
	screen.Fill(color.RGBA{0, 0, 0, 255})

	if g.introStep < len(g.introTexts) {
		text.Draw(screen, g.introTexts[g.introStep], g.fontSmall, 100, 800, color.White)
	}

	// Afficher le sprite du "professeur musique"
	if g.profSprite != nil {
		opts := &ebiten.DrawImageOptions{}
		scale := 0.3 // réduire à 30% de la taille originale
		opts.GeoM.Scale(scale, scale)

		sw, sh := screen.Size()
		w, h := g.profSprite.Size()
		scaledW := float64(w) * scale
		scaledH := float64(h) * scale

		// centré en bas de l’écran
		opts.GeoM.Translate((float64(sw)-scaledW)/2, float64(sh)-scaledH-50)
		screen.DrawImage(g.profSprite, opts)
	}

	// Texte de dialogue
	if g.introStep < len(g.introTexts) {
		msg := g.introTexts[g.introStep]

		if g.fontSmall != nil {
			text.Draw(screen, msg, g.fontSmall, 100, 850, color.White)
		} else {
			ebitenutil.DebugPrintAt(screen, msg, 100, 850)
		}

		// Indicateur "Appuie sur Entrée"
		if g.fontSmall != nil {
			text.Draw(screen, "▶ Appuie sur Entrée", g.fontSmall, 100, 900, color.RGBA{200, 200, 200, 255})
		} else {
			ebitenutil.DebugPrintAt(screen, "Appuie sur Entrée", 100, 900)
		}
	}
}
